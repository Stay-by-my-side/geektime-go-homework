package slidingwindow

import (
	"log"
	"time"
)

// 数据存储
type Datastore interface {
	Add(key string, start, delta int64) (int64, error)

	Get(key string, start int64) (int64, error)
}

type syncHelper struct {
	store        Datastore     // 数据存储
	syncInterval time.Duration // 同步的时间间隔

	inProgress bool      // 是否正在进行同步
	lastSynced time.Time // 最后同步时间
}

func newSyncHelper(store Datastore, syncInterval time.Duration) *syncHelper {
	return &syncHelper{store: store, syncInterval: syncInterval}
}

// 是否需要将数据同步到数据存储
func (h *syncHelper) IsTimeUp(now time.Time) bool {
	return !h.inProgress && now.Sub(h.lastSynced) >= h.syncInterval
}

func (h *syncHelper) InProgress() bool {
	return h.inProgress
}

func (h *syncHelper) Begin(now time.Time) {
	h.inProgress = true
	h.lastSynced = now
}

func (h *syncHelper) End() {
	h.inProgress = false
}

func (h *syncHelper) Sync(req SyncRequest) (resp SyncResponse, err error) {
	var newCount int64

	if req.Changes > 0 {
		newCount, err = h.store.Add(req.Key, req.Start, req.Changes)
	} else {
		newCount, err = h.store.Get(req.Key, req.Start)
	}

	if err != nil {
		return SyncResponse{}, err
	}

	return SyncResponse{
		OK:           true,
		Start:        req.Start,
		Changes:      req.Changes,
		OtherChanges: newCount - req.Count,
	}, nil
}

// 非阻塞模式下进行同步
type NonblockingSynchronizer struct {
	reqC  chan SyncRequest  // 请求通道
	respC chan SyncResponse // 响应通道

	stopC chan struct{} // 停止
	exitC chan struct{} // 退出

	helper *syncHelper
}

func NewNonblockingSynchronizer(store Datastore, syncInterval time.Duration) *NonblockingSynchronizer {
	return &NonblockingSynchronizer{
		reqC:   make(chan SyncRequest),
		respC:  make(chan SyncResponse),
		stopC:  make(chan struct{}),
		exitC:  make(chan struct{}),
		helper: newSyncHelper(store, syncInterval),
	}
}

func (s *NonblockingSynchronizer) Start() {
	go s.syncLoop()
}

func (s *NonblockingSynchronizer) Stop() {
	close(s.stopC)
	<-s.exitC
}

func (s *NonblockingSynchronizer) syncLoop() {
	for {
		select {
		case req := <-s.reqC:
			// 处理请求
			resp, err := s.helper.Sync(req)
			if err != nil {
				log.Printf("err: %v\n", err)
			}

			select {
			case s.respC <- resp:
			case <-s.stopC:
				goto exit
			}
		case <-s.stopC:
			goto exit
		}
	}

exit:
	close(s.exitC)
}

func (s *NonblockingSynchronizer) Sync(now time.Time, makeReq MakeFunc, handleResp HandleFunc) {
	if s.helper.IsTimeUp(now) {
		select {
		case s.reqC <- makeReq():
			s.helper.Begin(now)
		default:
		}
	}

	if s.helper.InProgress() {
		select {
		case resp := <-s.respC:
			handleResp(resp)
			s.helper.End()
		default:
		}
	}
}
