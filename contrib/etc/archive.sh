#!/bin/bash

mkdir -p dist
git archive --prefix=build-tools/ --format=tar HEAD | gzip >dist/build-tools.tgz
echo "NODE_VERSION=$NODE_VERSION" > dist/versions
echo "OS=$OS" >> dist/versions
echo "DOCKERFILE=$DOCKERFILE" >> dist/versions
echo "IMAGE_TAG=$IMAGE_TAG" >> dist/versions
echo "LATEST=$LATEST" >> dist/versions
echo "MAJOR_TAG=$MAJOR_TAG" >> dist/versions
echo "MINOR_TAG=$MINOR_TAG" >> dist/versions
echo "IMAGE_NAME=$IMAGE_NAME" >> dist/versions
echo "NPM_VERSION=$NPM_VERSION" >> dist/versions
git rev-parse HEAD >dist/build-tools.revision
cd src/node-v${NODE_VERSION}
git archive -o ../../dist/node-v${NODE_VERSION}.tgz v${NODE_VERSION}
cd ../..
shasum dist/* >checksum
cp -v checksum dist/dist.checksum
tar czvf ${ARCHIVE} dist/*