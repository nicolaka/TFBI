#!/bin/bash
helm uninstall tfbi-grafana -n tfbi
kubectl delete -f k8s/tfbi_exporter/tfbi-exporter-deployment.yaml -n tfbi
kubectl delete -f k8s/tfbi_exporter/tfbi-exporter-service.yaml -n tfbi
helm uninstall tfbi-prometheus -n tfbi
kubectl delete namespace tfbi
