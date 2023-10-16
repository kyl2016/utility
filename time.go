package utility

import (
	"fmt"
	"time"
)

func ParseDateTime(datePart, timePart int) time.Time {
	if datePart <= 0 {
		return time.Time{}
	}

	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic("时区加载失败")
	}
	t, err := time.Parse("2006-01-02 15:04:05", time.Unix(int64(datePart/1000), 0).In(loc).Format("2006-01-02")+" "+time.Unix(int64(timePart), 0).In(loc).Format("15:04:05"))
	if err != nil {
		return time.Time{}
	}
	fmt.Println(time.Unix(int64(timePart), 0).In(loc).Format("15:04:05"))
	fmt.Println(t.String())
	return t
}
