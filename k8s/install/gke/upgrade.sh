helm upgrade tfbi-prometheus prometheus-community/prometheus --values k8s/prometheus/helm/gke_values.yaml
helm upgrade tfbi-grafana grafana/grafana --values k8s/grafana/helm/gke_values.yaml
kubectl get pods
kubectl get services