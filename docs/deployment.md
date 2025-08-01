# Deployment Guide

## Prerequisites

- Docker Engine 20.10+
- Docker Compose 2.0+
- Python 3.8+
- Node.js 16+
- Go 1.19+

## Environment Setup

### Development Environment

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

### Production Environment

1. **Environment Variables**
   ```bash
   export MONGO_ADMIN_USER=admin
   export MONGO_ADMIN_PASSWORD=secure_password
   export APP_DB_USER=appuser
   export APP_DB_PASSWORD=app_password
   ```

2. **Security Configuration**
   - Update MongoDB keyfile permissions
   - Configure firewall rules
   - Set up SSL certificates

3. **Monitoring Setup**
   - Configure log aggregation
   - Set up health checks
   - Enable metrics collection

## Deployment Steps

### Step 1: Infrastructure Deployment
```bash
# Deploy MongoDB replica set
docker-compose -f mongo/docker-compose.yml up -d

# Wait for services to be healthy
docker-compose -f mongo/docker-compose.yml ps
```

### Step 2: Database Initialization
```bash
# Initialize replica set
python scripts/init_mongo_servers.py

# Create application user
python scripts/create_app_user.py

# Verify setup
python scripts/check_replicaset_status.py
```

### Step 3: Application Deployment
```bash
# Deploy Node.js application
cd app-node
docker build -t product-creator .
docker run -d --name product-creator product-creator

# Deploy Go application
cd app-go
docker build -t product-reader .
docker run -d --name product-reader product-reader
```

## Health Checks

### Database Health
```bash
# Check replica set status
python scripts/check_replicaset_status.py

# Check MongoDB logs
docker logs mongo-0
docker logs mongo-1
docker logs mongo-2
```

### Application Health
```bash
# Check application logs
docker logs product-creator
docker logs product-reader

# Test connectivity
curl http://localhost:8080/health
```

## Troubleshooting

### Common Issues

1. **Replica Set Not Initialized**
   ```bash
   python scripts/init_mongo_servers.py
   ```

2. **Connection Refused**
   - Check if containers are running
   - Verify port mappings
   - Check firewall settings

3. **Authentication Failed**
   - Verify user credentials
   - Check keyfile permissions
   - Ensure proper replica set configuration

### Log Analysis

```bash
# View all container logs
docker-compose logs

# Follow specific service logs
docker-compose logs -f mongo-0

# Check application logs
docker logs product-creator --tail 100
```

## Rollback Procedures

1. **Application Rollback**
   ```bash
   docker stop product-creator product-reader
   docker run -d --name product-creator product-creator:previous
   docker run -d --name product-reader product-reader:previous
   ```

2. **Database Rollback**
   ```bash
   # Stop all services
   docker-compose down
   
   # Restore from backup
   docker run --rm -v mongo_backup:/backup mongo:7.0.16 mongorestore /backup
   
   # Restart services
   docker-compose up -d
   ``` 