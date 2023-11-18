### Build with Docker
```
export DOCKER_BUILDKIT=0                                                                                                                                                    
export COMPOSE_DOCKER_CLI_BUILD=0
export DOCKER_DEFAULT_PLATFORM=linux/amd64
```
```
docker run -it -v "$PWD:/app" -w /app krakend/builder:2.5.0 go build -buildmode=plugin -o plugin/krakend-plugin.so .
```
### Run
```
docker run -it --rm --name nats -p "8080:8080" -v $PWD:/etc/krakend/ devopsfaith/krakend:2.5.0 run -c krakend.json
```

### Request example
```
curl --location --request GET 'http://localhost:8080/billing' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer XXXXXXXX' \
--data '{
    "client": 123,
    "payment": true
}'
```