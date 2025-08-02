# Application Setup

## Application Architecture

The project contains two applications:
- **app-node** - Node.js application (port 3000)
- **app-go** - Go application (port 8080)

## Node.js Application Setup

### Installing Dependencies

```bash
cd app-node
npm install
cd ..
```

### Checking package.json

Make sure `app-node/package.json` has the correct dependencies:
```json
{
  "dependencies": {
    "mongodb": "^4.0.0",
    "express": "^4.17.1"
  }
}
```

### Environment Variables

The application uses the following environment variables:
- `MONGO_HOST` - MongoDB host (default: mongo-1)
- `MONGO_PORT` - MongoDB port (default: 27017)
- `MONGO_DB` - database name (default: appdb)
- `MONGO_USER` - MongoDB user
- `MONGO_PASSWORD` - MongoDB password
- `MONGO_DIRECT_CONNECTION` - direct connection (default: true)

### Testing Connection

```bash
# Check environment variables
docker exec app-node env | grep MONGO

# Test MongoDB connection
docker exec app-node node -e "
const { MongoClient } = require('mongodb');
const client = new MongoClient('mongodb://mongo-1:27017/appdb');
client.connect()
  .then(() => console.log('Connected!'))
  .catch(e => console.log('Error:', e.message))
  .finally(() => client.close());
"
```

## Go Application Setup

### Installing Go

```bash
sudo snap install go --classic
export PATH=/snap/bin:$PATH
go version
```

### Installing Dependencies

```bash
cd app-go
export PATH=/snap/bin:$PATH
go mod tidy
cd ..
```

### Checking go.mod

Make sure `app-go/go.mod` has the correct dependencies:
```go
require (
    go.mongodb.org/mongo-driver v1.10.0
    github.com/gorilla/mux v1.8.0
)
```

### Environment Variables

The application uses the following environment variables:
- `MONGO_HOST` - MongoDB host (default: mongo-1)
- `MONGO_PORT` - MongoDB port (default: 27017)
- `MONGO_DB` - database name (default: appdb)
- `MONGO_USER` - MongoDB user
- `MONGO_PASSWORD` - MongoDB password
- `MONGO_REPLICA_SET` - replica set name (default: rs0)

### Testing Connection

```bash
# Check environment variables
docker exec app-go env | grep MONGO

# Test MongoDB connection
docker exec app-go go run read_products.go
```

## Starting Applications

### Via Docker Compose

```bash
# Start all applications
docker-compose up -d app-node app-go

# Rebuild and start
docker-compose up -d --build app-node app-go
```

### Individual Startup

```bash
# Node.js application
cd app-node
npm start

# Go application
cd app-go
export PATH=/snap/bin:$PATH
go run read_products.go
```

## Health Check

### Checking Logs

```bash
# Node.js application logs
docker logs app-node

# Go application logs
docker logs app-go

# Logs with recent lines
docker logs app-node --tail 20
docker logs app-go --tail 20
```

### Checking Endpoints

```bash
# Check Node.js application
curl http://localhost:3000/health

# Check Go application
curl http://localhost:8080/health
```

### Testing Functionality

```bash
# Create product via Node.js
docker exec app-node node create_product.js

# Read products via Go
docker exec app-go go run read_products.go
```

## Debugging Applications

### Entering Container

```bash
# Node.js container
docker exec -it app-node bash

# Go container
docker exec -it app-go bash
```

### Checking Processes

```bash
# Processes in container
docker exec app-node ps aux

# Network connections
docker exec app-node netstat -tlnp
```

### Checking Files

```bash
# Directory contents
docker exec app-node ls -la

# Check configuration
docker exec app-node cat package.json
```

## Performance Monitoring

### Resource Usage

```bash
# Container statistics
docker stats app-node app-go

# Detailed information
docker inspect app-node
docker inspect app-go
```

### Performance Logs

```bash
# Logs with timestamps
docker logs app-node --timestamps

# Logs from last 100 lines
docker logs app-go --tail 100
```

## Troubleshooting

### Application Not Starting

1. **Check logs:**
   ```bash
   docker logs app-node
   ```

2. **Check environment variables:**
   ```bash
   docker exec app-node env | grep MONGO
   ```

3. **Check MongoDB connection:**
   ```bash
   docker exec app-node ping mongo-1
   ```

### Dependency Issues

1. **Rebuild containers:**
   ```bash
   docker-compose build --no-cache app-node app-go
   ```

2. **Clear npm cache:**
   ```bash
   docker exec app-node npm cache clean --force
   ```

3. **Update dependencies:**
   ```bash
   docker exec app-node npm update
   docker exec app-go go mod tidy
   ```

### MongoDB Connection Issues

1. **Check MongoDB availability:**
   ```bash
   docker exec mongo-1 mongo --eval "db.adminCommand('ping')"
   ```

2. **Check replica set:**
   ```bash
   docker exec mongo-0 mongo --eval "rs.status()"
   ```

3. **Test connection from application:**
   ```bash
   docker exec app-node node -e "
   const { MongoClient } = require('mongodb');
   const client = new MongoClient('mongodb://mongo-1:27017/appdb');
   client.connect()
     .then(() => console.log('Connected!'))
     .catch(e => console.log('Error:', e.message))
     .finally(() => client.close());
   "
   ``` 