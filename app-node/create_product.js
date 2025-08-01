const { MongoClient } = require('mongodb');
const express = require('express');

// getMongoURI constructs MongoDB connection URI from environment variables
function getMongoURI() {
  const user = process.env.MONGO_USER || 'appuser';
  const password = process.env.MONGO_PASSWORD || 'appuserpassword';
  const host = process.env.MONGO_HOST || '127.0.0.1';
  const port = process.env.MONGO_PORT || '27032';
  const db = process.env.MONGO_DB || 'appdb';
  const directConnection = process.env.MONGO_DIRECT_CONNECTION || 'true';

  return `mongodb://${user}:${password}@${host}:${port}/${db}?directConnection=${directConnection}`;
}

// MongoDB connection URI for direct connection to replica set member
const uri = getMongoURI();

// Create Express app for health checks
const app = express();
const PORT = process.env.PORT || 3000;

// Health check endpoint
app.get('/health', (req, res) => {
  res.json({
    status: 'healthy',
    service: 'node-app',
    timestamp: new Date().toISOString()
  });
});

// Start HTTP server
app.listen(PORT, () => {
  console.log(`Health check server running on port ${PORT}`);
});

async function run() {
  const client = new MongoClient(uri, { useUnifiedTopology: true });
  try {
    // Connect to MongoDB
    await client.connect();
    const db = client.db('appdb');
    const products = db.collection('products');
    
    // Generate random product name for demonstration
    const randomName = 'Product_' + Math.random().toString(36).substring(2, 10);
    
    // Insert new product with current timestamp
    const result = await products.insertOne({ name: randomName, createdAt: new Date() });
    console.log('Inserted product:', result.insertedId, 'with name:', randomName);
  } catch (err) {
    console.error('Error:', err);
  } finally {
    // Always close the connection
    await client.close();
  }
}

// Execute the product creation every 5 seconds
setInterval(run, 5000);

// Initial run
run();
