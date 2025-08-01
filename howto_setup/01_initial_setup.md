# Initial Project Setup

## Prerequisites

### Installing Docker and Docker Compose

1. **Installing Docker:**
   ```bash
   sudo apt update
   sudo apt install docker.io -y
   sudo systemctl start docker
   sudo systemctl enable docker
   sudo usermod -aG docker $USER
   ```

2. **Installing Docker Compose:**
   ```bash
   sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
   sudo chmod +x /usr/local/bin/docker-compose
   sudo ln -sf /usr/local/bin/docker-compose /usr/bin/docker-compose
   ```

3. **Verifying installation:**
   ```bash
   docker --version
   docker-compose --version
   ```

### Installing Go

```bash
sudo snap install go --classic
export PATH=/snap/bin:$PATH
go version
```

### Installing Node.js (if not installed)

```bash
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs
node --version
npm --version
```

## Environment Setup

1. **Copying environment file:**
   ```bash
   cp env.example .env
   ```

2. **Installing Python dependencies:**
   ```bash
   pip install -r scripts/requirements.txt
   ```

3. **Installing Node.js dependencies:**
   ```bash
   cd app-node
   npm install
   cd ..
   ```

4. **Installing Go dependencies:**
   ```bash
   cd app-go
   export PATH=/snap/bin:$PATH
   go mod tidy
   cd ..
   ```

## Starting the Project

1. **Starting all services:**
   ```bash
   docker-compose up -d
   ```

2. **Checking container status:**
   ```bash
   docker ps
   ```

3. **Initializing MongoDB replica set:**
   ```bash
   docker exec mongo-0 mongo --eval "rs.initiate({_id: 'rs0', members: [{_id: 0, host: 'mongo-0:27017'}, {_id: 1, host: 'mongo-1:27017'}, {_id: 2, host: 'mongo-2:27017'}]})"
   ```

4. **Checking replica set status:**
   ```bash
   docker exec mongo-0 mongo --eval "rs.status().members.forEach(function(m) { print(m.name + ': ' + m.stateStr); });"
   ```

## Health Check

1. **Checking MongoDB:**
   ```bash
   docker exec mongo-0 mongo --eval "db.adminCommand('ping')"
   ```

2. **Checking applications:**
   ```bash
   docker logs app-node
   docker logs app-go
   ```

## Stopping the Project

```bash
docker-compose down
``` 