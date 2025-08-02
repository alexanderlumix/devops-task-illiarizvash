const { MongoClient } = require('mongodb');
const express = require('express');
const winston = require('winston');
const rateLimit = require('express-rate-limit');
const prometheus = require('prom-client');

// Initialize Prometheus metrics
const register = new prometheus.Registry();
prometheus.collectDefaultMetrics({ register });

// Custom metrics
const httpRequestDurationMicroseconds = new prometheus.Histogram({
  name: 'http_request_duration_seconds',
  help: 'Duration of HTTP requests in seconds',
  labelNames: ['method', 'route', 'status_code'],
  buckets: [0.1, 0.5, 1, 2, 5]
});

const httpRequestsTotal = new prometheus.Counter({
  name: 'http_requests_total',
  help: 'Total number of HTTP requests',
  labelNames: ['method', 'route', 'status_code']
});

const productsCreated = new prometheus.Counter({
  name: 'products_created_total',
  help: 'Total number of products created'
});

const mongoOperations = new prometheus.Counter({
  name: 'mongo_operations_total',
  help: 'Total number of MongoDB operations',
  labelNames: ['operation', 'status']
});

// Register metrics
register.registerMetric(httpRequestDurationMicroseconds);
register.registerMetric(httpRequestsTotal);
register.registerMetric(productsCreated);
register.registerMetric(mongoOperations);

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

// Add JSON body parser
app.use(express.json());

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

// Metrics middleware
app.use((req, res, next) => {
  const start = Date.now();
  
  res.on('finish', () => {
    const duration = Date.now() - start;
    const route = req.route ? req.route.path : req.path;
    
    httpRequestDurationMicroseconds
      .labels(req.method, route, res.statusCode)
      .observe(duration / 1000);
    
    httpRequestsTotal
      .labels(req.method, route, res.statusCode)
      .inc();
  });
  
  next();
});

// Metrics endpoint
app.get('/metrics', async (req, res) => {
  try {
    res.set('Content-Type', register.contentType);
    res.end(await register.metrics());
  } catch (err) {
    res.status(500).end(err);
  }
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
app.post('/products', (req, res) => {
  // Manual validation instead of express-validator
  const { name, price, description } = req.body;
  const errors = [];

  // Validate name
  if (!name || typeof name !== 'string' || name.trim().length === 0) {
    errors.push({
      type: 'field',
      msg: 'Product name is required',
      path: 'name',
      location: 'body'
    });
  } else if (name.length > 100) {
    errors.push({
      type: 'field',
      msg: 'Product name must be less than 100 characters',
      path: 'name',
      location: 'body'
    });
  }

  // Validate price
  if (price !== undefined && (isNaN(price) || price < 0)) {
    errors.push({
      type: 'field',
      msg: 'Price must be a non-negative number',
      path: 'price',
      location: 'body'
    });
  }

  // Validate description
  if (description && description.length > 500) {
    errors.push({
      type: 'field',
      msg: 'Description must be less than 500 characters',
      path: 'description',
      location: 'body'
    });
  }

  // Check for validation errors
  if (errors.length > 0) {
    logger.warn('Validation failed', {
      errors: errors,
      remoteAddress: req.ip
    });
    return res.status(400).json({ 
      error: 'Validation failed',
      errors: errors 
    });
  }

  // Process valid request
  logger.info('Product creation request', {
    productName: name,
    remoteAddress: req.ip
  });

  // Track product creation metric
  productsCreated.inc();

  // Here you would normally save to database
  res.json({
    message: 'Product created successfully',
    product: {
      name: name,
      price: price,
      description: description,
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
