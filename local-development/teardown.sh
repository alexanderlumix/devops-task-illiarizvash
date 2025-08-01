#!/bin/bash

# Local Development Teardown Script
# Clean up local project environment

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

# Stop and remove containers
stop_and_remove_containers() {
    log_info "Stopping and removing containers..."
    
    if docker-compose ps | grep -q "Up"; then
        log_info "Stopping containers..."
        docker-compose down
        log_success "Containers stopped"
    else
        log_success "No running containers"
    fi
    
    # Remove containers if they remain
    if docker ps -a | grep -q "devops-task-illiarizvash"; then
        log_info "Removing remaining containers..."
        docker ps -a | grep "devops-task-illiarizvash" | awk '{print $1}' | xargs -r docker rm -f
        log_success "Remaining containers removed"
    fi
}

# Remove Docker images
remove_docker_images() {
    log_info "Removing Docker images..."
    
    # Remove project images
    if docker images | grep -q "devops-task-illiarizvash"; then
        log_info "Removing project images..."
        docker images | grep "devops-task-illiarizvash" | awk '{print $3}' | xargs -r docker rmi -f
        log_success "Project images removed"
    else
        log_success "Project images not found"
    fi
    
    # Remove unused images
    if docker images -q | wc -l | grep -q -v "0"; then
        log_info "Removing unused images..."
        docker image prune -f
        log_success "Unused images removed"
    fi
}

# Remove Docker volumes
remove_docker_volumes() {
    log_info "Removing Docker volumes..."
    
    # Remove project volumes
    if docker volume ls | grep -q "devops-task-illiarizvash"; then
        log_info "Removing project volumes..."
        docker volume ls | grep "devops-task-illiarizvash" | awk '{print $2}' | xargs -r docker volume rm -f
        log_success "Project volumes removed"
    else
        log_success "Project volumes not found"
    fi
    
    # Remove unused volumes
    if docker volume ls -q | wc -l | grep -q -v "0"; then
        log_info "Removing unused volumes..."
        docker volume prune -f
        log_success "Unused volumes removed"
    fi
}

# Remove Docker networks
remove_docker_networks() {
    log_info "Removing Docker networks..."
    
    # Remove project networks
    if docker network ls | grep -q "devops-task-illiarizvash"; then
        log_info "Removing project networks..."
        docker network ls | grep "devops-task-illiarizvash" | awk '{print $1}' | xargs -r docker network rm
        log_success "Project networks removed"
    else
        log_success "Project networks not found"
    fi
}

# Clean up Docker system
cleanup_docker_system() {
    log_info "Cleaning up Docker system..."
    
    # Remove unused resources
    docker system prune -f
    log_success "Docker system cleaned"
}

# Remove local files
remove_local_files() {
    log_info "Removing local files..."
    
    # Remove .env file
    if [[ -f ".env" ]]; then
        log_info "Removing .env file..."
        rm -f .env
        log_success ".env file removed"
    fi
    
    # Remove MongoDB key
    if [[ -f "mongo/mongo-keyfile" ]]; then
        log_info "Removing MongoDB key..."
        rm -f mongo/mongo-keyfile
        log_success "MongoDB key removed"
    fi
    
    # Remove node_modules (optional)
    if [[ -d "app-node/node_modules" ]]; then
        log_info "Removing node_modules..."
        rm -rf app-node/node_modules
        log_success "node_modules removed"
    fi
    
    # Remove Go cache (optional)
    if [[ -d "app-go/go.sum" ]]; then
        log_info "Removing Go cache..."
        rm -f app-go/go.sum
        log_success "Go cache removed"
    fi
}

# Clean up logs
cleanup_logs() {
    log_info "Cleaning up logs..."
    
    # Clean Docker logs
    if [[ -d "/var/lib/docker/containers" ]]; then
        log_info "Cleaning Docker logs..."
        sudo find /var/lib/docker/containers -name "*.log" -delete
        log_success "Docker logs cleaned"
    fi
}

# Full cleanup (optional)
full_cleanup() {
    if [[ "$1" == "--full" ]]; then
        log_warning "Performing full cleanup..."
        
        # Stop Docker daemon
        log_info "Stopping Docker daemon..."
        sudo systemctl stop docker
        
        # Clean Docker data
        log_info "Cleaning Docker data..."
        sudo rm -rf /var/lib/docker
        sudo mkdir -p /var/lib/docker
        
        # Restart Docker daemon
        log_info "Restarting Docker daemon..."
        sudo systemctl start docker
        
        log_success "Full cleanup completed"
    fi
}

# Verify cleanup
verify_cleanup() {
    log_info "Verifying cleanup..."
    
    echo "=== Container Check ==="
    docker ps -a | grep -E "(devops-task-illiarizvash|mongo|app-node|app-go)" || echo "Containers not found"
    
    echo "=== Image Check ==="
    docker images | grep -E "(devops-task-illiarizvash|mongo)" || echo "Images not found"
    
    echo "=== Volume Check ==="
    docker volume ls | grep "devops-task-illiarizvash" || echo "Volumes not found"
    
    echo "=== Network Check ==="
    docker network ls | grep "devops-task-illiarizvash" || echo "Networks not found"
    
    echo "=== File Check ==="
    ls -la .env mongo/mongo-keyfile 2>/dev/null || echo "Files not found"
    
    log_success "Cleanup verification completed"
}

# Main function
main() {
    log_info "Starting local environment cleanup..."
    
    check_project_root
    
    # Request confirmation
    if [[ "$1" != "--force" && "$1" != "--full" ]]; then
        echo -e "${YELLOW}WARNING: This action will delete all project data!${NC}"
        read -p "Continue? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            log_info "Cleanup cancelled"
            exit 0
        fi
    fi
    
    stop_and_remove_containers
    remove_docker_images
    remove_docker_volumes
    remove_docker_networks
    cleanup_docker_system
    remove_local_files
    cleanup_logs
    full_cleanup "$1"
    verify_cleanup
    
    log_success "Local environment successfully cleaned!"
    log_info "To reinstall, run: ./local-development/setup.sh"
}

# Command line argument handling
case "${1:-}" in
    --help|-h)
        echo "Usage: $0 [options]"
        echo "Options:"
        echo "  --help, -h     Show this help"
        echo "  --force        Force cleanup without confirmation"
        echo "  --full         Full cleanup (including Docker data)"
        echo ""
        echo "Examples:"
        echo "  $0              Normal cleanup with confirmation"
        echo "  $0 --force      Force cleanup"
        echo "  $0 --full       Full Docker data cleanup"
        exit 0
        ;;
    --force)
        log_info "Force cleanup"
        ;;
    --full)
        log_info "Full cleanup"
        ;;
    "")
        ;;
    *)
        log_error "Unknown option: $1"
        echo "Use --help for help"
        exit 1
        ;;
esac

# Run main function
main "$@" 