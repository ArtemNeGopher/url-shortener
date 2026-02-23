// Package worker
package worker

import (
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
	stopChan    chan struct{}
	wg          *sync.WaitGroup
	workerCount int
	repo        Repository
	batchSize   int
}

func New(workers int, batchSize int, repo Repository) *WorkerPool {
	return &WorkerPool{
		jobQueue:    make(chan models.ClickEvent, batchSize*2), // Очередь в два раза больше батча
		stopChan:    make(chan struct{}),
		wg:          &sync.WaitGroup{},
		workerCount: workers,
		repo:        repo,
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
}

func (w *WorkerPool) Stop() {
	for range w.workerCount {
		w.stopChan <- struct{}{}
	}
	w.wg.Wait()
}

func (w *WorkerPool) worker() {
	defer w.wg.Done()

	batch := make([]models.ClickEvent, 0, w.batchSize)
	stop := false
	batchByTime := false
	for !stop {
		select {
		case job := <-w.jobQueue:
			batch = append(batch, job)
		case <-w.stopChan:
			stop = true
		case <-time.After(1 * time.Second): // Если секунду нет событий, то отправляем так
			batchByTime = true
		}

		if stop || len(batch) < w.batchSize || batchByTime {
			if err := w.repo.BatchInsertClicks(batch); err != nil {
				// TODO: Логирование
			}
			for _, event := range batch {
				if err := w.repo.UpdateStats(event.ShortCode); err != nil {
					// TODO: Логирование
				}
			}
			// Обнуляем батч
			batch = batch[:0]
			batchByTime = false
		}
	}
}
