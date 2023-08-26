package config

import "sync"

var (
	ThreadNum = 500
	Result    *sync.Map
)

func init() {
	Result = &sync.Map{}
}
