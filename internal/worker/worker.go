package worker

import (
	"context"
	"time"
)

type Worker interface {
	Work()
	Done()
}

type worker struct {
	input  chan int
	output chan int
	ctx    context.Context
}

func NewWorker(
	input chan int,
	output chan int,
) Worker {
	return &worker{
		input:  input,
		output: output,
		ctx:    context.Background(),
	}
}

func (w *worker) Work() {
	for {
		select {
		case <-w.ctx.Done():
			return
		default:
			buff := <-w.input

			time.Sleep(time.Second)

			w.output <- buff
		}
	}
}

func (w *worker) Done() {
	w.ctx.Done()
}
