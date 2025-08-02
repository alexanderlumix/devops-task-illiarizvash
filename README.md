# DevOps Engineer Task - MongoDB Replica Set with Applications

This project demonstrates a complete MongoDB replica set setup with Node.js and Go applications for product management, following DevOps best practices and enterprise-grade standards.

## 📋 Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Quick Start](#quick-start)
- [Documentation](#documentation)
- [Development](#development)
- [Security](#security)
- [Quality Assurance](#quality-assurance)
- [Contributing](#contributing)

## 🚀 Overview

This project implements a microservices-based architecture with MongoDB replica set for high availability and data consistency. It serves as a comprehensive example of modern DevOps practices including containerization, CI/CD, security scanning, and code quality tools.

## 🏗️ Architecture

The project consists of:
- **MongoDB Replica Set**: 3-node replica set with HAProxy load balancer
- **Node.js Application**: Creates products in the database
- **Go Application**: Reads and displays products from the database
- **Python Scripts**: Automation for replica set initialization and user management

For detailed architecture information, see [docs/architecture.md](docs/architecture.md).

## ⚡ Quick Start

### Prerequisites

- Docker Engine 20.10+
- Docker Compose 2.0+
- Python 3.8+
- Node.js 16+
- Go 1.19+

### Setup Instructions

1. **Clone Repository**
   ```bash
   git clone <repository-url>
   cd devops-task-illiarizvash
   ```

2. **Install Dependencies**
   ```bash
   # Python dependencies
   pip install -r scripts/requirements.txt
   
   # Node.js dependencies
   cd app-node && npm install
   
   # Go dependencies
   cd app-go && go mod tidy
   ```

3. **Start Infrastructure**
   ```bash
   cd mongo
   docker-compose up -d
   ```

4. **Initialize Database**
   ```bash
   cd scripts
   python init_mongo_servers.py
   python create_app_user.py
   ```

5. **Test Applications**
   ```bash
   # Create products
   cd app-node
   node create_product.js
   
   # Read products
   cd app-go
   go run read_products.go
   ```

For detailed deployment instructions, see [docs/deployment.md](docs/deployment.md).

## 📚 Documentation

### Core Documentation
- [**Architecture**](docs/architecture.md) - System architecture and component interactions
- [**Deployment Guide**](docs/deployment.md) - Detailed deployment instructions for different environments
- [**Architecture Decision Records**](docs/decisions.md) - Technical decisions and rationale

### Development
- [**Code Quality Tools**](code-quality/) - Linters, formatters, and quality checkers
- [**Pre-commit Hooks**](demo_pre_commit_hooks/) - Git hooks for automated quality checks
- [**Security Policy**](security/security-policy.md) - Security guidelines and procedures

### CI/CD
- [**GitHub Actions**](.github/workflows/) - Automated testing, building, and deployment
- [**Docker Best Practices**](app-go/Dockerfile) - Optimized container configurations

## 🔧 Development

### Code Quality

The project includes comprehensive code quality tools:

- **Python**: Black, isort, flake8, mypy, bandit
- **Go**: golangci-lint, gofmt, govet
- **JavaScript**: ESLint with security plugins
- **Docker**: Hadolint for Dockerfile validation

### Pre-commit Hooks

Install and configure pre-commit hooks for automated quality checks:

```bash
# Install pre-commit
pip install pre-commit

# Install hooks
pre-commit install

# Run manually
pre-commit run --all-files
```

For detailed configuration, see [demo_pre_commit_hooks/](demo_pre_commit_hooks/).

### Testing

```bash
# Python tests
python -m pytest tests/

# Go tests
cd app-go && go test ./...

# Node.js tests
cd app-node && npm test
```

## 🔒 Security

### Security Features

- **Vulnerability Scanning**: Automated scanning with Trivy and Bandit
- **Secret Detection**: Pre-commit hooks to prevent secret leaks
- **Container Security**: Non-root users and minimal base images
- **Network Security**: Isolated Docker networks and firewall rules

### Security Tools

- **Static Analysis**: Bandit, ESLint security, gosec
- **Container Scanning**: Trivy, Hadolint
- **Secret Management**: detect-secrets pre-commit hook

For detailed security policies, see [security/](security/).

## 🎯 Quality Assurance

### Automated Checks

The project includes comprehensive quality assurance:

- **Code Formatting**: Automated formatting with Black, gofmt, ESLint
- **Static Analysis**: Type checking, linting, security scanning
- **Testing**: Unit tests, integration tests, security tests
- **Documentation**: Automated documentation generation and validation

### Quality Metrics

- **Code Coverage**: Automated test coverage reporting
- **Security Score**: Vulnerability scanning and scoring
- **Performance**: Automated performance testing
- **Compliance**: Automated compliance checking

## 🤝 Contributing

### Development Workflow

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/amazing-feature`
3. **Make your changes** following the coding standards
4. **Run quality checks**: `pre-commit run --all-files`
5. **Write tests** for new functionality
6. **Update documentation** as needed
7. **Submit a pull request**

### Code Standards

- Follow the established code formatting rules
- Write comprehensive tests
- Update documentation for any changes
- Follow security best practices
- Use conventional commit messages

### Quality Gates

All contributions must pass:
- ✅ Code formatting checks
- ✅ Static analysis
- ✅ Security scanning
- ✅ Unit tests
- ✅ Integration tests
- ✅ Documentation validation

## 📊 Monitoring

### Health Checks

```bash
# Check replica set status
python scripts/check_replicaset_status.py

# Check application logs
docker logs product-creator
docker logs product-reader

# Monitor HAProxy
docker logs haproxy-lb
```

### Metrics

- **Database**: Replica set status, connection health
- **Applications**: Response times, error rates
- **Infrastructure**: Container health, resource usage

## 🚨 Troubleshooting

### Common Issues

1. **Replica Set Not Initialized**
   ```bash
   python scripts/init_mongo_servers.py
   ```

2. **Connection Refused**
   - Check if containers are running: `docker ps`
   - Verify port mappings
   - Check firewall settings

3. **Authentication Failed**
   - Verify user credentials
   - Check keyfile permissions
   - Ensure proper replica set configuration

For more troubleshooting information, see [docs/deployment.md](docs/deployment.md).

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- MongoDB for the excellent database technology
- Docker for containerization platform
- HAProxy for load balancing
- All open-source contributors to the tools used in this project