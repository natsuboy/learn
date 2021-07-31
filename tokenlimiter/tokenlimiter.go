package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	tokenLimiter := NewTokenLimiter(3, 5)
	for {
		n := 4
		for i := 0; i < n; i++ {
			go func(i int) {
				if !tokenLimiter.IsAllow() {
					fmt.Printf("forbid [%d]\n", i)
				} else {
					fmt.Printf("allow  [%d]\n", i)
				}
			}(i)
		}
		time.Sleep(time.Second)
		fmt.Println("==============================")
	}
}

type TokenLimiter struct {
	limit  int // 速率
	burst  int // 桶的大小
	tokens int // 桶里的令牌数量

	lock     sync.Mutex
	lastTime time.Time // 最后一次获取令牌的时间
}

func NewTokenLimiter(l, b int) *TokenLimiter {
	return &TokenLimiter{
		limit: l,
		burst: b,
	}
}

func (tl *TokenLimiter) IsAllow() bool {
	return tl.isAllowN(time.Now(), 1)
}

func (tl *TokenLimiter) isAllowN(t time.Time, n int) bool {
	tl.lock.Lock()
	defer tl.lock.Unlock()

	delta := int(t.Sub(tl.lastTime).Seconds()) * tl.limit
	tl.tokens += delta

	if tl.tokens > tl.burst {
		tl.tokens = tl.burst
	}

	if n > tl.tokens {
		return false
	}

	tl.tokens -= n
	tl.lastTime = t

	return true
}
