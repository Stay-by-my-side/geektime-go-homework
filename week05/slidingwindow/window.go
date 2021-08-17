package slidingwindow

import (
	"time"
)

type LocalWindow struct {
	// 窗口开始边界（纳秒）
	// [start, start + size)
	start int64

	// 窗口中的请求总数
	count int64
}

func NewLocalWindow() (*LocalWindow, StopFunc) {
	return &LocalWindow{}, func() {}
}

func (w *LocalWindow) Start() time.Time {
	return time.Unix(0, w.start)
}

func (w *LocalWindow) Count() int64 {
	return w.count
}

func (w *LocalWindow) AddCount(n int64) {
	w.count += n
}

func (w *LocalWindow) Reset(s time.Time, c int64) {
	w.start = s.UnixNano()
	w.count = c
}

func (w *LocalWindow) Sync(now time.Time) {}

type (
	// 非阻塞模式下，使用此请求来发起操作
	SyncRequest struct {
		Key     string
		Start   int64
		Count   int64
		Changes int64
	}

	// 非阻塞模式下的请求操作结果
	SyncResponse struct {
		// 同步是否成功
		OK    bool
		Start int64
		// 本地累计计数
		Changes int64
		// 同步结果与本地累计计数的偏差
		OtherChanges int64
	}

	MakeFunc   func() SyncRequest
	HandleFunc func(SyncResponse)
)

type Synchronizer interface {
	// 开启同步goroutine
	Start()

	// 停止同步goroutine并等待退出
	Stop()

	// 发送同步请求
	Sync(time.Time, MakeFunc, HandleFunc)
}

type SyncWindow struct {
	LocalWindow
	changes int64

	key    string
	syncer Synchronizer
}

func NewSyncWindow(key string, syncer Synchronizer) (*SyncWindow, StopFunc) {
	w := &SyncWindow{
		key:    key,
		syncer: syncer,
	}

	w.syncer.Start()
	return w, w.syncer.Stop
}

func (w *SyncWindow) AddCount(n int64) {
	w.changes += n
	w.LocalWindow.AddCount(n)
}

func (w *SyncWindow) Reset(s time.Time, c int64) {
	w.changes = 0

	w.LocalWindow.Reset(s, c)
}

func (w *SyncWindow) makeSyncRequest() SyncRequest {
	return SyncRequest{
		Key:     w.key,
		Start:   w.LocalWindow.start,
		Count:   w.LocalWindow.count,
		Changes: w.changes,
	}
}

func (w *SyncWindow) handleSyncResponse(resp SyncResponse) {
	if resp.OK && resp.Start == w.LocalWindow.start {
		w.LocalWindow.count += resp.OtherChanges

		// 减去已同步的数量
		w.changes -= resp.Changes
	}
}

func (w *SyncWindow) Sync(now time.Time) {
	w.syncer.Sync(now, w.makeSyncRequest, w.handleSyncResponse)
}
