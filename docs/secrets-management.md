# Secrets Management Documentation

## Overview

This project implements a secure secrets management system that supports both local development and production environments. The system provides a unified interface for accessing credentials while maintaining security best practices.

## Architecture

### Local Development
- **File-based**: `credentials.local.json` for local development
- **Security**: File excluded from version control via `.gitignore`
- **Validation**: Automatic validation of credentials structure
- **Fallback**: Default values for missing credentials

### Production Environment
- **AWS Secrets Manager**: Secure cloud-based secret storage
- **IAM Roles**: Least privilege access control
- **Encryption**: At-rest and in-transit encryption
- **Audit Logging**: Complete access audit trail

## Setup Instructions

### Local Development Setup

#### 1. Create Local Credentials File

```bash
# Copy the example file
cp env.credentials.example credentials.local.json

# Edit the file with your local credentials
nano credentials.local.json
```

#### 2. Configure Local Credentials

```json
{
  "mongodb": {
    "admin_user": "mongo-1",
    "admin_password": "mongo-1",
    "app_user": "appuser",
    "app_password": "appuserpassword",
    "host": "127.0.0.1",
    "port": "27017",
    "database": "appdb",
    "replica_set": "rs0"
  },
  "jwt": {
    "secret": "local-dev-secret-key-change-in-production"
  },
  "encryption": {
    "key": "local-encryption-key-32-chars"
  }
}
```

#### 3. Verify Configuration

```bash
# Test Node.js secrets
node -e "
const SecretsManager = require('./app-node/secrets-manager.js');
const sm = new SecretsManager();
sm.testSecrets().then(result => console.log('Test result:', result));
"

# Test Go secrets
cd app-go
go run secrets/secrets.go
```

### Production Setup

#### 1. AWS Secrets Manager Configuration

```bash
# Create MongoDB credentials secret
aws secretsmanager create-secret \
  --name "mongodb-credentials" \
  --description "MongoDB connection credentials" \
  --secret-string '{
    "admin_user": "admin",
    "admin_password": "secure-admin-password",
    "app_user": "appuser",
    "app_password": "secure-app-password",
    "host": "mongodb-cluster.example.com",
    "port": "27017",
    "database": "appdb",
    "replica_set": "rs0"
  }'
```

#### 2. IAM Role Configuration

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "secretsmanager:GetSecretValue"
      ],
      "Resource": "arn:aws:secretsmanager:region:account:secret:mongodb-credentials-*"
    }
  ]
}
```

#### 3. Environment Variables

```bash
# Set production environment
export NODE_ENV=production
export GO_ENV=production

# AWS credentials (if not using IAM roles)
export AWS_ACCESS_KEY_ID=your-access-key
export AWS_SECRET_ACCESS_KEY=your-secret-key
export AWS_REGION=us-east-1
```

## Usage Examples

### Node.js Application

```javascript
const SecretsManager = require('./secrets-manager.js');

async function connectToDatabase() {
  const sm = new SecretsManager();
  
  try {
    // Get secrets
    const secrets = await sm.getSecrets();
    
    // Get MongoDB URI
    const uri = sm.getMongoURI(secrets, false);
    
    // Connect to MongoDB
    const client = new MongoClient(uri);
    await client.connect();
    
    console.log('Connected to MongoDB');
    return client;
  } catch (error) {
    console.error('Failed to connect:', error);
    throw error;
  }
}
```

### Go Application

```go
package main

import (
    "log"
    "./secrets"
)

func main() {
    sm := secrets.NewSecretsManager()
    
    // Get secrets
    secrets, err := sm.GetSecrets()
    if err != nil {
        log.Fatal("Failed to get secrets:", err)
    }
    
    // Get MongoDB URI
    uri := sm.GetMongoURI(secrets, false)
    
    // Connect to MongoDB
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
    if err != nil {
        log.Fatal("Failed to connect to MongoDB:", err)
    }
    
    log.Println("Connected to MongoDB")
}
```

## Security Best Practices

### Local Development
- ✅ Credentials file in `.gitignore`
- ✅ Hardcoded passwords only for local testing
- ✅ Clear documentation about local vs production
- ✅ Validation of local credentials file
- ✅ Secure file permissions (0600)

### Production Environment
- ✅ AWS IAM roles with minimal permissions
- ✅ Secret rotation policies
- ✅ Audit logging for secret access
- ✅ Encryption at rest and in transit
- ✅ Network security (VPC, security groups)

### General Security
- ✅ No hardcoded credentials in code
- ✅ Environment-specific configurations
- ✅ Proper error handling for missing secrets
- ✅ Secure credential transmission
- ✅ Regular security audits

## Testing

### Local Testing

```bash
# Test Node.js secrets manager
node -e "
const SecretsManager = require('./app-node/secrets-manager.js');
const sm = new SecretsManager();
sm.testSecrets().then(result => console.log('Test result:', result));
"

# Test Go secrets manager
cd app-go
go run -c secrets/secrets.go
```

### Production Testing

```bash
# Test AWS Secrets Manager access
aws secretsmanager get-secret-value --secret-id mongodb-credentials

# Test application with production secrets
NODE_ENV=production node app-node/create_product.js
```

## Troubleshooting

### Common Issues

#### 1. Local Credentials File Not Found
```
Error: Local credentials file not found. Please create credentials.local.json
```

**Solution**: Copy the example file and configure it:
```bash
cp env.credentials.example credentials.local.json
# Edit the file with your credentials
```

#### 2. Invalid Credentials Structure
```
Error: Missing required field: admin_user
```

**Solution**: Ensure all required fields are present in the credentials file.

#### 3. AWS Secrets Manager Access Denied
```
Error: User: arn:aws:sts::account:assumed-role/role is not authorized
```

**Solution**: Configure proper IAM roles and policies.

#### 4. Production Environment Not Set
```
Error: Production secrets management not implemented
```

**Solution**: Set environment variables:
```bash
export NODE_ENV=production
export GO_ENV=production
```

### Debug Mode

Enable debug logging for troubleshooting:

```bash
# Node.js
DEBUG=secrets:* node app-node/create_product.js

# Go
LOG_LEVEL=debug go run app-go/read_products.go
```

## Monitoring and Alerting

### Secret Access Monitoring

```bash
# Monitor AWS Secrets Manager access
aws cloudtrail lookup-events \
  --lookup-attributes AttributeKey=EventName,AttributeValue=GetSecretValue

# Monitor application secret usage
grep "secrets" /var/log/application.log
```

### Alerting Rules

- **Secret Access Failures**: Alert on repeated access failures
- **Credential Rotation**: Alert when secrets are near expiration
- **Unauthorized Access**: Alert on suspicious access patterns
- **Application Errors**: Alert on secret-related application errors

## Backup and Recovery

### Backup Strategy

```bash
# Backup local credentials (encrypted)
gpg -e credentials.local.json

# Backup AWS secrets
aws secretsmanager get-secret-value --secret-id mongodb-credentials > backup.json
```

### Recovery Procedures

1. **Local Recovery**: Restore from encrypted backup
2. **Production Recovery**: Restore from AWS Secrets Manager
3. **Application Recovery**: Redeploy with new credentials
4. **Validation**: Test all connections after recovery

## Compliance

### Audit Requirements

- **Access Logging**: All secret access is logged
- **Rotation Policies**: Regular secret rotation
- **Encryption**: All secrets encrypted at rest
- **Access Control**: Least privilege access
- **Monitoring**: Continuous security monitoring

### Documentation Requirements

- **Configuration**: Document all secret configurations
- **Procedures**: Document access and rotation procedures
- **Recovery**: Document backup and recovery procedures
- **Compliance**: Document compliance requirements

## Future Enhancements

### Planned Features

1. **HashiCorp Vault Integration**: Alternative to AWS Secrets Manager
2. **Dynamic Secret Generation**: Automatic secret rotation
3. **Multi-Region Support**: Cross-region secret replication
4. **Advanced Monitoring**: Enhanced security monitoring
5. **Compliance Reporting**: Automated compliance reports

### Integration Roadmap

- **Phase 1**: AWS Secrets Manager (Current)
- **Phase 2**: HashiCorp Vault support
- **Phase 3**: Dynamic secret generation
- **Phase 4**: Advanced monitoring and alerting
- **Phase 5**: Compliance automation 