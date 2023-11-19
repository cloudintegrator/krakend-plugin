### Introduction
A demo plugin for the KrakenD API Gateway. The plugin simply intercepts the incoming request for the **/billing** endpoint and
send the payload to a NATS server.

## Prerequisites
- Go version: 1.20.11
- Krakend version: 2.5
- A local NATS server.

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
docker run -it -v "$PWD:/app" -w /app krakend/builder:2.5.0 go build -buildmode=plugin -o plugin/krakend-plugin.so .
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