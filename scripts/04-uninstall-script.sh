#!/bin/bash

set -e

echo "=== Uninstalling OpenYurt Components ==="

if helm list -n kube-system | grep -q yurt-manager; then
    echo "Uninstalling yurt-manager from kube-system..."
    helm uninstall yurt-manager -n kube-system
    echo "✓ yurt-manager uninstalled"
elif helm list -n openyurt-system 2>/dev/null | grep -q yurt-manager; then
    echo "Uninstalling yurt-manager from openyurt-system..."
    helm uninstall yurt-manager -n openyurt-system
    echo "✓ yurt-manager uninstalled"
else
    echo "yurt-manager not found in kube-system or openyurt-system"
fi

kubectl wait --for=delete pod -l app=yurt-manager -n kube-system --timeout=60s 2>/dev/null || true
kubectl wait --for=delete pod -l app=yurt-manager -n openyurt-system --timeout=60s 2>/dev/null || true

echo "Cleaning up leftover OpenYurt Deployments in kube-system..."
DEPLOYS=$(kubectl get deployment -n kube-system -o name | grep yurt || true)
if [ -n "$DEPLOYS" ]; then
    echo "$DEPLOYS" | xargs -r kubectl delete -n kube-system
    echo "✓ Leftover deployments removed"
else
    echo "No leftover deployments found"
fi

echo "Removing OpenYurt CRDs..."
CRDS=$(kubectl get crd -o name | grep -E "(openyurt\.io|raven\.openyurt\.io)" || true)
if [ -n "$CRDS" ]; then
    echo "$CRDS" | xargs -r kubectl delete --ignore-not-found=true
    echo "✓ All OpenYurt CRDs removed"
else
    echo "No OpenYurt CRDs found"
fi

for ns in openyurt-system yurt-dashboard; do
    if kubectl get ns "$ns" >/dev/null 2>&1; then
        echo "Deleting namespace: $ns"
        kubectl delete ns "$ns" --ignore-not-found
    fi
done

echo "Removing leftover RBAC & Webhook configurations..."
kubectl delete clusterrole,clusterrolebinding -l app.kubernetes.io/name=openyurt --ignore-not-found || true
kubectl delete mutatingwebhookconfiguration,validatingwebhookconfiguration \
    -l app.kubernetes.io/name=openyurt --ignore-not-found || true

REMAINING=$(kubectl get pods -A | grep yurt || true)
if [ -n "$REMAINING" ]; then
    echo "WARNING: Some OpenYurt pods still remain:"
    echo "$REMAINING"
else
    echo "✓ No remaining OpenYurt pods found"
fi

echo "=== OpenYurt Uninstallation Completed Successfully ==="
