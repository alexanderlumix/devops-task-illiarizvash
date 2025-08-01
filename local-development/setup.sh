#!/bin/bash

# Local Development Setup Script
# Automatic installation and initialization of local project environment

set -e  # Stop on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check that script is run from project root directory
check_project_root() {
    if [[ ! -f "docker-compose.yml" ]]; then
        log_error "Script must be run from project root directory"
        exit 1
    fi
    log_success "Project root directory found"
}

# Check and install Docker
install_docker() {
    log_info "Checking Docker..."
    
    if command -v docker &> /dev/null; then
        log_success "Docker is already installed"
        docker --version
    else
        log_warning "Docker not found. Installing..."
        sudo apt update
        sudo apt install -y docker.io
        sudo systemctl start docker
        sudo systemctl enable docker
        sudo usermod -aG docker $USER
        log_success "Docker installed. Restart system or log out/in"
    fi
}

# Check and install Docker Compose
install_docker_compose() {
    log_info "Checking Docker Compose..."
    
    if command -v docker-compose &> /dev/null; then
        log_success "Docker Compose is already installed"
        docker-compose --version
    else
        log_warning "Docker Compose not found. Installing..."
        sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
        sudo chmod +x /usr/local/bin/docker-compose
        sudo ln -sf /usr/local/bin/docker-compose /usr/bin/docker-compose
        log_success "Docker Compose installed"
        docker-compose --version
    fi
}

# Check and install Go
install_go() {
    log_info "Checking Go..."
    
    if command -v go &> /dev/null; then
        log_success "Go is already installed"
        go version
    else
        log_warning "Go not found. Installing..."
        sudo snap install go --classic
        export PATH=/snap/bin:$PATH
        log_success "Go installed"
        go version
    fi
}

# Check and install Node.js
install_nodejs() {
    log_info "Checking Node.js..."
    
    if command -v node &> /dev/null; then
        log_success "Node.js is already installed"
        node --version
        npm --version
    else
        log_warning "Node.js not found. Installing..."
        curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
        sudo apt-get install -y nodejs
        log_success "Node.js installed"
        node --version
        npm --version
    fi
}

# Check and install Python dependencies
install_python_deps() {
    log_info "Checking Python dependencies..."
    
    if [[ -f "scripts/requirements.txt" ]]; then
        log_info "Installing Python dependencies..."
        pip install -r scripts/requirements.txt
        log_success "Python dependencies installed"
    else
        log_warning "requirements.txt file not found"
    fi
}

# Setup environment variables
setup_environment() {
    log_info "Setting up environment variables..."
    
    if [[ ! -f ".env" ]]; then
        if [[ -f "env.example" ]]; then
            cp env.example .env
            log_success ".env file created from env.example"
        else
            log_warning "env.example file not found, creating basic .env"
            cat > .env << EOF
# MongoDB Configuration
MONGO_ADMIN_USER=admin
MONGO_ADMIN_PASSWORD=mongo-1
APP_DB_USER=appuser
APP_DB_PASSWORD=appuserpassword
MONGO_DB=appdb
MONGO_REPLICA_SET=rs0

# Application Configuration
NODE_ENV=development
GO_ENV=development
LOG_LEVEL=info
EOF
            log_success "Basic .env file created"
        fi
    else
        log_success ".env file already exists"
    fi
}

# Install application dependencies
install_app_dependencies() {
    log_info "Installing application dependencies..."
    
    # Node.js dependencies
    if [[ -d "app-node" ]]; then
        log_info "Installing Node.js dependencies..."
        cd app-node
        npm install
        cd ..
        log_success "Node.js dependencies installed"
    fi
    
    # Go dependencies
    if [[ -d "app-go" ]]; then
        log_info "Installing Go dependencies..."
        cd app-go
        export PATH=/snap/bin:$PATH
        go mod tidy
        cd ..
        log_success "Go dependencies installed"
    fi
}

# Create MongoDB key
setup_mongodb_key() {
    log_info "Setting up MongoDB key..."
    
    if [[ ! -f "mongo/mongo-keyfile" ]]; then
        log_info "Creating MongoDB key..."
        openssl rand -base64 756 > mongo/mongo-keyfile
        sudo chmod 400 mongo/mongo-keyfile
        log_success "MongoDB key created"
    else
        log_success "MongoDB key already exists"
    fi
}

# Stop existing containers
stop_existing_containers() {
    log_info "Stopping existing containers..."
    
    if docker-compose ps | grep -q "Up"; then
        docker-compose down
        log_success "Existing containers stopped"
    else
        log_success "No running containers"
    fi
}

# Start project
start_project() {
    log_info "Starting project..."
    
    docker-compose up -d
    
    # Wait for containers to start
    log_info "Waiting for containers to start..."
    sleep 30
    
    # Check container status
    log_info "Checking container status..."
    docker ps
    
    log_success "Project started"
}

# Initialize MongoDB replica set
init_mongodb_replica_set() {
    log_info "Initializing MongoDB replica set..."
    
    # Wait for MongoDB containers to be ready
    log_info "Waiting for MongoDB containers to be ready..."
    sleep 60
    
    # Check that containers are running
    if ! docker ps | grep -q "mongo-0"; then
        log_error "MongoDB containers are not running"
        return 1
    fi
    
    # Initialize replica set
    log_info "Initializing replica set..."
    docker exec mongo-0 mongo --eval "rs.initiate({_id: 'rs0', members: [{_id: 0, host: 'mongo-0:27017'}, {_id: 1, host: 'mongo-1:27017'}, {_id: 2, host: 'mongo-2:27017'}]})"
    
    # Wait for replica set to stabilize
    log_info "Waiting for replica set to stabilize..."
    sleep 30
    
    # Check replica set status
    log_info "Checking replica set status..."
    docker exec mongo-0 mongo --eval "rs.status().members.forEach(function(m) { print(m.name + ': ' + m.stateStr); });"
    
    log_success "MongoDB replica set initialized"
}

# Test applications
test_applications() {
    log_info "Testing applications..."
    
    # Wait for applications to be ready
    sleep 30
    
    # Check health checks
    log_info "Checking health checks..."
    
    # Node.js application
    if docker exec app-node wget --quiet --tries=1 --spider http://localhost:3000/health; then
        log_success "Node.js application is working"
    else
        log_warning "Node.js application not responding to health check"
    fi
    
    # Go application
    if docker exec app-go wget --quiet --tries=1 --spider http://localhost:8080/health; then
        log_success "Go application is working"
    else
        log_warning "Go application not responding to health check"
    fi
    
    # Check logs
    log_info "Recent application logs:"
    echo "=== Node.js logs ==="
    docker logs app-node --tail 5
    echo "=== Go logs ==="
    docker logs app-go --tail 5
}

# Final check
final_check() {
    log_info "Final system check..."
    
    echo "=== Container status ==="
    docker ps
    
    echo "=== MongoDB replica set status ==="
    docker exec mongo-0 mongo --eval "rs.status().members.forEach(function(m) { print(m.name + ': ' + m.stateStr); });"
    
    echo "=== Health checks ==="
    docker ps --format "table {{.Names}}\t{{.Status}}"
    
    log_success "Installation completed!"
    log_info "Project available at:"
    log_info "- Node.js application: http://localhost:3000"
    log_info "- Go application: http://localhost:8080"
    log_info "- MongoDB: localhost:27030-27032"
    log_info "- HAProxy: localhost:27034"
}

# Main function
main() {
    log_info "Starting local environment setup..."
    
    check_project_root
    install_docker
    install_docker_compose
    install_go
    install_nodejs
    install_python_deps
    setup_environment
    install_app_dependencies
    setup_mongodb_key
    stop_existing_containers
    start_project
    init_mongodb_replica_set
    test_applications
    final_check
    
    log_success "Local environment successfully configured!"
}

# Command line argument handling
case "${1:-}" in
    --help|-h)
        echo "Usage: $0 [options]"
        echo "Options:"
        echo "  --help, -h     Show this help"
        echo "  --skip-deps    Skip dependency installation"
        echo "  --skip-mongo   Skip MongoDB initialization"
        echo "  --force        Force reinstallation"
        exit 0
        ;;
    --skip-deps)
        log_info "Skipping dependency installation"
        SKIP_DEPS=true
        ;;
    --skip-mongo)
        log_info "Skipping MongoDB initialization"
        SKIP_MONGO=true
        ;;
    --force)
        log_info "Force reinstallation"
        FORCE=true
        ;;
esac

# Run main function
main "$@" 