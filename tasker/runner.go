package tasker

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"
)

type Job struct {
	Context context.Context
	Limit   int
	ChData  chan []int
	ChDone  chan bool
}

type Runner struct {
	status sync.Map
	// goroutine control
	wg sync.WaitGroup

	// generic config
	workers int

	chJob        chan Job
	shutdownCtx  context.Context
	shutdownFunc context.CancelCauseFunc

	lock sync.RWMutex

	logger *zap.SugaredLogger
}

func (r *Runner) Shutdown() {
	r.shutdownFunc(fmt.Errorf("shutting down"))
	r.wg.Wait()

}

func (r *Runner) SubmitJob(ctx context.Context, limit, workers int) []int {

	data := make(chan []int)
	done := make(chan bool, workers)

	// ensures the order
	// r.lock.Lock()
	for i := 0; i < workers; i++ {
		r.chJob <- Job{
			Context: ctx,
			Limit:   limit,
			ChData:  data,
			ChDone:  done,
		}
	}
	// r.lock.Unlock()

	nums := make([]int, 0)

	for {
		select {
		// closing due to context timeout
		case <-ctx.Done():
			close(data)
			close(done)
			return nil
			// routine check to see if jobs have finished
		case num, ok := <-data:
			if !ok {
				fmt.Println("channel closed")
			}
			nums = append(nums, num...)
		default:
			if cap(done) == len(done) {

				close(data)
				close(done)
				// r.inspectResults(done)
				return nums
			}

		}
	}

}

func (r *Runner) count(_ context.Context, x int, chData chan []int, chDone chan bool) {
	fmt.Println("count started with ", x)
	// catch closed channel panic
	defer func() {
		if recover() == nil {
			return
		}

	}()

	y := make([]int, 0)
	for i := 0; i < x; i++ {
		y = append(y, i)
	}
	chData <- y

	chDone <- true
}

func (r *Runner) consume(ID int) {
	r.status.Store(ID, true)
	for {
		select {
		case job := <-r.chJob:
			fmt.Println(ID, " received job: ", job)
			r.count(job.Context, job.Limit, job.ChData, job.ChDone)
			fmt.Println(ID, " done job: ", job)
		case <-r.shutdownCtx.Done():
			// r.status.Store(ID, false)
			fmt.Println(ID, " going down")
			r.wg.Done()
			return
			// default:
			// 	continue
		}

	}
}

func (r *Runner) start() {
	r.wg.Add(r.workers)

	for i := 0; i < r.workers; i++ {
		go r.consume(i)
	}

	// for range time.NewTicker(10 * time.Second).C {
	// 	fmt.Println("workers status")
	// 	r.status.Range(
	// 		func(key, value any) bool {
	// 			fmt.Printf("ID %v: state %v", key, value)
	// 			return true
	// 		},
	// 	)
	// 	fmt.Println("job queue length: ", len(r.chJob))
	// }

}

func NewRunner(workers int, logger *zap.SugaredLogger) *Runner {
	r := new(Runner)

	r.workers = workers
	r.logger = logger

	r.shutdownCtx, r.shutdownFunc = context.WithCancelCause(context.Background())
	r.chJob = make(chan Job, workers)

	go r.start()

	return r
}
