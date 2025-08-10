#!/bin/bash
set -e

GREEN="\033[0;32m"
YELLOW="\033[1;33m"
RED="\033[0;31m"
NC="\033[0m"

info()    { echo -e "${YELLOW}[INFO]${NC} $1"; }
success() { echo -e "${GREEN}[OK]${NC} $1"; }
error()   { echo -e "${RED}[ERROR]${NC} $1"; }

NAMESPACE="kube-system"

echo "=== Installing OpenYurt Components ==="

info "Adding OpenYurt Helm repo..."
helm repo add openyurt https://openyurtio.github.io/openyurt-helm
helm repo update

info "Installing yurt-manager..."
helm upgrade --install yurt-manager \
    openyurt/yurt-manager \
    -n "$NAMESPACE" \
    --wait --timeout=300s --debug

kubectl get pods -n "$NAMESPACE" -l app=yurt-manager
kubectl wait --for=condition=Ready pod -l app=yurt-manager -n "$NAMESPACE" --timeout=300s
kubectl get svc -n "$NAMESPACE" | grep yurt-manager

info "Installing yurthub artifacts..."
SERVER_ADDR=$(kubectl config view --minify -o jsonpath='{.clusters[0].cluster.server}')
helm upgrade --install yurt-hub \
    openyurt/yurthub \
    -n "$NAMESPACE" \
    --set kubernetesServerAddr="$SERVER_ADDR" \
    --wait --timeout=300s --debug

kubectl get yss -n "$NAMESPACE"

info "Installing raven agent..."
helm upgrade --install raven-agent \
    openyurt/raven-agent \
    -n "$NAMESPACE" \
    --wait --timeout=300s --debug

kubectl get pods -n "$NAMESPACE" | grep raven-agent

success "All OpenYurt components installed successfully!"
