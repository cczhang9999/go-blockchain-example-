# 区块链API服务器

这是一个基于Go语言开发的区块链API服务器，提供了钱包管理、转账和区块链信息查询等功能。

## 功能特性

- 🏦 **钱包管理**: 创建新钱包，查询余额
- 💸 **转账功能**: 支持钱包之间的转账操作
- ⛓️ **区块链信息**: 查看区块链状态和区块信息
- 🔒 **安全性**: 使用以太坊标准的加密算法
- 🚀 **高性能**: 基于Gin框架的高性能Web服务器
- 📊 **数据库支持**: 使用MySQL存储区块链数据

## 技术栈

- **后端框架**: Gin
- **数据库**: MySQL
- **区块链**: 自定义区块链实现
- **加密**: 以太坊加密库
- **ORM**: GORM

## 快速开始

### 1. 环境要求

- Go 1.24+
- MySQL 8.0+

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 配置数据库

确保MySQL服务正在运行，并创建相应的数据库和表。

### 4. 运行服务器

```bash
# 编译
go build -o blockchain-server

# 运行
./blockchain-server
```

服务器将在 `http://localhost:8080` 启动。

## API接口

### 基础信息

- **基础URL**: `http://localhost:8080`
- **API版本**: `v1`
- **API前缀**: `/api/v1`

### 接口列表

#### 1. 获取API信息
```
GET /
```
返回所有可用的API接口信息。

#### 2. 健康检查
```
GET /api/v1/health
```
检查服务器状态。

#### 3. 创建钱包
```
POST /api/v1/wallet
```

**响应示例**:
```json
{
  "success": true,
  "data": {
    "address": "0xB41F21D2999b3D2a6ad6e26D2864539Dd9186F0F",
    "private_key": "a93e264f9d1d64f81477e9d6798c7a432687d8fdb1b1860c8a90faa17e996580"
  }
}
```

#### 4. 查询余额
```
GET /api/v1/wallet/:address
```

**参数**:
- `address`: 钱包地址

**响应示例**:
```json
{
  "success": true,
  "data": {
    "address": "0xB41F21D2999b3D2a6ad6e26D2864539Dd9186F0F",
    "balance": 100.5
  }
}
```

#### 5. 转账
```
POST /api/v1/transfer
```

**请求体**:
```json
{
  "from_address": "0xB41F21D2999b3D2a6ad6e26D2864539Dd9186F0F",
  "to_address": "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
  "amount": 50.0
}
```

**响应示例**:
```json
{
  "success": true,
  "message": "Transfer completed successfully",
  "data": {
    "from_address": "0xB41F21D2999b3D2a6ad6e26D2864539Dd9186F0F",
    "to_address": "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
    "amount": 50.0
  }
}
```

#### 6. 获取区块链信息
```
GET /api/v1/blockchain
```

**响应示例**:
```json
{
  "success": true,
  "data": {
    "is_valid": true,
    "blocks": [...],
    "count": 6
  }
}
```

## 项目结构

```
hello-go/
├── main.go                 # 主程序入口
├── handlers/
│   └── api.go             # API处理函数
├── blockchain/
│   └── chain.go           # 区块链核心逻辑
├── models/
│   └── block.go           # 数据模型
├── database/
│   ├── mysql.go           # 数据库连接
│   └── blockchain_mysql.go # 区块链数据访问层
├── config/
│   └── config.go          # 配置管理
└── README.md              # 项目说明
```

## 开发说明

### 代码优化

1. **单例模式**: 使用单例模式管理区块链实例，避免重复初始化
2. **错误处理**: 统一的错误处理和响应格式
3. **中间件**: 添加了CORS支持和日志记录
4. **API版本化**: 使用版本化的API路径
5. **参数验证**: 使用Gin的绑定验证请求参数

### 安全特性

- 使用以太坊标准的加密算法
- 参数验证和输入检查
- 统一的错误处理，避免信息泄露

## 测试

使用curl测试API接口：

```bash
# 创建钱包
curl -X POST http://localhost:8080/api/v1/wallet

# 查询余额
curl -X GET "http://localhost:8080/api/v1/wallet/0xB41F21D2999b3D2a6ad6e26D2864539Dd9186F0F"

# 转账
curl -X POST http://localhost:8080/api/v1/transfer \
  -H "Content-Type: application/json" \
  -d '{"from_address":"0xB41F21D2999b3D2a6ad6e26D2864539Dd9186F0F","to_address":"0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6","amount":50}'

# 获取区块链信息
curl -X GET http://localhost:8080/api/v1/blockchain
```

## 许可证

MIT License
