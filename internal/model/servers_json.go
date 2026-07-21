package model

type ServersFile struct {
	Mode  string `json:"mode"`
	Pools []Pool `json:"pools"`
}
