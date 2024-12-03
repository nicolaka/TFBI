group "default" {
  targets = ["tfbi-exporter"]
}

variable "TAG" {
  default = "latest"
}

target "tfbi-exporter" {
  context = "."
  dockerfile = "Dockerfile"
  platforms = ["linux/amd64", "linux/arm64"]
  tags = ["nicolaka/tfbi-exporter:${TAG}"]
  output = ["type=image","type=docker"]
}

