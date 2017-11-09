# kubeletPod
When we want to monitor the Pods in a Kubernetes cluster, we need the pods information (labels, resource limit/requests), 
as wells as the resource(CPU/Memory/Network) usage by containers. 

Pods information can be got from Kubernetes API server, and the resource usage metrics can be got by scraping the kubelets.
However, different versions/platforms/flavors(such as openshift, GKE) of Kubernetes clusters have different settings.
For example, to access GKE Kubernetes API server, [it is best to use the in-cluster settings](https://github.com/kubernetes/client-go/issues/242),
otherwise, the client needs to install `gcloud`.

So the purpose of this project is to **help developer find the best way to get metrics from kubelets, and get Pods from Kubernetes API server**.


## kubelet Port for metrics
|-|Kubernetes|Openshfit|GKE|
|-|-|-|-|
|kubeletPort| 10255 | 10250|10255|
|http.scheme| HTTP | HTTPS| HTTP|

Note1: Openshift version is 3.4, Kubernetes is 1.7;

Note2: GKE (google container engine) 1.6+ uses [`http and port 10255`](https://github.com/prometheus/prometheus/issues/2606).

## Run it

#### Build it
```bash
make product

## build a docker image
sh build.sh
```

#### Deploy it in kubernetes as a Pod:
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
