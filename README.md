# 区块链API服务器

这是一个基于Go语言开发的区块链API服务器，提供了钱包管理、转账、交易记录查询和区块链信息查询等功能。

## 功能特性

- 🏦 **钱包管理**: 创建新钱包，查询余额
- 💸 **转账功能**: 支持钱包之间的转账操作
- 📊 **交易记录**: 完整的交易历史查询功能
- ⛓️ **区块链信息**: 查看区块链状态和区块信息
- 🔒 **安全性**: 使用以太坊标准的加密算法
- 🚀 **高性能**: 基于Gin框架的高性能Web服务器
- 📊 **数据库支持**: 使用MySQL存储区块链数据
- 📄 **分页支持**: 交易记录支持分页查询

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
  "message": "Wallet created successfully",
  "data": {
    "address": "0xfc33F29F4023E2B59B75BdbAaB27F87a3f7521D1",
    "private_key": "67249bfe32fbd2b5d96b3f7dea066913b85a833d871d2742f87172b023899070",
    "balance": 0
  },
  "timestamp": "2025-07-06T13:29:41.703465+08:00"
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
  "message": "Balance retrieved successfully",
  "data": {
    "address": "0xfc33F29F4023E2B59B75BdbAaB27F87a3f7521D1",
    "balance": 100.5
  },
  "timestamp": "2025-07-06T13:29:41.703465+08:00"
}
```

#### 5. 转账
```
POST /api/v1/transfer
```

**请求体**:
```json
{
  "from_address": "0xfc33F29F4023E2B59B75BdbAaB27F87a3f7521D1",
  "to_address": "0xd6e1EFbe8C8eE752a4B371D1e59D4a735d075557",
  "amount": 50.0
}
```

**响应示例**:
```json
{
  "success": true,
  "message": "Transfer completed successfully",
  "data": {
    "from_address": "0xfc33F29F4023E2B59B75BdbAaB27F87a3f7521D1",
    "to_address": "0xd6e1EFbe8C8eE752a4B371D1e59D4a735d075557",
    "amount": 50.0,
    "timestamp": "2025-07-06T13:29:41.703465+08:00"
  },
  "timestamp": "2025-07-06T13:29:41.703465+08:00"
}
```

#### 6. 获取所有交易记录
```
GET /api/v1/transactions
```

**查询参数**:
- `page`: 页码（默认1）
- `limit`: 每页数量（默认20，最大100）

**响应示例**:
```json
{
  "success": true,
  "message": "All transactions retrieved successfully",
  "data": {
    "transactions": [
      {
        "id": 5,
        "from_address": "0x9b71ee886C2f82AeF96F58448a6E1A1734b50437",
        "to_address": "0xd6e1EFbe8C8eE752a4B371D1e59D4a735d075557",
        "amount": 50,
        "timestamp": "2025-07-05T14:38:14+08:00"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 5,
      "total_pages": 1
    }
  },
  "timestamp": "2025-07-06T13:29:56.732163+08:00"
}
```

#### 7. 获取地址交易历史
```
GET /api/v1/transactions/history/:address
```

**参数**:
- `address`: 钱包地址

**响应示例**:
```json
{
  "success": true,
  "message": "Transaction history retrieved successfully",
  "data": {
    "address": "0x9b71ee886C2f82AeF96F58448a6E1A1734b50437",
    "transactions": [
      {
        "id": 5,
        "from_address": "0x9b71ee886C2f82AeF96F58448a6E1A1734b50437",
        "to_address": "0xd6e1EFbe8C8eE752a4B371D1e59D4a735d075557",
        "amount": 50,
        "timestamp": "2025-07-05T14:38:14+08:00",
        "transaction_type": "sent"
      }
    ],
    "total_count": 1
  },
  "timestamp": "2025-07-06T13:29:28.146415+08:00"
}
```

#### 8. 获取区块交易
```
GET /api/v1/transactions/block/:block_id
```

**参数**:
- `block_id`: 区块ID

#### 9. 获取区块链信息
```
GET /api/v1/blockchain
```

**响应示例**:
```json
{
  "success": true,
  "message": "Blockchain information retrieved successfully",
  "data": {
    "is_valid": true,
    "blocks": [...],
    "block_count": 6,
    "last_updated": "2025-07-06T13:29:56.732163+08:00"
  },
  "timestamp": "2025-07-06T13:29:56.732163+08:00"
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
├── format_json.py         # JSON格式化工具
└── README.md              # 项目说明
```

## 开发说明

### 代码优化

1. **单例模式**: 使用单例模式管理区块链实例，避免重复初始化
2. **错误处理**: 统一的错误处理和响应格式
3. **中间件**: 添加了CORS支持和日志记录
4. **API版本化**: 使用版本化的API路径
5. **参数验证**: 使用Gin的绑定验证请求参数
6. **分页支持**: 交易记录支持分页查询
7. **交易类型**: 区分发送和接收的交易类型

### 安全特性

- 使用以太坊标准的加密算法
- 参数验证和输入检查
- 统一的错误处理，避免信息泄露

## 测试

### 使用curl测试API接口：

```bash
# 创建钱包
curl -X POST http://localhost:8080/api/v1/wallet | python3 format_json.py

# 查询余额
curl -X GET "http://localhost:8080/api/v1/wallet/0xfc33F29F4023E2B59B75BdbAaB27F87a3f7521D1" | python3 format_json.py

# 转账
curl -X POST http://localhost:8080/api/v1/transfer \
  -H "Content-Type: application/json" \
  -d '{"from_address":"0xfc33F29F4023E2B59B75BdbAaB27F87a3f7521D1","to_address":"0xd6e1EFbe8C8eE752a4B371D1e59D4a735d075557","amount":50}' | python3 format_json.py

# 获取所有交易记录
curl -X GET "http://localhost:8080/api/v1/transactions" | python3 format_json.py

# 获取地址交易历史
curl -X GET "http://localhost:8080/api/v1/transactions/history/0xfc33F29F4023E2B59B75BdbAaB27F87a3f7521D1" | python3 format_json.py

# 分页查询交易记录
curl -X GET "http://localhost:8080/api/v1/transactions?page=1&limit=3" | python3 format_json.py

# 获取区块链信息
curl -X GET http://localhost:8080/api/v1/blockchain | python3 format_json.py
```

### JSON格式化

项目提供了 `format_json.py` 工具来格式化JSON输出：

```bash
# 使用Python格式化
curl -X GET "http://localhost:8080/api/v1/transactions" | python3 format_json.py

# 或使用Python内置工具
curl -X GET "http://localhost:8080/api/v1/transactions" | python3 -m json.tool
```

## 交易记录功能

### 主要特性

1. **完整交易历史**: 记录所有转账交易
2. **地址查询**: 按地址查询交易历史
3. **交易类型**: 区分发送和接收的交易
4. **分页支持**: 支持分页查询大量交易记录
5. **时间戳**: 记录详细的交易时间
6. **区块关联**: 交易与区块的关联关系

### 数据库表结构

```sql
-- 交易记录表
CREATE TABLE transactions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    block_id BIGINT,
    from_addr VARCHAR(42),
    to_addr VARCHAR(42),
    amount DECIMAL(20,8),
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_from_addr (from_addr),
    INDEX idx_to_addr (to_addr),
    INDEX idx_timestamp (timestamp)
);
```

## 许可证

MIT License
