package workerpool

import (
	"context"
	"github.com/cyril-jump/shortener/internal/app/dto"
	"github.com/cyril-jump/shortener/internal/app/storage"
	"log"
	"sync"
	"time"
)

type InputWorker struct {
	ch    chan dto.Task
	done  chan struct{}
	index int
	ctx   context.Context
}

type OutputWorker struct {
	ch   chan dto.Task
	done chan struct{}
	db   storage.DB
	ctx  context.Context
	mu   *sync.Mutex
}

func NewInputWorker(ch chan dto.Task, done chan struct{}, ctx context.Context) *InputWorker {
	index := 0
	return &InputWorker{
		ch:    ch,
		done:  done,
		index: index,
		ctx:   ctx,
	}
}

func NewOutputWorker(ch chan dto.Task, done chan struct{}, ctx context.Context, db storage.DB, mu *sync.Mutex) *OutputWorker {
	return &OutputWorker{
		ch:   ch,
		done: done,
		ctx:  ctx,
		db:   db,
		mu:   mu,
	}
}

func (w *InputWorker) Do(t dto.Task) {
	w.ch <- t
	w.index++
	log.Println(w.index)
	if w.index == 20 {
		w.done <- struct{}{}
		w.index = 0
	}
}

func (w *OutputWorker) Do() error {
	timer := time.NewTicker(15 * time.Second)
	models := make([]dto.Task, 0, 1)
	defer timer.Stop()
	for {
		select {
		case <-w.ctx.Done():
			return nil
		case <-w.done:
			log.Println("chReady")
			for task := range w.ch {
				models = append(models, task)
				if len(w.ch) == 0 {
					w.mu.Lock()
					if err := w.db.DelBatchShortURLs(models); err != nil {
						log.Println(err)
					}
					w.mu.Unlock()
					models = nil
					break
				}

			}
		case <-timer.C:
			log.Println("timer")
			for task := range w.ch {
				models = append(models, task)
				if len(w.ch) == 0 {
					w.mu.Lock()
					if err := w.db.DelBatchShortURLs(models); err != nil {
						log.Println(err)
					}
					w.mu.Unlock()
					models = nil
					break
				}

			}
		}
	}

}
