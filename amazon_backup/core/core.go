package core

import (
	"fmt"

	"github.com/sony/sonyflake"
)

var (
	_idGenerator *sonyflake.Sonyflake
)

func init() {

	// id generator
	_idGenerator = sonyflake.NewSonyflake(sonyflake.Settings{})
}

// GenerateID _
func GenerateID() string {
	a, _ := _idGenerator.NextID()
	return fmt.Sprintf("%x", a)
}
