package GenID

import (
	"errors"
	"time"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func Init(startTime string, mechineID int64) error {
	if len(startTime) == 0 || mechineID <= 0 {
		return errors.New("invalid param")
	}

	var st time.Time
	st, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		return err
	}
	snowflake.Epoch = st.UnixNano() / 1000000
	node, err = snowflake.NewNode(mechineID)
	return err
}

func Get() int64 {
	return node.Generate().Int64()
}
