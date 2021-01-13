# init-tools

init-tools is a set of tools to use in init-container or in docker-entrypoint.sh

## setsvc

This tool get the listen port from the process in container and update the kubernetes service passed in the parameter

**Prerequisite**

The serviceAccount running the pod sould have the update access on the service

```sh
#> setsvc myservice
```

