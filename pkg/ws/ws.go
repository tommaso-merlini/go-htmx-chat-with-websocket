package ws

import "roomate/types"

var WSConnections = types.Conns{}

func Init() {
	WSConnections = make(types.Conns)
}
