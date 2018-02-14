FROM alpine:3.7
# This image provides a Node.JS environment you can use to run your Node.JS applications.

## START s2i specific stuff
ENV \
    # Path to be used in other layers to place s2i scripts into
    STI_SCRIPTS_PATH=/usr/libexec/s2i \
    APP_ROOT=/opt/app-root \
    # The $HOME is not set by default, but some applications needs this variable
    HOME=/opt/app-root/src \
    PATH=/opt/app-root/src/bin:/opt/app-root/bin:$PATH

# This is the list of basic dependencies that all language Docker image can consume.
RUN apk update \
    && apk add --no-cache \
        libstdc++ \
    && apk add findutils \
        gettext \
        bash \
        clang \
        binutils-gold \
        g++ \
        gcc \
        linux-headers \
        make \
        python

# Copy extra files to the image.
COPY image/alpine3/root/ /

# Directory with the sources is set as the working directory so all STI scripts
# can execute relative to this path.
WORKDIR ${HOME}

ENTRYPOINT ["container-entrypoint"]
CMD ["base-usage"]

# Add default user
RUN adduser -u 1001 -S -G root -h ${HOME} -s /sbin/nologin default && \
    chown -R 1001:0 ${APP_ROOT}

RUN mkdir /licences && \
    mkdir -p /opt/app-root && \
    mkdir -p $STI_SCRIPTS_PATH

## END s2i specific stuff

EXPOSE 8080

# This image will be initialized with "npm run $NPM_RUN"
# See https://docs.npmjs.com/misc/scripts, and your repo's package.json
# file for possible values of NPM_RUN
ARG NODE_VERSION
ARG NPM_VERSION
ARG V8_VERSION
ARG PREBUILT

ENV NPM_RUN=start \
    NODE_VERSION=${NODE_VERSION} \
    NPM_VERSION=${NPM_VERSION} \
    V8_VERSION=${V8_VERSION} \
    NODE_LTS=false \
    NPM_CONFIG_LOGLEVEL=info \
    NPM_CONFIG_PREFIX=$HOME/.npm-global \
    PATH=$HOME/node_modules/.bin/:$HOME/.npm-global/bin/:$PATH \
    DEBUG_PORT=5858 \
    NODE_ENV=production \
    DEV_MODE=false \
    PREBUILT=${PREBUILT}

LABEL io.k8s.description="Platform for building and running Node.js applications" \
    io.k8s.display-name="Node.js $NODE_VERSION" \
    io.openshift.expose-services="8080:http" \
    io.openshift.tags="builder,nodejs,nodejs-$NODE_VERSION" \
    com.redhat.deployments-dir="/opt/app-root/src" \
    io.openshift.s2i.scripts-url=image:///usr/libexec/s2i \
    io.s2i.scripts-url=image:///usr/libexec/s2i \
    com.redhat.component="s2i-builder-image" \
    maintainer="nearForm Developer Support <devsupport@nearform.com>" \
    name="alpine3-s2i-nodejs" \
    vendor="nearForm Ltd" \
    summary="Alpine based s2i image for NodeJS applications" \
    description="Alpine based s2i image for NodeJS applications"

RUN mkdir /src

COPY ./src/ /src
RUN ls /src
COPY ./s2i/ $STI_SCRIPTS_PATH
COPY ./contrib/ /opt/app-root

### Atomic Help File - Write in Markdown, it will be converted to man format at build time.
### https://github.com/projectatomic/container-best-practices/blob/master/creating/help.adoc
COPY help/help.md /tmp/

### add licenses to this directory
COPY image/alpine3/licenses /licenses

RUN /opt/app-root/etc/install_node_source.sh

RUN apk del --purge clang \
    findutils \
    gettext \
    bash \
    binutils-gold

USER 1001

RUN ls $STI_SCRIPTS_PATH

# Set the default CMD to print the usage
CMD ${STI_SCRIPTS_PATH}/usage