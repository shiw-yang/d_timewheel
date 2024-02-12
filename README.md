# D-TimeWheel

## 1. 介绍

D-TimeWheel 是一个 go 语言实现的基于时间轮算法的分布式定时器, 用于处理定时任务, 同时支持 Quartz 的语法。

这个定时任务

## 2. 核心功能

- Quartz 语法注册定时任务
- 支持秒级别的定时任务
- 支持分布式部署
- 支持单次调度和循环调度

## 3. 功能演示

> TODO: 先用伪代码描述

- 定时中心启动，注册回调函数

```
err := d_time_wheel.ServiceStart(conf, func(data JobInfo){})
```

- 注册定时任务

```
jobInfo := &d_time_wheel.JobInfo{
    JobName: "job name",
    Cron: "0/5 * * * * ?",
    Desc: "job desc",
    Job:  "job", // 任务 这个似乎什么类型都可以
}
jobDto, err := d_time_wheel.RegisterJob(jobInfo)
```

- 取消定时任务

```
err := d_time_wheel.CancelJob(jobID)
```

- 查询定时任务信息

```
jobDto, err := d_time_wheel.QueryJob(jobID)
```

- 定时中心关闭

```
d_time_wheel.ServiceStop() // 优雅关闭
```

## 4. 部署依赖

使用 etcd 做分布式锁

## 5. 项目的目标

> 明确这个项目的目标才能更好的完成这个项目

1. 这个项目的主要目标是学习分布式定时任务中心的设计和实现。
2. 如果完成的好的话，可以整理成一个完整的文章或者教程，作为一个 go 语言项目

项目主要的 milestone 有：

1. 完成单机的时间轮算法，实现秒级别的定时任务系统
2. 引入语法解析层，实现 Quartz 语法等其他语法功能的引入
3. 引入 redis 或者 etcd，实现可分布式部署的定时任务系统
4. 引入关系型数据库 mysql 或者 postgresql，实现定时任务的持久化

在整个项目的过程中需要注意的东西有：

- 整个项目模块化设计，明确职责，降低耦合，接口化的方式实现模块之间的调用
- 每一个模块的能力边界，解决的问题，以及模块之间的交互
- 代码的可读性，可维护性，可测试性
- 项目的文档化，包括代码注释，项目文档，项目使用文档等

## 6. 项目的设计

> 按文件夹的目录结构将各个功能模块划分出来

```bash
.
├── cmd
├── conf
│   └── conf_sample.toml
├── config
│   └── config.go
├── distributed
│   └── lock.go
├── docs
├── errors
│   └── errors.go
├── job
├── parser
│   ├── parser.go
└── timewheel
    ├── timewheel.go
    └── timewheel_test.go
```

项目的主要划分如上：

- cmd: 项目的启动入口
- conf: 项目的配置文件
- config: 项目的配置文件读取
- distributed: 分布式锁的实现，lock.go 定义了标准的分布式锁接口
- docs: 项目的文档，所有的文档都在这里集合，如果有接口文档的话也会定义在这里
- errors: 项目统一封装的错误定义
- job: 任务本身的概念
- parser: 语法解析层，parser.go 定义了标准的语法解析接口
- timewheel: 时间轮算法的实现

每一个 impl，都会尽可能的有一个实现和测试两个文件

## 7. 项目的实现

每一个 milestone 都会有一个单独的分支，每一个实现完成的 milestone 都会合并到 master 分支，等到包括文档教程在内的所有东西都做完会打一个 release 版本，项目实现的过程尽可能的有详细的 commit 记录。

- 01-timewheel
- 02-parser
- 03-distributed
- 04-persistence

每一个 commit 都有规范：可以使用`better-commits`工具来完成 commit 提交，具体内容参考：[better-commits](https://github.com/Everduin94/better-commits)

```bash
# install
npm install -g better-commits
# usages
better-commits
```
