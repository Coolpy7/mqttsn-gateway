package mqttsn

import (
	"log"
	"sync/atomic"
	"time"
)

type StatisticsReporter struct {
	sendPublish  uint64
	recvPublish  uint64
	pubComp      uint64
	sessionCount uint64
	tatMax       int64
	tatMin       int64
	tatTotal     int64
	interval     time.Duration
}

func NewStatisticsReporter(interval time.Duration) *StatisticsReporter {
	r := &StatisticsReporter{0, 0, 0, 0, 0, 65535, 0, interval}
	return r
}

func (r *StatisticsReporter) loggingLoop() {
	t := time.NewTicker(r.interval * time.Second)
	for {
		select {
		case <-t.C:
			log.Println("[STATISTICS-REPORT] SESSION-COUNT : ",
				atomic.LoadUint64(&r.sessionCount),
				", PUBLISH-RECV-COUNT : ",
				atomic.LoadUint64(&r.recvPublish),
				", PUBLISH-SEND-COUNT : ",
				atomic.LoadUint64(&r.sendPublish),
				", PUBLISH-COMP-COUNT : ",
				atomic.LoadUint64(&r.pubComp))

			avg := int64(0)
			if atomic.LoadUint64(&r.pubComp) > 0 {
				avg = atomic.LoadInt64(&r.tatTotal) /
					int64(atomic.LoadUint64(&r.pubComp))
			}
			log.Println("[STATISTICS-REPORT] PUBLISH-TAT-AVG : ",
				avg,
				", PUBLISH-TAT-MIN : ",
				atomic.LoadInt64(&r.tatMin),
				", PUBLISH-TAT-MAX : ",
				atomic.LoadInt64(&r.tatMax))

			// init counter
			atomic.StoreUint64(&r.sendPublish, 0)
			atomic.StoreUint64(&r.recvPublish, 0)
			atomic.StoreUint64(&r.pubComp, 0)
			atomic.StoreInt64(&r.tatTotal, 0)
			atomic.StoreInt64(&r.tatMax, 0)
			atomic.StoreInt64(&r.tatMin, 65535)
		}
	}
}

func (r *StatisticsReporter) countUpSendPublish() {
	atomic.AddUint64(&r.sendPublish, 1)
}

func (r *StatisticsReporter) countUpRecvPublish() {
	atomic.AddUint64(&r.recvPublish, 1)
}

func (r *StatisticsReporter) countUpPubComp() {
	atomic.AddUint64(&r.pubComp, 1)
}

func (r *StatisticsReporter) storeSessionCount(c uint64) {
	atomic.StoreUint64(&r.sessionCount, c)
}

func (r *StatisticsReporter) storeTat(t int64) {
	atomic.AddInt64(&r.tatTotal, t)

	max := atomic.LoadInt64(&r.tatMax)
	if t > max {
		atomic.StoreInt64(&r.tatMax, t)
	}

	min := atomic.LoadInt64(&r.tatMin)
	if t < min {
		atomic.StoreInt64(&r.tatMin, t)
	}
}
