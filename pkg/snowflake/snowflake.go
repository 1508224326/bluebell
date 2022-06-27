package snowflake

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"time"
)

var node *snowflake.Node

func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-01", startTime)
	if err != nil {
		fmt.Println("时间解析错误")
		return
	}

	// 时间片偏移到指定的时间上
	snowflake.Epoch = st.UnixNano() / 1000000
	node, err = snowflake.NewNode(machineID)
	return
}

func GenID() int64 {
	return node.Generate().Int64()
}

func Demo() {
	if err := Init("2022-03-01", 1); err != nil {
		fmt.Println("init failed, err: ", err)
		return
	}
	id := GenID()
	fmt.Println(id)
}
