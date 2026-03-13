#!/usr/bin/env bash
# build-push.sh — Builds and pushes all Docker images tagged with the current git SHA.
# Usage: REGISTRY=registry.example.com/chaos-monkey ./scripts/build-push.sh

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

: "${REGISTRY:?REGISTRY must be set (e.g. docker.io/myuser)}"

TAG="${IMAGE_TAG:-$(git -C "$ROOT_DIR" rev-parse --short HEAD)}"

echo "==> Building and pushing images"
echo "    Registry : $REGISTRY"
echo "    Tag      : $TAG"

build_push() {
  local name="$1"
  local dir="$ROOT_DIR/$name"
  local image="$REGISTRY/$name:$TAG"

  echo ""
  echo "--- Building $image (linux/amd64) ---"
  docker buildx build --platform linux/amd64 -t "$image" --push "$dir"
  echo "--- Pushed $image ---"
}

build_push "victim-app"
build_push "health-checker"
build_push "chaos-monkey"

echo ""
echo "==> All images pushed with tag: $TAG"
echo ""
echo "Update your Kubernetes manifests with:"
echo "    REGISTRY=$REGISTRY TAG=$TAG"
echo ""
echo "Or use scripts/apply-kubeconfig.sh to update the manifests automatically."
