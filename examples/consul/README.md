# running consul tests on k8s
content here has been derived from https://github.com/kelseyhightower/consul-on-kubernetes

pl. follow these steps:
```bash
$ GOSSIP_ENCRYPTION_KEY=$(consul keygen)
$ kubectl create secret generic consul --from-literal="gossip-encryption-key=${GOSSIP_ENCRYPTION_KEY}"
$ 

```