# 🔍 Terraform Business Insights (TFBI) 
> Business, Operational, and Adoption Insights for Terraform Cloud & Enterprise 

## Summary

Terraform Cloud Business Insights (TFBI) is a tool that provides business, operational, and adoption insights for Terraform Cloud operators. It implements both custom Prometheus collectors and metrics to query the Terraform Cloud API using [go-tfe](https://pkg.go.dev/github.com/hashicorp/go-tfe) Go libary and a Grafana dashboard to easily explore common business, operational, and adoption metrics. 

## Dashboard

![dashboard](img/dashboard_1.png)
![dashboard](img/dashboard_2.png)

## Metrics

| API/Category | Metric Name | Type | Description | Implementation Status
| - | - | - | - | -| 
| Organization | Organization Summary | `Table` | Organization Details  |  ✅  | 
| Teams | Total # of Teams | `Gauge` | Current number of active teams in the organization  |  ✅  | 
| Teams | Teams Summary | `Table` | Team Summary Table  |  ✅  | 
| Projects | Projects Count | `Gauge` | Current number of active projects in the organization  |  ✅  | 
| Projects | Projects Summary | `Table` | Projects Summary  |  ✅  | 
| Projects | Projects Count Over Time | `Time Series Graph` | Time series graph showing of # number of active projects over time |  ✅  | 
| Users | Total # of Users | `Gauge` | Current number of active users in the organization  |  ✅  | 
| Workspaces | Workspace Count | `Gauge` | Current number of active workspaces in the organization  |  ✅  | 
| Workspaces | Workspaces Summary | `Table` | Workspaces Summary  |  ✅  | 
| Workspaces | Workspaces Status Overview | `Chart` | Workspaces status distribution chart |  ✅  | 
| Workspaces | Terraform Version Distribution | `Chart` | Terraform version distribution chart |  ✅  |
| Workspaces | Drift Detection & Continious Validation Enabled | `Chart` | Chart showing details on number / % of workspaces that enabled drift detection/continious validation |  ✅  |
| Workspaces | Workspaces Count Over Time | `Time Series Graph` | Time series graph showing of # number of active workspaces over time |  ✅  | 
| Workspaces | Workspaces Status History | `Time Series Graph` | Time series graph showing workspace status over time |  ✅  | 
| Runs | Total Runs | `Counter` | Total number of runs executed  |  ✅  | 
| Runs | Total Run Failures | `Counter` | Total number of failed runs  |  ✅  | 
| Resources  | Resources Under Management | `Gauge` | Number of Resources Under Management(RUM) |  ✅  |
| Resources  | Resources Under Management Over Time | `Time Series Graph` | Number of Resources Under Management(RUM) over time |  ❌ |
| Policy Sets | Policy Set Count | `Gauge` | Current number of active policy sets organization  |  ✅  | 
| Policy Sets | Total Policy Check Failures | `Counter` | Total number of policy check failures  |  ✅  | 
| Policy Sets | Policy Set Summary | `Table` | Policy Sets Summary  |  ✅  | 
| Policy Sets  | Policy Type Distribution | `Chart` | Policy type distribution chart |  ✅  |
| Modules  | Modules Count | `Gauge` | Number of Modules in the Private Module Registry |  ✅  |
| Modules  | Modules Count Over Time | `Time Series Graph` | Number of Modules in the Private Module Registry over Time |  ❌  |
| Modules  | No-Code Module Distribution | `Chart` | Percentage of modules that are no-code ready |  ✅  |




## Usage

0. Clone this repo. 
1. Create a [Terraform Cloud API Token](https://app.terraform.io/app/settings/tokens)
2. Export your token and the name of your TFC Org:

```
export TF_API_TOKEN="TOKEN"
export TF_ORGANIZATIONS="ORG_NAME"
```

3. Spin up the application using Docker Compose

```
$ docker compose up -d
[+] Running 3/0
 ✔ Container tfbi-exporter-1    Running                                                                                                                                                                                                                                       0.0s 
 ✔ Container tfbi-prometheus-1  Running                                                                                                                                                                                                                                       0.0s 
 ✔ Container tfbi-grafana-1     Running                  


$ docker compose ps   
NAME                IMAGE             COMMAND                                                                                                                                                                                                                 SERVICE      CREATED         STATUS         PORTS
tfbi-exporter-1     tfbi-exporter     "./hot-reload.sh '--log-level debug'"                                                                                                                                                                                   exporter     4 seconds ago   Up 3 seconds   0.0.0.0:9100->9100/tcp
tfbi-grafana-1      grafana/grafana   "/run.sh"                                                                                                                                                                                                               grafana      4 seconds ago   Up 3 seconds   0.0.0.0:3000->3000/tcp
tfbi-prometheus-1   prom/prometheus   "/bin/prometheus --config.file=/etc/prometheus/prometheus.yml --storage.tsdb.path=/prometheus --web.console.libraries=/usr/share/prometheus/console_libraries --web.console.templates=/usr/share/prometheus/consoles"   prometheus   4 seconds ago   Up 3 seconds   0.0.0.0:9090->9090/tcp

```

4. Now you can access the dashboard using http://localhost:3000












