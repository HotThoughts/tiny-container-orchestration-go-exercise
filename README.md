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

## Deploy the tool

Depending on how the application is being deployed to k8s cluster, we have many ways to deploy the tool or use other tools to achieve the same behavior. For example, if the application is packaged into one helm chart, we can add the tiny tool definition to the halm chart yaml file.

We might also prefer other well-developed CD tools if they can cover our problem. In this case, [ArgoCD](https://argo-cd.readthedocs.io/en/stable/) is a good choice for managing deployments in a Kubernetes cluster. We can achieve the same behavior by enabling `selfHeal` attribute, which tells ArgoCD to watch out for changes that diverge from the desired state, and overwrites any manual changes.

## References

[^1]: [Official documentation of `docker buildx bake`](https://docs.docker.com/engine/reference/commandline/buildx_build/)
