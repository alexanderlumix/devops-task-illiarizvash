const { MongoClient } = require('mongodb');

// MongoDB connection URI for direct connection to replica set member
const uri = 'mongodb://appuser:appuserpassword@127.0.0.1:27032/appdb?directConnection=true';

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

// Execute the product creation
run();
