package workerpool

import "context"

func newWorker[T any](ctx context.Context, job Job[T], tasks <-chan *Task[T]) *worker[T] {
	return &worker[T]{
		ctx,
		job,
		tasks,
		nil,
		nil,
	}
}

type worker[T any] struct {
	ctx context.Context
	job Job[T]
	tasks <-chan *Task[T]
	done chan struct{}
	stop chan struct{}
}

// Start - запуск в воркера на обработку. Не потокобезопасно TODO возможно надо сделать потокобезопасно
func (w *worker[T]) Start() error {
	if w.done != nil{
		return AlreadyRunningError
	}
	w.done = make(chan struct{})
	w.stop = make(chan struct{}) // закрытие канала в методе "Stop"
	go func(){
		defer close(w.done)
		for{
			if !w.attempt(){
				break
			}
		}
	}()
	return nil
}

func (w *worker[T]) attempt() bool{
	select {
	case task, ok := <-w.tasks:
		if !ok {
			return false
		}
		w.job(w.ctx, task.data) // TODО обработка ошибок
		return true
	case <-w.stop:
	case <-w.ctx.Done():
	}
	return false
}

func (w *worker[T]) Stop(){
	defer func(){
		w.done = nil
	}()
	close(w.stop)
	select{
	case <-w.done:
	case <-w.ctx.Done():
	}
}