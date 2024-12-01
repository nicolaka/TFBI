group "default" {
  targets = ["exporter"]
}

target "exporter" {
  context = "."
  dockerfile = "Dockerfile.dev"

  tags = ["dev"]
}

