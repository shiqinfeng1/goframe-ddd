package application

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/bench"
)

var benchmark *bench.Benchmark

// 基准测试无业务逻辑处理，不需要domian层参与，因此直接在app层实现测试逻辑
func (h *Application) PubSubBenchmark(ctx context.Context, in *PubSubBenchmarkInput) error {
	numTopics := len(in.Subjects)

	benchmark = bench.NewBenchmark("NATS", in.NumSubs*numTopics, in.NumPubs*numTopics)
	// 配置连接选项
	opts := []nats.Option{nats.Name("NATS Benchmark")}

	var (
		startwg sync.WaitGroup
		donewg  sync.WaitGroup
	)

	donewg.Add((in.NumPubs + in.NumSubs) * numTopics)

	// 先运行订阅者，一个订阅者使用一个连接
	startwg.Add(in.NumSubs * numTopics)
	for i := range numTopics {
		for range in.NumSubs {
			nc, err := nats.Connect(g.Cfg().MustGet(ctx, "nats.serverAddr").String(), opts...)
			if err != nil {
				g.Log().Fatalf(ctx, "Can't connect: %v", err)
			}
			defer nc.Close()

			go runSubscriber(nc, in.Subjects[i], &startwg, &donewg, in.NumMsgs, in.MsgSize)
		}
	}

	startwg.Wait()

	// 再运行发布者，一个发布者一个连接
	startwg.Add(in.NumPubs * numTopics)
	pubCounts := bench.MsgsPerClient(in.NumMsgs, in.NumPubs)
	for j := range numTopics {
		for i := range in.NumPubs {
			nc, err := nats.Connect(g.Cfg().MustGet(ctx, "nats.serverAddr").String(), opts...)
			if err != nil {
				g.Log().Fatalf(ctx, "Can't connect: %v\n", err)
			}
			defer nc.Close()

			go runPublisher(nc, in.Subjects[j], &startwg, &donewg, pubCounts[i], in.MsgSize)
		}
	}

	g.Log().Infof(ctx, "Starting benchmark [topics=%v msgs=%d, msgsize=%d, pubs=%d, subs=%d]", numTopics, in.NumMsgs, in.MsgSize, in.NumPubs, in.NumSubs)

	startwg.Wait()
	donewg.Wait()

	benchmark.Close()

	g.Log().Infof(ctx, "\n-----------%v\n-----------", benchmark.Report())

	csv := benchmark.CSV()
	csvFile := fmt.Sprintf("pubsub_topics%v_pubs%v_subs%v_msgs%v_size%v_%v.csv", numTopics, in.NumPubs, in.NumSubs, in.NumMsgs, in.MsgSize, time.Now().Format("20060102_150405"))
	os.WriteFile(csvFile, []byte(csv), 0o644)
	g.Log().Infof(ctx, "Saved metric data in csv file %s", csvFile)

	return nil
}

func runPublisher(nc *nats.Conn, subj string, startwg, donewg *sync.WaitGroup, numMsgs int, msgSize int) {
	startwg.Done()

	var msg []byte
	if msgSize > 0 {
		msg = make([]byte, msgSize)
	}

	start := time.Now()

	for range numMsgs {
		nc.Publish(subj, msg)
	}
	nc.Flush()
	benchmark.AddPubSample(bench.NewSample(numMsgs, msgSize, start, time.Now(), nc))

	donewg.Done()
}

func runSubscriber(nc *nats.Conn, subj string, startwg, donewg *sync.WaitGroup, numMsgs int, msgSize int) {

	received := 0
	ch := make(chan time.Time, 2)
	sub, _ := nc.Subscribe(subj, func(msg *nats.Msg) {
		received++
		if received == 1 {
			ch <- time.Now()
		}
		if received >= numMsgs {
			ch <- time.Now()
		}
	})
	sub.SetPendingLimits(-1, -1)
	nc.Flush()
	startwg.Done()

	start := <-ch
	end := <-ch
	benchmark.AddSubSample(bench.NewSample(numMsgs, msgSize, start, end, nc))
	nc.Close()
	donewg.Done()
}
