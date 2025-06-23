provider "tfe" {
  hostname = "app.terraform.io"
  organization = var.org
}

variable "org" {
  default = "test"
}

variable "project" {
  default = "test"
}

data "tfe_project" "project" {
    name = var.project

variable "workspace_count" {
  default = 0
}

variable "team_count" {
  default = 0
}

variable "project_count" {
  default = 0
}

resource "tfe_team" "team" {
   count = var.team_count
   name  = "scale-team-${count.index + 1}"
   organization = var.org
}

resource "tfe_workspace" "ws" {
  count = var.workspace_count
  name  = "scale-workspace-${count.index + 1}"
  organization = var.org
  project_id = data.tfe_project.project.id
}

resource "tfe_project" "project" {
  count = var.project_count
  name  = "scale-project-${count.index + 1}"
  organization = var.org
}
