# Chaos Monkey — UpCloud UKS

A chaos engineering system that continuously disrupts a Kubernetes workload while measuring its uptime, built entirely on [UpCloud](https://upcloud.com).

## What it does

Two Kubernetes clusters run on a shared private SDN network:

- **Cluster A (victim)** — runs a simple HTTP app that gets attacked
- **Cluster B (chaos + monitoring)** — runs the chaos monkey, uptime prober, Prometheus, and Grafana

The chaos monkey fires a random disruptive action every 30–60 seconds. A health-checker probes the victim app every second and records whether it responded. Grafana visualises uptime, response latency, and chaos activity in real time.

```
UpCloud SDN (shared private network)
┌─────────────────────────────────────────────────┐
│  Cluster A (victim)      Cluster B (chaos)       │
│  ├─ victim-app (×3) ◄─── ├─ chaos-monkey        │
│  └─ alpine DaemonSet     ├─ health-checker       │
│     (tc netem)           ├─ prometheus           │
│                          └─ grafana              │
└─────────────────────────────────────────────────┘
```

Monitoring lives on Cluster B so that if chaos kills Cluster A entirely, Grafana keeps recording the outage.

## Chaos actions

| Action | Weight | Effect |
|---|---|---|
| `pod_kill` | 40 | Deletes a random running pod in the victim namespace |
| `network_disrupt` | 20 | Adds 200ms network delay to a node for 30s via `tc netem` |
| `scale_zero` | 20 | Scales victim-app to 0 replicas for 30s, then restores |

Weights are relative — pod_kill fires ~50% of the time. All weights and intervals are configurable via environment variables on the Deployment.

> Node-level actions (reboot/stop) are not implemented — the UpCloud Go SDK requires username/password auth, which is incompatible with API token login.

## Components

**`victim-app/`** — Simple Go HTTP server (`/health`, `/api/data`, `/api/slow`). Deployed with 3 replicas and an UpCloud LoadBalancer service on Cluster A.

**`health-checker/`** — Go service that probes `/health` every second and exposes Prometheus metrics: `uptime_probe_success` (0/1 gauge), `probe_response_seconds` (histogram), `probe_errors_total`, `probe_consecutive_failures`.

**`chaos-monkey/`** — Go service on Cluster B with a kubeconfig for Cluster A. Weighted random scheduler fires chaos actions and exposes `chaos_actions_total` / `chaos_action_errors_total` metrics. Network disruption execs `tc netem` into a privileged alpine DaemonSet on Cluster A.

**`monitoring/`** — Prometheus scrapes health-checker and chaos-monkey every 15s. Grafana has a pre-provisioned dashboard.

**`terraform/`** — UpCloud infrastructure: two UKS clusters, shared SDN network, router, and NAT gateway. Uses the [UpCloudLtd/upcloud](https://registry.terraform.io/providers/UpCloudLtd/upcloud) provider (~> 5.0).

## Grafana dashboard

| Panel | What it shows |
|---|---|
| Uptime % (5m) | Rolling 5-minute average of probe success |
| Response Time p99 | 99th percentile probe latency over 1m |
| Current Probe Status | Live UP / DOWN indicator |
| Uptime Over Time | Raw probe success as a time series — shows outage shape |
| Response Time Histogram | p50 / p95 / p99 latency — spikes during network_disrupt |
| Chaos Actions (5m) | Count of each action type in rolling 5m window |
| Consecutive Failures | Current streak of failed probes |

## Deployment

```bash
# 1. Provision infrastructure
cd terraform && terraform apply -var upcloud_token=$UPCLOUD_TOKEN

# 2. Build and push Docker images (requires Docker Hub account)
REGISTRY=docker.io/<user> ./scripts/build-push.sh

# 3. Wire clusters together (injects kubeconfig + LB IP as secrets, applies all manifests)
UPCLOUD_TOKEN=<token> REGISTRY=docker.io/<user> ./scripts/apply-kubeconfig.sh
```

## Configuration

| Variable | Default | Description |
|---|---|---|
| `MIN_INTERVAL` | `30s` | Minimum time between actions |
| `MAX_INTERVAL` | `60s` | Maximum time between actions |
| `WEIGHT_POD_KILL` | `40` | Relative weight for pod_kill |
| `WEIGHT_NETWORK_DISRUPT` | `20` | Relative weight for network_disrupt |
| `WEIGHT_SCALE_ZERO` | `20` | Relative weight for scale_zero |
| `NETWORK_DISRUPT_DELAY` | `200ms` | Delay injected by tc netem |
| `SCALE_ZERO_DURATION` | `30s` | How long victim-app stays at 0 replicas |
