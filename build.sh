#!/bin/bash

tag=beekman9527/kubelet-pod

docker build -t $tag .
docker push $tag
