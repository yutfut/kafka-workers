package manager

import (
	"context"
	"kafka-workers/internal/worker"
	"log"
	"time"
)

type Manager interface {
	ManagerWorker(ctx context.Context)
}

type manager struct {
	workers     []worker.Worker
	add         int
	del         int
	input       chan int
	output      chan int
	lastBuffLen int
}

func NewManager(
	input chan int,
	output chan int,
) Manager {
	return &manager{
		workers:     make([]worker.Worker, 10),
		add:         0,
		del:         0,
		input:       input,
		output:      output,
		lastBuffLen: 0,
	}
}

func (m *manager) ManagerWorker(
	ctx context.Context,
) {
	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			log.Printf("--//-- ManagerWorker del: %d add: %d", m.del, m.add)
			log.Printf("lastBuffLen: %d, len(input): %d", m.lastBuffLen, len(m.input))

			if len(m.input) > m.lastBuffLen {
				newWorker := worker.NewWorker(
					m.input,
					m.output,
				)

				m.workers[m.add] = newWorker
				go m.workers[m.add].Work()
				m.add++
				if len(m.workers) == m.add {
					m.add = 0
				}
			} else if len(m.input) < m.lastBuffLen {
				m.workers[m.del].Done()
				m.workers[m.del] = nil
				m.del++
				if len(m.workers) == m.del {
					m.del = 0
				}
			} else {
				log.Printf("equal")
			}

			m.lastBuffLen = len(m.input)
		}
	}
}
