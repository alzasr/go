package workerpool

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWoker(t *testing.T) {
	t.Run("smoke", func(t *testing.T) {
		attemptCount := 100
		result := make([]bool, attemptCount)
		job := func(ctx context.Context, data int) error {
			result[data] = true
			return nil
		}
		var tasks = make(chan *Task[int])
		w := newWorker(
			t.Context(),
			job,
			tasks,
		)
		w.Start()
		for i := range result {
			tasks <- &Task[int]{
				data: i,
			}
		}
		w.Stop()
		assert.NotContains(t, result, false)	       		           
	})
}