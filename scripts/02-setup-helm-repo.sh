#!/bin/bash
set -e

GREEN="\033[0;32m"
YELLOW="\033[1;33m"
RED="\033[0;31m"
NC="\033[0m"

info()    { echo -e "${YELLOW}[INFO]${NC} $1"; }
success() { echo -e "${GREEN}[OK]${NC} $1"; }
error()   { echo -e "${RED}[ERROR]${NC} $1"; }

echo "=== OpenYurt Helm Repository Setup ==="

if helm repo list | grep -q "openyurt"; then
    info "OpenYurt repo already exists. Updating it..."
    helm repo update openyurt || {
        error "Failed to update OpenYurt repo."
        exit 1
    }
else
    info "Adding OpenYurt repository..."
    helm repo add openyurt https://openyurtio.github.io/openyurt-helm || {
        error "Failed to add OpenYurt repo."
        exit 1
    }
fi

info "Updating all Helm repositories..."
helm repo update

info "Listing OpenYurt charts..."
helm search repo openyurt || {
    error "Failed to search OpenYurt repo."
    exit 1
}

success "Helm repository setup completed!"
