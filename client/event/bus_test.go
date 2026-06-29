package event

import (
	"sync"
	"testing"
)

func TestNewBus(t *testing.T) {
	b := NewBus()
	if b == nil {
		t.Fatal("expected non-nil bus")
	}
}

func TestSubscribeAndPublish(t *testing.T) {
	b := NewBus()
	var got Event
	b.Subscribe("foo", func(e Event) { got = e })
	b.Publish(Event{Topic: "foo", Data: 42})
	if got.Topic != "foo" {
		t.Fatalf("expected topic foo, got %q", got.Topic)
	}
	if got.Data != 42 {
		t.Fatalf("expected data 42, got %v", got.Data)
	}
}

func TestPublishNoSubscribers(t *testing.T) {
	b := NewBus()
	// must not panic
	b.Publish(Event{Topic: "unknown"})
}

func TestMultipleSubscribers(t *testing.T) {
	b := NewBus()
	count := 0
	b.Subscribe("tick", func(_ Event) { count++ })
	b.Subscribe("tick", func(_ Event) { count++ })
	b.Subscribe("tick", func(_ Event) { count++ })
	b.Publish(Event{Topic: "tick"})
	if count != 3 {
		t.Fatalf("expected 3, got %d", count)
	}
}

func TestSubscribersIsolatedByTopic(t *testing.T) {
	b := NewBus()
	aCount, bCount := 0, 0
	b.Subscribe("a", func(_ Event) { aCount++ })
	b.Subscribe("b", func(_ Event) { bCount++ })
	b.Publish(Event{Topic: "a"})
	if aCount != 1 {
		t.Fatalf("expected aCount 1, got %d", aCount)
	}
	if bCount != 0 {
		t.Fatalf("expected bCount 0, got %d", bCount)
	}
}

func TestConcurrentPublish(t *testing.T) {
	b := NewBus()
	var mu sync.Mutex
	count := 0
	b.Subscribe("ping", func(_ Event) {
		mu.Lock()
		count++
		mu.Unlock()
	})
	var wg sync.WaitGroup
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			b.Publish(Event{Topic: "ping"})
		}()
	}
	wg.Wait()
	if count != 100 {
		t.Fatalf("expected 100, got %d", count)
	}
}

func TestConcurrentSubscribeAndPublish(_ *testing.T) {
	b := NewBus()
	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			b.Subscribe("x", func(_ Event) {})
		}()
	}
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			b.Publish(Event{Topic: "x"})
		}()
	}
	wg.Wait()
}
