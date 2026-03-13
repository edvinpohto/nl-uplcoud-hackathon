#!/usr/bin/env bash
# chaos-status.sh — Shows chaos monkey logs and uptime metrics
# Usage: ./scripts/chaos-status.sh

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
KC_A="${KUBECONFIG_A:-$ROOT_DIR/kc-a.yaml}"
KC_B="${KUBECONFIG_B:-$ROOT_DIR/kc-b.yaml}"

echo "=== Cluster A (Victim) ==="
if [ -f "$KC_A" ]; then
  echo "--- Nodes ---"
  kubectl --kubeconfig="$KC_A" get nodes -o wide 2>/dev/null || echo "(unavailable)"
  echo ""
  echo "--- Victim App Pods ---"
  kubectl --kubeconfig="$KC_A" get pods -n victim -o wide 2>/dev/null || echo "(unavailable)"
  echo ""
  echo "--- Victim App Service ---"
  kubectl --kubeconfig="$KC_A" get svc -n victim 2>/dev/null || echo "(unavailable)"
else
  echo "kc-a.yaml not found, skipping Cluster A status"
fi

echo ""
echo "=== Cluster B (Chaos + Monitor) ==="
if [ -f "$KC_B" ]; then
  echo "--- Nodes ---"
  kubectl --kubeconfig="$KC_B" get nodes -o wide 2>/dev/null || echo "(unavailable)"
  echo ""
  echo "--- Chaos Namespace Pods ---"
  kubectl --kubeconfig="$KC_B" get pods -n chaos -o wide 2>/dev/null || echo "(unavailable)"
  echo ""
  echo "--- Chaos Monkey Logs (last 50 lines) ---"
  kubectl --kubeconfig="$KC_B" logs -n chaos deployment/chaos-monkey --tail=50 2>/dev/null || echo "(no logs yet)"
  echo ""
  echo "--- Health Checker Logs (last 20 lines) ---"
  kubectl --kubeconfig="$KC_B" logs -n chaos deployment/health-checker --tail=20 2>/dev/null || echo "(no logs yet)"
else
  echo "kc-b.yaml not found, skipping Cluster B status"
fi

echo ""
echo "==> To watch chaos logs live:"
echo "    kubectl --kubeconfig=kc-b.yaml logs -f deployment/chaos-monkey -n chaos"
echo ""
echo "==> To open Grafana dashboard:"
echo "    ./scripts/port-forward-grafana.sh"
