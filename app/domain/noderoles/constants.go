package noderoles

import "github.com/codecrafters-io/redis-starter-go/app/domain"

const (
	Master domain.NodeRole = "master"
	Slave  domain.NodeRole = "slave"
)
