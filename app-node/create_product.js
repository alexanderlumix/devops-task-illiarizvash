const { MongoClient } = require('mongodb');
const express = require('express');
const winston = require('winston');
const { body, validationResult } = require('express-validator');
const rateLimit = require('express-rate-limit');

// Configure structured logging
const logger = winston.createLogger({
  level: process.env.LOG_LEVEL || 'info',
  format: winston.format.combine(
    winston.format.timestamp(),
    winston.format.errors({ stack: true }),
    winston.format.json()
  ),
  defaultMeta: { 
    service: 'node-app',
    version: '1.0.0',
    environment: process.env.NODE_ENV || 'development'
  },
  transports: [
    new winston.transports.Console({
      format: winston.format.combine(
        winston.format.colorize(),
        winston.format.simple()
      )
    })
  ]
});

// Configure rate limiting
const limiter = rateLimit({
  windowMs: 15 * 60 * 1000, // 15 minutes
  max: 100, // limit each IP to 100 requests per windowMs
  message: {
    error: 'Too many requests from this IP, please try again later.',
    retryAfter: '15 minutes'
  },
  standardHeaders: true,
  legacyHeaders: false,
});

// getMongoURI constructs MongoDB connection URI from environment variables
function getMongoURI() {
  const host = process.env.MONGO_HOST || 'mongo-1';
  const port = process.env.MONGO_PORT || '27017';
  const db = process.env.MONGO_DB || 'appdb';
  const directConnection = process.env.MONGO_DIRECT_CONNECTION || 'true';

  // Connect without authentication for development
  const uri = `mongodb://${host}:${port}/${db}?directConnection=${directConnection}`;
  
  logger.info('MongoDB URI constructed', {
    host,
    port,
    database: db,
    directConnection
  });

  return uri;
}

// MongoDB connection URI for direct connection to replica set member
const uri = getMongoURI();

// Create Express app for health checks
const app = express();
const PORT = process.env.PORT || 3000;

// Apply rate limiting to all requests
app.use(limiter);

// Add request logging middleware
app.use((req, res, next) => {
  logger.info('HTTP request', {
    method: req.method,
    url: req.url,
    userAgent: req.get('User-Agent'),
    remoteAddress: req.ip
  });
  next();
});

// Health check endpoint
app.get('/health', (req, res) => {
  const healthData = {
    status: 'healthy',
    service: 'node-app',
    timestamp: new Date().toISOString(),
    version: '1.0.0',
    environment: process.env.NODE_ENV || 'development'
  };
  
  logger.info('Health check requested', {
    remoteAddress: req.ip,
    userAgent: req.get('User-Agent')
  });
  
  res.json(healthData);
});

// Product creation endpoint with validation
app.post('/products', [
  body('name')
    .isLength({ min: 1, max: 100 })
    .trim()
    .escape()
    .withMessage('Product name must be between 1 and 100 characters'),
  body('price')
    .optional()
    .isNumeric()
    .withMessage('Price must be a number'),
  body('description')
    .optional()
    .isLength({ max: 500 })
    .trim()
    .escape()
    .withMessage('Description must be less than 500 characters')
], (req, res) => {
  // Check for validation errors
  const errors = validationResult(req);
  if (!errors.isEmpty()) {
    logger.warn('Validation failed', {
      errors: errors.array(),
      remoteAddress: req.ip
    });
    return res.status(400).json({ 
      error: 'Validation failed',
      errors: errors.array() 
    });
  }

  // Process valid request
  logger.info('Product creation request', {
    productName: req.body.name,
    remoteAddress: req.ip
  });

  // Here you would normally save to database
  res.json({
    message: 'Product created successfully',
    product: {
      name: req.body.name,
      price: req.body.price,
      description: req.body.description,
      createdAt: new Date()
    }
  });
});

// Start HTTP server
app.listen(PORT, () => {
  logger.info('Health check server started', { port: PORT });
});

async function run() {
  const client = new MongoClient(uri, { useUnifiedTopology: true });
  
  try {
    logger.info('Connecting to MongoDB', { uri: uri.replace(/\/\/[^:]+:[^@]+@/, '//***:***@') });
    
    // Connect to MongoDB
    await client.connect();
    logger.info('Successfully connected to MongoDB');
    
    const db = client.db('appdb');
    const products = db.collection('products');
    
    // Generate random product name for demonstration
    const randomName = 'Product_' + Math.random().toString(36).substring(2, 10);
    
    logger.info('Creating new product', { productName: randomName });
    
    // Insert new product with current timestamp
    const result = await products.insertOne({ 
      name: randomName, 
      createdAt: new Date() 
    });
    
    logger.info('Product created successfully', {
      productId: result.insertedId,
      productName: randomName,
      timestamp: new Date().toISOString()
    });
    
    console.log('Inserted product:', result.insertedId, 'with name:', randomName);
    
  } catch (err) {
    logger.error('Error in product creation', {
      error: err.message,
      stack: err.stack
    });
    console.error('Error:', err);
  } finally {
    // Always close the connection
    await client.close();
    logger.info('MongoDB connection closed');
  }
}

// Execute the product creation every 5 seconds
setInterval(run, 5000);

// Initial run
logger.info('Starting Node.js application');
run();
