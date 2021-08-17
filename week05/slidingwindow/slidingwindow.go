package slidingwindow

import (
	"sync"
	"time"
)

// 窗口
type Window interface {
	// 开始边界
	Start() time.Time

	// 累计计数
	Count() int64

	// 累计计数增加n
	AddCount(n int64)

	// 重置开始边界和计数
	Reset(s time.Time, c int64)

	// 保持计数是最新的
	Sync(now time.Time)
}

// 停止窗口同步
type StopFunc func()

// 创建窗口并返回停止同步的func
type NewWindow func() (Window, StopFunc)

// 限制器
type Limiter struct {
	size  time.Duration // 一个窗口的持续时间
	limit int64         // 一个窗口内允许的最大请求数

	mu sync.Mutex

	curr Window
	prev Window
}

// 创建限制器并返回停止当前窗口内同步的func
func NewLimiter(size time.Duration, limit int64, newWindow NewWindow) (*Limiter, StopFunc) {
	currWin, currStop := newWindow()

	prevWin, _ := NewLocalWindow()

	lim := &Limiter{
		size:  size,
		limit: limit,
		curr:  currWin,
		prev:  prevWin,
	}

	return lim, currStop
}

// 返回窗口持续时间
func (lim *Limiter) Size() time.Duration {
	return lim.size
}

// 返回窗口内允许的最大请求数
func (lim *Limiter) Limit() int64 {
	lim.mu.Lock()
	defer lim.mu.Unlock()
	return lim.limit
}

// 重新设置限制器的最大请求数
func (lim *Limiter) SetLimit(newLimit int64) {
	lim.mu.Lock()
	defer lim.mu.Unlock()
	lim.limit = newLimit
}

// 发生1个请求
func (lim *Limiter) Allow() bool {
	return lim.AllowN(time.Now(), 1)
}

// 发生n个请求
func (lim *Limiter) AllowN(now time.Time, n int64) bool {
	lim.mu.Lock()
	defer lim.mu.Unlock()

	lim.advance(now)

	elapsed := now.Sub(lim.curr.Start())
	weight := float64(lim.size-elapsed) / float64(lim.size)
	count := int64(weight*float64(lim.prev.Count())) + lim.curr.Count()

	// 同步
	defer lim.curr.Sync(now)

	if count+n > lim.limit {
		return false
	}

	lim.curr.AddCount(n)
	return true
}

// 提前更新当前和上一个时间窗口
func (lim *Limiter) advance(now time.Time) {
	// 计算now的窗口起始时间
	newCurrStart := now.Truncate(lim.size)

	diffSize := newCurrStart.Sub(lim.curr.Start()) / lim.size
	if diffSize >= 1 {
		// 当前窗口至少比预期窗口大一个窗口大小。

		newPrevCount := int64(0)
		if diffSize == 1 {
			// 新的前一个窗口将与旧的当前窗口重叠
			newPrevCount = lim.curr.Count()
		}
		lim.prev.Reset(newCurrStart.Add(-lim.size), newPrevCount)

		lim.curr.Reset(newCurrStart, 0)
	}
}
