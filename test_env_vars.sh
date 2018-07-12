#!/bin/bash
export ND_NODEVERSION=9.5.0
export ND_V8VERSION=6.2.414.46
export ND_OS=centos7
export ND_DOCKERFILE=./image/alpine3/Dockerfile
export ND_IMAGETAG=9.5.0
export ND_LATEST=T
export ND_PREBUILT=T
export ND_MAJORTAG=9
export ND_MINORTAG=9.5
export ND_IMAGENAME=nearform/alpine3-s2i-nodejs
export ND_NPMVERSION=5.6.0
export ND_FROMDATA="{ \"from\": { \"image\": \"centos/s2i-base-centos7\", \"tag\": \"latest\", \"last_updated\": \"2018-06-05T16:17:48.190766Z\", \"sha\":\"\" } }"
export ND_PREBUILT=T