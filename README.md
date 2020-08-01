# Kubernetes Graviton Scheduler Extender

This is a extension for Kubernetes that extends the default scheduler to support architecture-based filtering and prioritization for pod placement.
This extension is tested in EKS clusters, and designed to support mixed-architecture clusters (clusters that have both arm64 (Graviton2) and amd64 (x86) instances).

A Kubernetes Sheduler Extender is a built in feature of the default scheduler to support "last mile" webhook-based prioritization and filtering of candidate nodes for pod scheduling.

This extender will analyze every pod that's being scheduled to extract the container image names and tags. 
The extender then attempts to query the registry to determine the architectures supported.
If any containers in a pod do not have multi-architecture images available, the "filter" method of this extender will ensure that only the compatible nodes are included for scheduling.
If all images used by all containers in a pod have both amd64 and arm64 architectures available, the extender will reprioritize the node list to favor the arm64 instances over other instances.

## Motivation

This was built to support running Kubernetes Operators and Helm Charts (images that we didn't create) on a cluster with mixed architecture nodes. 
Our first use case was trying to run the Elasticsearch Operator (ECK) on EKS.
When deploying, we preferred to run on Graviton2 nodes to save cost, but adding these nodes to the cluster forced us to understand all images that the operator will deploy.
Operators (and Helm Charts) don't often contain a manifest of images and architectures.
This scheduler extender ensures that as upstream images start to include arm64 compatible images, the workload will automatically shift over to the Graviton2 nodes (on the next reschedule event).

## Installation

To install:

```shell
kubectl apply -f ./install/v<k8s_version>/graviton-scheduler-extender.yaml
```

The installer is specific to the version of Kubernetes you are running. 
For example, if you are are running `1.16.13`, the installation command would be:

```shell
kubectl apply -f ./install/v1.16.13/graviton-scheduler-extender.yaml
```

## Security

This is a scheduler extension.
It doesn't replace the default scheduler, but is a webhook that the Kubernetes scheduler will call into before finalizing the node placement.
As a result, this extension does require `cluster-admin` RBAC to install. 
The install scripts do this automatically.