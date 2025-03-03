package limiter

import (
	"context"
	"io"

	"golang.org/x/time/rate"
)

type LimitWriter struct {
	next             io.ReadWriteCloser
	limiter          *rate.Limiter
	defaultTokenNums int
	chunkSize        int // 根据发送速率计算: speed/defaultTokenNums
}

var (
	defaultTokenNums = 10   // 令牌桶每秒产生的token数量, 表示每秒的执行次数
	defaultTokenCap  = 1    // 令牌通容量
	minSpeed         = 1024 // 最低限速：1k
)

func NewLimitWriter(writer io.ReadWriteCloser, speed int) io.ReadWriteCloser {
	if speed < minSpeed {
		speed = minSpeed
	}
	limiter := rate.NewLimiter(rate.Limit(defaultTokenNums), defaultTokenCap)
	return &LimitWriter{
		next:             writer,
		limiter:          limiter,
		defaultTokenNums: defaultTokenNums,
		chunkSize:        speed / defaultTokenNums,
	}
}

func (w LimitWriter) Read(p []byte) (n int, err error) {
	return w.next.Read(p)
}

func (w LimitWriter) Close() error {
	return w.next.Close()
}

func (w LimitWriter) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	sendedSize := 0
	for {
		// 限速等待
		err := w.limiter.Wait(context.Background())
		if err != nil {
			return 0, err
		}
		// 发送完整分块
		if sendedSize+w.chunkSize < len(p) {
			if n, err := w.next.Write(p[sendedSize : sendedSize+w.chunkSize]); err != nil {
				return sendedSize + n, err
			}
			sendedSize += w.chunkSize
		} else { // 发送剩余数据
			if n, err := w.next.Write(p[sendedSize:]); err != nil {
				return sendedSize + n, err
			}
			sendedSize = len(p)
		}
		// 全部发送完成，直接返回
		if sendedSize >= len(p) {
			return sendedSize, nil
		}
	}
}
