package mqttsn

import (
	"log"
	"sync"
	"time"
)

type TatCalculator struct {
	mutex      sync.RWMutex
	pubTimeMap map[uint16]time.Time
}

func NewTatCalculator() *TatCalculator {
	return &TatCalculator{sync.RWMutex{}, make(map[uint16]time.Time)}
}

func (tc *TatCalculator) RegistPublishTime(msgId uint16, t time.Time) {
	tc.mutex.Lock()
	defer tc.mutex.Unlock()
	if _, ok := tc.pubTimeMap[msgId]; ok {
		log.Println("WARN : Publish time for ", msgId, " is already exists.")
	}
	tc.pubTimeMap[msgId] = t
}

func (tc *TatCalculator) CalculateTatAndRemoveMsgId(msgId uint16, t time.Time) int64 {
	tc.mutex.Lock()
	defer tc.mutex.Unlock()
	if _, ok := tc.pubTimeMap[msgId]; !ok {
		log.Println("WARN : Publish time for ", msgId, " is not exists.")
		return 0
	}
	tat := t.Sub(tc.pubTimeMap[msgId]) / time.Millisecond
	delete(tc.pubTimeMap, msgId)
	return int64(tat)
}
