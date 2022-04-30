# Tiny Container Orchestration

A tiny container orchestration tool written in golang. It watchers a list of containers, restarts stopped ones, as well as removes unknown ones.

## Build docker images

We use `docker-bake.hcl`[^1] to build two docker images

- `tinygoexercise:local-dev-api` is the `status-state-api`
- `tinygoexercise:local-dev-controller` is the `status-state-controller`

```shell
docker buildx bake
```

## Run container

### Run watcher

```sh
docker run \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -p 9001:9001 \
  tinygoexercise:local-dev-api
```

Note: In order to gain access to the Docker Engine API, we must gain access to the socket connect via mount `/var/run/docker.sock`.

### Run controller

```sh
docker run \
  -v /var/run/docker.sock:/var/run/docker.sock \
  --net=host \
  tinygoexercise:local-dev-controller
```

## References

[^1]: [Official documentation of `docker buildx bake`](https://docs.docker.com/engine/reference/commandline/buildx_build/)
