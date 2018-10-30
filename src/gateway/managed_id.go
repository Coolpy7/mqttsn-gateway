package mqttsn

import (
	"sync"
)

type ManagedId struct {
	mutex sync.RWMutex
	id    [0xffff + 1]bool
}

func (t *ManagedId) EnsureId(id uint16) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if !t.id[id] {
		t.id[id] = true
		return true
	}

	return false
}

func (t *ManagedId) NextId() uint16 {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	for i, v := range t.id {
		if v == false {
			t.id[i] = true
			return uint16(i)
		}
	}

	// TODO: error handling

	return 0
}

func (t *ManagedId) FreeId(i uint16) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.id[i] = false

	return
}
