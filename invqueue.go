package invqueue

import (
	"container/list"

	"github.com/monetas/bmutil/wire"
)

// Keep track of the unknown objects in both a map and a list so that
// we can iterate over the objects in the order they were received and
// check quickly whether the object is known about. Can't use a
// MruInventoryMap without redesigning it.
type InvQueue struct {
	objList *list.List
	objMap  map[wire.InvVect]struct{}
	next    *list.Element
	last    *list.Element
}

func (q *InvQueue) Len() int {
	return len(q.objMap)
}

// PushBack adds an inv to the back of the queue.
func (q *InvQueue) PushBack(iv *wire.InvVect) bool {
	if _, ok := q.objMap[*iv]; ok {
		return false
	}
	q.objMap[*iv] = struct{}{}
	q.objList.PushBack(iv)
	return true
}

// PushFront adds an inv to the front of the queue.
func (q *InvQueue) PushFront(iv *wire.InvVect) bool {
	if _, ok := q.objMap[*iv]; ok {
		return false
	}
	q.objMap[*iv] = struct{}{}
	q.objList.PushFront(iv)
	return true
}

// PushFront adds an inv to the front of the queue.
func (q *InvQueue) PushBackList(list *list.List) bool {
	valid := true
	for e := list.Front(); e != nil; e = e.Next() {
		valid = valid && q.PushBack(e.Value.(*wire.InvVect))
	}
	return valid
}

// Exists says whether an inv is in the queue.
func (q *InvQueue) Exists(iv *wire.InvVect) bool {
	_, ok := q.objMap[*iv]
	return ok
}

// Front returns the inv at the front of the queue for iterating.
func (q *InvQueue) Front() *wire.InvVect {
	if q.objList.Len() == 0 {
		return nil
	}
	q.next = nil
	q.last = q.objList.Front()
	/*fmt.Println("Front called. first element is ", *q.last.Value.(*wire.InvVect))
	if q.last.Next() != nil {
		fmt.Println("Also, the element after that is ", *q.last.Next().Value.(*wire.InvVect))
	}*/
	return q.last.Value.(*wire.InvVect)
}

// Next iterates through the queue and returns nil when it reaches the end.
func (q *InvQueue) Next() *wire.InvVect {
	// This happens if the current element in the list has been removed.
	if q.next != nil {
		/*fmt.Println("Case 1: next is ", *q.next.Value.(*wire.InvVect))
		if q.next.Next() != nil {
			fmt.Println("Also, the element after that is ", *q.next.Next().Value.(*wire.InvVect))
		}*/
		q.last = q.next
		q.next = nil
		return q.last.Value.(*wire.InvVect)
	}
	if q.last != nil {
		/*fmt.Println("Case 2: last is ", *q.last.Value.(*wire.InvVect))
		if q.last.Next() != nil {
			fmt.Println("Also, the element after that is ", *q.last.Next().Value.(*wire.InvVect))
		} else {
			fmt.Println("Next element is nil.")
		}*/
		q.last = q.last.Next()
		if q.last != nil {
			return q.last.Value.(*wire.InvVect)
		}
		return nil
	}
	return nil
}

// Remove removes an inv from the queue.
func (q *InvQueue) Remove(iv *wire.InvVect) bool {
	if _, ok := q.objMap[*iv]; !ok {
		return false
	}
	if q.last != nil {
		if *iv == *q.last.Value.(*wire.InvVect) {
			q.next = q.last.Next()
			q.objList.Remove(q.last)
			/*if q.next != nil {
				fmt.Println("Remove called. last =", *q.last.Value.(*wire.InvVect), ", next =", *q.next.Value.(*wire.InvVect))
				if q.next.Next() != nil {
					fmt.Println("Also, the element after that is ", *q.next.Next().Value.(*wire.InvVect))
				}
			} else {
				fmt.Println("Remove called. last =", *q.last.Value.(*wire.InvVect), ", next is nil.")
			}*/
			q.last = nil
			delete(q.objMap, *iv)
			return true
		}
	}
	for e := q.objList.Front(); e != nil; e = e.Next() {
		if *e.Value.(*wire.InvVect) == *iv {
			q.objList.Remove(e)
			delete(q.objMap, *iv)
			return true
		}
	}
	return false // This line should never actually happen.
}

// NewInvQueue creates a new inv queue.
func NewInvQueue() *InvQueue {
	return &InvQueue{
		objList: list.New(),
		objMap:  make(map[wire.InvVect]struct{}),
	}
}

func (q *InvQueue) CheckIntegrity() int {
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
