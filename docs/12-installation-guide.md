# ZeroWall — Installation Guide

This guide provides step-by-step instructions for installing ZeroWall on physical hardware or virtual environments.

---

## 1. System Requirements

### Minimum (Home/Small Lab)
- **CPU:** Dual-core 64-bit x86 (Intel/AMD)
- **RAM:** 2GB
- **Storage:** 8GB (SSD/HDD/NVMe)
- **Network:** 2x Gigabit Ethernet NICs (Intel recommended)

### Enterprise (High Throughput / IDS / IPS)
- **CPU:** Quad-core or higher (High clock speed for single-threaded packet processing)
- **RAM:** 8GB - 16GB
- **Storage:** 32GB+ SSD
- **Network:** 4x or more Gigabit/10G/25G NICs

## 2. Preparing the Installation Media
1. Download the latest ZeroWall `.iso` (for VMs) or `.img` (for USB) from the [Official Downloads](https://zerowall.io/download).
2. For USB installation, write the image using a tool like `dd`, `Rufus`, or `BalenaEtcher`.
   ```sh
   # Example for Linux/macOS
   sudo dd if=zerowall-24.11.img of=/dev/sdX bs=1M status=progress
   ```

## 3. Physical Hardware Installation
1. Insert the installation media into your appliance.
2. Boot from the media (adjust BIOS/UEFI settings if necessary).
3. The ZeroWall installer will load into memory and present the welcome screen.
4. Select **(I)nstall ZeroWall** from the menu.
5. Follow the prompts to:
   - Select the target disk.
   - Choose the filesystem type (ZFS recommended for data integrity and snapshots).
   - Set the root password.
6. Once complete, remove the installation media and **Reboot**.

## 4. Virtual Machine Installation
### Proxmox / KVM
- **OS Type:** Linux / Other (FreeBSD 14)
- **NIC:** VirtIO (Paravirtualized)
- **CPU:** Host type recommended.
- **Storage:** VirtIO SCSI with Discard enabled.

### VMware ESXi
- **Guest OS:** FreeBSD 13 or later (64-bit)
- **NIC:** VMXNET3
- **Disk Controller:** VMware Paravirtual

## 5. Initial Configuration (Console)
On the first boot, ZeroWall will prompt you to assign network interfaces:

1. **VLANs:** Opt-in to configure VLANs first if required.
2. **Assign WAN:** Select the NIC connected to your ISP (e.g., `em0`).
3. **Assign LAN:** Select the NIC connected to your internal switch (e.g., `em1`).
4. **Finalize:** Confirm the assignments to apply settings and start the web dashboard.

## 6. Accessing the Web Dashboard
1. Connect a computer to the LAN interface.
2. Ensure your computer is set to obtain an IP via DHCP (ZeroWall assigns `192.168.1.x` by default).
3. Open a browser and navigate to: `https://192.168.1.1`
4. Log in with the default credentials:
   - **Username:** `admin`
   - **Password:** `zerowall`
5. The **Setup Wizard** will automatically launch to guide you through initial host settings, NTP, and WAN configuration.

---

*Previous: [Database Design](11-database-design.md) | Next: [Admin Guide](13-configration-guide.md)*
