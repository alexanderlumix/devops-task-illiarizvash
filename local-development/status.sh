#!/bin/bash

# Local Development Status Script
# Quick status check of local project environment

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
}

# Check Docker
check_docker() {
    log_info "Checking Docker..."
    
    if command -v docker &> /dev/null; then
        log_success "Docker installed: $(docker --version)"
        
        if docker info &> /dev/null; then
            log_success "Docker daemon is running"
        else
            log_error "Docker daemon is not running"
            return 1
        fi
    else
        log_error "Docker not installed"
        return 1
    fi
}

# Check Docker Compose
check_docker_compose() {
    log_info "Checking Docker Compose..."
    
    if command -v docker-compose &> /dev/null; then
        log_success "Docker Compose installed: $(docker-compose --version)"
    else
        log_error "Docker Compose not installed"
        return 1
    fi
}

# Check Go
check_go() {
    log_info "Checking Go..."
    
    if command -v go &> /dev/null; then
        log_success "Go installed: $(go version)"
    else
        log_warning "Go not installed"
    fi
}

# Check Node.js
check_nodejs() {
    log_info "Checking Node.js..."
    
    if command -v node &> /dev/null; then
        log_success "Node.js installed: $(node --version)"
        log_success "npm installed: $(npm --version)"
    else
        log_warning "Node.js not installed"
    fi
}

# Check containers
check_containers() {
    log_info "Checking containers..."
    
    if docker ps &> /dev/null; then
        echo "=== Running containers ==="
        docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
        
        echo ""
        echo "=== All project containers ==="
        docker ps -a --filter "name=devops-task-illiarizvash" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
    else
        log_error "Failed to get container list"
        return 1
    fi
}

# Check MongoDB replica set
check_mongodb() {
    log_info "Checking MongoDB replica set..."
    
    if docker ps | grep -q "mongo-0"; then
        log_success "MongoDB containers are running"
        
        echo "=== Replica set status ==="
        docker exec mongo-0 mongo --eval "rs.status().members.forEach(function(m) { print(m.name + ': ' + m.stateStr); });" 2>/dev/null || log_warning "Failed to get replica set status"
    else
        log_warning "MongoDB containers are not running"
    fi
}

# Check applications
check_applications() {
    log_info "Checking applications..."
    
    # Node.js application
    if docker ps | grep -q "app-node"; then
        log_success "Node.js application is running"
        
        # Check health check
        if docker exec app-node wget --quiet --tries=1 --spider http://localhost:3000/health 2>/dev/null; then
            log_success "Node.js health check: OK"
        else
            log_warning "Node.js health check: FAILED"
        fi
    else
        log_warning "Node.js application is not running"
    fi
    
    # Go application
    if docker ps | grep -q "app-go"; then
        log_success "Go application is running"
        
        # Check health check
        if docker exec app-go wget --quiet --tries=1 --spider http://localhost:8080/health 2>/dev/null; then
            log_success "Go health check: OK"
        else
            log_warning "Go health check: FAILED"
        fi
    else
        log_warning "Go application is not running"
    fi
}

# Check HAProxy
check_haproxy() {
    log_info "Checking HAProxy..."
    
    if docker ps | grep -q "haproxy"; then
        log_success "HAProxy is running"
        
        # Check port availability
        if netstat -tlnp 2>/dev/null | grep -q ":27034"; then
            log_success "HAProxy port 27034 is available"
        else
            log_warning "HAProxy port 27034 is not available"
        fi
    else
        log_warning "HAProxy is not running"
    fi
}

# Check files
check_files() {
    log_info "Checking files..."
    
    # Check .env
    if [[ -f ".env" ]]; then
        log_success ".env file exists"
    else
        log_warning ".env file not found"
    fi
    
    # Check MongoDB key
    if [[ -f "mongo/mongo-keyfile" ]]; then
        log_success "MongoDB key exists"
        
        # Check access rights
        if [[ -r "mongo/mongo-keyfile" ]]; then
            log_success "MongoDB key is readable"
        else
            log_warning "MongoDB key is not readable"
        fi
    else
        log_warning "MongoDB key not found"
    fi
    
    # Check dependencies
    if [[ -d "app-node/node_modules" ]]; then
        log_success "Node.js dependencies installed"
    else
        log_warning "Node.js dependencies not installed"
    fi
}

# Check logs
check_logs() {
    log_info "Checking logs..."
    
    echo "=== Recent Node.js logs ==="
    docker logs app-node --tail 3 2>/dev/null || echo "Logs not available"
    
    echo ""
    echo "=== Recent Go logs ==="
    docker logs app-go --tail 3 2>/dev/null || echo "Logs not available"
    
    echo ""
    echo "=== Recent MongoDB logs ==="
    docker logs mongo-0 --tail 3 2>/dev/null || echo "Logs not available"
}

# Check resources
check_resources() {
    log_info "Checking resources..."
    
    echo "=== Resource usage ==="
    docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.MemPerc}}" 2>/dev/null || echo "Statistics not available"
    
    echo ""
    echo "=== Disk usage ==="
    docker system df 2>/dev/null || echo "Information not available"
}

# Check network
check_network() {
    log_info "Checking network..."
    
    echo "=== Docker networks ==="
    docker network ls --filter "name=devops-task-illiarizvash" --format "table {{.Name}}\t{{.Driver}}\t{{.Scope}}"
    
    echo ""
    echo "=== Volumes ==="
    docker volume ls --filter "name=devops-task-illiarizvash" --format "table {{.Name}}\t{{.Driver}}"
}

# Final summary
final_summary() {
    log_info "Final summary..."
    
    echo ""
    echo "=== STATUS SUMMARY ==="
    
    # Count components
    total_components=0
    healthy_components=0
    
    # Docker
    if command -v docker &> /dev/null && docker info &> /dev/null; then
        ((total_components++))
        ((healthy_components++))
        echo "✅ Docker: Working"
    else
        ((total_components++))
        echo "❌ Docker: Not working"
    fi
    
    # Docker Compose
    if command -v docker-compose &> /dev/null; then
        ((total_components++))
        ((healthy_components++))
        echo "✅ Docker Compose: Installed"
    else
        ((total_components++))
        echo "❌ Docker Compose: Not installed"
    fi
    
    # MongoDB
    if docker ps | grep -q "mongo-0"; then
        ((total_components++))
        ((healthy_components++))
        echo "✅ MongoDB: Running"
    else
        ((total_components++))
        echo "❌ MongoDB: Not running"
    fi
    
    # Node.js application
    if docker ps | grep -q "app-node"; then
        ((total_components++))
        ((healthy_components++))
        echo "✅ Node.js: Running"
    else
        ((total_components++))
        echo "❌ Node.js: Not running"
    fi
    
    # Go application
    if docker ps | grep -q "app-go"; then
        ((total_components++))
        ((healthy_components++))
        echo "✅ Go: Running"
    else
        ((total_components++))
        echo "❌ Go: Not running"
    fi
    
    # HAProxy
    if docker ps | grep -q "haproxy"; then
        ((total_components++))
        ((healthy_components++))
        echo "✅ HAProxy: Running"
    else
        ((total_components++))
        echo "❌ HAProxy: Not running"
    fi
    
    echo ""
    echo "=== OVERALL STATUS ==="
    echo "Components: $total_components"
    echo "Working: $healthy_components"
    
    if [[ $healthy_components -eq $total_components ]]; then
        log_success "All components are working correctly!"
    elif [[ $healthy_components -gt 0 ]]; then
        log_warning "Some components are working. Check logs for diagnostics."
    else
        log_error "System is not working. Run setup.sh for installation."
    fi
}

# Main function
main() {
    log_info "Checking local environment status..."
    
    check_project_root
    check_docker
    check_docker_compose
    check_go
    check_nodejs
    check_containers
    check_mongodb
    check_applications
    check_haproxy
    check_files
    check_logs
    check_resources
    check_network
    final_summary
}

# Command line argument handling
case "${1:-}" in
    --help|-h)
        echo "Usage: $0 [options]"
        echo "Options:"
        echo "  --help, -h     Show this help"
        echo "  --quick        Quick check (main components only)"
        echo "  --verbose      Detailed check (including logs and resources)"
        exit 0
        ;;
    --quick)
        log_info "Quick check..."
        check_project_root
        check_docker
        check_containers
        check_applications
        final_summary
        ;;
    --verbose)
        log_info "Detailed check..."
        main
        ;;
    "")
        main
        ;;
    *)
        log_error "Unknown option: $1"
        echo "Use --help for help"
        exit 1
        ;;
esac 