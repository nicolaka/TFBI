group "default" {
  targets = ["tfbi-exporter"]
}

variable "TAG" {
  default = "v0.3"
}

target "tfbi-exporter" {
  context = "."
  dockerfile = "Dockerfile"
  platforms = ["linux/amd64", "linux/arm64"]
  tags = ["nicolaka/tfbi-exporter:${TAG}","nicolaka/tfbi-exporter:latest"]
  output = ["type=image"]
}

