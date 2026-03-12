# ZeroWall — Database Design

While ZeroWall primarily uses a single XML file (`config.xml`) for its core configuration to enable portability and atomic updates, it also utilizes lightweight, high-performance databases for logging, telemetry, and runtime data. This document describes the database architecture.

---

## 1. Primary Configuration Store (XML)
- **File Location:** `/etc/zerowall/config.xml`
- **Role:** Source of truth for all persistent settings.
- **Why XML?**
  - Human-readable and auditable.
  - Easy to backup and restore.
  - Atomic writes ensure data integrity during configuration changes.
  - Version control friendly (Git).

## 2. Runtime Databases (SQLite / Berkeley DB)
For transient data that requires high-frequency access or structured queries during uptime, ZeroWall uses embedded databases:

- **DHCP Leases:** Stored in `/var/db/dhcpd.leases` (structured text/DB format).
- **DNS Cache:** Unbound maintains its cache in memory, optionally synced to disk.
- **State Table:** Managed directly in the kernel's RAM (pf state table).
- **User Sessions:** Short-lived JWT identifiers are cached in an in-memory SQLite database for rapid validation.

## 3. Logs and Telemetry (Structured JSON / SQLite)
- **Local Logs:** Stored as structured JSON files in `/var/log/zerowall/`.
- **System Metrics:** Time-series data (CPU, RAM, Bandwidth) is collected by `zwlogd` and stored in a rolling SQLite database supporting up to 30 days of retention.
- **IDS/IPS Alerts:** Suricata generates EVE JSON logs, which can be indexed locally in a dedicated SQLite instance for rapid dashboard retrieval.

## 4. Optional External Database (PostgreSQL)
For enterprise deployments requiring multi-node log aggregation, long-term auditing, or integration with external SIEM/reporting tools, ZeroWall supports exporting data to an external PostgreSQL database.

**Schemas supported for export:**
- **Netflow/IPFIX:** Detailed traffic flow records.
- **Audit Logs:** Full history of administrative configuration changes.
- **Firewall Events:** Normalized pass/block events with GeoIP enrichment.

## 5. Security and Integrity
- **Restricted Access:** Database files are owned by the specific service user (e.g., `zwlog`, `unbound`) with `0600` permissions.
- **Atomic Updates:** All database transitions use transaction-based writes to prevent corruption during power loss.
- **Backup:** The primary `config.xml` is automatically backed up nightly and before any major system update.

---

*Previous: [API Design](10-api-design.md) | Next: [Installation Guide](12-installation-guide.md)*