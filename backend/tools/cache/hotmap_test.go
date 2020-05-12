package cache

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHotMap_Has(t *testing.T) {
	hm := NewHotMap(100)
	k := 1
	hm.Push(k, 1)
	has := hm.Has(k)
	require.Equal(t, true, has)
	require.Equal(t, false, hm.Has(2))
}
func TestHotMap_Push(t *testing.T) {
	hm := NewHotMap(10)
	for i := 1; i <= 11; i++ {
		hm.Push(i, i)
	}
	if len(hm.k) > 10 {
		t.Fail()
	}
	t.Log(hm.k)
}
