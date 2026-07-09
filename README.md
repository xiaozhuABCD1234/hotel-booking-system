# 上海电力大学数据库原理课程设计 — 酒店预订管理系统

选题四：酒店预订管理系统。本系统模拟一个第三方酒店预订平台（如携程），支持用户注册、酒店搜索、客房预订、订单管理、评价及数据统计等功能。

## 技术栈

| 层级         | 技术                                                                                |
| ------------ | ----------------------------------------------------------------------------------- |
| **后端**     | Go + [Fiber v3](https://github.com/gofiber/fiber)                                   |
| **前端**     | Vue 3 + TypeScript + [Vite](https://vite.dev/)                                      |
| **UI 框架**  | [Tailwind CSS v4](https://tailwindcss.com/) + [Reka UI](https://reka-ui.com/)       |
| **状态管理** | [Pinia](https://pinia.vuejs.org/)                                                   |
| **路由**     | [Vue Router](https://router.vuejs.org/)                                             |
| **数据库**   | [PostgreSQL 18](https://www.postgresql.org/about/news/postgresql-18-released-3142/) |
| **ORM**      | [GORM](https://gorm.io/)                                                            |
| **API 文档** | [Swaggo / Swagger](https://github.com/swaggo/swag)                                  |
| **认证**     | JWT ([golang-jwt](https://github.com/golang-jwt/jwt))                               |
| **图标**     | [Lucide](https://lucide.dev/)                                                       |

## 项目结构

```
├── backend/          # Go 后端 (Fiber + GORM)
│   ├── cmd/          # 入口
│   ├── config/       # 配置
│   ├── database/     # 数据库连接 & 迁移
│   ├── handler/      # 请求处理器
│   ├── middleware/    # 中间件 (JWT, CORS 等)
│   ├── model/        # 数据模型
│   ├── repo/         # 数据访问层
│   ├── router/       # 路由注册
│   ├── service/      # 业务逻辑层
│   ├── auth/         # 认证相关
│   ├── test/         # 测试
│   └── docs/         # Swagger 文档
├── frontend/         # Vue 3 前端
│   ├── src/
│   │   ├── api/      # API 调用
│   │   ├── components/  # 通用组件
│   │   ├── composables/ # 组合式函数
│   │   ├── layouts/     # 布局组件
│   │   ├── lib/         # 工具库
│   │   ├── router/      # 前端路由
│   │   ├── stores/      # Pinia 状态
│   │   ├── types/       # TypeScript 类型
│   │   └── views/       # 页面视图
│   └── ...
├── sql/              # 数据库脚本
│   ├── table.sql     # 表、类型、索引
│   ├── trigger.sql   # 触发器
│   ├── function.sql  # 存储函数/过程
│   ├── view.sql      # 视图
│   ├── insert.sql    # 测试数据
│   └── select.sql    # 查询示例
└── doc/              # 课程设计文档
    ├── 需求分析.md
    ├── 数据库设计.md
    └── 图.drawio
```

## 功能需求

1. **用户管理** — 注册、登录、信息维护；积分体系与 VIP 升级及优惠策略
2. **客房预订** — 支持为自己/他人预订，一次多间，确认入住人信息
3. **酒店搜索** — 按城市、区域、酒店名、价格范围、入住日期等条件筛选
4. **房型选择** — 进入酒店后展示多种房型及价格供用户选择
5. **订单管理** — 查询全部订单，支持按状态（已预订/已入住等）分类查询
6. **用户评价** — 对入住信息进行评价
7. **退订** — 已预订未入住的订单可退订
8. **后台维护** — 系统后台信息（如客房情况）的维护
9. **汇总统计** —
    - 按城市、区域、酒店、价格等因素统计客房预订情况
    - 按入住人信息（年龄、性别、职业、学历、收入等）统计预订情况，分析客户偏好

## 数据库设计

按照数据库设计规范完成，包含：

- 各表主键、外键约束
- 默认约束（如性别、入住时间等）
- 非空约束（如客房类型名）
- CHECK 约束（如离店时间 ≥ 入住时间）
- 规则约束（如身份证号 18 位）
- 适当设计索引

## 数据库对象

- **自定义类型** — ENUM（`user_role`、`order_status`）
- **视图** — 按城市、区域、酒店查看所有客房详细信息
- **触发器** — 客房被预订时更新同类型可预订数量
- **存储过程** — 按用户信息查询订单；统计一段时间内各类客房入住率
- **JSON 字段** — 入住人信息、评价内容以 JSON 存储

## 快速开始

### 前置要求

- Go 1.26+
- Node.js 20+
- pnpm
- PostgreSQL 18+

### 后端

```bash
cd backend
cp .env.example .env    # 编辑数据库连接信息
go run cmd/main.go
```

### 前端

```bash
cd frontend
pnpm install
pnpm dev
```

### 数据库

```bash
psql -U <user> -d <dbname> -f sql/table.sql
psql -U <user> -d <dbname> -f sql/trigger.sql
psql -U <user> -d <dbname> -f sql/function.sql
psql -U <user> -d <dbname> -f sql/view.sql
psql -U <user> -d <dbname> -f sql/insert.sql
```

## API 文档

启动后端后访问 Swagger UI：

```
http://localhost:<port>/swagger/index.html
```
