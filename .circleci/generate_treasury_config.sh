#!/bin/sh

export DOCKER_USER=$1
export DOCKER_PASS=$2
export TREASURY_USER=$3
export TREASURY_PASS=$4

envsubst '${DOCKER_USER}:${DOCKER_PASS}:${TREASURY_USER}:${TREASURY_PASS}' < /opt/treasury-cli/config.yaml.example > /opt/treasury-cli/config.yaml
