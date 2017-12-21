#!/bin/bash

CONFIG_PATH=.config/config.json

isDebug() {
    if [ ! -z "$DEBUG_BUILD" ]; then
        return 0
    fi
    return 1
}

checkVersionIsSet() {
    if [ -z "$VERSION" ]; then
        echo "Error: VERSION not set"
        exit 1
    fi
    return $true
}

getProjectSecret() {
    checkVersionIsSet
    MAJOR=$(echo ${VERSION} | cut -d'.' -f1)
    secret_varname=$(jq -r --arg version ${MAJOR} '.products[0].projects[] | select(.version==$version) | .secret_env_name' ${CONFIG_PATH})
    if [ ! -z "${!secret_varname}" ]; then
        echo "${!secret_varname}"
    else
        echo "Error: ${secret_varname}: not set"
        exit 1
    fi
}

getProjectId() {
    checkVersionIsSet
    MAJOR=$(echo ${VERSION} | cut -d'.' -f1)
    jq -r --arg version ${MAJOR} '.products[0].projects[] | select(.version==$version) | .project_id' ${CONFIG_PATH}
}

shouldPublish() {
    checkVersionIsSet
    MAJOR=$(echo ${VERSION} | cut -d'.' -f1)
    RH_V=${RH_MIN_VERSION-8}
    [ "$RH_V" -le "$MAJOR" ] && return
    return
}