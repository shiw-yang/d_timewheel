package timewheel

import (
	"container/list"
	"sync"
	"time"
)

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

// taskElement 时间轮的任务元素
type taskElement struct {
	task  func() // 任务函数
	pos   int    // 任务在时间轮的槽的位置
	cycle int    // 任务的周期
	key   int64  // 任务的唯一标识
}

// NewTimeWheel 初始化时间轮，interval表示每个槽位的时间单位，默认1s；slotNums表示槽位数量；默认60s
func NewTimeWheel(interval time.Duration, slotNums int) *TimeWheel {
	// 判断并设置默认值
	if interval <= 0 {
		interval = time.Second
	}
	if slotNums <= 0 {
		slotNums = 60
	}

	// 初始化时间轮
	tw := &TimeWheel{
		interval:    interval,
		ticker:      time.NewTicker(interval),
		slots:       make([]*list.List, slotNums),
		stopc:       make(chan struct{}),
		key2ETask:   make(map[int64]*list.Element),
		addTaskc:    make(chan *taskElement),
		removeTaskc: make(chan int64),
	}
	// 初始化槽位
	for i := range tw.slots {
		tw.slots[i] = list.New()
	}

	// 异步goroutine启动时间轮
	go tw.run()
	return tw
}

// AddTask 添加任务到时间轮,taskF为定时任务函数，time为任务执行时间
func (tw *TimeWheel) AddTask(key int64, taskF func(), time time.Time) {
	// 计算任务的位置
	pos, cycle := tw.getPosAndCycle(time)
	// 创建任务
	task := &taskElement{
		task:  taskF,
		pos:   pos,
		cycle: cycle,
		key:   key,
	}
	// 发送添加任务的信号
	tw.addTaskc <- task
}

// RemoveTask 移除时间轮中的任务
func (tw *TimeWheel) RemoveTask(key int64) {
	tw.removeTaskc <- key
}

// getPosAndCycle 根据入参t计算任务的位置和周期
func (tw *TimeWheel) getPosAndCycle(t time.Time) (int, int) {
	delay := int(time.Until(t))

	// 定时任务从属的环状数组 pos
	pos := (tw.currentSlot + delay/int(tw.interval)) % len(tw.slots)
	// 定时触发的延迟轮次
	cycle := delay / (int(tw.interval) * len(tw.slots))

	return pos, cycle
}

// run 启动时间轮
func (tw *TimeWheel) run() {
	// 捕获时间轮生命周期的panic
	// 防止因执行函数错误而破坏时间轮的生命周期
	defer func() {
		if err := recover(); err != nil {
			print(err)
		}
	}()

	// 启动时间轮,for select结构监听时间轮生命周期功能
	for {
		select {
		case <-tw.stopc: // 停止时间轮的信号
			return
		case <-tw.ticker.C: // 时间轮的刻度间隔触发器
			tw.tickHandler()
		case task := <-tw.addTaskc: // 添加任务的信号
			tw.addTaskHandler(task)
		case key := <-tw.removeTaskc: // 移除任务的信号
			tw.removeTaskHandler(key)
		}
	}
}

// tickHandler 时间轮的刻度间隔触发器的处理函数
func (tw *TimeWheel) tickHandler() {
	list := tw.slots[tw.currentSlot]
	defer tw.circularIncr()
	tw.execute(list)
}

// addTaskHandler 添加任务的处理函数
func (tw *TimeWheel) addTaskHandler(task *taskElement) {
	taskList := tw.slots[task.pos]
	// 判断任务是否已经存在,存在删除后更新
	if _, ok := tw.key2ETask[task.key]; ok {
		tw.removeTaskHandler(task.key)
	}

	// 添加任务到时间轮
	ele := taskList.PushBack(task)
	tw.key2ETask[task.key] = ele
}

// removeTaskHandler 移除时间轮中的任务
func (tw *TimeWheel) removeTaskHandler(key int64) {
	// 根据key找到任务的链表节点内容
	ele, ok := tw.key2ETask[key]
	if !ok {
		return
	}

	// 从链表节点复原出任务信息
	task := ele.Value.(*taskElement)

	// 删除在key2ETask中的映射和时间轮内的任务
	delete(tw.key2ETask, task.key)
	tw.slots[task.pos].Remove(ele)
}

// circularIncr 时间轮的槽的循环增加
func (tw *TimeWheel) circularIncr() {
	tw.currentSlot = (tw.currentSlot + 1) % len(tw.slots)
}

func (tw *TimeWheel) execute(taskList *list.List) {
	for ele := taskList.Front(); ele != nil; {
		// 获取节点任务信息
		taskEle := ele.Value.(*taskElement)
		// 判断任务是否存在延迟周期
		if taskEle.cycle > 0 {
			taskEle.cycle--
			ele = ele.Next()
			continue
		}

		// 达到执行条件，开始执行
		go func() {
			// 捕获任务函数带来的panic
			defer func() {
				if err := recover(); err != nil {
					print(err)
				}
			}()
			taskEle.task()
		}()

		// 执行任务下达，删除任务
		next := ele.Next()
		taskList.Remove(ele)
		delete(tw.key2ETask, taskEle.key)
		ele = next
	}
}

func (tw *TimeWheel) Stop() {
	// 保证并发安全
	tw.once.Do(func() {
		// 关闭时间轮的触发器
		tw.ticker.Stop()
		// 关闭时间轮的信号
		close(tw.stopc)
	})
}
