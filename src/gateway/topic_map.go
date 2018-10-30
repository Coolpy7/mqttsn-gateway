package mqttsn

import (
	"log"
	"sync"
)

type TopicMap struct {
	mutex   sync.RWMutex
	topics  map[uint16]string
	topicId *ManagedId
}

func NewTopicMap() *TopicMap {
	t := &TopicMap{
		mutex:   sync.RWMutex{},
		topics:  make(map[uint16]string),
		topicId: &ManagedId{}}

	// reserve special id
	t.topicId.EnsureId(0x0000)
	t.topicId.EnsureId(0xffff)
	return t
}

func (t *TopicMap) StoreTopicWithId(topicName string, id uint16) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	ok := t.topicId.EnsureId(id)
	if !ok {
		log.Println("ERROR : failed to ensure topic id.")
		return false
	}

	t.topics[id] = topicName
	return true
}

func (t *TopicMap) StoreTopic(topicName string) uint16 {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	for k, v := range t.topics {
		if v == topicName {
			return k
		}
	}

	// get next id
	id := t.topicId.NextId()

	t.topics[id] = topicName
	return id
}

func (t *TopicMap) LoadTopic(topicId uint16) (string, bool) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	topicName, ok := t.topics[topicId]
	return topicName, ok
}

func (t *TopicMap) LoadTopicId(topicName string) (uint16, bool) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	for k, v := range t.topics {
		if v == topicName {
			return k, true
		}
	}
	return 0, false
}
