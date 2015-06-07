package invqueue_test

import (
	"testing"

	"github.com/DanielKrawisz/invqueue"
	"github.com/monetas/bmutil/wire"
)

func TestDuplicate(t *testing.T) {
	hashA, _ := wire.NewShaHashFromStr("0b4448198e83b24ccfb332e0b5b3788244e3c324ddaf0f9f6cdd67962a673992")
	hashB, _ := wire.NewShaHashFromStr("0b4448198e83b24ccfb332e0b5b3788244e3c324ddaf0f9f6cdd67962a673992")

	a := &wire.InvVect{Hash: *hashA}
	b := &wire.InvVect{Hash: *hashB}

	invQueue := invqueue.NewInvQueue()

	lenth := invQueue.Len()
	if lenth != 0 {
		t.Error("Len is %d, expected 0", lenth)
	}

	if invQueue.Front() != nil {
		t.Error("Front failed to return nil.")
	}

	if invQueue.Remove(a) {
		t.Error("Should be nothing to remove.")
	}
	lenth = invQueue.Len()
	if lenth != 0 {
		t.Error("Len is %d, expected 0", lenth)
	}

	invQueue.PushBack(a)

	lenth = invQueue.Len()
	if lenth != 1 {
		t.Errorf("Len is %d, expected 1", lenth)
	}

	invQueue.PushBack(b)

	lenth = invQueue.Len()
	if lenth != 1 {
		t.Error("Len is %d, expected 1", lenth)
	}

	check := invQueue.CheckIntegrity()
	if check != 0 {
		t.Errorf("Integrity check fails after adding duplicate item to the back: ", check)
	}

	invQueue.PushFront(b)

	check = invQueue.CheckIntegrity()
	if check != 0 {
		t.Error("Integrity check fails after adding duplicate item to the front: ", check)
	}

	if !invQueue.Exists(a) {
		t.Error("Queue should contain object a.")
	}

	if !invQueue.Exists(b) {
		t.Error("Queue should contain object b.")
	}

	if !invQueue.Remove(b) {
		t.Error("Should have removed object b.")
	}
	lenth = invQueue.Len()
	if lenth != 0 {
		t.Error("Len is %d, expected 0", lenth)
	}

	if invQueue.Exists(a) {
		t.Error("Queue should not contain object a.")
	}

	if invQueue.Exists(b) {
		t.Error("Queue should not contain object b.")
	}

	check = invQueue.CheckIntegrity()
	if check != 0 {
		t.Error("Integrity check fails after removing item: ", check)
	}
}

func TestRemove(t *testing.T) {
	for i := 0; i < 3; i++ {
		hashA, _ := wire.NewShaHashFromStr("000048198e83b24ccfb332e0b5b3788244e3c324ddaf0f9f6cdd67962a673992")
		hashB, _ := wire.NewShaHashFromStr("1111f020493cb3374894433c03c7e5f671ac39fa7443d25217e61b693074927d")
		hashC, _ := wire.NewShaHashFromStr("22227c88d21e5933f0a63b98084b1201e67530ad874ffa5aad2bb88a20bca36c")

		a := &wire.InvVect{Hash: *hashA}
		b := &wire.InvVect{Hash: *hashB}
		c := &wire.InvVect{Hash: *hashC}

		invQueue := invqueue.NewInvQueue()

		invQueue.PushBack(a)
		invQueue.PushFront(b)
		invQueue.PushBack(c)

		lenth := invQueue.Len()
		if lenth != 3 {
			t.Error("Len is %d, expected 3", lenth)
		}

		if !invQueue.Exists(a) {
			t.Error("Queue should contain object a.")
		}

		if !invQueue.Exists(b) {
			t.Error("Queue should contain object b.")
		}

		if !invQueue.Exists(c) {
			t.Error("Queue should contain object c.")
		}

		check := invQueue.CheckIntegrity()
		if check != 0 {
			t.Error("Integrity check fails after populating list: ", check)
		}

		expectedOrder := []*wire.InvVect{b, a, c}

		invQueue.Remove(expectedOrder[i])

		lenth = invQueue.Len()
		if lenth != 2 {
			t.Error("Len is %d, expected 0", lenth)
		}

		k := 0
		for x := invQueue.Front(); x != nil; x = invQueue.Next() {
			if i == k {
				k++
			}
			t.Log("k = ", k, ", x =", *x)
			if *x != *expectedOrder[k] {
				t.Errorf("Ordering is incorrect after iterating through the list. Expected %s got %s", *expectedOrder[k], *x)
			}

			k++
		}
	}
}

func TestIterateRemove(t *testing.T) {
	hashA, _ := wire.NewShaHashFromStr("000048198e83b24ccfb332e0b5b3788244e3c324ddaf0f9f6cdd67962a673992")
	hashB, _ := wire.NewShaHashFromStr("1111f020493cb3374894433c03c7e5f671ac39fa7443d25217e61b693074927d")
	hashC, _ := wire.NewShaHashFromStr("22227c88d21e5933f0a63b98084b1201e67530ad874ffa5aad2bb88a20bca36c")

	a := &wire.InvVect{Hash: *hashA}
	b := &wire.InvVect{Hash: *hashB}
	c := &wire.InvVect{Hash: *hashC}

	var reorderList func(newList, oldList []*wire.InvVect, swap, elem int)
	reorderList = func(newList, oldList []*wire.InvVect, swap, elem int) {
		if elem == 2 {
			newList[elem] = oldList[swap]
			return
		}
		if elem >= swap {
			newList[elem] = oldList[elem+1]
			reorderList(newList, oldList, swap, elem+1)
			return
		}
		newList[elem] = oldList[elem]
		reorderList(newList, oldList, swap, elem+1)
	}

	testFunc := func(invQueue *invqueue.InvQueue, expectedOrderBefore []*wire.InvVect,
		expectedOrderAfter []*wire.InvVect, i int) {
		j := 0
		for x := invQueue.Front(); x != nil; x = invQueue.Next() {
			t.Log("j = ", j)
			if j >= 3 {
				t.Fatal("Iteration went too far!")
			}
			if *x != *expectedOrderBefore[j] {
				t.Error("Ordering is incorrect.")
			}

			if i == j {
				if !invQueue.Remove(x) {
					t.Error("Should have removed element ", j)
				}
			}

			check := invQueue.TstCheckIntegrity()
			if check != 0 {
				t.Error("Integrity check fails after iteration ", j)
			}

			j++
		}

		lenth := invQueue.Len()
		if lenth != 2 {
			t.Error("Len is %d, expected 2", lenth)
		}

		invQueue.PushBack(expectedOrderBefore[i])

		k := 0
		for x := invQueue.Front(); x != nil; x = invQueue.Next() {

			t.Log("k = ", k, ", x =", *x)
			if *x != *expectedOrderAfter[k] {
				t.Errorf("Ordering is incorrect after iterating through the list. Expected %s got %s", *expectedOrderAfter[k], *x)
			}

			k++
		}
	}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {

			invQueue := invqueue.NewInvQueue()

			invQueue.PushBack(a)
			invQueue.PushFront(b)
			invQueue.PushBack(c)

			lenth := invQueue.Len()
			if lenth != 3 {
				t.Error("Len is %d, expected 3", lenth)
			}

			if !invQueue.Exists(a) {
				t.Error("Queue should contain object a.")
			}

			if !invQueue.Exists(b) {
				t.Error("Queue should contain object b.")
			}

			if !invQueue.Exists(c) {
				t.Error("Queue should contain object c.")
			}

			check := invQueue.CheckIntegrity()
			if check != 0 {
				t.Error("Integrity check fails after populating list: ", check)
			}

			expectedOrderBefore := []*wire.InvVect{b, a, c}

			expectedOrderAfter := make([]*wire.InvVect, 3)
			reorderList(expectedOrderAfter, expectedOrderBefore, i, 0)

			testFunc(invQueue, expectedOrderBefore, expectedOrderAfter, i)

			expectedOrderBefore = expectedOrderAfter

			expectedOrderAfter = make([]*wire.InvVect, 3)
			reorderList(expectedOrderAfter, expectedOrderBefore, j, 0)

			testFunc(invQueue, expectedOrderBefore, expectedOrderAfter, j)

		}
	}
}

func TestBreakRemove(t *testing.T) {
	hashA, _ := wire.NewShaHashFromStr("000048198e83b24ccfb332e0b5b3788244e3c324ddaf0f9f6cdd67962a673992")
	hashB, _ := wire.NewShaHashFromStr("1111f020493cb3374894433c03c7e5f671ac39fa7443d25217e61b693074927d")
	hashC, _ := wire.NewShaHashFromStr("22227c88d21e5933f0a63b98084b1201e67530ad874ffa5aad2bb88a20bca36c")

	a := &wire.InvVect{Hash: *hashA}
	b := &wire.InvVect{Hash: *hashB}
	c := &wire.InvVect{Hash: *hashC}

	invQueue := invqueue.NewInvQueue()

	invQueue.PushBack(a)
	invQueue.PushFront(b)
	invQueue.PushBack(c)

	invQueue.TstBreakRemove(b) // b is removed from the list but not the hash.

	invQueue.Remove(b) // Remove should now execute a line that should never really happen.

	check := invQueue.TstCheckIntegrity()
	if check != 0 {
		t.Error("Integrity check fails removing element.")
	}
}

func TestIterateRemove2(t *testing.T) {
	for i := 0; i < 3; i++ {
		for j := i + 1; j < 4; j++ {
			hashA, _ := wire.NewShaHashFromStr("000048198e83b24ccfb332e0b5b3788244e3c324ddaf0f9f6cdd67962a673992")
			hashB, _ := wire.NewShaHashFromStr("1111f020493cb3374894433c03c7e5f671ac39fa7443d25217e61b693074927d")
			hashC, _ := wire.NewShaHashFromStr("22227c88d21e5933f0a63b98084b1201e67530ad874ffa5aad2bb88a20bca36c")
			hashD, _ := wire.NewShaHashFromStr("33336db1db7268e09ee4c3f51e54ae2f0b548541e4fd37d9f42fdb1cabe29266")

			a := &wire.InvVect{Hash: *hashA}
			b := &wire.InvVect{Hash: *hashB}
			c := &wire.InvVect{Hash: *hashC}
			d := &wire.InvVect{Hash: *hashD}

			invQueue := invqueue.NewInvQueue()

			invQueue.PushBack(a)
			invQueue.PushBack(b)
			invQueue.PushBack(c)
			invQueue.PushBack(d)

			expectedOrder := []*wire.InvVect{a, b, c, d}

			k := 0
			for x := invQueue.Front(); x != nil; x = invQueue.Next() {
				t.Log("k = ", k)
				if *x != *expectedOrder[k] {
					t.Errorf("Ordering is incorrect; expected %s, got %s", *expectedOrder[k], *x)
				}

				if k == j || k == i {
					if !invQueue.Remove(x) {
						t.Error("Should have removed element ", k)
					}
				}

				check := invQueue.TstCheckIntegrity()
				if check != 0 {
					t.Error("Integrity check fails after iteration ", k)
				}

				k++
			}

			lenth := invQueue.Len()
			if lenth != 2 {
				t.Error("Len is %d, expected 2", lenth)
			}

			invQueue.PushBack(expectedOrder[i])
			invQueue.PushBack(expectedOrder[j])

			k = 0
			var z int
			for x := invQueue.Front(); x != nil; x = invQueue.Next() {
				if k == i {
					k++
				}
				if k == j {
					k++
				}
				z = k
				if k == 4 {
					z = i
				}
				if k == 5 {
					z = j
				}
				t.Log("k = ", k, ", z = ", z, ", x =", *x)
				if *x != *expectedOrder[z] {
					t.Errorf("Ordering is incorrect after iterating through the list. Expected %s got %s", *expectedOrder[z], *x)
				}

				k++
			}
		}
	}
}

func TestPushList(t *testing.T) {

}
