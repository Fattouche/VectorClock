package main

import (
	"fmt"
)

// This is an api that allows distributed systems to manage vector clocks.

// vectorClock struct to control who owns the clock
type vectorClock struct {
	myNode string
	vector map[string]int
}

// Merges two vector clocks together and updates the self vectorClock
// Parameters:
//			self,peer(vectorclock): The two vector clocks to be merged together
// Returns:
//			err(error): The error if the vector clocks cannot be merged

func (self vectorClock) merge(peer vectorClock) (err error) {
	if len(self.vector) != len(peer.vector) {
		err = fmt.Errorf("vector clocks differ in length a: %d  b: %d", len(self.vector), len(peer.vector))
		return
	}
	for node := range self.vector {
		if _, ok := peer.vector[node]; !ok {
			err = fmt.Errorf("vector clocks differ in contents, found %s in a but not in b", node)
			return
		}
		self.vector[node] = max(self.vector[node], peer.vector[node])
	}
	return
}

// Increments the value for a single node within a vector clock
//	Returns:
//			err(error): Errors if the node to be incremented does not exist within the vector clock

func (vc vectorClock) increment() {
	vc.vector[vc.myNode]++
	return
}

// Creates a new vector clock
// Parameters:
//			myIndex(int): The index of the caller node within the slice of nodes
// 			nodes([]string): A slice of the nodes that will exist in the distributed system
// Returns:
//			vc(vectorClock): The newly created vector clock
//			err(error): The error during creation if any

func newVectorClock(myIndex int, nodes []string) (vc vectorClock, err error) {
	if myIndex >= len(nodes) || myIndex < 0 {
		err = fmt.Errorf("myIndex parameter invalid: Expected between 0 and %d, got %d", len(nodes), myIndex)
		return
	}
	vector := make(map[string]int)
	vc = vectorClock{nodes[myIndex], vector}
	for _, node := range nodes {
		vc.vector[node] = 0
	}
	return
}

// Returns the max between two numbers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
