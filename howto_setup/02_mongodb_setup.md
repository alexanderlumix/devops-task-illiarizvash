# MongoDB Setup

## MongoDB Architecture

The project uses a MongoDB replica set with three nodes:
- **mongo-0** (port 27030)
- **mongo-1** (port 27031) 
- **mongo-2** (port 27032)

## MongoDB 6.0 Issues

MongoDB 6.0 requires a CPU with AVX instruction support. If your system doesn't support AVX, use MongoDB 4.4.

### Solving AVX Issues

1. **Changing MongoDB version in docker-compose.yml:**
   ```yaml
   # Replace all occurrences of
   image: mongo:6.0
   # with
   image: mongo:4.4
   ```

2. **Changing healthcheck commands:**
   ```yaml
   # Replace
   test: ["CMD", "mongosh", "--eval", "db.adminCommand('ping')"]
   # with
   test: ["CMD", "mongo", "--eval", "db.adminCommand('ping')"]
   ```

## Authentication Setup

### Option 1: Without Authentication (for development)

Remove from docker-compose.yml:
- `environment` sections with `MONGO_INITDB_ROOT_USERNAME` and `MONGO_INITDB_ROOT_PASSWORD`
- `--keyFile /etc/mongo-keyfile` from commands
- mounting `./mongo/mongo-keyfile:/etc/mongo-keyfile`

### Option 2: With Authentication

1. **Creating replica set key:**
   ```bash
   openssl rand -base64 756 > mongo/mongo-keyfile
   sudo chmod 400 mongo/mongo-keyfile
   ```

2. **Checking access rights:**
   ```bash
   ls -la mongo/mongo-keyfile
   ```

## Replica Set Initialization

### Automatic Initialization

```bash
docker exec mongo-0 mongo --eval "rs.initiate({_id: 'rs0', members: [{_id: 0, host: 'mongo-0:27017'}, {_id: 1, host: 'mongo-1:27017'}, {_id: 2, host: 'mongo-2:27017'}]})"
```

### Status Check

```bash
# Check all replica set members
docker exec mongo-0 mongo --eval "rs.status().members.forEach(function(m) { print(m.name + ': ' + m.stateStr); });"

# Check primary node
docker exec mongo-0 mongo --eval "rs.isMaster()"
```

### Manual Replica Set Management

```bash
# Freeze node (prevent primary election)
docker exec mongo-0 mongo --eval "rs.freeze(0)"

# Unfreeze node
docker exec mongo-0 mongo --eval "rs.freeze(0)"

# Force primary step down
docker exec mongo-0 mongo --eval "rs.stepDown()"
```

## Connecting to MongoDB

### From Container

```bash
# Connect to primary
docker exec -it mongo-0 mongo

# Connect to specific node
docker exec -it mongo-1 mongo
```

### From Host

```bash
# Connect to primary via HAProxy
mongo localhost:27017

# Direct connection to node
mongo localhost:27030
```

## MongoDB Monitoring

### Checking Logs

```bash
# Specific node logs
docker logs mongo-0

# All nodes logs
docker logs mongo-0 mongo-1 mongo-2
```

### Performance Check

```bash
# Operation statistics
docker exec mongo-0 mongo --eval "db.serverStatus()"

# Replica set statistics
docker exec mongo-0 mongo --eval "rs.status()"
```

## Backup

```bash
# Create dump
docker exec mongo-0 mongodump --out /tmp/backup

# Restore from dump
docker exec mongo-0 mongorestore /tmp/backup
```

## Troubleshooting

### Containers Not Starting

1. **Check logs:**
   ```bash
   docker logs mongo-0
   ```

2. **Check key access rights:**
   ```bash
   ls -la mongo/mongo-keyfile
   sudo chmod 400 mongo/mongo-keyfile
   ```

3. **Recreate key:**
   ```bash
   openssl rand -base64 756 > mongo/mongo-keyfile
   sudo chmod 400 mongo/mongo-keyfile
   ```

### Replica Set Not Initializing

1. **Check node availability:**
   ```bash
   docker exec mongo-0 ping mongo-1
   docker exec mongo-0 ping mongo-2
   ```

2. **Check ports:**
   ```bash
   docker exec mongo-0 netstat -tlnp | grep 27017
   ```

3. **Restart with data cleanup:**
   ```bash
   docker-compose down -v
   docker-compose up -d
   ``` 