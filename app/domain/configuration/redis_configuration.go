package configuration

type ReplicaOf struct {
	Host string
	Port string
}

type RedisConfiguration struct {
	Port      string
	Role      string
	ReplicaOf *ReplicaOf
}
