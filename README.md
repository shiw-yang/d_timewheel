# D-TimeWheel

## 1. 介绍

D-TimeWheel是一个go语言实现的基于时间轮算法的分布式定时器, 用于处理定时任务, 同时支持Quartz的语法。

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

使用etcd做分布式锁

