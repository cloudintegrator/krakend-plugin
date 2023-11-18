# TODO
FROM devopsfaith/krakend:2.5
COPY krakend.json /etc/krakend/krakend.json
COPY plugin/krakend-plugin.so /etc/krakend/plugin/krakend-plugin.so
ENTRYPOINT ["krakend", "run", "-c", "/etc/krakend/krakend.json"]
