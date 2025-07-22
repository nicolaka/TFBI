#!/bin/bash

# Create namespace
kubectl create namespace tfbi

# Install Prometheus
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install tfbi-prometheus prometheus-community/prometheus -f k8s/prometheus/helm/gke_values.yaml -n tfbi

# Create config map for TFBI exporter
kubectl create configmap tfbi-config \
  --from-literal=TF_ORGANIZATIONS=$TFE_ORGANIZATIONS \
  --from-literal=TFE_ADDRESS=https://app.terraform.io \
  -n tfbi

# Create secret for TFBI exporter
kubectl create secret generic tfbi-token --from-literal=TF_API_TOKEN=$TFE_API_TOKEN -n tfbi

# Install TFBI exporter
kubectl apply -f k8s/tfbi_exporter/tfbi-exporter-deployment.yaml -n tfbi
kubectl apply -f k8s/tfbi_exporter/tfbi-exporter-service.yaml -n tfbi

# Create config map for Grafana dashboard
kubectl create configmap tfbi-grafana-dashboard-config --from-file=grafana/dashboards/general.json -n tfbi

# Install Grafana
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update
helm install tfbi-grafana grafana/grafana -f k8s/grafana/helm/gke_values.yaml -n tfbi

kubectl config set-context --current --namespace=tfbi

kubectl get pods
kubectl get services
