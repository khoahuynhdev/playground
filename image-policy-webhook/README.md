# ImagePolicyWebhook POC for Kubernetes

## ImagePolicyWebhook

- An ImagePolicyWebhook is a type of image that can be used in kubernetes cluster to enforce security on images before they are deployed in the cluster.
- This can be done by leverage the Admission Webhook in Kubernetes and the Admission Plugin "ImagePolicyWebhook"
- This code is for a POC on ImagePolicyWebhook where you write your own webhook and deploy it inside the cluster.

## Prior knowledge

- Linux
- Kubernetes (knowing how to use deployment, service, etc is enough)
- Openssl
- Go
