package application

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/bench"
)

var benchmark *bench.Benchmark

func (h *Application) PubSubBenchmark(ctx context.Context, in *PubSubBenchmarkInput) error {
	benchmark := bench.NewBenchmark("NATS", in.NumSubs, in.NumPubs)
	// Connect Options.
	opts := []nats.Option{nats.Name("NATS Benchmark")}
	var startwg sync.WaitGroup
	var donewg sync.WaitGroup

	donewg.Add(in.NumPubs + in.NumSubs)

	// Run Subscribers first
	startwg.Add(in.NumSubs)
	for i := 0; i < in.NumSubs; i++ {
		nc, err := nats.Connect(g.Cfg().MustGet(ctx, "nats.serverAddr").String(), opts...)
		if err != nil {
			log.Fatalf("Can't connect: %v\n", err)
		}
		defer nc.Close()

		go runSubscriber(nc, &startwg, &donewg, in.NumMsgs, in.MsgSize)
	}
	startwg.Wait()

	// Now Publishers
	startwg.Add(in.NumPubs)
	pubCounts := bench.MsgsPerClient(in.NumMsgs, in.NumPubs)
	for i := 0; i < in.NumPubs; i++ {
		nc, err := nats.Connect(g.Cfg().MustGet(ctx, "nats.serverAddr").String(), opts...)
		if err != nil {
			log.Fatalf("Can't connect: %v\n", err)
		}
		defer nc.Close()

		go runPublisher(nc, &startwg, &donewg, pubCounts[i], in.MsgSize)
	}

	log.Printf("Starting benchmark [msgs=%d, msgsize=%d, pubs=%d, subs=%d]\n", in.NumMsgs, in.MsgSize, in.NumPubs, in.NumSubs)

	startwg.Wait()
	donewg.Wait()

	benchmark.Close()

	fmt.Print(benchmark.Report())

	csv := benchmark.CSV()
	csvFile := "pubsub_nats_benchmark_" + time.Now().Format("2006-01-02_150405") + ".csv"
	os.WriteFile(csvFile, []byte(csv), 0o644)
	fmt.Printf("Saved metric data in csv file %s\n", csvFile)

	return nil
}

func runPublisher(nc *nats.Conn, startwg, donewg *sync.WaitGroup, numMsgs int, msgSize int) {
	startwg.Done()

	args := flag.Args()
	subj := args[0]
	var msg []byte
	if msgSize > 0 {
		msg = make([]byte, msgSize)
	}

	start := time.Now()

	for i := 0; i < numMsgs; i++ {
		nc.Publish(subj, msg)
	}
	nc.Flush()
	benchmark.AddPubSample(bench.NewSample(numMsgs, msgSize, start, time.Now(), nc))

	donewg.Done()
}

func runSubscriber(nc *nats.Conn, startwg, donewg *sync.WaitGroup, numMsgs int, msgSize int) {
	args := flag.Args()
	subj := args[0]

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
