package pgoutbox

import (
	"context"
	"fmt"
)

type outbox struct {
	store  *store
	worker *worker
}

func (o *outbox) Start(ctx context.Context) {
	o.worker.Start(ctx)
}

func (o *outbox) Send(ctx context.Context, tx Tx, msg *Message) error {
	err := o.store.CreateTx(ctx, tx, msg)
	if err != nil {
		return fmt.Errorf("save msg: %w", err)
	}
	return nil
}

func (o *outbox) Flush(ctx context.Context, msg *Message) {
	ctx = context.WithoutCancel(ctx)
	go o.worker.Send(ctx, msg)
}
