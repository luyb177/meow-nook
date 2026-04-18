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