# nearForm Docker distribution for Node.js Applications

[![CircleCI](https://circleci.com/gh/nearform/nodejs-distribution.svg?style=svg)](https://circleci.com/gh/nearform/nodejs-distribution)

This repository contains sources for an [s2i](https://github.com/openshift/source-to-image) builder image for Node.js releases from nodejs.org.

## Support
Commercial support for LTS versions is available by [contacting](https://www.nearform.com/contact/) nearForm directly.

## Flavors
We build and images based on various Operating Systems:
* Centos 7

  [![docker hub stats](http://dockeri.co/image/nearform/centos7-s2i-nodejs)](https://hub.docker.com/r/nearform/centos7-s2i-nodejs/)

* Alpine 3

  [![docker hub stats](http://dockeri.co/image/nearform/alpine3-s2i-nodejs)](https://hub.docker.com/r/nearform/alpine3-s2i-nodejs/)

The image can be used like any other image specifying it in your Dockerfile like this:
```
FROM nearform/centos7-s2i-nodejs
...
WORKDIR /opt/app-root/src
...
```
The Images are also prepared for use with [s2i](https://github.com/openshift/source-to-image/), a clean way to run your Node.js code in a controlled, squashed and secure image.

For more information about using these images with OpenShift, please see the
official [OpenShift Documentation](https://docs.openshift.org/latest/using_images/s2i_images/nodejs.html).

### Building instructions for CircleCI ###
The configuration can be found in /.circleci/config.json
The build is configured using the following [build parameters](https://circleci.com/docs/2.0/env-vars/#injecting-environment-variables-with-the-api).
* OS, Operating System, i.e. "centos7"
* VERSION, node.js version i.e."8.9.3"
* V8, V8 version i.e. "6.1.534.48"
* NPM, npm version i.e. "5.5.1"
* TAG, image tag i.e. "8.9.3"
* MAJOR, major version i.e. "8"
* MINOR, major and minor version i.e. "8.9"
* LTS, Long Time Support string i.e. "Carbon"
* PREBUILT, use prebuilt binaries if set to "T", otherwise build from sources

For the Red Hat images, there is a configurationfile at `.config/config.json` to map node.js versions to Red Hat projects. In the [Red Hat Catalog](https://access.redhat.com/containers/#/vendor/nearform) the different versions are organised in their own repository.
In order to push images to the Red Hat cetification registry a secret has to be provided for each project.
The configurationfile provides the ENV variables used to obtain a secret for each project, i.e. `NODEJS_8_SECRET`.

## Versions

Node.js versions currently provided:

<!-- versions.start -->
* **`8`**: (8.x, latest, LTS, [Red Hat Catalog](https://access.redhat.com/containers/?tab=overview#/registry.connect.redhat.com/nearform/nearform-s2i-nodejs8), supported)
* **`10`**: (10.x, latest, [Red Hat Catalog](https://access.redhat.com/containers/?tab=overview#/registry.connect.redhat.com/nearform/nearform-s2i-nodejs10), unsupported)
<!-- versions.end -->

## Source2image Usage

Using this image with OpenShift `oc` command line tool, or with `s2i` directly, will
assemble your application source with any required dependencies to create a new image.
This resulting image contains your Node.js application and all required dependencies,
and can be run either by OpenShift Origin or by Docker.

The [`oc` command-line tool](https://github.com/openshift/origin/releases) can be used to start a build, layering your desired nodejs `REPO_URL` sources into a centos7 image with your selected `RELEASE` of Node.js via the following command format:

```
oc new-app nearform/centos7-s2i-nodejs:RELEASE~REPO_URL
```

For example, you can run a build (including `npm install` steps), using  [`s2i-nodejs`](http://github.com/bucharest-gold/s2i-nodejs) example repo, and the `latest` release of
Node.js with:

```
oc new-app nearform/centos7-s2i-nodejs:latest~https://github.com/bucharest-gold/s2i-nodejs
```

### Environment variables

Use the following environment variables to configure the runtime behavior of the
application image created from this builder image.

NAME        | Description
------------|-------------
NPM_RUN     | Select an alternate / custom runtime mode, defined in your `package.json` file's [`scripts`](https://docs.npmjs.com/misc/scripts) section (default: npm run "start")
NPM_MIRROR  | Sets the npm registry URL
NODE_ENV    | Node.js runtime mode (default: "production")
HTTP_PROXY  | use an npm proxy during assembly
HTTPS_PROXY | use an npm proxy during assembly

One way to define a set of environment variables is to include them as key value pairs
in a `.s2i/environment` file in your source repository.

Example: `DATABASE_USER=sampleUser`

### Using Docker's exec

To change your source code in a running container, use Docker's [exec](http://docker.io) command:

```
docker exec -it <CONTAINER_ID> /bin/bash
```

After you [Docker exec](http://docker.io) into the running container, your current directory is set to `/opt/app-root/src`, where the source code for your application is located.

### Using OpenShift's rsync

If you have deployed the container to OpenShift, you can use [oc rsync](https://docs.openshift.org/latest/dev_guide/copy_files_to_container.html) to copy local files to a remote container running in an OpenShift pod.

## Builds

The [Source2Image cli tools](https://github.com/openshift/source-to-image/releases) are available as a standalone project, allowing you to run builds outside of OpenShift.

This example will produce a new docker image named `webapp`:

```
s2i build https://github.com/bucharest-gold/s2i-nodejs nearform/centos7-s2i-nodejs:current webapp
```

## Installation

There are several ways to make this base image and the full list of tagged Node.js releases available to users during OpenShift's web-based "Add to Project" workflow.

### For OpenShift Online
Those without admin privileges can install the latest Node.js releases within their project context with:

```
oc create -f https://s3.amazonaws.com/nodejs-distro-imagestreams/nearform_nodejs_8_imagestream.json
```

To ensure that each of the latest Node.js release tags are available and displayed correctly in the web UI, try upgrading / reinstalling the image stream:

```
oc delete is/centos7-s2i-nodejs ; oc create -f https://s3.amazonaws.com/nodejs-distro-imagestreams/nearform_nodejs_8_imagestream.json
```

If you've (automatically) imported this image using the [`oc new-app` example command](#usage), then you may need to clear the auto-imported image stream reference and re-install it.

### For Administrators

Administrators can make these Node.js releases available globally (visible in all projects, by all users) by adding them to the `openshift` namespace:

```
oc create -n openshift -f https://s3.amazonaws.com/nodejs-distro-imagestreams/nearform_nodejs_8_imagestream.json
```

To replace [the default SCL-packaged `openshift/nodejs` image](https://hub.docker.com/r/openshift/nodejs-010-centos7/) (admin access required), run:

```
oc delete is/nodejs -n openshift ; oc create -n openshift -f https://raw.githubusercontent.com/nearform/centos7-s2i-nodejs/master/centos7-s2i-nodejs.json
```

## Building your own Builder images

Clone a copy of this repo to fetch the build sources:

```
git clone https://github.com/nearform/centos7-s2i-nodejs.git
cd centos7-s2i-nodejs
```

### Requirements - docker-squash

`pip install docker-squash`

To build your own S2I Node.js builder images from scratch, run:

```
make all
```
