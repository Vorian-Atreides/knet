# knet

Kubernetes services will always do a round robin load balancing between your pods,
which is fine in most of the time but isn't enough in many other.

Knet is a simple library, taking advantage of the Kubernetes' headless services
to give you more tools to build your services.

## Resolver

Resolver helps to retrieve Pod's hostnames from the Kubernetes DNS.

## Consistent hashing

Will help you to define a consistent routing from a given resource to a given
destination.

## Headless service

To help you prototype your application, you can use the below instructions to
setup a dummy application with its headless service.

First, create the deployment of your service:
> $ kubectl create deployment hello-node --image=gcr.io/hello-minikube-zero-install/hello-node

Then create the headless service with the selector on the previous deployment:
> $ kubectl expose deployment hello-node --cluster-ip=None --port=8080

