package main

import (
	"context"
	"testing"
	"time"
)

func TestSync(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := Sync(ctx); err != nil {
		t.Error(err)
	}
}
