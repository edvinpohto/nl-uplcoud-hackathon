#!/usr/bin/env bash
# apply-kubeconfig.sh — Extracts Cluster A kubeconfig + victim-app LB IP from Terraform output
# and injects them into Cluster B as Kubernetes Secrets/ConfigMaps.
#
# Run AFTER:
#   1. terraform apply (clusters are up)
#   2. kubectl apply on Cluster A (victim-app LoadBalancer created)
#
# Usage: REGISTRY=... TAG=... ./scripts/apply-kubeconfig.sh

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
TERRAFORM_DIR="$ROOT_DIR/terraform"
MANIFESTS_DIR="$ROOT_DIR/manifests"

: "${UPCLOUD_TOKEN:?UPCLOUD_TOKEN must be set}"
REGISTRY="${REGISTRY:-docker.io/chaosmonkey}"
TAG="${IMAGE_TAG:-$(git -C "$ROOT_DIR" rev-parse --short HEAD 2>/dev/null || echo latest)}"

echo "==> Extracting Terraform outputs"

cd "$TERRAFORM_DIR"

CLUSTER_A_UUID=$(terraform output -raw cluster_a_id)
CLUSTER_B_UUID=$(terraform output -raw cluster_b_id)

echo "    Cluster A UUID: $CLUSTER_A_UUID"
echo "    Cluster B UUID: $CLUSTER_B_UUID"

# Fetch kubeconfigs via UpCloud API (response is JSON: {"kubeconfig": "..."})
curl -sf -H "Authorization: Bearer $UPCLOUD_TOKEN" \
  "https://api.upcloud.com/1.3/kubernetes/$CLUSTER_A_UUID/kubeconfig" \
  | python3 -c "import json,sys; print(json.load(sys.stdin)['kubeconfig'])" > /tmp/kc-a.yaml
curl -sf -H "Authorization: Bearer $UPCLOUD_TOKEN" \
  "https://api.upcloud.com/1.3/kubernetes/$CLUSTER_B_UUID/kubeconfig" \
  | python3 -c "import json,sys; print(json.load(sys.stdin)['kubeconfig'])" > /tmp/kc-b.yaml

export KUBECONFIG=/tmp/kc-b.yaml

# Ensure chaos namespace exists
kubectl apply -f "$MANIFESTS_DIR/cluster-b/chaos-monkey/namespace.yaml"

# Encode Cluster A kubeconfig as base64
KCA_B64=$(base64 -i /tmp/kc-a.yaml | tr -d '\n')

echo "==> Creating/updating victim-cluster-config Secret in Cluster B"
kubectl create secret generic victim-cluster-config \
  --namespace=chaos \
  --from-literal=kubeconfig="$KCA_B64" \
  --from-literal=cluster_uuid="$CLUSTER_A_UUID" \
  --dry-run=client -o yaml | kubectl apply -f -

# Create UpCloud credentials secret
echo "==> Creating/updating upcloud-credentials Secret in Cluster B"
kubectl create secret generic upcloud-credentials \
  --namespace=chaos \
  --from-literal=token="$UPCLOUD_TOKEN" \
  --dry-run=client -o yaml | kubectl apply -f -

# Get victim-app LoadBalancer IP from Cluster A
echo "==> Waiting for victim-app LoadBalancer IP on Cluster A..."
export KUBECONFIG=/tmp/kc-a.yaml

for i in $(seq 1 30); do
  VICTIM_LB_IP=$(kubectl get svc victim-app -n victim \
    -o jsonpath='{.status.loadBalancer.ingress[0].ip}{.status.loadBalancer.ingress[0].hostname}' 2>/dev/null || true)
  if [ -n "$VICTIM_LB_IP" ]; then
    echo "    victim-app LB IP: $VICTIM_LB_IP"
    break
  fi
  echo "    Waiting for LB IP... ($i/30)"
  sleep 10
done

if [ -z "${VICTIM_LB_IP:-}" ]; then
  echo "ERROR: Could not get victim-app LoadBalancer IP after 5 minutes"
  exit 1
fi

# Update health-checker ConfigMap on Cluster B
export KUBECONFIG=/tmp/kc-b.yaml

echo "==> Updating health-checker-config ConfigMap on Cluster B"
kubectl create configmap health-checker-config \
  --namespace=chaos \
  --from-literal=target_url="http://$VICTIM_LB_IP:8080/health" \
  --dry-run=client -o yaml | kubectl apply -f -

# Apply manifests with correct image tags
echo "==> Applying Cluster B manifests with image tag: $TAG"

# Substitute REGISTRY and TAG in manifests and apply
for dir in chaos-monkey health-checker; do
  for f in "$MANIFESTS_DIR/cluster-b/$dir"/*.yaml; do
    sed "s|REGISTRY|$REGISTRY|g; s|:TAG|:$TAG|g" "$f" | kubectl apply -f -
  done
done

kubectl apply -f "$MANIFESTS_DIR/cluster-b/rbac.yaml"

echo ""
echo "==> Cluster A manifest deployment"
export KUBECONFIG=/tmp/kc-a.yaml

for f in "$MANIFESTS_DIR/cluster-a/victim-app"/*.yaml; do
  sed "s|REGISTRY|$REGISTRY|g; s|:TAG|:$TAG|g" "$f" | kubectl apply -f -
done

kubectl apply -f "$MANIFESTS_DIR/cluster-a/network-disrupt/daemonset.yaml"

echo ""
echo "==> Deploy monitoring stack on Cluster B"
export KUBECONFIG=/tmp/kc-b.yaml

kubectl apply -f "$ROOT_DIR/monitoring/prometheus/"
kubectl apply -f "$ROOT_DIR/monitoring/grafana/"

echo ""
echo "==> Deployment complete!"
echo "    Victim app LB: http://$VICTIM_LB_IP:8080/health"
echo "    Run: ./scripts/port-forward-grafana.sh"
echo "    Or: ./scripts/chaos-status.sh"

# Save kubeconfigs for convenience
cp /tmp/kc-a.yaml "$ROOT_DIR/kc-a.yaml"
cp /tmp/kc-b.yaml "$ROOT_DIR/kc-b.yaml"
echo "    Kubeconfigs saved: kc-a.yaml, kc-b.yaml"
