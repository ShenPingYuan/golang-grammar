# API 文档

## Base URL

`http://localhost:8080`

## 公开接口

| Method | Path              | 说明     |
|--------|-------------------|----------|
| GET    | /healthz          | 健康检查 |
| POST   | /api/v1/register  | 用户注册 |
| POST   | /api/v1/login     | 用户登录 |

## 认证接口（需 Bearer Token）

| Method | Path                    | 说明         |
|--------|-------------------------|--------------|
| GET    | /api/v1/users/me        | 当前用户信息 |
| GET    | /api/v1/users           | 用户列表     |
| GET    | /api/v1/users/{id}      | 用户详情     |
| POST   | /api/v1/orders          | 创建订单     |
| GET    | /api/v1/orders          | 我的订单     |
| GET    | /api/v1/orders/{id}     | 订单详情     |
| POST   | /api/v1/orders/{id}/pay | 支付订单     |