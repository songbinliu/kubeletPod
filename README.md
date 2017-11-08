# kubeletPod
Tester to get Metrics from Kubelets and Pods from K8s API server.

## kubelet Port for metrics
|-|Kubernetes|Openshfit|GKE|
|-|-|-|-|
|kubeletPort| 10255 | 10250|10255|
|http.scheme| HTTP | HTTPS| HTTP|

Note1: Openshift version is 3.4, Kubernetes is 1.7;

Note2: GKE (google container engine) 1.6+ use [`http and port 10255`](https://github.com/prometheus/prometheus/issues/2606).

## Run it

#### Build it
```bash
make product

## build a docker image
sh build.sh
```

####Deploy it in kubernetes as a Pod:
```yaml
apiVersion: v1
kind: Pod
metadata:
    name: test-get-pods 
spec:
    containers:
    - name: test-pod
      image: beekman9527/kubelet-pod:latest 
      imagePullPolicy: Always
      args:
      - --v=4
      - --kubeletHttps=false
      - --kubeletPort=10255
```
