package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// ./nats_benchmark -subject pubsub.station1.1.IED.1.point.1 -size 1024 -rate 1000 -duration 30s
// ./nats_benchmark -subject test -size 1024 -rate 1000 -duration 30s -js -stream TEST_STREAM -stream-subjects test

type BenchmarkConfig struct {
	ServerURL     string        // NATS服务器地址
	Subject       string        // 消息主题
	MessageSize   int           // 消息大小(字节)
	Duration      time.Duration // 测试持续时间
	RatePerSecond int           // 每秒发送消息数
	JetStream     bool          // 是否使用JetStream
	StreamName    string        // JetStream流名称
}

var (
	totalMessages uint64 = 0
	config               = parseFlags() // 解析命令行参数
)

func main() {

	// 连接到NATS服务器
	nc, err := nats.Connect(config.ServerURL)
	if err != nil {
		log.Fatalf("无法连接到NATS服务器: %v", err)
	}
	defer nc.Close()

	// 如果启用JetStream，创建流
	var js jetstream.JetStream
	if config.JetStream {
		js, err = jetstream.New(nc)
		if err != nil {
			log.Fatalf("无法创建JetStream: %v", err)
		}

		// 检查流是否存在，不存在则创建
		stream, err := js.Stream(context.Background(), config.StreamName)
		if err != nil || stream == nil {
			log.Fatalf("未找到JetStream流 %v: %v", config.StreamName, err)
		}
	}

	// 生成测试消息
	message := make([]byte, config.MessageSize)
	for i := range message {
		message[i] = 'a' // 填充测试数据
	}

	// 启动基准测试
	log.Printf("开始基准测试: 主题=%s, 大小=%d字节, 速率=%d消息/秒, 持续时间=%v, JetStream=%v",
		config.Subject, config.MessageSize, config.RatePerSecond, config.Duration, config.JetStream)

	ctx, cancel := context.WithTimeout(context.Background(), config.Duration)
	defer cancel()

	go reportStats(ctx)

	start := time.Now()
	// 根据配置选择发布方法
	if config.JetStream {
		runJetStreamBenchmark(ctx, js, config.Subject, message, config.RatePerSecond, &totalMessages)
	} else {
		runRegularBenchmark(ctx, nc, config.Subject, message, config.RatePerSecond, &totalMessages)
	}
	// 计算结果
	duration := time.Since(start)
	messagesPerSecond := float64(totalMessages) / duration.Seconds()
	bytesPerSecond := messagesPerSecond * float64(config.MessageSize)

	// 打印结果
	log.Printf("基准测试完成:")
	log.Printf("  总发送消息数: %d", totalMessages)
	log.Printf("  测试持续时间: %.2f秒", duration.Seconds())
	log.Printf("  消息速率: %.2f 消息/秒", messagesPerSecond)
	log.Printf("  数据速率: %.2f KB/秒 (%.2f MB/秒)",
		bytesPerSecond/1024, bytesPerSecond/(1024*1024))

}
func reportStats(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	prevCount := uint64(0)
	prevTime := time.Now()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			current := atomic.LoadUint64(&totalMessages)
			now := time.Now()
			elapsed := now.Sub(prevTime).Seconds()
			rate := float64(current-prevCount) / elapsed
			bytesPerSecond := rate * float64(config.MessageSize)
			fmt.Printf("Published: %d (%.2f msg/sec)\n", current, rate)
			// 打印结果
			log.Printf("  |总发送消息数: %d  |统计间隔: %.2f秒 |速率: %.2f 消息/秒  %.4f MB/秒",
				totalMessages, elapsed, rate, bytesPerSecond/(1024*1024))
			prevCount = current
			prevTime = now
		}
	}
}

// 解析命令行参数
func parseFlags() BenchmarkConfig {
	config := BenchmarkConfig{}

	flag.StringVar(&config.ServerURL, "s", "nats://localhost:4222", "NATS服务器地址")
	flag.StringVar(&config.Subject, "subject", "test", "消息主题")
	flag.IntVar(&config.MessageSize, "size", 1024, "消息大小(字节)")
	flag.DurationVar(&config.Duration, "duration", 10*time.Second, "测试持续时间")
	flag.IntVar(&config.RatePerSecond, "rate", 100, "每秒发送消息数")
	flag.BoolVar(&config.JetStream, "js", false, "是否使用JetStream")
	flag.StringVar(&config.StreamName, "stream", "TEST_STREAM", "JetStream流名称")

	flag.Parse()

	return config
}

// 普通消息基准测试
func runRegularBenchmark(ctx context.Context, nc *nats.Conn, subject string, message []byte, ratePerSecond int, counter *uint64) {
	// 计算发送间隔
	interval := time.Duration(1e9 / ratePerSecond)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := nc.Publish(subject, message)
			if err != nil {
				log.Printf("发布消息失败: %v", err)
				continue
			}
			atomic.AddUint64(counter, 1)
		}
	}
}

// JetStream消息基准测试
func runJetStreamBenchmark(ctx context.Context, js jetstream.JetStream, subject string, message []byte, ratePerSecond int, counter *uint64) {
	// 计算发送间隔
	interval := time.Duration(1e9 / ratePerSecond)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_, err := js.Publish(ctx, subject, message)
			if err != nil {
				log.Printf("发布JetStream消息失败: %v", err)
				continue
			}
			atomic.AddUint64(counter, 1)
		}
	}
}
