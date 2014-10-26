package server

// TODO expand functionality.
import (
	"container/list"
)

// Keeps a pool of integers.
func NewIdPool(totNum uint32) *IdPool {
	idPool := &IdPool{}
	idPool.init(totNum)
	return idPool
}

type IdPool struct {
	ids *list.List
}

// We start from 1, so that 0 is not used as an id.
func (idp *IdPool) init(totNum uint32) {
	idp.ids = list.New()
	for i := uint32(1); i <= totNum; i++ {
		idp.ids.PushBack(i)
	}
}

func (idp *IdPool) GetId() uint32 {
	val := idp.ids.Front()
	idp.ids.Remove(val)
	num, _ := val.Value.(uint32)
	return num
}

func (idp *IdPool) ReturnId(id uint32) {
	idp.ids.PushBack(id)
}