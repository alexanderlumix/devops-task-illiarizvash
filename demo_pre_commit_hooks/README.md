# Pre-commit Hooks Configuration

This directory contains pre-commit hooks configuration for the MongoDB Replica Set project to ensure code quality and security before each commit.

## Overview

Pre-commit hooks are automated checks that run before each git commit to ensure:
- Code quality standards
- Security best practices
- Consistent formatting
- No hardcoded secrets

## Features

### 🔒 Security Hooks
- **Password Detection**: Custom hook to detect hardcoded passwords and credentials
- **Secret Scanning**: detect-secrets integration for comprehensive secret detection
- **Security Linting**: Bandit for Python security vulnerabilities
- **Container Security**: Hadolint for Dockerfile security

### 🎨 Code Quality Hooks
- **Python**: Black (formatting), isort (imports), flake8 (linting), mypy (type checking)
- **Go**: golangci-lint, go vet, go fmt
- **JavaScript**: ESLint with security plugins
- **Markdown**: markdownlint for documentation quality

### 🔧 Utility Hooks
- **File Checks**: trailing whitespace, end-of-file, large files
- **YAML/JSON**: validation and linting
- **Merge Conflicts**: detection and prevention

## Installation

### Prerequisites
```bash
# Install pre-commit
pip install pre-commit

# Install language-specific tools
pip install black isort flake8 bandit mypy
npm install -g eslint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Setup
```bash
# Install pre-commit hooks
pre-commit install

# Install all hooks
pre-commit install --hook-type pre-commit --hook-type commit-msg
```

## Custom Password Detection Hook

### Overview
The custom password detection hook (`scripts/check_passwords.py`) scans files for:
- Hardcoded passwords and credentials
- MongoDB connection strings with passwords
- Common weak passwords
- API keys and tokens

### Patterns Detected
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
# Run manually
python scripts/check_passwords.py file1.py file2.js

# Run via pre-commit
pre-commit run password-check --all-files
```

### Configuration
The hook is configured in `.pre-commit-config.yaml`:
```yaml
- repo: local
  hooks:
    - id: password-check
      name: Password Detection
      entry: python scripts/check_passwords.py
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

## Usage

### Running Hooks
```bash
# Run all hooks on staged files
pre-commit run

# Run specific hook
pre-commit run black
pre-commit run password-check

# Run on all files
pre-commit run --all-files

# Run specific hook on all files
pre-commit run password-check --all-files
```

### Skipping Hooks
```bash
# Skip specific hook
git commit -m "message" --no-verify

# Skip all hooks (not recommended)
git commit -m "message" --no-verify
```

## Hook Details

### Password Detection Hook
- **Purpose**: Detect hardcoded credentials
- **Files**: Python, JavaScript, Go, YAML, Markdown, Text files
- **Excludes**: Documentation, baseline files, node_modules
- **Action**: Blocks commit if credentials found

### Security Hooks
- **detect-secrets**: Comprehensive secret scanning
- **bandit**: Python security vulnerabilities
- **hadolint**: Dockerfile security

### Code Quality Hooks
- **Black**: Python code formatting
- **isort**: Python import sorting
- **flake8**: Python linting
- **mypy**: Python type checking
- **golangci-lint**: Go linting
- **ESLint**: JavaScript linting

## Troubleshooting

### Common Issues

#### Hook Fails on Password Detection
```bash
# Check what passwords were detected
python scripts/check_passwords.py $(git diff --cached --name-only)

# Fix by using environment variables
# Instead of: password = "secret123"
# Use: password = os.getenv("DB_PASSWORD")
```

#### Hook Fails on Code Formatting
```bash
# Auto-fix formatting issues
pre-commit run black --all-files
pre-commit run isort --all-files
```

#### Hook Fails on Linting
```bash
# Check specific linting errors
pre-commit run flake8
pre-commit run golangci-lint
```

### Updating Hooks
```bash
# Update all hooks to latest versions
pre-commit autoupdate

# Update specific hook
pre-commit autoupdate --freeze
```

## Best Practices

### 1. Use Environment Variables
```python
# ❌ Bad
password = "secret123"

# ✅ Good
import os
password = os.getenv("DB_PASSWORD")
```

### 2. Create .env Files
```bash
# Create .env for local development
echo "DB_PASSWORD=your_password_here" > .env

# Add to .gitignore
echo ".env" >> .gitignore
```

### 3. Use Secret Management
```python
# For production, use proper secret management
# AWS Secrets Manager, HashiCorp Vault, etc.
```

### 4. Regular Updates
```bash
# Keep hooks updated
pre-commit autoupdate

# Review baseline regularly
detect-secrets audit .secrets.baseline
```

## Configuration Files

- `.pre-commit-config.yaml`: Main configuration
- `scripts/check_passwords.py`: Custom password detection
- `.secrets.baseline`: Known secrets baseline
- `code-quality/`: Language-specific configurations

## Integration with CI/CD

The same hooks can be run in CI/CD pipelines:

```yaml
# GitHub Actions example
- name: Run pre-commit hooks
  run: |
    pip install pre-commit
    pre-commit run --all-files
```

## Contributing

When adding new hooks:
1. Update `.pre-commit-config.yaml`
2. Test with `pre-commit run --all-files`
3. Update this README
4. Commit changes

## Support

For issues with pre-commit hooks:
1. Check the hook documentation
2. Review error messages
3. Test individual hooks
4. Update hook versions if needed 