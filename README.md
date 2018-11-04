# VectorClock API
This is a golang api for vector clocks to allow for partial ordering of events.

## API
#####   newVectorClock(myIndex int, nodes []string)

**Creates a new vector clock**
Parameters:
* _myIndex(int)_: The index of the caller node within the slice of nodes
* _nodes([]string)_: A slice of the nodes that will exist in the distributed system

Returns:
* _vc(vectorClock)_: The newly created vector clock
* _err(error)_: The error during creation if any

##### (vc vectorClock) increment()
**Increments the value for a single node within a vector clock** 
Returns:
* err(error): Errors if the node to be incremented does not exist within the vector clock

##### (self vectorClock) merge(peer vectorClock)

**Merges two vector clocks together and updates the self vectorClock**
Parameters:
* _self,peer(vectorclock)_: The two vector clocks to be merged together

Returns:
* _err(error)_: The error if the vector clocks cannot be merged


## Design
The code is designed to function as a standalone api that can be utilized by any distributed system that needs partial ordering of node events. The core of the structure stems from the _vectorClock_ struct type which is the vector clock for a specific node. The respective node can then call functions like _increment()_ and _merge(peer vectorClock)_ on their own vectorClock object. In order for the node to get a new vectorClock it must pass in the list of all nodes within the distributed system as well as the index of itself within the list. This information is used to keep track of who owns what clock and allows the code to have a more object oriented feel. When a node recieves a new message from another node, it must call _merge()_ with its own vectorClock while passing in the newly recieved vectorClock from the peer. This will then update the recieing nodes vector clock by reference.