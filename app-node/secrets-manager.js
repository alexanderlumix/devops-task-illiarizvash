const fs = require('fs');
const path = require('path');

class SecretsManager {
  constructor() {
    this.isProduction = process.env.NODE_ENV === 'production';
    this.credentialsPath = path.join(process.cwd(), 'credentials.local.json');
    this.examplePath = path.join(process.cwd(), 'env.credentials.example');
  }

  /**
   * Get secrets based on environment
   * @returns {Promise<Object>} Secrets object
   */
  async getSecrets() {
    if (this.isProduction) {
      return await this.getProductionSecrets();
    } else {
      return this.getLocalSecrets();
    }
  }

  /**
   * Get production secrets from AWS Secrets Manager
   * @returns {Promise<Object>} Production secrets
   */
  async getProductionSecrets() {
    try {
      // In production, this would use AWS SDK
      // For now, we'll throw an error to indicate this needs implementation
      throw new Error('Production secrets management not implemented. Please configure AWS Secrets Manager.');
      
      // Example AWS Secrets Manager implementation:
      // const AWS = require('aws-sdk');
      // const secretsManager = new AWS.SecretsManager();
      // const data = await secretsManager.getSecretValue({
      //   SecretId: 'mongodb-credentials'
      // }).promise();
      // return JSON.parse(data.SecretString);
    } catch (error) {
      console.error('Error fetching production secrets:', error.message);
      throw error;
    }
  }

  /**
   * Get local secrets from credentials file
   * @returns {Object} Local secrets
   */
  getLocalSecrets() {
    try {
      if (!fs.existsSync(this.credentialsPath)) {
        this.createLocalCredentialsFile();
        throw new Error(`Local credentials file not found. Please create ${this.credentialsPath} based on env.credentials.example`);
      }
      
      const credentials = JSON.parse(fs.readFileSync(this.credentialsPath, 'utf8'));
      this.validateCredentials(credentials);
      return credentials.mongodb;
    } catch (error) {
      console.error('Error reading local secrets:', error.message);
      throw error;
    }
  }

  /**
   * Create local credentials file from example
   */
  createLocalCredentialsFile() {
    try {
      if (fs.existsSync(this.examplePath)) {
        const exampleContent = fs.readFileSync(this.examplePath, 'utf8');
        fs.writeFileSync(this.credentialsPath, exampleContent);
        console.log(`Created ${this.credentialsPath} from example. Please update with your local credentials.`);
      } else {
        console.error('Example credentials file not found:', this.examplePath);
      }
    } catch (error) {
      console.error('Error creating local credentials file:', error.message);
    }
  }

  /**
   * Validate credentials structure
   * @param {Object} credentials - Credentials object
   */
  validateCredentials(credentials) {
    const requiredFields = ['mongodb'];
    const requiredMongoFields = ['admin_user', 'admin_password', 'app_user', 'app_password', 'host', 'port', 'database'];

    // Check top-level structure
    for (const field of requiredFields) {
      if (!credentials[field]) {
        throw new Error(`Missing required field: ${field}`);
      }
    }

    // Check MongoDB credentials
    const mongoCreds = credentials.mongodb;
    for (const field of requiredMongoFields) {
      if (!mongoCreds[field]) {
        throw new Error(`Missing required MongoDB field: ${field}`);
      }
    }

    // Validate port is numeric
    if (isNaN(parseInt(mongoCreds.port))) {
      throw new Error('MongoDB port must be a number');
    }
  }

  /**
   * Get MongoDB connection URI from secrets
   * @param {Object} secrets - MongoDB secrets
   * @param {boolean} directConnection - Whether to use direct connection
   * @returns {string} MongoDB connection URI
   */
  getMongoURI(secrets, directConnection = false) {
    const { app_user, app_password, host, port, database, replica_set } = secrets;
    
    if (directConnection) {
      return `mongodb://${app_user}:${app_password}@${host}:${port}/${database}?directConnection=true`;
    } else {
      return `mongodb://${app_user}:${app_password}@${host}:${port}/${database}?replicaSet=${replica_set}`;
    }
  }

  /**
   * Get admin MongoDB URI for user management
   * @param {Object} secrets - MongoDB secrets
   * @returns {string} Admin MongoDB connection URI
   */
  getAdminMongoURI(secrets) {
    const { admin_user, admin_password, host, port, replica_set } = secrets;
    return `mongodb://${admin_user}:${admin_password}@${host}:${port}/admin?replicaSet=${replica_set}&authSource=admin`;
  }

  /**
   * Test secrets configuration
   * @returns {Promise<boolean>} True if secrets are valid
   */
  async testSecrets() {
    try {
      const secrets = await this.getSecrets();
      this.validateCredentials({ mongodb: secrets });
      console.log('✅ Secrets configuration is valid');
      return true;
    } catch (error) {
      console.error('❌ Secrets configuration error:', error.message);
      return false;
    }
  }

  /**
   * Get JWT secret
   * @returns {string} JWT secret
   */
  getJWTSecret() {
    try {
      if (fs.existsSync(this.credentialsPath)) {
        const credentials = JSON.parse(fs.readFileSync(this.credentialsPath, 'utf8'));
        return (credentials.jwt && credentials.jwt.secret) || 'default-jwt-secret-change-in-production';
      }
      return 'default-jwt-secret-change-in-production';
    } catch (error) {
      console.error('Error reading JWT secret:', error.message);
      return 'default-jwt-secret-change-in-production';
    }
  }

  /**
   * Get encryption key
   * @returns {string} Encryption key
   */
  getEncryptionKey() {
    try {
      if (fs.existsSync(this.credentialsPath)) {
        const credentials = JSON.parse(fs.readFileSync(this.credentialsPath, 'utf8'));
        return (credentials.encryption && credentials.encryption.key) || 'default-encryption-key-32-chars';
      }
      return 'default-encryption-key-32-chars';
    } catch (error) {
      console.error('Error reading encryption key:', error.message);
      return 'default-encryption-key-32-chars';
    }
  }
}

module.exports = SecretsManager; 