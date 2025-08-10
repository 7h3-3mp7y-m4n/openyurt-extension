#!/bin/bash
set -e

GREEN="\033[0;32m"
RED="\033[0;31m"
YELLOW="\033[1;33m"
NC="\033[0m"

info()    { echo -e "${YELLOW}[INFO]${NC} $1"; }
success() { echo -e "${GREEN}[OK]${NC} $1"; }
error()   { echo -e "${RED}[ERROR]${NC} $1"; }

echo "=== OpenYurt Prerequisites Check ==="

info "Checking for kubectl/k8s..."
if ! command -v kubectl &>/dev/null; then
    error "kubectl not found. Please install it before proceeding."
    exit 1
fi

KUBECTL_VERSION=$(kubectl version --client --short 2>/dev/null || kubectl version --client)
success "kubectl found: $KUBECTL_VERSION"

info "Checking Kubernetes cluster access..."
if ! kubectl cluster-info &>/dev/null; then
    error "Cannot connect to Kubernetes cluster."
    exit 1
fi

success "Kubernetes cluster is accessible."

info "Checking for Ready nodes..."
READY_NODES=$(kubectl get nodes --no-headers 2>/dev/null | grep -c ' Ready')
if [ "$READY_NODES" -eq 0 ]; then
    error "No nodes are in 'Ready' state."
    exit 1
fi

success "$READY_NODES nodes are Ready."
kubectl get nodes

if command -v helm &>/dev/null; then
    HELM_VERSION=$(helm version --short)
    success "Helm found: $HELM_VERSION"
else
    info "Helm not found Install Helm for required for OpenYurt setup."
fi

echo -e "\n${GREEN}✔ All prerequisite checks passed.${NC}"
