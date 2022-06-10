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
	ch     chan dto.Task
	done   chan struct{}
	index  int
	ticker *time.Ticker
	ctx    context.Context
}

type OutputWorker struct {
	id   int
	ch   chan dto.Task
	done chan struct{}
	db   storage.DB
	ctx  context.Context
	mu   *sync.Mutex
}

func NewInputWorker(ch chan dto.Task, done chan struct{}, ctx context.Context) *InputWorker {
	index := 0
	ticker := time.NewTicker(10 * time.Second)
	return &InputWorker{
		ch:     ch,
		done:   done,
		index:  index,
		ticker: ticker,
		ctx:    ctx,
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
			return nil
		case <-w.ticker.C:
			log.Println("timer")
			w.done <- struct{}{}
			w.index = 0
		}
	}
}

func (w *OutputWorker) Do() error {
	models := make([]dto.Task, 0, 20)
	for {
		select {
		case <-w.ctx.Done():
			return nil
		case <-w.done:
			if len(w.ch) == 0 {
				break
			}
			log.Println("chReady")
			for task := range w.ch {
				log.Println(w.id, "id out worker")
				models = append(models, task)
				if len(w.ch) == 0 {
					if err := w.db.DelBatchShortURLs(models); err != nil {
						log.Println(err, "error del")
					}
					models = nil
					break
				}
			}
		}
	}
}
