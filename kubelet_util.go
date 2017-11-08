package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	kubeclient "k8s.io/client-go/rest"
)

const (
	specPath    = "/spec"
	summaryPath = "/stats/summary"

	// Note: openshift cluster.kubeletPort is 10250; kubernetes is 10255
	defaultKubeletPort = 10250
	defaultTimeOut     = 10
)

type kubeletConfig struct {
	enableHttps bool
	port        int
	timeout     int //timeout in seconds
}

type KubeletClient struct {
	port    int
	scheme  string
	hclient *http.Client
}

func NewDefaultKubeletConfig() *kubeletConfig {
	return &kubeletConfig{
		enableHttps: false,
		port: defaultKubeletPort,
		timeout: defaultTimeOut,
	}
}

func NewKubeletClient(c *kubeletConfig, kc *kubeclient.Config) (*KubeletClient, error) {
	transport, err := MakeTransport(kc, c.enableHttps)
	if err != nil {
		return nil, err
	}

	hclient := &http.Client{
		Timeout:   time.Second * time.Duration(c.timeout),
		Transport: transport,
	}

	scheme := "http"
	if c.enableHttps {
		scheme = "https"
	}

	kclient := &KubeletClient{
		port:    c.port,
		scheme:  scheme,
		hclient: hclient,
	}

	return kclient, nil
}

func (c *KubeletClient) GetMachineInfo(host string) error {

	requestURL := url.URL{
		Scheme: c.scheme,
		Host:   fmt.Sprintf("%s:%d", host, c.port),
		Path:   specPath,
	}

	req, err := http.NewRequest("GET", requestURL.String(), nil)
	if err != nil {
		glog.Errorf("failed to build request[%v]: %v", requestURL, err)
		return err
	}
	glog.V(2).Infof("[request]: %s", req.URL.String())

	return SendRequestGetValue(c.hclient, req, nil)
}

func (c *KubeletClient) GetSummary(host string) error {
	requestURL := url.URL{
		Scheme: c.scheme,
		Host:   fmt.Sprintf("%s:%d", host, c.port),
		Path:   summaryPath,
	}

	req, err := http.NewRequest("GET", requestURL.String(), nil)
	if err != nil {
		glog.Errorf("failed to build Summary request[%v]: %v", requestURL, err)
		return err
	}
	glog.V(2).Infof("[request]: %s", req.URL.String())

	return SendRequestGetValue(c.hclient, req, nil)
}

func SendRequestGetValue(client *http.Client, req *http.Request, value interface{}) error {
	response, err := client.Do(req)
	if err != nil {
		e := fmt.Errorf("failed to send http request[%+v]: \n%v", req, err)
		glog.Error(e)
		return e
	}

	if response.StatusCode == http.StatusNotFound {
		e := fmt.Errorf("%s was not found", req.URL.String())
		glog.Error(e)
		return e
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if response.StatusCode != http.StatusOK {
		e := fmt.Errorf("request[%s] failed: %d-%q, body: %s", req.URL.String(), response.StatusCode, response.Status, string(body))
		glog.Error(e)
		return e
	}

	if value == nil {
		glog.V(2).Infof("body:\n%s", string(body))
		return nil
	}

	if err = json.Unmarshal(body, value); err != nil {
		e := fmt.Errorf("")
		glog.Error(e)
		return e
	}

	return nil
}
