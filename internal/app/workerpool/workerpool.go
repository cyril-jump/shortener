package workerpool

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/cyril-jump/shortener/internal/app/dto"
	"github.com/cyril-jump/shortener/internal/app/storage"
)

type InputWorker struct {
	ch     chan dto.Task
	done   chan struct{}
	index  int
	ticker *time.Ticker
	ctx    context.Context
	mu     *sync.Mutex
}

type OutputWorker struct {
	id   int
	ch   chan dto.Task
	done chan struct{}
	db   storage.DB
	ctx  context.Context
	mu   *sync.Mutex
}

func NewInputWorker(ch chan dto.Task, done chan struct{}, ctx context.Context, mu *sync.Mutex) *InputWorker {
	index := 0
	ticker := time.NewTicker(10 * time.Second)
	return &InputWorker{
		ch:     ch,
		done:   done,
		index:  index,
		ticker: ticker,
		ctx:    ctx,
		mu:     mu,
	}
}

func NewOutputWorker(id int, ch chan dto.Task, done chan struct{}, ctx context.Context, db storage.DB, mu *sync.Mutex) *OutputWorker {
	return &OutputWorker{
		id:   id,
		ch:   ch,
		done: done,
		ctx:  ctx,
		db:   db,
		mu:   mu,
	}
}

func (w *InputWorker) Do(t dto.Task) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.ch <- t
	w.index++
	log.Println(w.index)
	if w.index == 20 {
		w.done <- struct{}{}
		w.index = 0
	}
}

func (w *InputWorker) Loop() error {
	for {
		select {
		case <-w.ctx.Done():
			w.ticker.Stop()
			return nil
		case <-w.ticker.C:
			w.mu.Lock()
			w.done <- struct{}{}
			w.index = 0
			w.mu.Unlock()
		}
	}
}

func (w *OutputWorker) Do() error {
	models := make([]dto.Task, 0, 200)
	for {
		select {
		case <-w.ctx.Done():
			return nil
		case <-w.done:
			if len(w.ch) == 0 {
				break
			}
			for task := range w.ch {
				models = append(models, task)
				if len(w.ch) == 0 {
					if err := w.db.DelBatchShortURLs(models); err != nil {
						return err
					}
					models = nil
					break
				}
			}
		}
	}
}
