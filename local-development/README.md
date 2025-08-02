# Local Development

This folder contains scripts for automating local setup and cleanup of the project environment.

## Scripts

### 📦 `setup.sh` - Installation and Initialization

Automatic script for installing and configuring the local project environment.

#### Features:
- ✅ Check and install Docker and Docker Compose
- ✅ Check and install Go and Node.js
- ✅ Install Python dependencies
- ✅ Configure environment variables
- ✅ Install application dependencies
- ✅ Create MongoDB key
- ✅ Start project via Docker Compose
- ✅ Initialize MongoDB replica set
- ✅ Test applications
- ✅ Check health checks

#### Usage:

```bash
# Basic installation
./local-development/setup.sh

# Skip dependency installation
./local-development/setup.sh --skip-deps

# Skip MongoDB initialization
./local-development/setup.sh --skip-mongo

# Force reinstallation
./local-development/setup.sh --force

# Show help
./local-development/setup.sh --help
```

#### What the script does:

1. **Environment Check**
   - Checks for Docker, Docker Compose, Go, Node.js
   - Installs missing components

2. **Project Setup**
   - Creates `.env` file from `env.example`
   - Installs application dependencies
   - Creates MongoDB key

3. **System Startup**
   - Stops existing containers
   - Starts project via Docker Compose
   - Initializes MongoDB replica set

4. **Testing**
   - Checks application health checks
   - Tests MongoDB connections
   - Outputs system status

### 🧹 `teardown.sh` - Environment Cleanup

Script for complete cleanup of the local project environment.

### 📊 `status.sh` - Status Check

Script for quick status check of the local project environment.

#### Features:
- ✅ Stop and remove containers
- ✅ Remove Docker images
- ✅ Remove Docker volumes
- ✅ Remove Docker networks
- ✅ Clean Docker system
- ✅ Remove local files
- ✅ Clean logs
- ✅ Complete Docker data cleanup (optional)

#### Usage:

```bash
# Normal cleanup with confirmation
./local-development/teardown.sh

# Force cleanup without confirmation
./local-development/teardown.sh --force

# Complete cleanup (including Docker data)
./local-development/teardown.sh --full

# Show help
./local-development/teardown.sh --help
```

#### Features:
- ✅ Check Docker and Docker Compose
- ✅ Check containers and their status
- ✅ Check MongoDB replica set
- ✅ Check application health checks
- ✅ Check files and dependencies
- ✅ Check logs and resources
- ✅ Final status summary

#### Usage:

```bash
# Complete check
./local-development/status.sh

# Quick check
./local-development/status.sh --quick

# Detailed check
./local-development/status.sh --verbose

# Show help
./local-development/status.sh --help
```

#### What the script removes:

1. **Docker Resources**
   - Project containers
   - Project images
   - Project volumes
   - Project networks

2. **Local Files**
   - `.env` file
   - MongoDB key
   - `node_modules` (optional)
   - Go cache (optional)

3. **System Resources**
   - Unused Docker resources
   - Docker logs

## Quick Start

### First Installation:

```bash
# Make scripts executable
chmod +x local-development/setup.sh
chmod +x local-development/teardown.sh

# Run installation
./local-development/setup.sh
```

### Reinstallation:

```bash
# Clean environment
./local-development/teardown.sh

# Install again
./local-development/setup.sh
```

### Status Check:

```bash
# Quick check
./local-development/status.sh --quick

# Complete check
./local-development/status.sh

# Check containers
docker ps

# Check logs
docker logs app-node
docker logs app-go

# Check MongoDB
docker exec mongo-0 mongo --eval "rs.status()"
```

## Troubleshooting

### Installation Issues:

1. **Docker not installed**
   ```bash
   sudo apt update
   sudo apt install docker.io
   sudo systemctl start docker
   sudo usermod -aG docker $USER
   ```

2. **Permission issues**
   ```bash
   sudo chown $USER:$USER -R .
   chmod +x local-development/*.sh
   ```

3. **Port issues**
   ```bash
   # Check occupied ports
   sudo netstat -tlnp | grep -E "(3000|8080|27030|27031|27032)"
   ```

### MongoDB Issues:

1. **Replica set not initializing**
   ```bash
   docker exec mongo-0 mongo --eval "rs.initiate({_id: 'rs0', members: [{_id: 0, host: 'mongo-0:27017'}, {_id: 1, host: 'mongo-1:27017'}, {_id: 2, host: 'mongo-2:27017'}]})"
   ```

2. **Key issues**
   ```bash
   openssl rand -base64 756 > mongo/mongo-keyfile
   sudo chmod 400 mongo/mongo-keyfile
   ```

### Application Issues:

1. **Health checks failing**
   ```bash
   # Check logs
   docker logs app-node --tail 20
   docker logs app-go --tail 20
   
   # Restart applications
   docker-compose restart app-node app-go
   ```

2. **MongoDB connection issues**
   ```bash
   # Check environment variables
   docker exec app-node env | grep MONGO
   docker exec app-go env | grep MONGO
   ```

## Logs and Debugging

### Viewing Logs:

```bash
# All container logs
docker-compose logs

# Specific container logs
docker logs app-node
docker logs app-go
docker logs mongo-0

# Logs with last lines
docker logs app-node --tail 50
```

### Debugging:

```bash
# Enter container
docker exec -it app-node bash
docker exec -it app-go bash
docker exec -it mongo-0 bash

# Check processes
docker exec app-node ps aux
docker exec app-go ps aux
```

## Automation

### Adding to .bashrc:

```bash
# Add aliases to ~/.bashrc
echo 'alias dev-setup="./local-development/setup.sh"' >> ~/.bashrc
echo 'alias dev-clean="./local-development/teardown.sh"' >> ~/.bashrc
echo 'alias dev-status="./local-development/status.sh"' >> ~/.bashrc
source ~/.bashrc
```

### Using Aliases:

```bash
# Installation
dev-setup

# Cleanup
dev-clean

# Cleanup with confirmation
dev-clean --force

# Status check
dev-status --quick
```

## Security

### Important Notes:

1. **setup.sh script requires sudo for package installation**
2. **teardown.sh script removes all project data**
3. **Complete cleanup (--full) removes all Docker data**
4. **Always make backups before cleanup**

### Recommendations:

1. **Use a virtual machine for development**
2. **Regularly backup important data**
3. **Check logs before cleanup**
4. **Use --force only when necessary** 