# meow-nook
## 开发文档
### 划分线

单实线，这里比较帅气，之后的划分都可以用单实线
```api
// ──────────────────────────────────────────────
// Error codes
// ──────────────────────────────────────────────

```

### 版本
- protoc
```bash
protoc --version
libprotoc 34.1
```

- goctl 
```bash
goctl --version                                                
goctl version 1.10.1 windows/amd64 
```

- protoc-gen-go-grpc 
```bash
protoc-gen-go-grpc --version
protoc-gen-go-grpc 1.6.1
```

- protoc-gen-go 
```bash
protoc-gen-go --version     
protoc-gen-go.exe v1.36.11
```

### 命令行

#### 生成代码

进入对应的服务目录，更换具体的proto文件路径和proto_path路径，其中third_party目前没有引入，位置暂定
- user服务
```bash
goctl rpc protoc pb/user/v1/user.proto --proto_path=. --proto_path=..\..\third_party --go_out=. --go-grpc_out=. --zrpc_out=. --client=true  
```

- task 服务
```bash
goctl rpc protoc pb/task/v1/task.proto --proto_path=. --proto_path=..\..\third_party --go_out=. --go-grpc_out=. --zrpc_out=. --client=true
```

- post 服务
```bash
goctl rpc protoc pb/post/v1/post.proto --proto_path=. --proto_path=..\..\third_party --go_out=. --go-grpc_out=. --zrpc_out=. --client=true
```

- cat 服务
```bash
goctl rpc protoc pb/cat/v1/cat.proto --proto_path=. --proto_path=..\..\third_party --go_out=. --go-grpc_out=. --zrpc_out=. --client=true
```

- adoption 服务
```bash
goctl rpc protoc pb/adoption/v1/adoption.proto --proto_path=. --proto_path=..\..\third_party --go_out=. --go-grpc_out=. --zrpc_out=. --client=true
```

- gateway 服务
```bash
goctl api go -api gateway.api -dir . --style go_zero
```

- 生成 swag 文档
```bash
goctl api swagger -api gateway.api -dir ./docs -filename swagger
```

### impot 排序
对 import 进行排序的命令 可以直接在根目录里执行
```bash
gci.exe write -s standard -s "prefix(github.com/luyb177/meow-nook)" -s default   --skip-generated .
```
下面是基于你这段说明、并按“新建一个 grpc 微服务从 0 到能用”的路径做的结构化优化版（更像团队内部 Wiki/模板），同时把关键链路（grpc → gateway → handler）和“必须做的初始化点”强调清楚，方便复制到新服务里直接照做。

---

## 1) logger 使用规范

### 1.1 引入
在需要打印日志的地方引入 "github.com/luyb177/meow-nook/common/logger"。

### 1.2 初始化（关键：在服务入口创建实例 + 注入到上下文/中间件链路）
参考 **user 服务的 `user.go`**：

- 在入口（main/启动文件）创建 logger 实例
- 把 logger 相关中间件加入 grpc server 中

这样后续在 **logic / task / 任何下游** 都可以直接使用统一的 logger（并且能带上 traceId / reqId 等上下文信息）。

### 1.3 使用
初始化完成后，在业务逻辑中直接调用 logger 即可。

> 约定建议：
> - handler 只做参数解析 / 返回；日志尽量放在 logic（入口日志除外）。
> - 重要错误日志要带关键字段（uid、orderId、bizKey），避免只打 `err`。

---

## 2) errorx 使用规范（grpc ↔ gateway ↔ handler 全链路）

### 2.1 引入
- grpc 服务端：引入 `common/errorx`
- gateway 服务：同样引入 `common/errorx`（用于把 grpc err 转为可返回的业务错误）

### 2.2 grpc 服务端：只返回两类错误（强约束）
参考 **user 服务的 `test.go`**：

1) **直接返回 errorx 单例错误**
- 适用于“明确的业务错误码 + 明确的错误文案”，不需要额外包裹原始 error 的场景

2) **返回 `errorx.Wrap(...)` 包装后的 error**
- 适用于“需要保留原始 err（用于日志/排查），同时对外返回业务码和友好文案”的场景

> 建议规则：
> - 对外文案（msg）要稳定、可读；原始 err 只用于内部排查，通过 Wrap 保留。
> - 不要在 grpc 层直接返回 fmt.Errorf 之类的裸 error（否则 gateway 侧不好统一转换）。

### 2.3 gateway：遇到 grpc 错误，统一转换
在 gateway 调用 grpc 后，如果拿到 err：

- **统一走** `errorx.FromGRPC(err)`
- 得到的就是标准化的 errorx（包含 code/msg/detail 等）

### 2.4 handler：统一返回（你这段示例可以作为模板）
handler 里不要自己拼 JSON，统一用的 `httpresp.JsonBaseResponseCtx(...)`（或等价函数）：

```go
if err := httpx.Parse(r, &req); err != nil {
    httpresp.JsonBaseResponseCtx(r.Context(), w,
        errorx.Wrap(errorx.CodeBadRequest, "请求参数错误", err),
    )
    return
}

l := auth.NewTestLogic(r.Context(), svcCtx)
resp, err := l.Test(&req)
if err != nil {
    httpresp.JsonBaseResponseCtx(r.Context(), w, err)
    return
}
httpresp.JsonBaseResponseCtx(r.Context(), w, resp)
```

> 这里的关键点：
> - Parse 参数错误：属于 handler 自己的错误，直接 Wrap 成 CodeBadRequest。
> - logic 返回 error：可能已经是 errorx 或 gateway 转换后的 errorx，直接原样返回即可（保持一致性）。

---

## 3) 新建 grpc 微服务必须加的“入口中间件”（很关键）

> 你原文提到“每次新建 grpc 微服务要在入口程序里加入中间件”，这里建议把它写成强制 checklist。

参考 **user 服务 `user.go`**，新服务启动时必须完成：

- 注册 logger 中间件（让后续逻辑都能拿到 logger/trace）
- 注册 errorx 相关中间件（如果你们有统一 panic recover / 错误转换）
  -（如果有）注册 tracing / metrics / recovery / auth 等通用中间件

完成后：
- grpc 服务端才能规范产出 errorx
- 业务逻辑里才能随处可用 logger
- gateway 才能稳定把 grpc error 映射为标准响应

---

## 4) 消息队列（Kafka + Redis 延迟队列）使用规范

### 4.1 依赖与初始化（每个服务都要有）
每个微服务如果要用任务系统，必须同时初始化：

- **Kafka client/producer/consumer**
- **Redis client**（用于延迟队列/定时触发）

参考：
- **user 服务的 `user.go`**：如何创建 kafka/redis 实例
- **svcCtx 创建**：把这些实例放进 ServiceContext，供 logic/task 使用

> 建议约定：
> - svcCtx 里统一字段命名（Redis、Kafka、TaskManager…），避免每个服务一个叫法。
> - 初始化失败要直接 fail-fast（启动失败），不要带病运行。

### 4.2 任务处理（Task Worker）
参考 **user 服务 `pkg/task/task.go`**：

- 创建一个 task 实例（worker）
- 在 task 里注册/绑定你要处理的 topic / job 类型
- 在 handler/logic 只负责投递，不写消费逻辑

### 4.3 投递任务（Producer）
参考 **user 服务 `/logic/test.go`**：

- 在 logic 内投递任务到 kafka（立即任务）或 redis（延迟任务的调度/占位）
- 投递成功后返回业务响应
- 任务真正执行在 task worker 内进行（异步）

### 4.4 延迟队列的语义（建议写清楚，避免误用）
由于 Redis 只是“延迟触发/时间轮”的角色，你们的真实消费处理最终仍会落到 Kafka（或某个可消费队列）上。建议在文档里明确：

- **Redis：负责延迟、重试调度、到期触发**
- **Kafka：负责可扩展的实际消费处理**

---

## 5) mac构建镜像
- gateway
```bash
docker buildx build \
  --platform linux/amd64 \
  -t crpi-u5azhs6neq326bz0.cn-hangzhou.personal.cr.aliyuncs.com/yub_lu/memo-nook-gateway:0.0.1 \
  -f service/gateway/Dockerfile \
  . \
  --push
```
- cat
```bash
docker buildx build \
  --platform linux/amd64 \
  -t crpi-u5azhs6neq326bz0.cn-hangzhou.personal.cr.aliyuncs.com/yub_lu/memo-nook-cat:0.0.1 \
  -f service/cat/Dockerfile \
  . \
  --push
```
-user
```bash
docker buildx build \
  --platform linux/amd64 \
  -t crpi-u5azhs6neq326bz0.cn-hangzhou.personal.cr.aliyuncs.com/yub_lu/memo-nook-user:0.0.1 \
  -f service/user/Dockerfile \
  . \
  --push
```

# TODO
- [ ] casbin 的引入
- [ ] 消息队列的err使用处理