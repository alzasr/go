package pgoutbox

import (
	"context"
	"errors"
	"fmt"
)

type worker struct {
	logger Logger
	locker Locker
	sender Sender
	store  store
}

func (w *worker) Start(ctx context.Context) {
	go func() {
		for {
			if ctx.Err() != nil {
				return
			}
			err := w.sendRemaining(ctx)
			if err != nil {
				w.logger.Error(err)
			}
		}
	}()
}

func (w *worker) Send(ctx context.Context, msg *Message) (err error) {
	locked, unlock, err := w.tryLockMsg(ctx, msg)
	if err != nil {
		return fmt.Errorf("try lock msg: %w", err)
	}
	if !locked {
		return nil
	}
	defer unlock()

	msgs, err := w.listMessageChain(ctx, msg)
	if err != nil {
		return fmt.Errorf("list message chain: %w", err)
	}

	err = w.sender.Send(ctx, msgs)
	if err != nil {
		saveErr := w.saveSendError(ctx, msgs, err)
		return errors.Join(fmt.Errorf("send message: %w", err), saveErr)
	}

	err = w.store.DeleteBatch(ctx, msgs)
	if err != nil {
		return fmt.Errorf("delete messages: %w", err)
	}
	return nil
}

func (w *worker) saveSendError(ctx context.Context, msgs []*Message, sendErr error) error {
	for i := range msgs {
		msgs[i].Attempt++
		msgs[i].LastError = sendErr.Error()
	}
	err := w.store.UpdateBatch(ctx, msgs)
	if err != nil {
		return fmt.Errorf("update messages: %w", err)
	}
	return nil
}

func (w *worker) listMessageChain(ctx context.Context, msg *Message) ([]*Message, error) {
	if len(msg.Key) == 0 {
		return []*Message{msg}, nil
	}
	msgs, err := w.store.ListByKey(ctx, msg.Key)
	if err != nil {
		return nil, fmt.Errorf("list messages by key: %w", err)
	}
	return msgs, nil
}

func (w *worker) tryLockMsg(ctx context.Context, msg *Message) (locked bool, unlock func() error, err error) {
	if len(msg.Key) != 0 {
		locked, unlock, err = w.locker.TryLock(ctx, msg.Key)
		if err != nil {
			return false, nil, fmt.Errorf("try lock: %w", err)
		}
		if !locked {
			return false, nil, nil
		}
		return true, unlock, nil

	}
	// TODO lock message without key by ID
	return true, noop, nil
}

func (w *worker) sendRemaining(ctx context.Context) error {
	ids, err := w.store.ListIDsGroupingByKey(ctx)
	if err != nil {
		return fmt.Errorf("list ids grouping by key: %w", err)
	}
	for i := range ids {
		msg, err := w.store.Get(ctx, ids[i])
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				continue
			}
			return fmt.Errorf("get: %w", err)
		}
		err = w.Send(ctx, msg)
		if err != nil {
			return fmt.Errorf("send: %w", err)
		}
	}
	return nil
}
