## nearForm Node.js 8.x RHEL image - Getting Started

This guide is designed to quickly help you get a sample NodeJS application running using s2i and this rhrel7-s2i-nodejs image

### prerequisites
* [docker](https://docs.docker.com/engine/installation/)
* [s2i](https://github.com/openshift/source-to-image)

### usage
`s2i build https://github.com/sclorg/s2i-nodejs-container.git --context-dir=8/test/test-app/ nearform/rhel7-s2i-nodejs:8 nodejs-sample-app`

the output should look something like this:
```
Your branch is up-to-date with 'origin/master'.
Submodule 'common' (https://github.com/sclorg/container-common-scripts.git) registered for path 'common'
Cloning into '/private/var/folders/rg/kht1wrdx45s2d33y_t8vr6fw0000gn/T/s2i689433094/upload/tmp/common'...
Submodule path 'common': checked out 'a8d2885de09982496ffd7dd3e9a06296bd95e040'
---> Installing application source
---> Building your Node application from source
---> Using 'npm install'
up to date in 0.107s
Build completed successfully
```
Now you have a new docker image:
```
docker images | grep nodejs-sample-app
nodejs-sample-app                                                                                 latest              61621789b294        About a minute ago   520MB
```
which you can run:
```
docker run -p 8080:8080 nodejs-sample-app
npm info it worked if it ends with ok
npm info using npm@5.5.1
npm info using node@v8.9.1
npm info lifecycle node-echo@0.0.1~prestart: node-echo@0.0.1
npm info lifecycle node-echo@0.0.1~start: node-echo@0.0.1

> node-echo@0.0.1 start /opt/app-root/src
> node server.js

Server running on 0.0.0.0:8080
```
Then in a second terminal, verify that it actually works:
```
curl localhost:8080
This is a node.js echo service
Host: localhost:8080

node.js Production Mode: yes

HTTP/1.1
Request headers:
{ host: 'localhost:8080',
  'user-agent': 'curl/7.54.0',
  accept: '*/*' }
Request query:
{}
Request body:
{}

Host: 1d50fd173ff8
OS Type: Linux
OS Platform: linux
OS Arch: x64
OS Release: 4.9.49-moby
OS Uptime: 425
OS Free memory: 1627.421875mb
OS Total memory: 1999.0546875mb
OS CPU count: 2
OS CPU model: Intel(R) Core(TM) i5-6267U CPU @ 2.90GHz
OS CPU speed: 2900mhz
```
