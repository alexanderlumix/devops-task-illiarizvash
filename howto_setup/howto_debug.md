# Debugging Guide

## General Debugging Commands

### System Status Check

```bash
# Status of all containers
docker ps -a

# Status of running containers
docker ps

# Resource usage statistics
docker stats

# Network information
docker network ls
docker network inspect devops-task-illiarizvash_mongo-network
```

### Log Checking

```bash
# Logs of all containers
docker-compose logs

# Logs of specific container
docker logs mongo-0
docker logs app-node
docker logs app-go

# Logs with recent lines
docker logs mongo-0 --tail 50

# Logs with timestamps
docker logs app-node --timestamps

# Follow logs in real time
docker logs -f app-node
```

## MongoDB Debugging

### Replica Set Status Check

```bash
# Status of all replica set members
docker exec mongo-0 mongo --eval "rs.status().members.forEach(function(m) { print(m.name + ': ' + m.stateStr); });"

# Primary information
docker exec mongo-0 mongo --eval "rs.isMaster()"

# Replica set configuration
docker exec mongo-0 mongo --eval "rs.conf()"
```

### Connection Check

```bash
# Check node availability
docker exec mongo-0 ping mongo-1
docker exec mongo-0 ping mongo-2

# Check ports
docker exec mongo-0 netstat -tlnp | grep 27017

# Test MongoDB connection
docker exec mongo-0 mongo --eval "db.adminCommand('ping')"
```

### Problem Diagnostics

```bash
# Check key access rights
ls -la mongo/mongo-keyfile

# Check key contents
file mongo/mongo-keyfile

# Check MongoDB configuration
docker exec mongo-0 cat /etc/mongod.conf
```

### Fix Commands

```bash
# Recreate key
openssl rand -base64 756 > mongo/mongo-keyfile
sudo chmod 400 mongo/mongo-keyfile

# Restart MongoDB containers
docker-compose restart mongo-0 mongo-1 mongo-2

# Initialize replica set
docker exec mongo-0 mongo --eval "rs.initiate({_id: 'rs0', members: [{_id: 0, host: 'mongo-0:27017'}, {_id: 1, host: 'mongo-1:27017'}, {_id: 2, host: 'mongo-2:27017'}]})"
```

## Application Debugging

### Environment Variables Check

```bash
# Node.js application environment variables
docker exec app-node env | grep MONGO

# Go application environment variables
docker exec app-go env | grep MONGO

# Check specific variable
docker exec app-node sh -c 'echo $MONGO_HOST'
```

### Connection Testing

```bash
# Test Node.js connection to MongoDB
docker exec app-node node -e "
const { MongoClient } = require('mongodb');
const client = new MongoClient('mongodb://mongo-1:27017/appdb');
client.connect()
  .then(() => console.log('Connected!'))
  .catch(e => console.log('Error:', e.message))
  .finally(() => client.close());
"

# Test Go connection to MongoDB
docker exec app-go go run read_products.go
```

### Process Check

```bash
# Processes in container
docker exec app-node ps aux

# Network connections
docker exec app-node netstat -tlnp

# Memory and CPU usage
docker stats app-node app-go
```

### Entering Containers for Debugging

```bash
# Enter Node.js container
docker exec -it app-node bash

# Enter Go container
docker exec -it app-go bash

# Enter MongoDB container
docker exec -it mongo-0 bash
```

## Network Debugging

### Network Connection Check

```bash
# Docker networks list
docker network ls

# Network information
docker network inspect devops-task-illiarizvash_mongo-network

# Container IP addresses
docker inspect mongo-0 | grep IPAddress
docker inspect app-node | grep IPAddress
```

### Network Connectivity Testing

```bash
# Ping between containers
docker exec app-node ping mongo-1
docker exec mongo-0 ping mongo-1

# Telnet for port checking
docker exec app-node telnet mongo-1 27017

# DNS resolution
docker exec app-node nslookup mongo-1
```

## HAProxy Debugging

### Configuration Check

```bash
# Check configuration syntax
docker exec haproxy haproxy -c -f /usr/local/etc/haproxy/haproxy.cfg

# Check HAProxy status
docker exec haproxy haproxy -c -f /usr/local/etc/haproxy/haproxy.cfg

# HAProxy logs
docker logs haproxy
```

### HAProxy Testing

```bash
# Check MongoDB availability via HAProxy
telnet localhost 27017

# Check HAProxy statistics
docker exec haproxy haproxy -c -f /usr/local/etc/haproxy/haproxy.cfg
```

## Docker Compose Debugging

### Configuration Check

```bash
# Check docker-compose.yml syntax
docker-compose config

# Check environment variables
docker-compose config --env-file .env
```

### Service Management

```bash
# Restart specific service
docker-compose restart app-node

# Rebuild and start
docker-compose up -d --build app-node

# Stop all services
docker-compose down

# Stop with volume removal
docker-compose down -v
```

## Performance Monitoring

### System Resources

```bash
# Container resource usage
docker stats

# Detailed container information
docker inspect app-node

# Disk usage
docker system df
```

### Performance Logs

```bash
# Logs with timestamps
docker logs --timestamps app-node

# Logs for specific period
docker logs --since "2024-01-01T00:00:00" app-node

# Logs from last N lines
docker logs --tail 100 app-node
```

## Emergency Debugging Commands

### Complete System Restart

```bash
# Stop all containers
docker-compose down

# Remove all volumes
docker-compose down -v

# Clean Docker
docker system prune -f

# Restart with rebuild
docker-compose up -d --build
```

### Image Problem Diagnostics

```bash
# Check images
docker images

# Remove unused images
docker image prune -f

# Rebuild without cache
docker-compose build --no-cache
```

### File System Check

```bash
# Check access rights
ls -la mongo/mongo-keyfile

# Check file contents
cat docker-compose.yml
cat .env

# Check file sizes
du -sh *
```

## Automated Debugging Scripts

### System Status Check Script

```bash
#!/bin/bash
echo "=== Docker Status ==="
docker ps -a

echo "=== MongoDB Replica Set Status ==="
docker exec mongo-0 mongo --eval "rs.status().members.forEach(function(m) { print(m.name + ': ' + m.stateStr); });"

echo "=== Application Logs ==="
docker logs app-node --tail 10
docker logs app-go --tail 10

echo "=== Network Connectivity ==="
docker exec app-node ping -c 1 mongo-1
```

### Full Diagnostics Script

```bash
#!/bin/bash
echo "=== Full System Diagnostics ==="

echo "1. Docker Status:"
docker ps -a

echo "2. Network Status:"
docker network ls

echo "3. MongoDB Status:"
docker exec mongo-0 mongo --eval "rs.status()" 2>/dev/null || echo "MongoDB not accessible"

echo "4. Application Environment:"
docker exec app-node env | grep MONGO

echo "5. Recent Logs:"
docker logs app-node --tail 20
docker logs app-go --tail 20
``` 