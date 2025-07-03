group "default" {
  targets = ["tfbi-exporter"]
}

variable "VERSION" {
  default = "v0.5"
}

target "tfbi-exporter" {
  context = "."
  dockerfile = "Dockerfile"
  platforms = ["linux/amd64", "linux/arm64"]
  tags = ["nicolaka/tfbi-exporter:${VERSION}","nicolaka/tfbi-exporter:latest"]
  output = ["type=image","type=docker"]
  args = {
    TAG = "${VERSION}"
  }
}

