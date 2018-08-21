#!/bin/bash

set -ex

# Node.js binaries don't run on alpine, because glibc is missing
if [ x"${PREBUILT}" = "xT" ] && [ ! -f /etc/alpine-release ]; then
    echo "Installing from prebuilt binary"
    tar -zxf /src/node-v${NODE_VERSION}-linux-x64.tar.gz -C /usr/local --strip-components=1
    npm install -g npm@${NPM_VERSION} -s &>/dev/null
else
    echo "INFO: Building from source"
    if [ -d /src/node-v${NODE_VERSION}.tar.gz ]; then
        tar -zxf /src/node-v${NODE_VERSION}.tar.gz -C /tmp/ --strip-components=1
        cd /tmp/
    elif [ -d /src/node-v${NODE_VERSION} ]; then
        cd /src/node-v${NODE_VERSION}
    else
        echo "No source files found, exiting..."
        exit 1
    fi

    ./configure
    make -j$(getconf _NPROCESSORS_ONLN)
    if [ cat Makefile | grep "test-only" ]; then
        make -j$(getconf _NPROCESSORS_ONLN) test-only
    fi
    make install
fi

if [ $(node --version | grep "6.0.0") ] && [ -f /etc/alpine-release ]; then
    echo 'skip yarn install'
else
    # Install yarn
    npm install -g yarn -s &>/dev/null
fi

# Fix permissions for the npm update-notifier
if [ ! -d /opt/app-root/src/.config ] ; then
  mkdir -p /opt/app-root/src/.config
fi

chmod -R 777 /opt/app-root/src/.config

# Delete NPM things that we don't really need (like tests) from node_modules
find /usr/local/lib/node_modules/npm -name test -o -name .bin -type d | xargs rm -rf

# Clean up the stuff we downloaded
rm -rf /tmp/node-v${NODE_VERSION} ~/.npm ~/.node-gyp ~/.gnupg /usr/share/man /tmp/* /usr/local/lib/node_modules/npm/man /usr/local/lib/node_modules/npm/doc /usr/local/lib/node_modules/npm/html
