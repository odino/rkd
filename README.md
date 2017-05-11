# rkd

> Think Dockerfile & docker-compose for rkt containers

![up](https://github.com/odino/rkd/blob/master/bin/up.png?raw=true)

`rkd` (aka *rock-it dev*) is a simple tool to build
and run [rkt containers](https://coreos.com/rkt) locally.

```
$ rkd
Development environments powered by rkt containers, with ease.

Usage:
  rkd [command]

Available Commands:
  build       Build the containers
  help        Help about any command
  up          Run the container
  version     Print the version number of rkd

Flags:
  -h, --help   help for rkd

Use "rkd [command] --help" for more information about a command.

$ sudo rkd up
/root/.rkd/5290facf0b502d01ba15b7de9a1b9633.aci already built
/root/.rkd/a023872855269062eca818f2ea8c0b32.aci already built
rkt --insecure-options=image --net=host run --interactive --volume src,kind=host,source=/home/odino/projects/go/src/github.com/odino/rkd/example/src /root/.rkd/a023872855269062eca818f2ea8c0b32.aci
[nodemon] 1.11.0
[nodemon] to restart at any time, enter `rs`
[nodemon] watching: *.*
[nodemon] starting `node index.js`
server started...
```

## Usage

Suppose you have a [NodeJS webserver](https://github.com/odino/rkd/blob/a0038678bed322e4d48f8340c3597ac211c60463/example/src/index.js)
running locally with docker-compose and want to convert it to rkt, without
using conversion tool like [docker2aci](https://github.com/appc/docker2aci) & the likes.

First, create your `prod.rkd` which is basically docker's `Dockerfile`:

```
set-name example.com/node-hello
dep add quay.io/coreos/alpine-sh
run -- apk add --update nodejs
copy src /src
set-working-directory /src
set-exec -- node index.js
port add www tcp 8080
```

and create a `dev.rkd` which is your new `docker-compose.yml`, where you can
define commands dependencies you only need on your local machine:

```
run -- npm install -g nodemon
mount add src src
set-exec -- nodemon index.js
```

That's it, now run `rkd up` in the [current folder](https://github.com/odino/rkd/tree/a0038678bed322e4d48f8340c3597ac211c60463/example):

```
sudo rkd up   
Building /root/.rkd/prod-5290facf0b502d01ba15b7de9a1b9633.aci
acbuild begin
acbuild set-name example.com/node-hello
acbuild dep add quay.io/coreos/alpine-sh
acbuild run -- apk add --update nodejs
Downloading quay.io/coreos/alpine-sh: [========================] 2.65 MB/2.65 MB
fetch http://dl-4.alpinelinux.org/alpine/v3.2/main/x86_64/APKINDEX.tar.gz
(1/4) Installing libgcc (4.9.2-r6)
(2/4) Installing libstdc++ (4.9.2-r6)
(3/4) Installing libuv (1.5.0-r0)
(4/4) Installing nodejs (0.12.10-r0)
Executing busybox-1.23.2-r0.trigger
OK: 28 MiB in 19 packages
acbuild copy src /src
acbuild set-working-directory /src
acbuild set-exec -- node index.js
acbuild port add www tcp 8080
acbuild write /root/.rkd/prod-5290facf0b502d01ba15b7de9a1b9633.aci
acbuild end
Building /root/.rkd/dev-a023872855269062eca818f2ea8c0b32.aci
acbuild begin ./prod.aci
acbuild run -- npm install -g nodemon
Downloading quay.io/coreos/alpine-sh: [========================] 2.65 MB/2.65 MB
npm WARN optional dep failed, continuing fsevents@1.1.1
/usr/bin/nodemon -> /usr/lib/node_modules/nodemon/bin/nodemon.js
nodemon@1.11.0 /usr/lib/node_modules/nodemon
├── ignore-by-default@1.0.1
├── undefsafe@0.0.3
├── es6-promise@3.3.1
├── debug@2.6.6 (ms@0.7.3)
├── minimatch@3.0.4 (brace-expansion@1.1.7)
├── touch@1.0.0 (nopt@1.0.10)
├── ps-tree@1.1.0 (event-stream@3.3.4)
├── lodash.defaults@3.1.2 (lodash.restparam@3.6.1, lodash.assign@3.2.0)
├── update-notifier@0.5.0 (is-npm@1.0.0, semver-diff@2.1.0, chalk@1.1.3, string-length@1.0.1, repeating@1.1.3, configstore@1.4.0, latest-version@1.0.1)
└── chokidar@1.7.0 (path-is-absolute@1.0.1, inherits@2.0.3, async-each@1.0.1, glob-parent@2.0.0, is-binary-path@1.0.1, is-glob@2.0.1, readdirp@2.1.0, anymatch@1.3.0)
acbuild mount add src src
acbuild set-exec -- nodemon index.js
acbuild write /root/.rkd/dev-a023872855269062eca818f2ea8c0b32.aci
acbuild end
rkt --insecure-options=image --net=host run --interactive --volume src,kind=host,source=/home/odino/projects/go/src/github.com/odino/rkd/example/src /root/.rkd/dev-a023872855269062eca818f2ea8c0b32.aci
[nodemon] 1.11.0
[nodemon] to restart at any time, enter `rs`
[nodemon] watching: *.*
[nodemon] starting `node index.js`
server started...
```

and you can head to [localhost:8080](http://localhost:8080).

## Installation

> Make sure [acbuild](https://github.com/containers/build) is installed in your system.

Builds for a few systems are available in the [releases](https://github.com/odino/rkd/releases).

Alternatively, you can compile it straight away:

```
git clone git@github.com:odino/rkd.git
cd rkd

go build -o rkd main.go
mv rkd /usr/local/bin
```

and then you have the `rkd` executable up & running.

## Why this?

One of the arguments against rkt is that building and running containers seems
generally more complicated than using docker, so I decided to figure out a way
to replicate docker's simplicity on dev environments -- 2 files, one command,
running app.

The `*.rkd` files are basically a list or `acbuild` instructions used for
building 2 ACIs (`prod.aci` & `dev.aci`): `rkd` scans them, building the ACIs,
and bases `dev.aci` off of what it build in `prod.aci`.

## Troubleshooting

There's a plethora of stuff that could / needs to be done here as this is an early stage weekend
project. There's very less error handling etc in the codebase and that's
something I wish to work on granted that (1) I can find the time and (2) there's
some interest here.
