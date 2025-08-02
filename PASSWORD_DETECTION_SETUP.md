# Password Detection Setup

## Overview

This document describes the password detection and pre-commit hook setup added to the MongoDB Replica Set project.

## 🔒 Password Detection Hook

### Features
- **Custom Python Script**: `scripts/check_passwords.py`
- **Pre-commit Integration**: Automatically runs before each commit
- **Comprehensive Scanning**: Detects various types of hardcoded credentials
- **Smart Exclusions**: Ignores documentation and baseline files

### Detected Patterns
```python
# Common password patterns
password = "secret123"
pass = "admin123"
secret = "mysecret"

# MongoDB connection strings
mongodb://user:password@host:port/db
mongodb+srv://user:password@cluster.mongodb.net/db

# Weak passwords
"password", "123456", "qwerty", "admin", "root"
```

### Usage
```bash
# Manual testing
python3 scripts/check_passwords.py file1.py file2.js

# Via pre-commit
pre-commit run password-check --all-files
```

## 🚀 Demo GitHub Actions

### Features
- **Complete CI/CD Pipeline**: Simulates full build process
- **Multiple Jobs**: Code quality, security, testing, build, integration, deploy
- **Echo Commands**: All steps use echo for demonstration
- **Realistic Workflow**: Follows actual CI/CD best practices

### Jobs Included
1. **Code Quality Analysis**
   - Python, Node.js, Go linting
   - Security scanning
   - Type checking

2. **Security Analysis**
   - Container security scanning
   - Dependency vulnerability checks

3. **Unit Tests**
   - Python, Node.js, Go test simulation

4. **Build and Push**
   - Docker image building
   - Container registry push

5. **Integration Tests**
   - MongoDB infrastructure setup
   - Application testing

6. **Deploy**
   - Production deployment simulation
   - Health checks

### Usage
```bash
# Trigger manually
gh workflow run demo-build.yml

# Or push to trigger automatically
git push origin dev
```

## 📁 Files Added

### Pre-commit Configuration
- `demo_pre_commit_hooks/.pre-commit-config.yaml` - Updated with password detection
- `scripts/check_passwords.py` - Custom password detection script
- `.secrets.baseline` - Baseline for detect-secrets
- `env.example` - Example environment variables

### GitHub Actions
- `.github/workflows/demo-build.yml` - Complete demo pipeline
- Updated `.github/workflows/ci-cd.yml` - Enhanced existing pipeline

### Documentation
- Updated `demo_pre_commit_hooks/README.md` - Comprehensive documentation
- `PASSWORD_DETECTION_SETUP.md` - This summary document

## 🔧 Configuration

### Pre-commit Hook
```yaml
- repo: local
  hooks:
    - id: password-check
      name: Password Detection
      entry: python3 scripts/check_passwords.py
      language: python
      files: \.(py|js|go|yml|yaml|md|txt)$
      exclude: |
        (?x)^(
            \.secrets\.baseline|
            \.env\.example|
            docs/.*\.md|
            README\.md|
            node_modules/|
            \.git/
        )$
```

### Environment Variables
```bash
# Copy example file
cp env.example .env

# Fill in your values
MONGO_ADMIN_PASSWORD=your_secure_password
APP_DB_PASSWORD=your_secure_password
```

## 🧪 Testing

### Test Password Detection
```bash
# Test with current hardcoded credentials
python3 scripts/check_passwords.py app-go/read_products.go

# Expected output: Should detect hardcoded passwords
```

### Test Pre-commit Hook
```bash
# Install pre-commit
pip install pre-commit
pre-commit install

# Try to commit (should be blocked)
git add .
git commit -m "test commit"
```

## 🎯 Benefits

### Security
- **Prevents credential leaks**: Blocks commits with hardcoded passwords
- **Comprehensive scanning**: Detects various credential patterns
- **Baseline tracking**: Tracks known secrets over time

### Development
- **Automated checks**: No manual password hunting needed
- **Early detection**: Catches issues before they reach production
- **Team consistency**: Ensures all developers follow security practices

### CI/CD
- **Complete pipeline**: Full build and deployment simulation
- **Quality gates**: Multiple checkpoints before deployment
- **Realistic workflow**: Mirrors actual production processes

## 🚨 Current Issues Detected

The password detection hook currently finds these hardcoded credentials:

1. **app-go/read_products.go:18**
   - `appuserpassword` in MongoDB connection string

2. **app-node/create_product.js:4**
   - `appuserpassword` in MongoDB connection string

3. **scripts/create_app_user.py:9,14**
   - `mongo-1` and `appuserpassword` hardcoded

## 💡 Recommendations

### Immediate Actions
1. **Replace hardcoded credentials** with environment variables
2. **Create .env file** for local development
3. **Update applications** to use environment variables

### Example Fix
```python
# Before (❌ Bad)
password = "secret123"

# After (✅ Good)
import os
password = os.getenv("DB_PASSWORD")
```

### Production Setup
1. **Use secret management** (AWS Secrets Manager, HashiCorp Vault)
2. **Environment-specific configs** for different deployments
3. **Regular security audits** of the baseline file

## 🔄 Next Steps

1. **Fix hardcoded credentials** in the codebase
2. **Set up environment variables** for development
3. **Configure production secrets** management
4. **Test the complete pipeline** with real builds
5. **Monitor and update** the secrets baseline regularly 