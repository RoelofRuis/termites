package termites

type NodeId Identifier
type InPortId Identifier
type OutPortId Identifier
type ConnectionId Identifier

type NodeStatus uint8

const (
	NodeActive    NodeStatus = 0 // The node is active.
	NodeSuspended NodeStatus = 1 // The node is temporarily suspended and might be activated.
	NodeError     NodeStatus = 2 // The node has encountered an error.
)

type NodeRunningStatus uint8

const (
	NodePreStarted NodeRunningStatus = 0 // The node is not yet started.
	NodeRunning    NodeRunningStatus = 1 // The node is running.
	NodeTerminated NodeRunningStatus = 2 // The node has completed running.
)

type NodeRef struct {
	Id            NodeId
	Name          string
	Status        NodeStatus
	RunningStatus NodeRunningStatus
	InPorts       map[InPortId]InPortRef
	OutPorts      map[OutPortId]OutPortRef
	RunInfo       *FunctionInfo
	ShutdownInfo  *FunctionInfo
}

type InPortRef struct {
	Id   InPortId
	Name string
}

type OutPortRef struct {
	Id          OutPortId
	Name        string
	Connections []ConnectionRef
}

type ConnectionRef struct {
	Id      ConnectionId
	Adapter *AdapterRef
	In      *InPortRef
}

type AdapterRef struct {
	Name          string
	TransformInfo *FunctionInfo
}

type FunctionInfo struct {
	File string
	Line int
}