### Introduction
A demo plugin for the KrakenD API Gateway. The plugin simply intercepts the incoming request for the **/billing** endpoint and
send the payload to a NATS server.

## Prerequisites
- Go version: 1.20.11
- Krakend version: 2.5
- A local NATS server.

## GO env for this plugin
```
GO111MODULE="on"
GOARCH="arm64"
GOBIN=""
GOCACHE="/Users/anupam.gogoi.br/Library/Caches/go-build"
GOENV="/Users/anupam.gogoi.br/Library/Application Support/go/env"
GOEXE=""
GOEXPERIMENT=""
GOFLAGS=""
GOHOSTARCH="arm64"
GOHOSTOS="darwin"
GOINSECURE=""
GOMODCACHE="/Users/anupam.gogoi.br/go/pkg/mod"
GONOPROXY=""
GONOSUMDB=""
GOOS="darwin"
GOPATH="/Users/anupam.gogoi.br/go"
GOPRIVATE=""
GOPROXY="https://proxy.golang.org,direct"
GOROOT="/usr/local/go"
GOSUMDB="sum.golang.org"
GOTMPDIR=""
GOTOOLDIR="/usr/local/go/pkg/tool/darwin_arm64"
GOVCS=""
GOVERSION="go1.20.11"
GCCGO="gccgo"
AR="ar"
CC="clang"
CXX="clang++"
CGO_ENABLED="1"
GOMOD="/Users/anupam.gogoi.br/github/anupamgogoi/sonic/krakend-plugin/go.mod"
GOWORK=""
CGO_CFLAGS="-O2 -g"
CGO_CPPFLAGS=""
CGO_CXXFLAGS="-O2 -g"
CGO_FFLAGS="-O2 -g"
CGO_LDFLAGS="-O2 -g"
PKG_CONFIG="pkg-config"
GOGCCFLAGS="-fPIC -arch arm64 -pthread -fno-caret-diagnostics -Qunused-arguments -fmessage-length=0 -fdebug-prefix-map=/var/folders/x6/p1455tb119335zj9_nc4gc300000gp/T/go-build3659829989=/tmp/go-build -gno-record-gcc-switches -fno-common"

```

### Configuration file
Here is the content of the [krakend.json](https://github.com/cloudintegrator/krakend-plugin/blob/main/krakend.json) file.

### Build
Build the plugin using the docker image provided by KrakenD.
```
export DOCKER_BUILDKIT=0                                                                                                                                                    
export COMPOSE_DOCKER_CLI_BUILD=0
export DOCKER_DEFAULT_PLATFORM=linux/amd64
```
```
docker run -it -v "$PWD:/app" -w /app krakend/builder:2.5.0 go build -v -buildmode=plugin -o plugin/krakend-plugin.so .
```
### Run
Run the plugin using the KrakenD docker image.
```
docker run -it --rm --name nats -p "8080:8080" -v $PWD:/etc/krakend/ devopsfaith/krakend:2.5.0 run -c krakend.json
```

### Request example
```
curl --location 'http://localhost:8080/billing' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer XXXXXXXX' \
--data '{
    "client": 123,
    "payment": false
}'
```

### Problem
When tried to use the NATS Go client in this [function](https://github.com/cloudintegrator/krakend-plugin/blob/592711144fecf356017b78784e9cfc56b7007f1d/plugin.go#L77), the compiled SO of the plugin could
not be run. It generated errors like,
```
 plugin was built with a different version of package github.com/nats-io/nats.go/encoders/builtin
```
