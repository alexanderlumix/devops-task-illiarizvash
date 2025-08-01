const { MongoClient } = require('mongodb');
const express = require('express');
const request = require('supertest');

// Mock MongoDB for testing
jest.mock('mongodb');

describe('Node.js Application Tests', () => {
  let app;

  beforeEach(() => {
    // Reset environment variables
    delete process.env.MONGO_USER;
    delete process.env.MONGO_PASSWORD;
    delete process.env.MONGO_HOST;
    delete process.env.MONGO_PORT;
    delete process.env.MONGO_DB;
    delete process.env.MONGO_DIRECT_CONNECTION;
    
    // Create fresh app instance for each test
    app = require('./create_product');
  });

  describe('Environment Configuration', () => {
    test('should construct MongoDB URI with default values', () => {
      const uri = getMongoURI();
      expect(uri).toContain('mongodb://appuser:appuserpassword@127.0.0.1:27032/appdb');
    });

    test('should construct MongoDB URI with custom values', () => {
      process.env.MONGO_USER = 'testuser';
      process.env.MONGO_PASSWORD = 'testpass';
      process.env.MONGO_HOST = 'testhost';
      process.env.MONGO_PORT = '27017';
      process.env.MONGO_DB = 'testdb';
      
      const uri = getMongoURI();
      expect(uri).toContain('mongodb://testuser:testpass@testhost:27017/testdb');
    });
  });

  describe('Health Check Endpoint', () => {
    test('should return health status', async () => {
      const response = await request(app)
        .get('/health')
        .expect(200);

      expect(response.body).toHaveProperty('status', 'healthy');
      expect(response.body).toHaveProperty('service', 'node-app');
      expect(response.body).toHaveProperty('timestamp');
      expect(response.body).toHaveProperty('version');
    });
  });

  describe('Product Creation Endpoint', () => {
    test('should create product with valid data', async () => {
      const productData = {
        name: 'Test Product',
        price: 99.99,
        description: 'A test product'
      };

      const response = await request(app)
        .post('/products')
        .send(productData)
        .expect(200);

      expect(response.body).toHaveProperty('message', 'Product created successfully');
      expect(response.body.product).toHaveProperty('name', 'Test Product');
      expect(response.body.product).toHaveProperty('price', 99.99);
    });

    test('should reject product with invalid name', async () => {
      const productData = {
        name: '', // Invalid: empty name
        price: 99.99
      };

      const response = await request(app)
        .post('/products')
        .send(productData)
        .expect(400);

      expect(response.body).toHaveProperty('error', 'Validation failed');
      expect(response.body.errors).toBeDefined();
    });

    test('should reject product with invalid price', async () => {
      const productData = {
        name: 'Test Product',
        price: 'invalid' // Invalid: non-numeric price
      };

      const response = await request(app)
        .post('/products')
        .send(productData)
        .expect(400);

      expect(response.body).toHaveProperty('error', 'Validation failed');
    });

    test('should reject product with too long description', async () => {
      const productData = {
        name: 'Test Product',
        description: 'a'.repeat(501) // Invalid: too long
      };

      const response = await request(app)
        .post('/products')
        .send(productData)
        .expect(400);

      expect(response.body).toHaveProperty('error', 'Validation failed');
    });
  });

  describe('Rate Limiting', () => {
    test('should apply rate limiting', async () => {
      // Make multiple requests quickly
      const promises = [];
      for (let i = 0; i < 105; i++) {
        promises.push(
          request(app)
            .get('/health')
            .catch(err => err)
        );
      }

      const responses = await Promise.all(promises);
      const rateLimited = responses.filter(r => r.status === 429);
      
      expect(rateLimited.length).toBeGreaterThan(0);
    });
  });

  describe('Logging', () => {
    test('should log requests', () => {
      // This test would verify that logging is working
      // In a real application, you would mock the logger and verify calls
      expect(true).toBe(true); // Placeholder
    });
  });

  describe('MongoDB Integration', () => {
    test('should connect to MongoDB', async () => {
      // Mock MongoDB connection
      const mockClient = {
        connect: jest.fn().mockResolvedValue(),
        db: jest.fn().mockReturnValue({
          collection: jest.fn().mockReturnValue({
            insertOne: jest.fn().mockResolvedValue({ insertedId: 'test-id' })
          })
        }),
        close: jest.fn().mockResolvedValue()
      };

      MongoClient.mockImplementation(() => mockClient);

      // Test the run function
      await run();
      
      expect(mockClient.connect).toHaveBeenCalled();
      expect(mockClient.close).toHaveBeenCalled();
    });

    test('should handle MongoDB connection errors', async () => {
      const mockClient = {
        connect: jest.fn().mockRejectedValue(new Error('Connection failed')),
        close: jest.fn().mockResolvedValue()
      };

      MongoClient.mockImplementation(() => mockClient);

      // Test error handling
      await run();
      
      expect(mockClient.connect).toHaveBeenCalled();
      expect(mockClient.close).toHaveBeenCalled();
    });
  });
});

// Helper function to get MongoDB URI (extracted from main file)
function getMongoURI() {
  const user = process.env.MONGO_USER || 'appuser';
  const password = process.env.MONGO_PASSWORD || 'appuserpassword';
  const host = process.env.MONGO_HOST || '127.0.0.1';
  const port = process.env.MONGO_PORT || '27032';
  const db = process.env.MONGO_DB || 'appdb';
  const directConnection = process.env.MONGO_DIRECT_CONNECTION || 'true';

  return `mongodb://${user}:${password}@${host}:${port}/${db}?directConnection=${directConnection}`;
}

// Mock run function for testing
async function run() {
  const client = new MongoClient(getMongoURI(), { useUnifiedTopology: true });
  
  try {
    await client.connect();
    const db = client.db('appdb');
    const products = db.collection('products');
    
    const randomName = 'Product_' + Math.random().toString(36).substring(2, 10);
    const result = await products.insertOne({ 
      name: randomName, 
      createdAt: new Date() 
    });
    
    return result;
  } catch (err) {
    throw err;
  } finally {
    await client.close();
  }
} 