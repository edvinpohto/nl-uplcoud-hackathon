#!/usr/bin/env bash
# port-forward-grafana.sh — Forwards Grafana to localhost:3000
# Usage: ./scripts/port-forward-grafana.sh

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
KC_B="${KUBECONFIG:-$ROOT_DIR/kc-b.yaml}"

if [ ! -f "$KC_B" ]; then
  echo "ERROR: Kubeconfig not found at $KC_B"
  echo "Run: ./scripts/apply-kubeconfig.sh first"
  exit 1
fi

echo "==> Forwarding Grafana to http://localhost:3000"
echo "    Credentials: admin / chaos-monkey"
echo "    Press Ctrl+C to stop"
echo ""

kubectl --kubeconfig="$KC_B" port-forward \
  svc/grafana 3000:3000 \
  -n chaos
