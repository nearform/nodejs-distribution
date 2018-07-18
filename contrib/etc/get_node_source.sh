#!/bin/bash

set -ex

NODE_VERSION="${1}"
SRCDIR="${2}"
NODEDIR="node-v${NODE_VERSION}"
NODE_REPO=${${3}:-https://github.com/nodejs/node.git}
COMMIT=${4}

mkdir -p "${SRCDIR}" || exit 1

# Download and install a binary from nodejs.org
# Add the gpg keys listed at https://github.com/nodejs/node
for key in \
        94AE36675C464D64BAFA68DD7434390BDBE9B9C5 \
        B9AE9905FFD7803F25714661B63B535A4C206CA9 \
        77984A986EBC2AA786BC0F66B01FBB92821C587A \
        56730D5401028683275BD23C23EFEFE93C4CFFFE \
        71DCFD284A79C3B38668286BC97EC7A07EDE3FC1 \
        FD3A5288F042B6850C66B31F09FE44734EB7990E \
        C4F0DFFF4E8C1A8236409D08E73BC641CC11F4C8 \
        DD8F2338BAE7501E3DD5AC78C273792F7D83545D \
        8FCCA13FEF1D0C2E91008E09770F7A9A5AE15600 \
    ; do
    gpg -q --keyserver ipv4.pool.sks-keyservers.net --recv-keys "$key" || \
    gpg -q --keyserver pgp.mit.edu --recv-keys "$key" || \
    gpg -q --keyserver keyserver.pgp.com --recv-keys "$key" ; \
    echo "$key:6" | gpg --import-ownertrust
done

# Get the node binary and it's shasum
cd "${SRCDIR}"
if [[ x"${PREBUILT}" == "xT" ]] && [ "${OS}" != "alpine3" ]; then

    if command -v sha256sum; then
        SHACMD=sha256sum
    elif command -v shasum; then
        SHACMD='shasum -a 256 '
    else
        echo "sha256sum or shasum required, exiting.."
        exit 1
    fi
    curl -O -sSL https://nodejs.org/dist/v${NODE_VERSION}/SHASUMS256.txt.asc
    gpg --verify SHASUMS256.txt.asc || exit 1
    curl -O -sSL https://nodejs.org/dist/v${NODE_VERSION}/node-v${NODE_VERSION}-linux-x64.tar.gz
    grep " node-v${NODE_VERSION}-linux-x64.tar.gz" SHASUMS256.txt.asc | ${SHACMD} -c -
else
    if [ -d ${NODEDIR}/.git ]; then
        cd ${NODEDIR}
        git fetch --all
    else
        rm -Rf ${NODEDIR}
        git clone ${NODE_REPO} ${NODEDIR}
        cd ${NODEDIR}
    fi
    if [[ x"${COMMIT}" == "x" ]]; then
        git verify-tag v${NODE_VERSION} || exit 1
        git checkout tags/v${NODE_VERSION}
    else
        git checkout ${COMMIT} || exit 1
    fi
    cd "${SRCDIR}"
fi
