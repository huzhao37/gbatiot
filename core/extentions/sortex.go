package extensions

import (
	"github.com/influxdata/influxdb/client"
	"sort"
)

// ByMotorId implements sort.Interface for []client.Point based on
// the MotorId field.

type PointWrapper struct {
	p [] client.Point
	by func(p, q * client.Point) bool
}

type SortBy func(p, q *client.Point) bool

func (pw PointWrapper) Len() int {         // 重写 Len() 方法
	return len(pw.p)
}
func (pw PointWrapper) Swap(i, j int){     // 重写 Swap() 方法
	pw.p[i], pw.p[j] = pw.p[j], pw.p[i]
}
func (pw PointWrapper) Less(i, j int) bool {    // 重写 Less() 方法
	return pw.by(&pw.p[i], &pw.p[j])
}

// 封装成 SortPoint 方法
func SortPoint(p []client.Point, by SortBy){
	sort.Sort(PointWrapper{p, by})
}


