# Known Issues and Solutions

## MongoDB Issues

### 1. MongoDB 6.0 requires AVX instructions

**Symptoms:**
```
FATAL: kernel too old
The kernel version is 4.19.0-6-amd64, but MongoDB 6.0 requires kernel version 4.19.0-10 or later
```

**Cause:** MongoDB 6.0 requires a CPU with AVX instruction support.

**Solution:**
1. Change MongoDB version in `docker-compose.yml`:
   ```yaml
   image: mongo:4.4  # instead of mongo:6.0
   ```

2. Update healthcheck commands:
   ```yaml
   test: ["CMD", "mongo", "--eval", "db.adminCommand('ping')"]  # instead of mongosh
   ```

### 2. MongoDB key file permission issues

**Symptoms:**
```
Error parsing keyFile /etc/mongo-keyfile: Invalid key file
```

**Cause:** Incorrect permissions or corrupted key file.

**Solution:**
```bash
# Recreate key
openssl rand -base64 756 > mongo/mongo-keyfile
sudo chmod 400 mongo/mongo-keyfile

# Restart containers
docker-compose restart mongo-0 mongo-1 mongo-2
```

### 3. Replica set not initializing

**Symptoms:**
```
"errmsg" : "No configuration found. Not yet initialized."
```

**Cause:** Replica set was not initialized.

**Solution:**
```bash
docker exec mongo-0 mongo --eval "rs.initiate({_id: 'rs0', members: [{_id: 0, host: 'mongo-0:27017'}, {_id: 1, host: 'mongo-1:27017'}, {_id: 2, host: 'mongo-2:27017'}]})"
```

### 4. All nodes in SECONDARY state

**Symptoms:**
```
mongo-0: SECONDARY
mongo-1: SECONDARY
mongo-2: SECONDARY
```

**Cause:** No primary node in replica set.

**Solution:**
```bash
# Unfreeze nodes
docker exec mongo-0 mongo --eval "rs.freeze(0); rs.freeze(1); rs.freeze(2);"

# Check status
docker exec mongo-0 mongo --eval "rs.status().members.forEach(function(m) { print(m.name + ': ' + m.stateStr); });"
```

## Docker Issues

### 1. Docker Compose not found

**Symptoms:**
```
docker-compose: command not found
```

**Cause:** Docker Compose is not installed or not in PATH.

**Solution:**
```bash
# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
sudo ln -sf /usr/local/bin/docker-compose /usr/bin/docker-compose
```

### 2. Docker daemon issues

**Symptoms:**
```
Cannot connect to the Docker daemon
```

**Cause:** Docker daemon is not running or user is not in the docker group.

**Solution:**
```bash
# Start Docker daemon
sudo systemctl start docker
sudo systemctl enable docker

# Add user to docker group
sudo usermod -aG docker $USER
# Reboot or log out/in
```

### 3. Docker network issues

**Symptoms:**
```
Error response from daemon: network not found
```

**Cause:** Docker network not created or corrupted.

**Solution:**
```bash
# Recreate network
docker-compose down
docker network prune -f
docker-compose up -d
```

## Application Issues

### 1. Node.js application cannot connect to MongoDB

**Symptoms:**
```
MongoNetworkError: connect ECONNREFUSED
```

**Cause:** Incorrect host or port in environment variables.

**Solution:**
1. Check environment variables:
   ```bash
   docker exec app-node env | grep MONGO
   ```

2. Make sure the correct host is used:
   ```bash
   # In docker-compose.yml
   environment:
     - MONGO_HOST=mongo-1  # instead of haproxy
   ```

3. Rebuild container:
   ```bash
   docker-compose up -d --build app-node
   ```

### 2. Go application cannot connect to MongoDB

**Symptoms:**
```
no reachable servers
```

**Cause:** Issues connecting to replica set.

**Solution:**
1. Check replica set:
   ```bash
   docker exec mongo-0 mongo --eval "rs.status()"
   ```

2. Make sure the correct replica set is specified:
   ```bash
   # In docker-compose.yml
   environment:
     - MONGO_REPLICA_SET=rs0
   ```

### 3. Dependency issues

**Symptoms:**
```
Cannot find module 'mongodb'
```

**Cause:** Dependencies not installed.

**Solution:**
```bash
# Rebuild container
docker-compose build --no-cache app-node

# Or install dependencies manually
docker exec app-node npm install
```

## HAProxy Issues

### 1. HAProxy cannot connect to MongoDB

**Symptoms:**
```
Connection refused
```

**Cause:** Incorrect backend server configuration.

**Solution:**
1. Use IP addresses instead of container names:
   ```conf
   server mongo1 172.18.0.2:27017 check
   server mongo2 172.18.0.3:27017 check
   server mongo3 172.18.0.4:27017 check
   ```

2. Or remove HAProxy from the chain and connect applications directly to MongoDB.

### 2. HAProxy does not start

**Symptoms:**
```
haproxy exited with code 1
```

**Cause:** Error in HAProxy configuration.

**Solution:**
```bash
# Check configuration
docker exec haproxy haproxy -c -f /usr/local/etc/haproxy/haproxy.cfg

# View logs
docker logs haproxy
```

## Go Issues

### 1. Go not found

**Symptoms:**
```
go: command not found
```

**Cause:** Go is not installed or not in PATH.

**Solution:**
```bash
# Install Go
sudo snap install go --classic

# Update PATH
export PATH=/snap/bin:$PATH

# Check
go version
```

### 2. Go version issues

**Symptoms:**
```
go: go.mod requires go >= 1.16
```

**Cause:** Old version of Go installed.

**Solution:**
```bash
# Update Go
sudo snap refresh go

# Check version
go version
```

## Network Issues

### 1. Containers cannot communicate with each other

**Symptoms:**
```
ping: unknown host mongo-1
```

**Cause:** DNS resolution issues in Docker network.

**Solution:**
1. Check network:
   ```bash
   docker network ls
   docker network inspect devops-task-illiarizvash_mongo-network
   ```

2. Use IP addresses instead of names:
   ```bash
   docker inspect mongo-1 | grep IPAddress
   ```

### 2. Port already in use

**Symptoms:**
```
Bind for 0.0.0.0:27017 failed: port is already allocated
```

**Cause:** Port is used by another process.

**Solution:**
```bash
# Find process using the port
sudo netstat -tlnp | grep 27017

# Stop the process or change port in docker-compose.yml
```

## Environment Variable Issues

### 1. Environment variables not loaded

**Symptoms:**
```
undefined environment variable
```

**Cause:** .env file not created or misconfigured.

**Solution:**
```bash
# Create .env file
cp env.example .env

# Check contents
cat .env

# Restart with rebuild
docker-compose up -d --build
```

### 2. Environment variables not updating

**Symptoms:** Old variable values in container.

**Cause:** Container caches environment variables.

**Solution:**
```bash
# Rebuild container
docker-compose up -d --build app-node

# Or restart with forced update
docker-compose down
docker-compose up -d
```

## General Issues

### 1. Not enough disk space

**Symptoms:**
```
no space left on device
```

**Solution:**
```bash
# Clean Docker
docker system prune -f

# Remove unused images
docker image prune -f

# Clean volumes
docker volume prune -f
```

### 2. Not enough memory

**Symptoms:**
```
Cannot allocate memory
```

**Solution:**
1. Increase Docker memory:
   ```bash
   # In /etc/docker/daemon.json
   {
     "memory": "2g"
   }
   ```

2. Restart Docker:
   ```bash
   sudo systemctl restart docker
   ```

### 3. Permission issues

**Symptoms:**
```
Permission denied
```

**Solution:**
```bash
# Check access rights
ls -la mongo/mongo-keyfile

# Fix permissions
sudo chmod 400 mongo/mongo-keyfile
sudo chown 999:999 mongo/mongo-keyfile
```

## Preventive Measures

### Regular Cleanup

```bash
# Weekly cleanup
docker system prune -f
docker image prune -f
docker volume prune -f
```

### Resource Monitoring

```bash
# Check resource usage
docker stats

# Check disk space
df -h
```

### Configuration Backup

```bash
# Create backup
cp docker-compose.yml docker-compose.yml.backup
cp .env .env.backup
``` 