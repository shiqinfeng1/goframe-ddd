package limiter

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

// 自定义的 MockWriter 结构体，用于模拟 Writer 接口
type MockWriter struct {
	// 用于存储写入的数据
	Buffer bytes.Buffer
	// 模拟写入时可能返回的错误
	Err error
}

func (w MockWriter) Read(p []byte) (n int, err error) {
	return
}

func (w MockWriter) Close() error {
	return nil
}

// 实现 Writer 接口的 Write 方法
func (m *MockWriter) Write(p []byte) (n int, err error) {
	if m.Err != nil {
		return 0, m.Err
	}
	return m.Buffer.Write(p)
}

func TestLimit(t *testing.T) {
	speed := 10 // 单位10k
	data := make([]byte, 0, 1024*1000)
	for range 1024 * 3 {
		data = append(data, []byte("1234567890")...)
	}
	// 创建一个 MockWriter 实例
	mockWriter := &MockWriter{}
	lw := NewLimitWriter(mockWriter, 1024*speed)
	start := time.Now()
	n, err := lw.Write(data)

	assert.Assert(t, func() (success bool, message string) {
		sp := int64(len(data)) * 1000000 / 1024 / time.Since(start).Microseconds()
		if sp-int64(speed) > 2 {
			return false, fmt.Sprintf("exp:%v act:%v", speed, sp)
		} else {
			return true, fmt.Sprintf("exp:%v act:%v", speed, sp)
		}
	})

	assert.NilError(t, err)
	assert.Equal(t, n, len(data))
	assert.DeepEqual(t, mockWriter.Buffer.Bytes(), data)

	// t.Errorf("speed= %v KB/s\n", int64(len(data))*1000000/1024/time.Since(start).Microseconds())
}
