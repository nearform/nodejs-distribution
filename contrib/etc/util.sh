#!/bin/bash

CONFIG_PATH=.config/config.json

osSupported() {
    declare -a arr=("rhel7" "centos7" "alpine3")

    local retVal=1
    for i in "${arr[@]}"; do
        if [[ $i = $1 ]]; then
            retVal=0
            break
        fi
    done
    return ${retVal}
}

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

getBaseImageForOs() {
    # echo ${OS}
    FROM=$(cat image/${OS}/Dockerfile | grep FROM | awk -F " " '{print $2}')
    SLASHES=$(echo $FROM | awk -F"/" '{print NF-1}')
    if [[ $SLASHES = 1 ]]; then
        FROM_OWNER=$(echo $FROM | awk -F '/' '{print $1}')
        FROM=$(echo $FROM | awk -F '/' '{print $2}')
        # echo $FROM
    elif [[ $SLASHES = 2 ]]; then
        FROM_OWNER=$(echo $FROM | awk -F '/' '{print $1}')
        FROM=$(echo $FROM | awk -F '/' '{print $2"%252F"$3}')
        # echo $FROM
    else
        FROM_OWNER="library"
    fi
    FROM_IMAGE=$(echo $FROM | awk -F ':' '{print $1}')
    if [[ $FROM = *":"* ]]; then
        FROM_TAG=$(echo $FROM | awk -F ':' '{print $2}')
    else
        FROM_TAG="latest"
    fi
    if [[ "$OS" = "alpine3" || "$OS" = "centos7" ]]; then
        URL="https://hub.docker.com/v2/repositories/${FROM_OWNER}/${FROM_IMAGE}/tags/${FROM_TAG}/"
        FROM_DATETIME=$(curl -s $URL | jq -r '.last_updated')
    fi
    if [[ "$OS" = "rhel7" ]]; then
        URL="https://www.redhat.com//wapps/containercatalog/rest/v1/repository/registry.access.redhat.com/${FROM_IMAGE}"
        JSON=$(curl -s $URL)
        FROM_DATETIME=$(echo $JSON | jq -r '.processed[0].images[0].repositories[0].push_date')
        FROM_TAG=$(echo $JSON | jq -r '.processed[0].images[0].repositories[0].tags[] | select(.name | contains("1-")).name')
        FROM_IMAGE=$(echo $FROM_IMAGE | sed 's#%252F#/#')
    fi
    # echo $URL
    echo "{ \"from\": { \"image\": \"$FROM_OWNER/$FROM_IMAGE\", \"tag\": \"$FROM_TAG\", \"last_updated\": \"$FROM_DATETIME\" } }"
}