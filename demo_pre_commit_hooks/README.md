# Pre-commit Hooks Demo

This directory contains examples and templates for setting up pre-commit hooks to ensure code quality and consistency across the project.

## Overview

Pre-commit hooks are scripts that run automatically before each commit to check code quality, formatting, and security issues.

## Setup

### 1. Install pre-commit

```bash
# Install pre-commit globally
pip install pre-commit

# Or install in a virtual environment
python -m pip install pre-commit
```

### 2. Install the hooks

```bash
# Install hooks from .pre-commit-config.yaml
pre-commit install

# Install commit-msg hook for commit message validation
pre-commit install --hook-type commit-msg
```

### 3. Run hooks manually

```bash
# Run all hooks on all files
pre-commit run --all-files

# Run specific hook
pre-commit run black --all-files

# Run hooks on staged files only
pre-commit run
```

## Available Hooks

### Code Formatting
- **Black**: Python code formatter
- **isort**: Python import sorter
- **go-fmt**: Go code formatter
- **go-imports**: Go import organizer
- **ESLint**: JavaScript/TypeScript linter

### Code Quality
- **flake8**: Python code quality checker
- **golangci-lint**: Go code quality checker
- **mypy**: Python type checker
- **bandit**: Python security linter

### Security
- **detect-secrets**: Find secrets in code
- **hadolint**: Dockerfile linter
- **gosec**: Go security scanner

### Documentation
- **markdownlint**: Markdown linter
- **yamllint**: YAML linter
- **license-eye**: License header checker

### Git
- **commitizen**: Conventional commit messages
- **commitlint**: Commit message validation
- **trailing-whitespace**: Remove trailing whitespace
- **end-of-file-fixer**: Ensure files end with newline

## Configuration Files

### .pre-commit-config.yaml
Main configuration file defining all hooks and their settings.

### .eslintrc.js
ESLint configuration for JavaScript/TypeScript code quality.

### .golangci.yml
GolangCI-lint configuration for Go code quality.

### pyproject.toml
Python tool configurations (Black, isort, Pylint, Bandit, MyPy).

## Custom Hooks

### Creating Custom Hooks

```yaml
# Example custom hook in .pre-commit-config.yaml
- repo: local
  hooks:
    - id: custom-check
      name: Custom Check
      entry: python scripts/custom_check.py
      language: python
      files: \.py$
```

### Local Scripts

Create custom validation scripts in the `scripts/` directory:

```python
# scripts/custom_check.py
#!/usr/bin/env python3
import sys

def main():
    # Your custom validation logic here
    print("Running custom check...")
    # Return 0 for success, 1 for failure
    return 0

if __name__ == "__main__":
    sys.exit(main())
```

## Best Practices

### 1. Incremental Adoption
- Start with basic formatting hooks
- Gradually add more strict quality checks
- Configure hooks to be warnings initially

### 2. Team Communication
- Document hook requirements
- Provide setup instructions
- Share configuration files

### 3. CI/CD Integration
- Run pre-commit hooks in CI/CD pipeline
- Fail builds on hook failures
- Generate reports for quality metrics

### 4. Performance
- Use language-specific hooks
- Exclude unnecessary files
- Cache hook results

## Troubleshooting

### Common Issues

1. **Hook fails on existing code**
   ```bash
   # Run hooks with --allow-unmatched
   pre-commit run --all-files --allow-unmatched
   ```

2. **Slow hook execution**
   ```bash
   # Run specific hooks only
   pre-commit run black --all-files
   ```

3. **Hook conflicts**
   - Review hook order in configuration
   - Check for conflicting rules
   - Adjust hook settings

### Debug Mode

```bash
# Run with verbose output
pre-commit run --verbose

# Run specific hook with debug
pre-commit run black --verbose --all-files
```

## Migration Guide

### From Manual Checks
1. Identify current manual checks
2. Find equivalent pre-commit hooks
3. Configure hooks to match current standards
4. Gradually increase strictness

### From Other Tools
1. Export current configurations
2. Map to pre-commit hook equivalents
3. Test configurations
4. Update team documentation

## Resources

- [Pre-commit Documentation](https://pre-commit.com/)
- [Hook Repository](https://github.com/pre-commit/pre-commit-hooks)
- [Black Documentation](https://black.readthedocs.io/)
- [ESLint Documentation](https://eslint.org/)
- [GolangCI-lint Documentation](https://golangci-lint.run/) 