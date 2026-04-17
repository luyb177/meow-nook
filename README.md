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
```bash
goctl rpc protoc pb/user.proto --proto_path=. --proto_path=..\..\third_party --go_out=. --go-grpc_out=. --zrpc_out=. --client=true
```