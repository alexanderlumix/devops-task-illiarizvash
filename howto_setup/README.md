# Project Setup Guide

This folder contains detailed documentation for setting up and debugging the project.

## Documentation Structure

### 📋 Main Instructions

1. **[01_initial_setup.md](01_initial_setup.md)** - Initial project setup
   - Installing Docker and Docker Compose
   - Installing Go and Node.js
   - Environment setup
   - Project startup

2. **[02_mongodb_setup.md](02_mongodb_setup.md)** - MongoDB setup
   - MongoDB replica set architecture
   - Solving AVX instruction issues
   - Authentication setup
   - Replica set initialization
   - Monitoring and backup

3. **[03_applications_setup.md](03_applications_setup.md)** - Application setup
   - Node.js application setup
   - Go application setup
   - Environment variables
   - Connection testing
   - Application debugging

### 🔧 Debugging and Diagnostics

4. **[howto_debug.md](howto_debug.md)** - Debugging guide
   - Common debugging commands
   - MongoDB debugging
   - Application debugging
   - Network and HAProxy debugging
   - Performance monitoring
   - Automated debugging scripts

### ⚠️ Known Issues

5. **[known_issues.md](known_issues.md)** - Known issues and their solutions
   - MongoDB issues
   - Docker issues
   - Application issues
   - Network issues
   - Common issues and preventive measures

## Quick Start

### Minimal Setup

1. **Install dependencies:**
   ```bash
   # Docker and Docker Compose
   sudo apt install docker.io docker-compose -y
   sudo systemctl start docker
   sudo usermod -aG docker $USER
   
   # Go
   sudo snap install go --classic
   export PATH=/snap/bin:$PATH
   ```

2. **Set up environment:**
   ```bash
   cp env.example .env
   pip install -r scripts/requirements.txt
   cd app-node && npm install && cd ..
   cd app-go && go mod tidy && cd ..
   ```

3. **Start the project:**
   ```bash
   docker-compose up -d
   docker exec mongo-0 mongo --eval "rs.initiate({_id: 'rs0', members: [{_id: 0, host: 'mongo-0:27017'}, {_id: 1, host: 'mongo-1:27017'}, {_id: 2, host: 'mongo-2:27017'}]})"
   ```

### Health Check

```bash
# Container status
docker ps

# MongoDB replica set status
docker exec mongo-0 mongo --eval "rs.status().members.forEach(function(m) { print(m.name + ': ' + m.stateStr); });"

# Application logs
docker logs app-node --tail 10
docker logs app-go --tail 10
```

## Frequently Used Commands

### Project Management

```bash
# Start all services
docker-compose up -d

# Stop all services
docker-compose down

# Rebuild and start
docker-compose up -d --build

# View logs
docker-compose logs -f
```

### MongoDB Debugging

```bash
# Replica set status
docker exec mongo-0 mongo --eval "rs.status()"

# Connect to MongoDB
docker exec -it mongo-0 mongo

# Restart MongoDB
docker-compose restart mongo-0 mongo-1 mongo-2
```

### Application Debugging

```bash
# Application logs
docker logs app-node
docker logs app-go

# Enter containers
docker exec -it app-node bash
docker exec -it app-go bash

# Environment variables
docker exec app-node env | grep MONGO
```

## Troubleshooting

### If something doesn't work

1. **Check system status:**
   ```bash
   docker ps -a
   docker logs mongo-0
   docker logs app-node
   ```

2. **Use the debugging guide:**
   - See [howto_debug.md](howto_debug.md) for detailed debugging commands

3. **Check known issues:**
   - See [known_issues.md](known_issues.md) for typical problem solutions

### Emergency Restart

```bash
# Complete system restart
docker-compose down -v
docker system prune -f
docker-compose up -d --build
```

## Project Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   app-node      │    │   app-go        │    │   MongoDB       │
│   (Node.js)     │    │   (Go)          │    │   Replica Set   │
│   Port: 3000    │    │   Port: 8080    │    │   Ports: 27030- │
└─────────────────┘    └─────────────────┘    │   27032         │
         │                      │              └─────────────────┘
         └──────────────────────┼────────────────────┘
                                │
                    ┌─────────────────┐
                    │   HAProxy       │
                    │   (Load Balancer)│
                    │   Port: 27017   │
                    └─────────────────┘
```

## Support

If you encounter problems:

1. Check the [known_issues.md](known_issues.md) section
2. Use commands from [howto_debug.md](howto_debug.md)
3. Check container logs: `docker logs <container-name>`
4. Make sure all dependencies are installed and updated

## Documentation Updates

This documentation is updated when:
- Project architecture changes
- New components are added
- New issues are discovered and their solutions
- Setup procedures change 