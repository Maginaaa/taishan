

# 泰山接口压测平台 🏔️

![Go Version](https://img.shields.io/badge/go-1.20-blue)
![Node Version](https://img.shields.io/badge/node-18.16.0-green)
![Vue Version](https://img.shields.io/badge/vue-3.x-brightgreen)

## 🌟 项目简介

泰山是一款高性能、易用的接口压测平台，致力于帮助开发者快速评估系统性能与稳定性。提供实时监控、压测报告生成、异常告警等功能，并可通过分布式部署实现高并发场景模拟。用户可通过可视化界面轻松配置压测任务，结合数据分析工具，精准定位性能瓶颈，为系统优化提供可靠依据。泰山压测平台适用于API、微服务及全链路压测，助力企业提升系统可靠性，保障业务平稳运行。

## 🛠️ 技术栈

| 领域        | 技术选型                                                              |
|-------------|-------------------------------------------------------------------|
| **前端**    | Node.js 18.16.0 + Vue 3 + TypeScript + Element Plus + ECharts     |
| **后端**    | Golang 1.20 + Gin + GORM Gen                                      |
| **存储**    | MySQL (业务数据) + Redis(缓存) + InfluxDB (施压数据) + Kafka + OSS(参数化文件存储) |
| **压测引擎**| 自研Go压测引擎                                                          |


## ⚡ 分布式能力
支持百万级并发压力

动态资源调度算法

多地域联合压测

## 🏗️ 系统架构

![img.png](readme-img/architecture.png)

## 📂 目录结构
```bash
泰山             
├── scene           # 场景管理主服务
├── engine          # 施压引擎
├── machine         # 机器管理(开源版仅施压机信息)
├── collector       # 压测数据聚合
├── report          # 压测报告
├── data            # 数据服务(开源版仅前置场景执行数据和采样数据)
├── fe              # 前端
```

## 🖼️ 核心功能界面预览
场景管理
![img.png](readme-img/scene-mng.png)
接口编辑
![img.png](readme-img/case-edit.png)
压测模式
![img.png](readme-img/press-type.png)
测试报告
![img.png](readme-img/report.png)

## 🚦 快速开始
1. 准备相应中间件
   1. 必须: mysql、redis、kafka、influxDB 
   2. 非必须: MongoDB (影响采样日志查询)、OSS(影响参数化文件使用)
2. 执行scene/sql文件内所有文件，生成对应mysql表
3. 配置各项目的conf.yml
4. 前端启动 node install && pnpm dev
5. 服务端启动，各项目跟目录go run main.go