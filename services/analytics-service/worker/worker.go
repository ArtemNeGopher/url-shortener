// Package worker
package worker

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/ArtemNeGopher/url-shortener/services/analytics-service/models"
)

type Repository interface {
	BatchInsertClicks(events []models.ClickEvent) error
	UpdateStats(shortCode string) error
	GetStats(shortCode string) (*models.Stats, error)
}

type WorkerPool struct {
	jobQueue    chan models.ClickEvent
	batchQueue  chan []models.ClickEvent
	stopCtx     context.Context
	stopFn      context.CancelFunc
	wg          *sync.WaitGroup
	workerCount int
	repo        Repository
	batchSize   int
	log         *slog.Logger
}

func New(workers int, batchSize int, repo Repository, log *slog.Logger) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	return &WorkerPool{
		jobQueue:    make(chan models.ClickEvent, batchSize),
		batchQueue:  make(chan []models.ClickEvent),
		stopCtx:     ctx,
		stopFn:      cancel,
		wg:          &sync.WaitGroup{},
		workerCount: workers,
		repo:        repo,
		batchSize:   batchSize,
		log:         log,
	}
}

func (w *WorkerPool) Submit(event models.ClickEvent) {
	w.jobQueue <- event
}

func (w *WorkerPool) Start() {
	for range w.workerCount {
		w.wg.Add(1)
		go w.worker()
	}

	go func() {
		stop := false
		var job models.ClickEvent
		batch := make([]models.ClickEvent, w.batchSize)

		timer := time.NewTimer(1 * time.Second)
		defer timer.Stop()

		i := 0
		for !stop {
			select {
			case job = <-w.jobQueue:
				batch[i] = job
				if i++; i >= w.batchSize {
					// Делаем копию
					local := make([]models.ClickEvent, i)
					copy(batch[:i], local)
					w.batchQueue <- local

					// Сбрасываем счётчик
					i = 0
				}
				timer.Reset(1 * time.Second)
			case <-timer.C:
				if i != 0 {
					// Делаем копию
					local := make([]models.ClickEvent, i)
					copy(batch[:i], local)
					w.batchQueue <- local

					// Сбрасываем счётчик
					i = 0
				}
			case <-w.stopCtx.Done():
				stop = true
			}
		}
	}()
}

func (w *WorkerPool) Stop() {
	w.stopFn()
	w.wg.Wait()
}

func (w *WorkerPool) worker() {
	defer w.wg.Done()

	for {
		var batch []models.ClickEvent
		select {
		case batch = <-w.batchQueue:
		case <-w.stopCtx.Done():
			return
		}
		if err := w.repo.BatchInsertClicks(batch); err != nil {
			w.log.Error("failed to insert batch", "error", err)
		}
	}
}
