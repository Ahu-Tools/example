package v1

import (
	"github.com/Ahu-Tools/example/edge/asynq"
	//@ahum: imports
)

const version = "v1"

// A list of task types.
const (
	ModuleName = "hello"
	TypeWorld  = "world"
	//@ahum: types
)

func GetPattern(handler string) string {
	return ModuleName + ":" + handler
}

func init() {
	asynq.RegisterHandler(version, GetPattern(TypeWorld), HandleWorld)
	// @ahum: registers
}
