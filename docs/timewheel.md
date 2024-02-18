# TimeWheel 时间轮原理及单极时间轮的实现

## 为什么要学习时间轮

这里涉及到一个实际的场景，我们在做一个定时任务的系统，会在里面有一个很大数量级的定时任务，随着项目的发展，使用 go 语言调用操作系统的定时器会出现内存占用过多，性能下降的问题，因此需要一个更高性能的定时器来完成我们这个问题。

使用最经典的 go 语言定时器模块[robfig/cron](https://github.com/robfig/cron), 随着定时任务的增多，每插入新任务所需要的时间越来越大。

```go
func benchmarkCron(t *testing.B, n int) {
	e := cron.New()
	e.Start()
	defer e.Stop()
	for i := 0; i < n; i++ {
		e.AddFunc("@every 1s", func() {
			i = i
		})
	}

}
```

```bash
$ g go test -bench . -benchtime=30s -benchmem
goos: darwin
goarch: arm64
pkg: d_timewheel/timewheel/day1
BenchmarkCron1000-8                 8834            3888343 ns/op          536893 B/op      11781 allocs/op
BenchmarkCron10000-8                 120          299165856 ns/op         5789890 B/op     127336 allocs/op
BenchmarkCron20000-8                  30         1167397692 ns/op        12788149 B/op     282568 allocs/op
BenchmarkCron50000-8                   5         7119547075 ns/op        51400086 B/op    1230223 allocs/op
BenchmarkCron100000-8                  2        28832153354 ns/op       186321972 B/op    4700672 allocs/op
PASS
ok      d_timewheel/timewheel/day1      295.136s
```

可以从表现和源码分析，这个定时器的时间空间复杂度，都没法在一个数量级内收敛，随着业务越来越大，这个定时器的性能会越来越差。

但在我们预期的定时任务系统中，操作和存储的时间复杂度空间复杂度都希望尽可能收敛到 O(1) 的范围，以应对大量的定时任务，因此需要一个新的定时器实现方式，在现有的功能中获得更好的性能。

## 时间轮原理

> 通过阅读论文 [Hashed and Hierarchical Timing Wheels: Data Structures for the Efficient Implementation of a Timer Facility ](http://www.cs.columbia.edu/~nahum/w6998/papers/sosp87-timing-wheels.pdf) 来学习和了解时间轮算法。

在文章的第二章：`2 Model and Performance Measures` 中，作者就提出了时间轮的主要功能：

1. START_TIMER(Interval, Request_ID, Expiry_Action): 新增一个定时任务，包括到期的时间 interval，请求的 ID（用于区分其他的任务），到期的动作。
2. STOP_TIMER(Request_ID): 传入任务的 id，将该任务在定时服务端中删除
3. PER_TICK_BOOKKEEPING：计时器的时间单位
4. PER_TICK_PROCESSING：定时任务服务端的处理 Expiry_Action 的操作器

其中，1，2 是通过服务端对外的接口，由客户端来调用，3，4 是服务端通过操作系统的 ticker 来协助实现的。

对于定时器模块的性能，有两个评价的维度：空间&延迟，我们也会在实现完成后来进行一波性能测试，跟操作系统常见的优先队列实现的定时任务做一个性能测试对比。

## 时间轮的实现 milestone

1. 实现一个简单的单极时间轮，支持定时任务的添加和删除，此时的任务并不是一个循环的任务，只是一个单次执行的任务，传入的参数是一个时间点，第一个阶段的难点在于基于论文实现一个时间轮的基本雏型，并保留时间轮的可扩展性。
2. 实现一个时分秒级别的多级时间轮。第二阶段的难点是怎么基于单极时间轮扩展出多级时间轮，涉及到的编码习惯和设计模式的问题。
3. 实现一个年月日时分秒的多级时间轮。第三阶段的难点在于，年月日的变换多种多样，如何设计出来一个足够准确的时间轮。
4. 改变定时任务的数据结构，使其可以循环在时间轮中被单次、循环执行，设置过期时间，第四阶段的难点在于之前几个阶段的发展过程中如何保证可扩展性，学习的是阶段性的小型重构。

## 单极时间轮的实现

> 基于以上的论文分析，我们得知需要有一个时间轮的数据结构，一个定时任务的数据结构。

1. 时间轮类定义

```go
// TimeWheel 时间轮数据结构
type TimeWheel struct {
	once        sync.Once               // 单例工具，保证时间轮生命周期函数并发安全
	interval    time.Duration           // 时间轮的刻度间隔
	ticker      *time.Ticker            // 时间轮的刻度间隔触发器
	slots       []*list.List            // 时间轮的槽
	currentSlot int                     // 当前槽
	stopc       chan struct{}           // 停止时间轮的信号
	addTaskc    chan *taskElement       // 添加任务的信号
	removeTaskc chan int64              // 移除任务的信号, 任务ID
	key2ETask   map[int64]*list.Element // 任务key到任务的映射,用与快速删除查找链表节点
}
```

2. 任务类定义

> 由于我们只所有的任务都只触发一次，因此一个任务里面只需要记录在时间轮内的相对位置即可

```go
// taskElement 时间轮的任务元素
type taskElement struct {
	task  func() // 任务函数
	pos   int    // 任务在时间轮的槽的位置
	cycle int    // 任务的周期
	key   int64  // 任务的唯一标识
}
```
