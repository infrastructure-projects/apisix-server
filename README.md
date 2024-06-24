# apisix-server

基于APISIX实现的微服务网关

该项目基于[apisix-go-plugin-runner](https://github.com/apache/apisix-go-plugin-runner)
实现的[APISIX](https://apisix.apache.org/)的扩展程序.

程序通过读取`/config/application.yaml`加载配置信息.当该配置文件发生变更时，程序会自动重新加载配置文件，无需重启.

### 一、SecurityFilter

> 对请求进行安全检查

配置内容:
```yaml
application:
  security:
    block-list:
      paths: # 阻止访问指定的path
        - /xxxx
        - /xxxx
      headers: # 阻止请求头中携带指定的key
        - x-app-username
      xss: # 阻止请求头中携带能够发起xss攻击的代码
        - $
```

### 二、构建命令:
```shell
make build
mv go-runner ./ci/go-runner
cd ./ci
# docker build command...
```

### 版本信息:

- APISIX: 3.x
- apisix-go-plugin-runner: 0.5.0
- golang: 1.22