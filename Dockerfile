FROM devopsfaith/krakend:2.5
COPY krakend.json /gw/krakend.json
COPY plugin/krakend-nats-plugin.so /gw/krakend-nats-plugin.so


