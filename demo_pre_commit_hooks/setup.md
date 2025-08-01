# Pre-commit Hooks Setup Guide

## Quick Setup

### 1. Install pre-commit

```bash
# Using pip
pip install pre-commit

# Using conda
conda install -c conda-forge pre-commit

# Using Homebrew (macOS)
brew install pre-commit
```

### 2. Install Hooks

```bash
# Install all hooks from .pre-commit-config.yaml
pre-commit install

# Install commit-msg hook for commit message validation
pre-commit install --hook-type commit-msg
```

### 3. Verify Installation

```bash
# Check installed hooks
pre-commit --version

# List installed hooks
pre-commit run --all-files
```

## Configuration

### 1. Basic Configuration

Copy the `.pre-commit-config.yaml` file to your project root:

```bash
cp demo_pre_commit_hooks/.pre-commit-config.yaml .
```

### 2. Customize Configuration

Edit `.pre-commit-config.yaml` to match your project needs:

```yaml
# Example: Disable specific hooks
repos:
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.5.0
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
      # Comment out hooks you don't need
      # - id: check-yaml
```

### 3. Language-Specific Setup

#### Python

```bash
# Install Python tools
pip install black isort flake8 bandit mypy

# Create pyproject.toml for tool configuration
cp code-quality/pyproject.toml .
```

#### Go

```bash
# Install Go tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Copy Go configuration
cp code-quality/.golangci.yml .
```

#### JavaScript/Node.js

```bash
# Install ESLint and plugins
npm install --save-dev eslint eslint-plugin-security

# Copy ESLint configuration
cp code-quality/.eslintrc.js .
```

## Usage

### 1. Automatic Execution

Hooks run automatically on:
- `git commit`
- `git push` (if configured)

### 2. Manual Execution

```bash
# Run all hooks on all files
pre-commit run --all-files

# Run specific hook
pre-commit run black --all-files

# Run hooks on staged files only
pre-commit run

# Run with verbose output
pre-commit run --verbose
```

### 3. Skip Hooks (Emergency)

```bash
# Skip all hooks
git commit -m "Emergency fix" --no-verify

# Skip specific hook
SKIP=black git commit -m "Fix formatting later"
```

## Troubleshooting

### 1. Hook Failures

```bash
# See detailed error messages
pre-commit run --verbose

# Run specific failing hook
pre-commit run black --verbose --all-files

# Fix auto-fixable issues
pre-commit run --all-files
```

### 2. Performance Issues

```bash
# Run hooks in parallel
pre-commit run --all-files --jobs 4

# Skip slow hooks during development
SKIP=bandit git commit -m "Quick fix"
```

### 3. Configuration Issues

```bash
# Validate configuration
pre-commit validate-config

# Update hook versions
pre-commit autoupdate
```

## Advanced Configuration

### 1. Custom Hooks

Create custom hooks in `.pre-commit-config.yaml`:

```yaml
- repo: local
  hooks:
    - id: custom-check
      name: Custom Validation
      entry: python scripts/custom_check.py
      language: python
      files: \.py$
      args: [--strict]
```

### 2. Conditional Hooks

```yaml
- repo: https://github.com/pycqa/bandit
  rev: 1.7.5
  hooks:
    - id: bandit
      args: [-r, scripts/]
      exclude: ^tests/
      # Only run on Python files
      types: [python]
```

### 3. Environment-Specific Hooks

```yaml
# Development hooks
- repo: local
  hooks:
    - id: dev-check
      name: Development Check
      entry: python scripts/dev_check.py
      language: python
      stages: [commit]
      # Only run in development
      additional_dependencies: [pytest]
```

## Integration with CI/CD

### 1. GitHub Actions

```yaml
# .github/workflows/pre-commit.yml
name: Pre-commit

on: [push, pull_request]

jobs:
  pre-commit:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-python@v4
        with:
          python-version: '3.11'
      - name: Install pre-commit
        run: pip install pre-commit
      - name: Run pre-commit
        run: pre-commit run --all-files
```

### 2. GitLab CI

```yaml
# .gitlab-ci.yml
pre-commit:
  stage: test
  image: python:3.11
  before_script:
    - pip install pre-commit
  script:
    - pre-commit run --all-files
```

### 3. Jenkins Pipeline

```groovy
// Jenkinsfile
pipeline {
    agent any
    stages {
        stage('Pre-commit') {
            steps {
                sh 'pip install pre-commit'
                sh 'pre-commit run --all-files'
            }
        }
    }
}
```

## Best Practices

### 1. Team Adoption

- Start with basic formatting hooks
- Gradually add stricter quality checks
- Provide clear documentation
- Offer training sessions

### 2. Performance Optimization

- Use language-specific hooks
- Exclude unnecessary files
- Cache hook results
- Run hooks in parallel

### 3. Maintenance

- Regular hook updates
- Configuration reviews
- Performance monitoring
- Team feedback collection

## Migration Guide

### From Manual Checks

1. **Identify Current Checks**
   ```bash
   # List current manual checks
   ls scripts/check_*.py
   ```

2. **Map to Pre-commit Hooks**
   - Code formatting → Black, isort
   - Linting → flake8, golangci-lint
   - Security → bandit, gosec

3. **Gradual Migration**
   ```bash
   # Start with warnings
   pre-commit run --all-files --show-diff-on-failure
   ```

### From Other Tools

1. **Export Configurations**
   ```bash
   # Export ESLint config
   cp .eslintrc.js code-quality/
   
   # Export Go config
   cp .golangci.yml code-quality/
   ```

2. **Test Integration**
   ```bash
   # Test with existing code
   pre-commit run --all-files
   ```

3. **Update Documentation**
   - Update README
   - Create setup guides
   - Train team members

## Resources

- [Pre-commit Documentation](https://pre-commit.com/)
- [Hook Repository](https://github.com/pre-commit/pre-commit-hooks)
- [Black Documentation](https://black.readthedocs.io/)
- [ESLint Documentation](https://eslint.org/)
- [GolangCI-lint Documentation](https://golangci-lint.run/) 