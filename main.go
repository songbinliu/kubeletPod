
package main

import (
	"flag"
	"github.com/golang/glog"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	client "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	api "k8s.io/client-go/pkg/api/v1"
)

var (
	masterUrl = ""
	kubeConfig = ""
)

func setFlags(c *kubeletConfig) {
	flag.StringVar(&masterUrl, "masterUrl", "", "kubernetes master url")
	flag.StringVar(&kubeConfig, "kubeConfig", "", "path to the kubeconfig(admin.kubeconfig) file.")
	flag.IntVar(&(c.timeout), "timeout", c.timeout, "the timeout when connecting to kubelet to get metrics.")
	flag.IntVar(&(c.port), "kubeletPort", c.port, "the Port of kubelet to get metrics.")
	flag.BoolVar(&(c.enableHttps), "kubeletHttps", c.enableHttps, "Whether to access the kubelet with https or not(http).")

	flag.Set("alsologtostderr", "true")
	flag.Parse()
}

func test_pod(kclient *client.Clientset, stop chan struct{}) {

	for {
		pods, err := kclient.CoreV1().Pods("").List(metav1.ListOptions{})
		if err != nil {
			glog.Error(err.Error)
			panic(err.Error())
		}
		glog.Infof("There are %d pods in the cluster\n", len(pods.Items))

		// Examples for error handling:
		// - Use helper functions like e.g. errors.IsNotFound()
		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
		podName := "kubeturbo-gke-max"
		_, err = kclient.CoreV1().Pods("default").Get(podName, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			glog.Errorf("Pod(%v) not found\n", podName)
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			glog.Errorf("Error getting pod %v: %v\n", podName, statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			glog.Infof("Found pod %v\n", podName)
		}

		glog.V(2).Infof("sleeping 20 seconds")
		time.Sleep(20 * time.Second)
	}
}

func test_kubelet(config *rest.Config, kletConfig *kubeletConfig, hosts []*api.Node) {
	kletClient, err := NewKubeletClient(kletConfig, config)
	if err != nil {
		glog.Errorf("Failed to create kubeletClient: %v", err)
		return
	}

	i := 0
	host := hosts[i % len(hosts)]
	for {
		glog.V(2).Infof("Get stats for host: %v", host.Name)
		kletClient.GetMachineInfo(host.Name)
		kletClient.GetSummary(host.Name)
		glog.V(2).Infof("sleeping 20 seconds")
		time.Sleep(30*time.Second)
		i ++
		i = i % len(hosts)
	}
}

func getNodes(kclient *client.Clientset) ([]*api.Node, error) {
	nodeList, err := kclient.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		glog.Errorf("failed to list all nodes in the cluster: %s", err)
		return nil, err
	}

	glog.V(2).Infof("There are %d nodes", len(nodeList.Items))

	result := make([]*api.Node, len(nodeList.Items))
	for i := 0; i < len(nodeList.Items); i ++ {
		result[i] = &nodeList.Items[i]
		glog.V(2).Infof("node: %v", result[i].Name)
	}

	return result, nil
}

func main() {
	kletConfig := NewDefaultKubeletConfig()
	setFlags(kletConfig)

	//1. create kclient
	config := GetKubeConfig(masterUrl, kubeConfig)
	if config == nil {
		glog.Errorf("Failed to create InCluster config")
		return
	}

	kclient, err := client.NewForConfig(config)
	if err != nil {
		glog.Errorf("Failed to create kubeClient: %v", err)
		return
	}

	//2. get all the nodes
	nodes, err := getNodes(kclient)
	if err != nil || len(nodes) == 0 {
		glog.Errorf("Failed to get nodes from API server.")
		return
	}

	//3. start the workers
	stop := make(chan struct{})
	go test_pod(kclient, stop)
	go test_kubelet(config, kletConfig, nodes)

	select {}
}
