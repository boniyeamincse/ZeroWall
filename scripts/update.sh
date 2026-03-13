#!/bin/bash

set -e

VERSION="1.0.0"
INSTALL_DIR="/usr/local"
CONFIG_DIR="/etc/zerowall"
LOG_FILE="/var/log/zerowall/update.log"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1" | tee -a "$LOG_FILE"
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ERROR:${NC} $1" | tee -a "$LOG_FILE" >&2
}

warn() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] WARNING:${NC} $1" | tee -a "$LOG_FILE"
}

info() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')] INFO:${NC} $1" | tee -a "$LOG_FILE"
}

show_usage() {
    cat << EOF
ZeroWall Update Script v${VERSION}

Usage: $0 [OPTIONS]

OPTIONS:
    -b, --backend         Update only backend API
    -f, --frontend        Update only frontend web UI
    -c, --config          Update configuration only
    -a, --all             Update all components (default)
    -d, --dry-run         Show what would be updated
    -s, --skip-build      Skip build step
    -h, --help            Show this help message

EXAMPLES:
    $0                      # Update all components
    $0 -b                   # Update only backend
    $0 -f                  # Update only frontend
    $0 -d                  # Dry run to see what would update

EOF
}

check_dependencies() {
    log "Checking dependencies..."
    
    local deps=("git" "make" "node" "npm")
    local missing=()
    
    for dep in "${deps[@]}"; do
        if ! command -v "$dep" &> /dev/null; then
            missing+=("$dep")
        fi
    done
    
    if [ ${#missing[@]} -gt 0 ]; then
        warn "Missing dependencies: ${missing[*]}"
        info "Some features may not work without these"
    fi
    
    log "Dependency check complete"
}

update_backend() {
    log "========================================="
    log "Updating ZeroWall Backend API"
    log "========================================="
    
    cd "$INSTALL_DIR/zerowall/api"
    
    if [ -f "go.mod" ]; then
        log "Updating Go dependencies..."
        go mod download
        go mod tidy
    fi
    
    log "Building backend..."
    go build -o zwapi .
    
    if [ $? -eq 0 ]; then
        log "Backend build successful"
        
        if [ -f "/usr/local/sbin/zwapi" ]; then
            log "Restarting zwapi service..."
            service zwapi restart 2>/dev/null || info "Could not restart service"
        fi
    else
        error "Backend build failed"
        return 1
    fi
}

update_frontend() {
    log "========================================="
    log "Updating ZeroWall Frontend Web UI"
    log "========================================="
    
    cd "$INSTALL_DIR/zerowall/web-ui"
    
    log "Installing frontend dependencies..."
    npm install
    
    log "Building frontend..."
    npm run build
    
    if [ $? -eq 0 ]; then
        log "Frontend build successful"
        
        log "Deploying frontend..."
        rm -rf /usr/local/www/zerowall
        cp -r dist /usr/local/www/zerowall
    else
        error "Frontend build failed"
        return 1
    fi
}

update_config() {
    log "========================================="
    log "Updating ZeroWall Configuration"
    log "========================================="
    
    if [ -d "$CONFIG_DIR" ]; then
        log "Backing up current configuration..."
        cp -r "$CONFIG_DIR" "${CONFIG_DIR}.bak.$(date +%Y%m%d)"
    fi
    
    log "Configuration update complete"
}

dry_run() {
    log "========================================="
    log "ZeroWall Update - Dry Run"
    log "========================================="
    
    info "Would update backend API from: $INSTALL_DIR/zerowall/api"
    info "Would update frontend from: $INSTALL_DIR/zerowall/web-ui"
    info "Would update configuration in: $CONFIG_DIR"
    
    log "Dry run complete - no changes made"
}

main() {
    UPDATE_BACKEND=true
    UPDATE_FRONTEND=true
    UPDATE_CONFIG=false
    DRY_RUN=false
    SKIP_BUILD=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            -b|--backend)
                UPDATE_FRONTEND=false
                UPDATE_CONFIG=false
                shift
                ;;
            -f|--frontend)
                UPDATE_BACKEND=false
                UPDATE_CONFIG=false
                shift
                ;;
            -c|--config)
                UPDATE_BACKEND=false
                UPDATE_FRONTEND=false
                shift
                ;;
            -a|--all)
                UPDATE_BACKEND=true
                UPDATE_FRONTEND=true
                UPDATE_CONFIG=true
                shift
                ;;
            -d|--dry-run)
                DRY_RUN=true
                shift
                ;;
            -s|--skip-build)
                SKIP_BUILD=true
                shift
                ;;
            -h|--help)
                show_usage
                exit 0
                ;;
            *)
                error "Unknown option: $1"
                show_usage
                exit 1
                ;;
        esac
    done
    
    mkdir -p "$(dirname $LOG_FILE)"
    
    log "========================================="
    log "ZeroWall Update Script v${VERSION}"
    log "========================================="
    
    check_dependencies
    
    if [ "$DRY_RUN" = true ]; then
        dry_run
        exit 0
    fi
    
    if [ "$UPDATE_CONFIG" = true ]; then
        update_config
    fi
    
    if [ "$UPDATE_BACKEND" = true ]; then
        update_backend
    fi
    
    if [ "$UPDATE_FRONTEND" = true ]; then
        update_frontend
    fi
    
    log "========================================="
    log "Update completed successfully!"
    log "========================================="
}

main "$@"
