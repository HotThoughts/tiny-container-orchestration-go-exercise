variable "IMAGE_NAME" {
  default = "tinygoexercise"
}

group "default" {
  targets = ["watcher", "controller"]
}

target "watcher" {
  dockerfile = "Dockerfile"
  tags = ["${IMAGE_NAME}:local-dev-api"]

  args = {
    GO_APP_NAME = "watcher"
  }
  cache-to = ["type=inline"]
  cache-from = ["type=registry,ref=${IMAGE_NAME}:local-dev-api"]
}

target "controller" {
  dockerfile = "Dockerfile"
  tags = ["${IMAGE_NAME}:local-dev-controller"]

  args = {
    GO_APP_NAME = "controller"
  }
  cache-to = ["type=inline"]
  cache-from = ["type=registry,ref=${IMAGE_NAME}:local-dev-controller"]
}
