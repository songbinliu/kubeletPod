package main

import (
	netutil "k8s.io/apimachinery/pkg/util/net"
	"net/http"

	"github.com/golang/glog"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/transport"
	"time"
)

func GetKubeConfig(masterUrl, configfile string) *rest.Config {

	// 1. in-cluster config: will use the Service Account Token mounted inside the Pod at
	//                    /var/run/secretes/kubernetes.io/serviceaccount
	if configfile == "" {
		config, err := rest.InClusterConfig()
		if err != nil {
			glog.Errorf("Failed to create config: %v", err)
			panic(err.Error())
		}

		return config
	}

	// 2. out-of-cluster config
	config, err := clientcmd.BuildConfigFromFlags(masterUrl, configfile)
	if err != nil {
		return nil
	}

	return config
}

func transportConfig(config *rest.Config, enableHttps bool) *transport.Config {
	cfg := &transport.Config{
		TLS: transport.TLSConfig{
			CAFile:   config.CAFile,
			CAData:   config.CAData,
			CertFile: config.CertFile,
			CertData: config.CertData,
			KeyFile:  config.KeyFile,
			KeyData:  config.KeyData,
		},
		BearerToken: config.BearerToken,
	}

	if enableHttps && !cfg.HasCA() {
		cfg.TLS.Insecure = true
		glog.Warning("insecure TLS transport.")
	}

	return cfg
}

// Generate a http.Transport based on rest.Config
func MakeTransport(config *rest.Config, enableHttps bool) (http.RoundTripper, error) {
	//1. get transport.config
	cfg := transportConfig(config, enableHttps)
	tlsConfig, err := transport.TLSConfigFor(cfg)
	if err != nil {
		return nil, err
	}
	if tlsConfig == nil {
		glog.Warningf("tlsConfig is nil.")
	}

	//2. http client
	rt := http.DefaultTransport
	if tlsConfig != nil {
		rt = netutil.SetOldTransportDefaults(&http.Transport{
			TLSClientConfig:     tlsConfig,
			TLSHandshakeTimeout: time.Second * 10,
		})
	}

	return transport.HTTPWrappersForConfig(cfg, rt)
}
