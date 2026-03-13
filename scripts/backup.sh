#!/bin/bash

set -e

VERSION="1.0.0"
BACKUP_DIR="/var/backups/zerowall"
CONFIG_DIR="/etc/zerowall"
LOG_FILE="/var/log/zerowall/backup.log"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
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

show_usage() {
    cat << EOF
ZeroWall Backup Script v${VERSION}

Usage: $0 [OPTIONS]

OPTIONS:
    -d, --dir          Backup directory (default: ${BACKUP_DIR})
    -c, --compress     Compress backup with gzip
    -r, --retention   Days to keep backups (default: 30)
    -f, --filename    Custom backup filename
    -v, --verbose     Verbose output
    -h, --help        Show this help message

EXAMPLES:
    $0                      # Run backup with defaults
    $0 -c -r 7             # Compress and keep 7 days
    $0 -f my-backup        # Custom filename
    $0 -d /backup/custom   # Custom backup directory

EOF
}

create_backup_dir() {
    if [ ! -d "$BACKUP_DIR" ]; then
        mkdir -p "$BACKUP_DIR"
        log "Created backup directory: $BACKUP_DIR"
    fi
    
    if [ ! -d "$(dirname $LOG_FILE)" ]; then
        mkdir -p "$(dirname $LOG_FILE)"
    fi
}

backup_config() {
    local config_files=(
        "${CONFIG_DIR}/config.xml"
        "${CONFIG_DIR}/config.xml.bak"
        "${CONFIG_DIR}/certs/"
        "${CONFIG_DIR}/aliases.conf"
    )
    
    log "Backing up configuration files..."
    
    for file in "${config_files[@]}"; do
        if [ -e "$file" ]; then
            cp -r "$file" "$BACKUP_DIR/temp/" 2>/dev/null || true
            log "  Backed up: $file"
        else
            warn "  Not found: $file"
        fi
    done
}

backup_firewall_rules() {
    log "Backing up firewall rules..."
    
    if [ -f "/etc/pf.conf" ]; then
        cp /etc/pf.conf "$BACKUP_DIR/temp/"
        log "  Backed up: /etc/pf.conf"
    else
        warn "  /etc/pf.conf not found"
    fi
    
    if [ -f "/etc/pf.conf.old" ]; then
        cp /etc/pf.conf.old "$BACKUP_DIR/temp/" 2>/dev/null || true
    fi
}

backup_network_config() {
    log "Backing up network configuration..."
    
    local network_files=(
        "/etc/rc.conf"
        "/etc/hosts"
        "/etc/resolv.conf"
        "/etc/dhcpd.conf"
    )
    
    for file in "${network_files[@]}"; do
        if [ -f "$file" ]; then
            cp "$file" "$BACKUP_DIR/temp/"
            log "  Backed up: $file"
        fi
    done
}

backup_vpn_config() {
    log "Backing up VPN configuration..."
    
    local vpn_dirs=(
        "/etc/openvpn"
        "/etc/wireguard"
        "/etc/ipsec.d"
    )
    
    for dir in "${vpn_dirs[@]}"; do
        if [ -d "$dir" ]; then
            cp -r "$dir" "$BACKUP_DIR/temp/"
            log "  Backed up: $dir"
        fi
    done
}

backup_dns_dhcp() {
    log "Backing up DNS/DHCP configuration..."
    
    if [ -d "/etc/unbound" ]; then
        cp -r /etc/unbound "$BACKUP_DIR/temp/"
        log "  Backed up: /etc/unbound"
    fi
    
    if [ -f "/etc/dhcpd.conf" ]; then
        cp /etc/dhcpd.conf "$BACKUP_DIR/temp/"
    fi
    
    if [ -d "/var/db/dhcpd" ]; then
        cp -r /var/db/dhcpd "$BACKUP_DIR/temp/" 2>/dev/null || true
    fi
}

backup_logs() {
    log "Backing up recent logs..."
    
    if [ -d "/var/log/zerowall" ]; then
        find /var/log/zerowall -type f -mtime -7 -exec cp {} "$BACKUP_DIR/temp/logs/" \; 2>/dev/null || true
        log "  Backed up recent logs"
    fi
}

create_archive() {
    local timestamp=$(date +%Y%m%d_%H%M%S)
    local archive_name="${BACKUP_DIR}/zerowall_backup_${timestamp}"
    
    if [ -n "$CUSTOM_FILENAME" ]; then
        archive_name="${BACKUP_DIR}/${CUSTOM_FILENAME}"
    fi
    
    log "Creating archive: ${archive_name}.tar.gz"
    
    cd "$BACKUP_DIR/temp"
    tar -czf "${archive_name}.tar.gz" . 2>/dev/null
    
    if [ $? -eq 0 ]; then
        log "Backup created successfully: ${archive_name}.tar.gz"
        local size=$(du -h "${archive_name}.tar.gz" | cut -f1)
        log "Backup size: $size"
        
        echo "${archive_name}.tar.gz"
    else
        error "Failed to create backup archive"
        return 1
    fi
}

cleanup_old_backups() {
    if [ -n "$RETENTION" ] && [ "$RETENTION" -gt 0 ]; then
        log "Cleaning up backups older than $RETENTION days..."
        find "$BACKUP_DIR" -name "zerowall_backup_*.tar.gz" -mtime +$RETENTION -delete 2>/dev/null || true
        log "Cleanup complete"
    fi
}

verify_backup() {
    local archive=$1
    
    log "Verifying backup..."
    
    if tar -tzf "$archive" > /dev/null 2>&1; then
        log "Backup verification: OK"
        return 0
    else
        error "Backup verification: FAILED"
        return 1
    fi
}

cleanup_temp() {
    if [ -d "$BACKUP_DIR/temp" ]; then
        rm -rf "$BACKUP_DIR/temp"
    fi
}

main() {
    COMPRESS=false
    RETENTION=30
    CUSTOM_FILENAME=""
    VERBOSE=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            -d|--dir)
                BACKUP_DIR="$2"
                shift 2
                ;;
            -c|--compress)
                COMPRESS=true
                shift
                ;;
            -r|--retention)
                RETENTION="$2"
                shift 2
                ;;
            -f|--filename)
                CUSTOM_FILENAME="$2"
                shift 2
                ;;
            -v|--verbose)
                VERBOSE=true
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
    
    log "========================================="
    log "ZeroWall Backup Script v${VERSION}"
    log "========================================="
    
    create_backup_dir
    
    mkdir -p "$BACKUP_DIR/temp"
    mkdir -p "$BACKUP_DIR/temp/logs"
    
    backup_config
    backup_firewall_rules
    backup_network_config
    backup_vpn_config
    backup_dns_dhcp
    backup_logs
    
    local archive=$(create_archive)
    
    if [ -n "$archive" ]; then
        verify_backup "$archive"
        cleanup_old_backups
    fi
    
    cleanup_temp
    
    log "========================================="
    log "Backup completed successfully!"
    log "========================================="
}

main "$@"
