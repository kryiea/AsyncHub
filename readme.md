# 异步任务调度框架 AsyncHub

## 简介
- flowsvr：任务流服务，对外提供任务处理、查询接口
- worker：处理某种/多种任务，其中集成了tasksdk，提供自动调度，部署在客户端


## 编译&启动
1. 提前安装MySQL
2. 在MySQL中执行create.sql脚本，创建相关表
3. 修改flowsvr的配置文件：包括服务地址，服务治理参数
 