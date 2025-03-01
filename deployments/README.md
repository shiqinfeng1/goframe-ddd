# app应用容器化部署指南

## 依赖检查

- Kubernetes: `>= 1.16.0-0`
- Helm: `>= 3.0`

假设 app项目根目录路径为 `APP_ROOT`

进入app项目根目录

```bash
$cd ${APP_ROOT}
```

## 容器化安装

具体安装步骤如下：

1) 生成配置文件

```bash
$export APP_APISERVER_INSECURE_BIND_ADDRESS=0.0.0.0
$export APP_AUTHZ_SERVER_INSECURE_BIND_ADDRESS=0.0.0.0
$./scripts/genconfig.sh scripts/install/environment.sh configs/app-apiserver.yaml > deployments/app/configs/app-apiserver.yaml
$./scripts/genconfig.sh scripts/install/environment.sh configs/app-authz-server.yaml > deployments/app/configs/app-authz-server.yaml
$./scripts/genconfig.sh scripts/install/environment.sh configs/app-pump.yaml > deployments/app/configs/app-pump.yaml
```

```bash
$kubectl -n app create configmap app --from-file=/etc/app/
$kubectl create configmap app-cert --from-file=/etc/app/cert
```

1) 使用Helm模板生成部署yaml文件: `app.yaml`

```bash
$helm template deployments/app > deployments/app.yaml
```

1) 安装app应用

```bash
$kubectl -n app apply -f deployments/app.yaml
```

1) 检查安装是否成功

检查app-apiserver

```bash
$export APP_APISERVER_HOST=x.x.x.x
$export APP_APISERVER_INSECURE_BIND_PORT=30080
$./scripts/install/test.sh app::test::apiserver
```

检查app-authz-server

```bash
$export APP_APISERVER_HOST=x.x.x.x
$export APP_APISERVER_INSECURE_BIND_PORT=30080
$./scripts/install/test.sh app::test::authzserver
```

检查app-pump

```bash
$export APP_APISERVER_HOST=x.x.x.x
$export APP_APISERVER_INSECURE_BIND_PORT=30080
$./scripts/install/test.sh app::test::pump
```

## Helm安装

1) 生成配置文件

```bash
$export APP_APISERVER_INSECURE_BIND_ADDRESS=0.0.0.0
$export APP_AUTHZ_SERVER_INSECURE_BIND_ADDRESS=0.0.0.0
$./scripts/genconfig.sh scripts/install/environment.sh configs/app-apiserver.yaml > deployments/app/configs/app-apiserver.yaml
$./scripts/genconfig.sh scripts/install/environment.sh configs/app-authz-server.yaml > deployments/app/configs/app-authz-server.yaml
$./scripts/genconfig.sh scripts/install/environment.sh configs/app-pump.yaml > deployments/app/configs/app-pump.yaml
```

1) Helm install

```bash
$helm install app deployments/app
```

1) 测试
