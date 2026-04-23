# 🔍 Port Scanner

A concurrent TCP port scanner written in Go. Supports single port connectivity checks and full port range sweeps across one or multiple hosts.

---

## Features

- **Single port check** — verify if a specific port is open on a host
- **Sweep scan** — scan all 65535 ports on a host concurrently
- **Multi-host support** — scan multiple hosts in one command
- **Configurable via CLI flags** — clean interface using Go's `flag` package

---

## Usage

### Build

```bash
go build -o port-scanner .
```

### Check if a specific port is open

```bash
./port-scanner -hosts=example.com -port=80
```

### Sweep scan all ports on a single host

```bash
./port-scanner -hosts=example.com
```

### Sweep scan across multiple hosts

```bash
./port-scanner -hosts=example.com,google.com,localhost
```

---

## Flags

| Flag      | Type   | Description                              |
|-----------|--------|------------------------------------------|
| `-hosts`  | string | Comma-separated list of hostnames or IPs |
| `-port`   | int    | Target port (1–65535)                    |

---

## Caveats

- **No rate limiting** — all 65535 goroutines are spawned simultaneously, which can be aggressive on some networks and may trigger firewalls or rate limits.
- **TCP only** — UDP ports are not scanned.
- **1-second timeout per connection** — full scans can take a while on slow or filtered networks.

---

## Example Output

```
port 22 is active
port 80 is active
port 443 is active
```

---

## Roadmap

- [ ] **Worker pool** — cap concurrent goroutines to avoid overwhelming the network or target host
- [ ] **SYN scan** — raw packet half-open scanning for faster, stealthier discovery (requires root/CAP_NET_RAW)