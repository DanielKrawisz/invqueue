package invqueue

import (
	"github.com/monetas/bmutil/wire"
)

func (q *InvQueue) TstCheckIntegrity() int {
	result := 0
	if len(q.objMap) != q.objList.Len() {
		result &= 1
	}

	testMap := make(map[wire.InvVect]struct{})
	for e := q.objList.Front(); e != nil; e = e.Next() {
		iv := e.Value.(*wire.InvVect)
		if _, ok := testMap[*iv]; ok {
			result &= 2
		}
		if _, ok := q.objMap[*iv]; !ok {
			result &= 4
		}
		testMap[*iv] = struct{}{}
	}

	for iv, _ := range q.objMap {
		if _, ok := testMap[iv]; !ok {
			result &= 8
		}
	}

	return result
}

// A special function to break an InvQueue and put it in an invalid state
// so as to test a line of code that shouldn't ever actually happen in real
// life.
func (q *InvQueue) TstBreakRemove(iv *wire.InvVect) bool {
	if _, ok := q.objMap[*iv]; !ok {
		return false
	}
	for e := q.objList.Front(); e != nil; e = e.Next() {
		if *e.Value.(*wire.InvVect) == *iv {
			q.objList.Remove(e)
			return true
		}
	}
	return true
}
