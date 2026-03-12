#!/bin/sh
# ZeroWall - Driver Verification Script
# Checks for compatible network interfaces and system hardware.

echo "--- ZeroWall Hardware Verification ---"

# Detect OS
OS_TYPE=$(uname)
if [ "$OS_TYPE" != "FreeBSD" ]; then
    echo "[WARNING] This script is designed for FreeBSD. Results may be inaccurate on $OS_TYPE."
fi

# Detect CPUs
CPU_COUNT=$(sysctl -n hw.ncpu)
echo "[CPU] Detected $CPU_COUNT logical processing cores."

# Detect Memory
MEM_TOTAL=$(sysctl -n hw.physmem)
MEM_GB=$(expr $MEM_TOTAL / 1024 / 1024 / 1024)
echo "[MEM] Detected ${MEM_GB}GB of physical memory."

# Scan Interface Drivers
echo "[NIC] Scanning network interfaces..."
INTERFACES=$(ifconfig -l)

for iface in $INTERFACES; do
    # Skip loopback
    if [ "$iface" = "lo0" ]; then continue; fi
    
    DRIVER=$(echo $iface | sed 's/[0-9]*//g')
    DESCRIPTION="Unknown Driver"
    
    case $DRIVER in
        em|igb|ix|ixl)
            DESCRIPTION="Intel PRO/1000 or 10GbE (Excellent Compatibility)"
            ;;
        re)
            DESCRIPTION="Realtek PCIe (Standard Compatibility)"
            ;;
        vtnet)
            DESCRIPTION="VirtIO Paravirtualized (Recommended for VMs)"
            ;;
        hn)
            DESCRIPTION="Microsoft Hyper-V Network (Recommended for Azure/Hyper-V)"
            ;;
        vmx)
            DESCRIPTION="VMWare VMXNET3 (Recommended for ESXi)"
            ;;
        *)
            DESCRIPTION="General Driver: $DRIVER"
            ;;
    esac
    
    echo "  - $iface: $DESCRIPTION"
done

echo "--- Verification Complete ---"
