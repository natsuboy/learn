package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	for {
		limitIPFreq(&gin.Context{}, 5, 8)
		time.Sleep(time.Millisecond * 300)
	}
}

var LimitQueue map[string][]int64
var ok bool

// 时间窗口限流，单机版
func LimitFreqSingle(queueName string, count uint, timeWindow int64) bool {
	curTime := time.Now().Unix()
	if LimitQueue == nil {
		LimitQueue = make(map[string][]int64)
	}
	if _, ok = LimitQueue[queueName]; !ok {
		LimitQueue[queueName] = make([]int64, 0)
	}
	if uint(len(LimitQueue[queueName])) < count {
		LimitQueue[queueName] = append(LimitQueue[queueName], curTime)
		fmt.Println("pass: ", curTime, LimitQueue[queueName])
		return true
	}
	earlyTime := LimitQueue[queueName][0]
	if curTime-earlyTime <= timeWindow {
		fmt.Println("msg: error Current IP frequently visited.", curTime, LimitQueue[queueName])
		return false
	}
	LimitQueue[queueName] = LimitQueue[queueName][1:]
	LimitQueue[queueName] = append(LimitQueue[queueName], curTime)
	fmt.Println("pass: ", curTime, LimitQueue[queueName])
	return true
}

func limitIPFreq(c *gin.Context, timeWindow int64, count uint) bool {
	// ip := c.ClientIP()
	ip := "1.2.3.4"
	key := "limit:" + ip
	if !LimitFreqSingle(key, count, timeWindow) {
		// c.JSON(200, gin.H{
		//     "code": 400,
		//     "msg":  "error Current IP frequently visited",
		// })
		return false
	}
	return true
}
