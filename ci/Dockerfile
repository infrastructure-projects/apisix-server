FROM docker.io/apache/apisix:3.9.1-redhat
COPY go-runner /plugins/go-runner
COPY ngx_tpl.lua /usr/local/apisix/apisix/cli/ngx_tpl.lua
RUN ln -s /config/apisix.yaml /usr/local/apisix/conf/apisix.yaml