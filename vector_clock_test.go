package main

import (
	"sync"
	"testing"
	"time"
)

func TestMerge(t *testing.T) {
	nodes := []string{"node1", "node2", "node3"}
	vc1, _ := newVectorClock(0, nodes)
	vc2, _ := newVectorClock(1, nodes)
	vcBad, _ := newVectorClock(0, []string{"node_random"})
	vc1.vector[nodes[0]] = 6
	vc1.vector[nodes[1]] = 3
	vc1.vector[nodes[2]] = 19

	vc2.vector[nodes[0]] = 8
	vc2.vector[nodes[1]] = 7
	vc2.vector[nodes[2]] = 4

	err := vc1.merge(vc2)
	if err != nil {
		t.Error(err)
	}
	if vc1.vector[nodes[0]] != 8 && vc1.vector[nodes[1]] != 7 && vc1.vector[nodes[2]] != 19 {
		t.Errorf("Expected values of %d,%d,%d for node1,node2,node3 respectiively but got %d,%d,%d", 8, 7, 19, vc1.vector[nodes[0]], vc1.vector[nodes[1]], vc1.vector[nodes[2]])
	}
	err = vc1.merge(vcBad)
	if err == nil {
		t.Errorf("Expected error for invalid merge but none was returned")
	}
}

func TestIncrement(t *testing.T) {
	nodes := []string{"node1", "node2", "node3"}
	vc, _ := newVectorClock(0, nodes)
	vc.increment()
	vc.increment()
	vc.increment()
	if vc.vector[nodes[0]] != 3 {
		t.Errorf("Unexpected value for node1 in vector clock: Expected %d but got %d", 3, vc.vector[nodes[0]])
	}
	if vc.vector[nodes[1]] != 0 {
		t.Errorf("Unexpected value for node1 in vector clock: Expected %d but got %d", 0, vc.vector[nodes[1]])
	}
	if vc.vector[nodes[2]] != 0 {
		t.Errorf("Unexpected value for node2 in vector clock: Expected %d but got %d", 0, vc.vector[nodes[2]])
	}
}

func TestNewVectorClock(t *testing.T) {
	nodes := []string{"node1", "node2", "node3"}
	vc, err := newVectorClock(0, nodes)
	if err != nil {
		t.Error(err)
	}
	if _, ok := vc.vector[nodes[0]]; !ok {
		t.Error("Failed to find node1 in vector clock")
	}
	val := vc.vector[nodes[0]]
	if val != 0 {
		t.Errorf("Unexpected value for node1 in vector clock: Expected %d but got %d", 0, val)
	}
	vc, err = newVectorClock(3, nodes)
	if err == nil {
		t.Errorf("Expected error for invalid index but none was returned")
	}
}

func TestMax(t *testing.T) {
	a, b := 0, 1
	c := max(a, b)
	if c != 1 {
		t.Errorf("Incorrect max: Expected %d but got %d", b, c)
	}
}

func clockCycle() {
	time.Sleep(time.Millisecond * 100)
}

// This test is taken from https://en.wikipedia.org/wiki/Vector_clock
// In a real distributed system this would be done in parallel, however this test is only concerned with the general correctness of the api.
func TestFullOne(t *testing.T) {
	nodes := []string{"node1", "node2", "node3"}
	vc1, _ := newVectorClock(0, nodes)
	vc2, _ := newVectorClock(1, nodes)
	vc3, _ := newVectorClock(2, nodes)
	wg := sync.WaitGroup{}
	wg.Add(3)
	//Node 3
	go func() {
		vc3.increment()
		clockCycle()
		clockCycle()
		clockCycle()
		clockCycle()
		clockCycle()
		clockCycle()
		clockCycle()
		vc3.increment()
		vc3.merge(vc2)
		clockCycle()
		clockCycle()
		vc3.increment()
		clockCycle()
		clockCycle()
		clockCycle()
		vc3.increment()
		vc3.merge(vc2)
		vc3.increment()
		clockCycle()
		clockCycle()
		wg.Done()
	}()
	//Node 2
	go func() {
		clockCycle()
		vc2.increment()
		vc2.merge(vc3)
		vc2.increment()
		clockCycle()
		clockCycle()
		clockCycle()
		vc2.increment()
		clockCycle()
		clockCycle()
		vc2.increment()
		vc2.merge(vc1)
		clockCycle()
		vc2.increment()
		clockCycle()
		clockCycle()
		clockCycle()
		clockCycle()
		clockCycle()
		clockCycle()
		clockCycle()
		wg.Done()
	}()
	//Node 1
	go func() {
		clockCycle()
		clockCycle()
		clockCycle()
		clockCycle()
		vc1.merge(vc2)
		vc1.increment()
		vc1.increment()
		clockCycle()
		clockCycle()
		clockCycle()
		clockCycle()
		clockCycle()
		clockCycle()
		clockCycle()
		vc1.increment()
		vc1.merge(vc3)
		clockCycle()
		clockCycle()
		clockCycle()
		vc1.increment()
		vc1.merge(vc3)
		wg.Done()
	}()
	wg.Wait()
	if vc1.vector[nodes[0]] != 4 || vc1.vector[nodes[1]] != 5 || vc1.vector[nodes[2]] != 5 {
		t.Errorf("Incorrect values: Expected 4,5,5 but got %d,%d,%d", vc1.vector[nodes[0]], vc1.vector[nodes[1]], vc1.vector[nodes[2]])
	}
}

func TestFullTwo(t *testing.T) {
	nodes := []string{"P0", "P1", "P2"}
	vc1, _ := newVectorClock(0, nodes)
	vc2, _ := newVectorClock(1, nodes)
	vc3, _ := newVectorClock(2, nodes)
	vc1.increment()
	vc2.increment()
	vc3.increment()
	vc1.increment()
	vc2.merge(vc1)
	vc1.increment()
	vc1.merge(vc2)
	vc2.increment()
	vc1.increment()
	vc3.increment()
	vc1.merge(vc3)
	vc3.increment()
	vc3.merge(vc1)
	vc1.increment()
	vc1.increment()
	vc2.merge(vc1)
	vc2.increment()
	vc1.increment()
	if vc1.vector[nodes[0]] != 7 || vc1.vector[nodes[1]] != 1 || vc1.vector[nodes[2]] != 2 {
		t.Errorf("Incorrect values: Expected 7,1,2 but got %d,%d,%d", vc1.vector[nodes[0]], vc1.vector[nodes[1]], vc1.vector[nodes[2]])
	}
}
