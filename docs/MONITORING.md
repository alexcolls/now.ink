# 📊 Monitoring & Analytics Setup

Complete guide for setting up production monitoring, analytics, error tracking, and alerting for now.ink.

---

## 📋 Table of Contents

1. [Overview](#overview)
2. [Error Tracking (Sentry)](#error-tracking-sentry)
3. [Metrics (Prometheus + Grafana)](#metrics-prometheus--grafana)
4. [User Analytics](#user-analytics)
5. [Logging](#logging)
6. [Alerting](#alerting)
7. [Dashboard Examples](#dashboard-examples)

---

## 🎯 Overview

### Monitoring Stack

```
┌─────────────────────────────────────────────────┐
│                 Applications                    │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐     │
│  │  Backend │  │  Mobile  │  │ Blockchain│     │
│  │   (Go)   │  │   (RN)   │  │  Scripts  │     │
│  └────┬─────┘  └────┬─────┘  └────┬──────┘     │
│       │             │             │             │
└───────┼─────────────┼─────────────┼─────────────┘
        │             │             │
        ▼             ▼             ▼
   ┌────────────────────────────────────┐
   │         Error Tracking             │
   │         (Sentry)                   │
   └────────────────────────────────────┘
        │
        ▼
   ┌────────────────────────────────────┐
   │         Metrics Collection         │
   │         (Prometheus)               │
   └────────┬───────────────────────────┘
            │
            ▼
   ┌────────────────────────────────────┐
   │         Visualization              │
   │         (Grafana)                  │
   └────────────────────────────────────┘
            │
            ▼
   ┌────────────────────────────────────┐
   │         Alerting                   │
   │   (Grafana + Sentry + Email)       │
   └────────────────────────────────────┘
```

### Tools Overview

| Tool | Purpose | Cost | Hosting |
|------|---------|------|---------|
| **Sentry** | Error tracking | Free (5K events/mo) | Cloud |
| **Prometheus** | Metrics collection | Free | Self-hosted |
| **Grafana** | Dashboards & viz | Free | Self-hosted |
| **Loki** | Log aggregation | Free | Self-hosted |
| **Plausible** | Privacy-friendly analytics | $9/mo | Cloud |

---

## 🔴 Error Tracking (Sentry)

### 1. Setup Sentry Account

1. Sign up at https://sentry.io
2. Create a new organization: `now-ink`
3. Create projects:
   - `backend-api` (Go)
   - `mobile-app` (React Native)
   - `blockchain-scripts` (Node.js)

### 2. Backend Integration (Go)

**Install Sentry SDK:**
```bash
cd backend
go get github.com/getsentry/sentry-go
go get github.com/getsentry/sentry-go/http
```

**Add to `backend/main.go`:**
```go
package main

import (
    "log"
    "time"
    
    "github.com/getsentry/sentry-go"
    sentryfiber "github.com/getsentry/sentry-go/fiber"
    "github.com/gofiber/fiber/v2"
)

func main() {
    // Initialize Sentry
    err := sentry.Init(sentry.ClientOptions{
        Dsn: os.Getenv("SENTRY_DSN"),
        Environment: os.Getenv("ENV"), // "production" or "development"
        Release: "now-ink-backend@0.3.2",
        TracesSampleRate: 1.0, // Capture 100% of transactions
        BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
            // Filter sensitive data
            if event.Request != nil {
                event.Request.Headers = nil // Remove headers
            }
            return event
        },
    })
    if err != nil {
        log.Fatalf("Sentry initialization failed: %v\n", err)
    }
    defer sentry.Flush(2 * time.Second)

    app := fiber.New()
    
    // Add Sentry middleware
    app.Use(sentryfiber.New(sentryfiber.Options{
        Repanic: true,
        WaitForDelivery: false,
    }))
    
    // Your routes...
    
    app.Listen(":8080")
}
```

**Capture errors manually:**
```go
func (s *StreamService) SaveStream(...) error {
    if err := s.db.Create(&stream).Error; err != nil {
        sentry.CaptureException(err)
        return err
    }
    return nil
}
```

**Add context:**
```go
sentry.ConfigureScope(func(scope *sentry.Scope) {
    scope.SetUser(sentry.User{
        ID: userID,
        Email: email,
        Username: username,
    })
    scope.SetTag("wallet", walletAddress)
    scope.SetContext("nft", map[string]interface{}{
        "mint_address": mintAddress,
        "arweave_tx": arweaveTxId,
    })
})
```

### 3. Mobile Integration (React Native)

**Install Sentry:**
```bash
cd mobile
npx expo install @sentry/react-native
```

**Configure `mobile/app.json`:**
```json
{
  "expo": {
    "hooks": {
      "postPublish": [
        {
          "file": "sentry-expo/upload-sourcemaps",
          "config": {
            "organization": "now-ink",
            "project": "mobile-app"
          }
        }
      ]
    }
  }
}
```

**Add to `mobile/App.tsx`:**
```typescript
import * as Sentry from '@sentry/react-native';

Sentry.init({
  dsn: process.env.EXPO_PUBLIC_SENTRY_DSN,
  environment: __DEV__ ? 'development' : 'production',
  release: 'now-ink-mobile@0.3.2',
  tracesSampleRate: 1.0,
  beforeSend(event, hint) {
    // Filter sensitive data
    if (event.user) {
      delete event.user.email;
    }
    return event;
  },
});

function App() {
  // Your app code...
}

export default Sentry.wrap(App);
```

**Capture errors:**
```typescript
try {
  await mintNFT(videoUri, metadata);
} catch (error) {
  Sentry.captureException(error, {
    tags: {
      action: 'nft_mint',
      wallet: walletAddress,
    },
    contexts: {
      video: {
        duration: videoDuration,
        size: videoSize,
      },
    },
  });
  throw error;
}
```

### 4. Environment Variables

Add to `.env`:
```env
# Sentry
SENTRY_DSN=https://YOUR_DSN@o123456.ingest.sentry.io/7891011
SENTRY_AUTH_TOKEN=YOUR_AUTH_TOKEN
SENTRY_ORG=now-ink
SENTRY_PROJECT_BACKEND=backend-api
SENTRY_PROJECT_MOBILE=mobile-app
```

Add to `mobile/.env`:
```env
EXPO_PUBLIC_SENTRY_DSN=https://YOUR_MOBILE_DSN@o123456.ingest.sentry.io/7891012
```

---

## 📈 Metrics (Prometheus + Grafana)

### 1. Prometheus Setup

**Create `monitoring/prometheus/prometheus.yml`:**
```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s
  external_labels:
    monitor: 'now-ink'
    environment: 'production'

scrape_configs:
  # Backend API metrics
  - job_name: 'backend-api'
    static_configs:
      - targets: ['backend:8080']
    metrics_path: '/metrics'
    
  # PostgreSQL metrics
  - job_name: 'postgres'
    static_configs:
      - targets: ['postgres-exporter:9187']
    
  # Node exporter (system metrics)
  - job_name: 'node'
    static_configs:
      - targets: ['node-exporter:9100']
    
  # Nginx metrics
  - job_name: 'nginx'
    static_configs:
      - targets: ['nginx-exporter:9113']

# Alerting rules
rule_files:
  - '/etc/prometheus/rules/*.yml'

alerting:
  alertmanagers:
    - static_configs:
        - targets: ['alertmanager:9093']
```

**Create `monitoring/prometheus/rules/alerts.yml`:**
```yaml
groups:
  - name: backend_alerts
    interval: 30s
    rules:
      # High error rate
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate detected"
          description: "Error rate is {{ $value }} (threshold: 0.05)"
      
      # API latency
      - alert: HighAPILatency
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High API latency"
          description: "95th percentile latency: {{ $value }}s"
      
      # Database connections
      - alert: HighDatabaseConnections
        expr: pg_stat_database_numbackends > 80
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High database connections"
          description: "{{ $value }} active connections"
      
      # Disk space
      - alert: LowDiskSpace
        expr: (node_filesystem_avail_bytes / node_filesystem_size_bytes) < 0.1
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "Low disk space"
          description: "Only {{ $value | humanizePercentage }} available"
      
      # NFT minting failures
      - alert: HighMintFailureRate
        expr: rate(nft_minting_failures_total[10m]) > 0.1
        for: 10m
        labels:
          severity: critical
        annotations:
          summary: "High NFT minting failure rate"
          description: "{{ $value }} failures per second"

  - name: system_alerts
    interval: 30s
    rules:
      # CPU usage
      - alert: HighCPUUsage
        expr: 100 - (avg by(instance) (rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100) > 80
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "High CPU usage"
          description: "CPU usage: {{ $value }}%"
      
      # Memory usage
      - alert: HighMemoryUsage
        expr: (1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100 > 85
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High memory usage"
          description: "Memory usage: {{ $value }}%"
```

### 2. Backend Prometheus Instrumentation

**Install Prometheus client:**
```bash
cd backend
go get github.com/prometheus/client_golang/prometheus
go get github.com/prometheus/client_golang/prometheus/promhttp
go get github.com/prometheus/client_golang/prometheus/promauto
```

**Create `backend/internal/metrics/metrics.go`:**
```go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    // HTTP metrics
    HTTPRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    HTTPRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint"},
    )
    
    // NFT minting metrics
    NFTMintingTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "nft_minting_total",
            Help: "Total NFT minting attempts",
        },
        []string{"status"}, // success, failure
    )
    
    NFTMintingDuration = promauto.NewHistogram(
        prometheus.HistogramOpts{
            Name: "nft_minting_duration_seconds",
            Help: "NFT minting duration",
            Buckets: []float64{1, 5, 10, 30, 60, 120},
        },
    )
    
    // Video upload metrics
    VideoUploadTotal = promauto.NewCounter(
        prometheus.CounterOpts{
            Name: "video_upload_total",
            Help: "Total video uploads",
        },
    )
    
    VideoUploadSize = promauto.NewHistogram(
        prometheus.HistogramOpts{
            Name: "video_upload_size_bytes",
            Help: "Video upload size in bytes",
            Buckets: prometheus.ExponentialBuckets(1024*1024, 2, 10), // 1MB to 1GB
        },
    )
    
    // Database metrics
    DatabaseQueriesTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "database_queries_total",
            Help: "Total database queries",
        },
        []string{"query_type", "status"},
    )
    
    DatabaseQueryDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "database_query_duration_seconds",
            Help: "Database query duration",
            Buckets: prometheus.DefBuckets,
        },
        []string{"query_type"},
    )
    
    // Active users
    ActiveUsers = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "active_users",
            Help: "Number of currently active users",
        },
    )
    
    // Arweave metrics
    ArweaveUploadsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "arweave_uploads_total",
            Help: "Total Arweave uploads",
        },
        []string{"status"},
    )
    
    ArweaveUploadDuration = promauto.NewHistogram(
        prometheus.HistogramOpts{
            Name: "arweave_upload_duration_seconds",
            Help: "Arweave upload duration",
            Buckets: []float64{1, 5, 10, 30, 60, 120, 300},
        },
    )
)
```

**Add metrics middleware in `backend/main.go`:**
```go
import (
    "time"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "your-project/internal/metrics"
)

func metricsMiddleware(c *fiber.Ctx) error {
    start := time.Now()
    
    // Process request
    err := c.Next()
    
    // Record metrics
    duration := time.Since(start).Seconds()
    status := c.Response().StatusCode()
    
    metrics.HTTPRequestsTotal.WithLabelValues(
        c.Method(),
        c.Path(),
        fmt.Sprintf("%d", status),
    ).Inc()
    
    metrics.HTTPRequestDuration.WithLabelValues(
        c.Method(),
        c.Path(),
    ).Observe(duration)
    
    return err
}

func main() {
    app := fiber.New()
    
    // Metrics middleware
    app.Use(metricsMiddleware)
    
    // Expose metrics endpoint
    app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
    
    // Your routes...
}
```

### 3. Grafana Setup

**Add to `docker-compose.monitoring.yml`:**
```yaml
version: '3.8'

services:
  prometheus:
    image: prom/prometheus:latest
    container_name: prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./monitoring/prometheus:/etc/prometheus
      - prometheus-data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--web.console.libraries=/etc/prometheus/console_libraries'
      - '--web.console.templates=/etc/prometheus/consoles'
      - '--web.enable-lifecycle'
    restart: unless-stopped
  
  grafana:
    image: grafana/grafana:latest
    container_name: grafana
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_USER=admin
      - GF_SECURITY_ADMIN_PASSWORD=${GRAFANA_PASSWORD}
      - GF_INSTALL_PLUGINS=grafana-piechart-panel
    volumes:
      - grafana-data:/var/lib/grafana
      - ./monitoring/grafana/provisioning:/etc/grafana/provisioning
      - ./monitoring/grafana/dashboards:/var/lib/grafana/dashboards
    depends_on:
      - prometheus
    restart: unless-stopped
  
  node-exporter:
    image: prom/node-exporter:latest
    container_name: node-exporter
    ports:
      - "9100:9100"
    command:
      - '--path.procfs=/host/proc'
      - '--path.sysfs=/host/sys'
      - '--collector.filesystem.mount-points-exclude=^/(sys|proc|dev|host|etc)($$|/)'
    volumes:
      - /proc:/host/proc:ro
      - /sys:/host/sys:ro
      - /:/rootfs:ro
    restart: unless-stopped
  
  postgres-exporter:
    image: prometheuscommunity/postgres-exporter:latest
    container_name: postgres-exporter
    ports:
      - "9187:9187"
    environment:
      - DATA_SOURCE_NAME=postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable
    restart: unless-stopped
  
  alertmanager:
    image: prom/alertmanager:latest
    container_name: alertmanager
    ports:
      - "9093:9093"
    volumes:
      - ./monitoring/alertmanager:/etc/alertmanager
      - alertmanager-data:/data
    command:
      - '--config.file=/etc/alertmanager/alertmanager.yml'
      - '--storage.path=/data'
    restart: unless-stopped

volumes:
  prometheus-data:
  grafana-data:
  alertmanager-data:
```

---

## 📱 User Analytics

### Option 1: Plausible (Recommended - Privacy-Friendly)

**Setup:**
1. Sign up at https://plausible.io
2. Add domain: `now.ink`
3. Add tracking script to landing page

**Add to `marketing/landing-page.html`:**
```html
<script defer data-domain="now.ink" src="https://plausible.io/js/script.js"></script>
```

**Track custom events:**
```html
<script>
function trackEvent(eventName, props) {
  if (window.plausible) {
    window.plausible(eventName, {props: props});
  }
}

// Track beta signup
document.getElementById('beta-form').addEventListener('submit', function() {
  trackEvent('Beta Signup');
});

// Track app download clicks
document.querySelectorAll('.download-btn').forEach(function(btn) {
  btn.addEventListener('click', function() {
    trackEvent('App Download Click', {
      platform: btn.dataset.platform
    });
  });
});
</script>
```

### Option 2: Self-Hosted (Free)

**Umami Analytics:**
```bash
# Clone Umami
git clone https://github.com/umami-software/umami.git
cd umami

# Add to docker-compose
docker-compose up -d
```

---

## 📝 Logging

### Loki + Promtail

**Add to `docker-compose.monitoring.yml`:**
```yaml
  loki:
    image: grafana/loki:latest
    container_name: loki
    ports:
      - "3100:3100"
    volumes:
      - ./monitoring/loki:/etc/loki
      - loki-data:/loki
    command: -config.file=/etc/loki/config.yml
    restart: unless-stopped
  
  promtail:
    image: grafana/promtail:latest
    container_name: promtail
    volumes:
      - ./monitoring/promtail:/etc/promtail
      - /var/log:/var/log:ro
      - /var/lib/docker/containers:/var/lib/docker/containers:ro
    command: -config.file=/etc/promtail/config.yml
    restart: unless-stopped

volumes:
  loki-data:
```

**Create `monitoring/loki/config.yml`:**
```yaml
auth_enabled: false

server:
  http_listen_port: 3100

ingester:
  lifecycler:
    ring:
      kvstore:
        store: inmemory
      replication_factor: 1
  chunk_idle_period: 15m
  chunk_retain_period: 30s

schema_config:
  configs:
    - from: 2024-01-01
      store: boltdb-shipper
      object_store: filesystem
      schema: v11
      index:
        prefix: index_
        period: 24h

storage_config:
  boltdb_shipper:
    active_index_directory: /loki/index
    cache_location: /loki/cache
    shared_store: filesystem
  filesystem:
    directory: /loki/chunks

limits_config:
  enforce_metric_name: false
  reject_old_samples: true
  reject_old_samples_max_age: 168h

chunk_store_config:
  max_look_back_period: 0s

table_manager:
  retention_deletes_enabled: true
  retention_period: 168h # 7 days
```

---

## 🚨 Alerting

### Configure Alertmanager

**Create `monitoring/alertmanager/alertmanager.yml`:**
```yaml
global:
  resolve_timeout: 5m

route:
  group_by: ['alertname', 'severity']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 12h
  receiver: 'email'
  
  routes:
    # Critical alerts go to PagerDuty (if configured)
    - match:
        severity: critical
      receiver: 'pagerduty'
      continue: true
    
    # All alerts go to email
    - match_re:
        severity: warning|critical
      receiver: 'email'

receivers:
  - name: 'email'
    email_configs:
      - to: 'alerts@now.ink'
        from: 'alertmanager@now.ink'
        smarthost: 'smtp.gmail.com:587'
        auth_username: 'alerts@now.ink'
        auth_password: '${SMTP_PASSWORD}'
        headers:
          Subject: '🚨 {{ .GroupLabels.alertname }} - now.ink'
  
  - name: 'pagerduty'
    pagerduty_configs:
      - service_key: '${PAGERDUTY_KEY}'
        description: '{{ .GroupLabels.alertname }}'

inhibit_rules:
  - source_match:
      severity: 'critical'
    target_match:
      severity: 'warning'
    equal: ['alertname']
```

### Slack Integration (Optional)

```yaml
  - name: 'slack'
    slack_configs:
      - api_url: '${SLACK_WEBHOOK_URL}'
        channel: '#alerts'
        title: '🚨 {{ .GroupLabels.alertname }}'
        text: |
          {{ range .Alerts }}
          *Alert:* {{ .Labels.alertname }}
          *Severity:* {{ .Labels.severity }}
          *Description:* {{ .Annotations.description }}
          {{ end }}
```

---

## 📊 Dashboard Examples

### Backend API Dashboard (Grafana JSON)

Save as `monitoring/grafana/dashboards/backend-api.json` (excerpt):

```json
{
  "dashboard": {
    "title": "now.ink Backend API",
    "panels": [
      {
        "title": "Request Rate",
        "targets": [{
          "expr": "rate(http_requests_total[5m])"
        }],
        "type": "graph"
      },
      {
        "title": "Error Rate",
        "targets": [{
          "expr": "rate(http_requests_total{status=~\"5..\"}[5m])"
        }],
        "type": "graph"
      },
      {
        "title": "95th Percentile Latency",
        "targets": [{
          "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))"
        }],
        "type": "graph"
      },
      {
        "title": "NFT Minting Rate",
        "targets": [{
          "expr": "rate(nft_minting_total{status=\"success\"}[5m])"
        }],
        "type": "stat"
      },
      {
        "title": "Active Users",
        "targets": [{
          "expr": "active_users"
        }],
        "type": "stat"
      }
    ]
  }
}
```

---

## ✅ Setup Checklist

### Initial Setup
- [ ] Create Sentry account and projects
- [ ] Install Sentry SDKs (backend + mobile)
- [ ] Deploy Prometheus + Grafana stack
- [ ] Configure Prometheus scrape targets
- [ ] Create Grafana dashboards
- [ ] Setup alerting rules
- [ ] Configure Alertmanager
- [ ] Test alert delivery

### Environment Variables
```env
# Sentry
SENTRY_DSN=https://YOUR_DSN@sentry.io/PROJECT_ID
SENTRY_AUTH_TOKEN=YOUR_TOKEN

# Grafana
GRAFANA_PASSWORD=YOUR_SECURE_PASSWORD

# Alerting
SMTP_PASSWORD=YOUR_SMTP_PASSWORD
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/YOUR/WEBHOOK/URL

# Database (for postgres-exporter)
POSTGRES_USER=nowink
POSTGRES_PASSWORD=YOUR_DB_PASSWORD
POSTGRES_DB=nowink
```

### Testing
```bash
# Test Prometheus targets
curl http://localhost:8080/metrics

# Test Grafana
open http://localhost:3000

# Test alert (simulate high error rate)
curl -X POST http://localhost:8080/test-error -H "Authorization: Bearer test"

# Check Alertmanager
open http://localhost:9093
```

---

## 📚 Resources

- **Sentry Docs**: https://docs.sentry.io
- **Prometheus Docs**: https://prometheus.io/docs
- **Grafana Docs**: https://grafana.com/docs
- **Alertmanager**: https://prometheus.io/docs/alerting/latest/alertmanager/
- **Plausible**: https://plausible.io/docs

---

**Created**: November 5, 2025  
**Version**: 1.0.0  
**Status**: Ready for Implementation
