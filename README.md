# Round-Robin Scheduler

This is a very simple scheduler demonstrating how works a Kubernetes scheduler.

## Build

```shell
docker build -t <your_registry>/scheduler-round-robin .
docker push <your_registry>/scheduler-round-robin
```

## Deploy

Replace <your_registry> in `scheduler-round-robin.yaml`, then:

```shell
kubectl apply -f scheduler-round-robin.yaml
```
