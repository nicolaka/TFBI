#!/bin/bash
# IMPORTANT: Run k8s/local_testing/gke/config.sh first to configure kubectl for your GKE cluster.

# Create namespace
kubectl create namespace tfbi

# Install Prometheus
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts --force-update
echo "Installing Prometheus"
helm install tfbi-prometheus prometheus-community/prometheus -f k8s/prometheus/helm/gke_values.yaml -n tfbi

# Create config map for TFBI exporter
kubectl create configmap tfbi-config \
  --from-literal=TF_ORGANIZATIONS=$TF_ORGANIZATIONS \
  --from-literal=TFE_ADDRESS=$TFE_ADDRESS \
  -n tfbi

# Create secret for TFBI exporter
kubectl create secret generic tfbi-token --from-literal=TF_API_TOKEN=$TF_API_TOKEN -n tfbi


# Install TFBI exporter
echo "Installing TFBI exporter"
kubectl apply -f k8s/tfbi_exporter/tfbi-exporter-deployment.yaml -n tfbi
kubectl apply -f k8s/tfbi_exporter/tfbi-exporter-service.yaml -n tfbi

# Create config map for Grafana dashboard
kubectl create configmap tfbi-grafana-dashboard-config --from-file=grafana/dashboards/general.json -n tfbi

# Create or update TLS secret for Grafana
kubectl create secret tls grafana-tls \
  --cert="$GRAFANA_TLS_CERT" \
  --key="$GRAFANA_TLS_KEY" \
  -n tfbi --dry-run=client -o yaml | kubectl apply -f -

# Install Grafana
helm repo add grafana https://grafana.github.io/helm-charts --force-update
echo "Installing Grafana"
helm install tfbi-grafana grafana/grafana -f k8s/grafana/helm/gke_values.yaml -n tfbi

kubectl config set-context --current --namespace=tfbi

kubectl get pods
kubectl get services
