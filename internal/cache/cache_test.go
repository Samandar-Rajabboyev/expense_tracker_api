package cache

import (
	"context"
	"testing"
	"time"
)

func TestCache_SetGet(t *testing.T) {
	c := NewCache("localhost:6379")
	ctx := context.Background()

	err := c.Set(ctx, "test_key", "hello", 5*time.Second)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	val, err := c.Get(ctx, "test_key")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if val != "hello" {
		t.Errorf("expected 'hello', got %q", val)
	}

	c.Delete(ctx, "test_key")
}
