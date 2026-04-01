# 架构说明

## 分层架构

HTTP: router → middleware → handler → service → repository → DB
gRPC: server → interceptor → grpc/svc → service → repository → DB
Event: service → event.Publish() → bus → event/handler
MQ: service → mq.Producer → Kafka → mq.Consumer → event/handler

## 原则

- **单向依赖**：handler 只调 service，service 只调 repository
- **接口解耦**：每层通过 interface 交互，方便 mock 测试
- **事件驱动**：跨领域副作用通过事件总线解耦