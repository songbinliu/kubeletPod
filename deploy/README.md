KubeletPod needs `view` access to Kubernetes Pods and Nodes from API. 
To make sure it has the right to access these resources, a [service account](https://kubernetes.io/docs/admin/service-accounts-admin/) is defined, 
and assign enought privileges for the pod.

1. Create a service account
In a designed namespace, create a service account.

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: turbo-user
```

create it by `oc` or `kubectl`
```console
# oc create -f my-sa.yaml -n <my-namespace>
```

checkt it by
```console
oc get sa -n <my-namespace> turbo-user
```

2. Assign role to the service account 
Assign `cluster-admin` role to the service account by:
```console
oadm policy add-cluster-role-to-user cluster-admin system:serviceaccount:<my-namespace>:turbo-user
```

or only assign `cluster-reader` role to the service account by:
```console
oadm policy add-cluster-role-to-user cluster-reader system:serviceaccount:<my-namespace>:turbo-user
```

3. Create pod with the service account
Create the pod with `turbo-user` service account, instead of the `default`service account:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: test-get-pods 
spec:
  serviceAccount: turbo-user
  containers:
  - name: test-pod
    image: beekman9527/kubelet-pod:latest 
    imagePullPolicy: Always
    args:
    - --v=4
    - --kubeletHttps=true
    - --kubeletPort=10250
```
