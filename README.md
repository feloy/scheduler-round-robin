# Round-Robin Scheduler

This is a very simple scheduler demonstrating how works a Kubernetes scheduler.

## Deploy

```shell
kubectl apply -f scheduler-round-robin.yaml
```

## Build

```shell
docker build -t <your_registry>/scheduler-round-robin .
docker push <your_registry>/scheduler-round-robin
```
 Then replace the `deployment.spec.container.image` value with `<your_registry>/scheduler-round-robin`.