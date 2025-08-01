# Local Development Setup

## Overview

This guide explains how to set up the MongoDB Replica Set project for local development using the local passwords file.

## 🔐 Local Passwords File

### File: `passwords_and_credentials.txt`

This file contains all the passwords and credentials needed for local development. It's automatically excluded from:
- Git version control (via .gitignore)
- Pre-commit hooks
- Security scans

### Contents
```bash
# MongoDB Admin Credentials
MONGO_ADMIN_USER=admin
MONGO_ADMIN_PASSWORD=mongo-1

# Application Database Credentials
APP_DB_USER=appuser
APP_DB_PASSWORD=appuserpassword

# MongoDB Connection Settings
MONGO_HOST=127.0.0.1
MONGO_PORT=27017
MONGO_DB=appdb
MONGO_REPLICA_SET=rs0
```

## 🚀 GitHub Actions Triggers

### Automatic Triggers
- **Push to branches**: `main`, `dev`, `feature/*`, `bugfix/*`
- **Pull Requests**: to `main` or `dev` branches

### Manual Triggers
```bash
# Trigger workflow manually via GitHub CLI
gh workflow run demo-build.yml

# Or via GitHub web interface:
# 1. Go to Actions tab
# 2. Select "Demo Build Pipeline"
# 3. Click "Run workflow"
# 4. Choose environment and options
```

### Manual Trigger Options
- **Environment**: development, staging, production
- **Force Deploy**: Skip test failures (not recommended)

## 📋 Setup Instructions

### 1. Initial Setup
```bash
# Clone the repository
git clone <repository-url>
cd devops-task-illiarizvash

# The passwords file should already exist
ls passwords_and_credentials.txt
```

### 2. Install Pre-commit Hooks
```bash
# Install pre-commit
pip install pre-commit

# Install hooks
pre-commit install

# Install all hook types
pre-commit install --hook-type pre-commit --hook-type commit-msg
```

### 3. Verify Password Detection
```bash
# Test password detection (should exclude passwords file)
python3 scripts/check_passwords.py app-go/read_products.go

# Expected: Should detect hardcoded passwords in code files
# But should NOT scan passwords_and_credentials.txt
```

### 4. Start Development Environment
```bash
# Start MongoDB replica set
cd mongo
docker-compose up -d

# Initialize replica set
cd ../scripts
python3 init_mongo_servers.py

# Create application user
python3 create_app_user.py

# Test applications
cd ../app-node
node create_product.js

cd ../app-go
go run read_products.go
```

## 🔧 Configuration

### Environment Variables
The `passwords_and_credentials.txt` file can be used to set environment variables:

```bash
# Source the file (bash/zsh)
source passwords_and_credentials.txt

# Or export variables manually
export MONGO_ADMIN_PASSWORD=mongo-1
export APP_DB_PASSWORD=appuserpassword
```

### Application Configuration
Update applications to use environment variables:

```python
# Python example
import os
password = os.getenv('APP_DB_PASSWORD', 'default_password')
```

```javascript
// Node.js example
const password = process.env.APP_DB_PASSWORD || 'default_password';
```

```go
// Go example
password := os.Getenv("APP_DB_PASSWORD")
if password == "" {
    password = "default_password"
}
```

## 🧪 Testing

### Test Pre-commit Hooks
```bash
# Run all hooks
pre-commit run --all-files

# Run specific hook
pre-commit run password-check

# Test commit (should be blocked if hardcoded passwords found)
git add .
git commit -m "test commit"
```

### Test GitHub Actions Locally
```bash
# Install act (GitHub Actions local runner)
# https://github.com/nektos/act

# Run workflow locally
act push -W .github/workflows/demo-build.yml
```

## 🚨 Security Best Practices

### 1. Never Commit Passwords
```bash
# ❌ Bad - Don't do this
git add passwords_and_credentials.txt
git commit -m "add passwords"

# ✅ Good - File is already in .gitignore
git status  # Should not show passwords file
```

### 2. Use Environment Variables
```python
# ❌ Bad
password = "hardcoded_password"

# ✅ Good
password = os.getenv("DB_PASSWORD")
```

### 3. Regular Security Checks
```bash
# Check for hardcoded passwords
python3 scripts/check_passwords.py $(find . -name "*.py" -o -name "*.js" -o -name "*.go")

# Update secrets baseline
detect-secrets scan --baseline .secrets.baseline
```

## 🔄 Workflow

### Development Workflow
1. **Start development**: Use `passwords_and_credentials.txt`
2. **Write code**: Use environment variables
3. **Test locally**: Run pre-commit hooks
4. **Commit**: Hooks will check for hardcoded passwords
5. **Push**: Triggers GitHub Actions automatically
6. **Review**: Check Actions results

### Production Workflow
1. **Create PR**: From feature branch to main
2. **CI/CD**: GitHub Actions run automatically
3. **Review**: Code quality and security checks
4. **Merge**: Deploy to production
5. **Monitor**: Check deployment status

## 🛠️ Troubleshooting

### Pre-commit Hook Issues
```bash
# Hook fails on password detection
python3 scripts/check_passwords.py $(git diff --cached --name-only)

# Fix by using environment variables
# Replace hardcoded passwords with os.getenv()
```

### GitHub Actions Issues
```bash
# Check workflow status
gh run list --workflow=demo-build.yml

# View logs
gh run view <run-id>
```

### Local Development Issues
```bash
# Check if passwords file exists
ls -la passwords_and_credentials.txt

# Verify .gitignore
git check-ignore passwords_and_credentials.txt

# Test password detection exclusion
python3 scripts/check_passwords.py passwords_and_credentials.txt
# Should show: No hardcoded passwords detected
```

## 📚 Additional Resources

- [Pre-commit Documentation](https://pre-commit.com/)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [MongoDB Security Best Practices](https://docs.mongodb.com/manual/security/)
- [Environment Variables Best Practices](https://12factor.net/config) 