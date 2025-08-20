package reflectx

import (
	"sync"
	"testing"
)

func TestPool(t *testing.T) {
	if err := RegisterTypedPool[float64](&sync.Pool{
		New: func() any { return new(float64) },
	}); err != nil {
		t.Fatal(err)
	}
}
