package termites

type NodeId = Identifier
type InPortId = Identifier
type OutPortId = Identifier
type ConnectionId = Identifier

type NodeRef struct {
	Id           NodeId
	Version      uint
	Name         string
	InPorts      map[InPortId]InPortRef
	OutPorts     map[OutPortId]OutPortRef
	RunInfo      FunctionInfo
	ShutdownInfo FunctionInfo
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
	TransformInfo FunctionInfo
}

type FunctionInfo struct {
	File string
	Line int
}
