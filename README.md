# 目录
```
magicmirror
│   README.md   -- 本文件
│   接口.md     -- 接口说明
│
└───MagicMirrorServer
    │   MainServer.go   -- 服务端主程序
    │   upload.gtpl     -- 上传文件用的表单
    │
    └───Data
    |   │   Comment.go  -- 评论
    |   │   Picture.go  -- 图片处理
    |   |   Post.go     -- 帖子
    |   |   Relation.go -- 用户关联
    |   |   User.go     -- 用户
    |
    └───uploaded        -- 存放上传的图片
    |   ...
```
## MainServer.go
需要以下全部：
- 已经安装 mysql 数据库以及 beego 相关驱动
- 相关的常量配置无误：数据库信息、服务端地址等
- 运行时对应的数据库存在（数据表不存在时会自行创建）
