package cache

func NewHotMap(capacity int) *HotMap {
	return &HotMap{
		c: capacity,
		m: make(map[interface{}]interface{}),
		k: make([]interface{}, 0),
	}
}

type HotMap struct {
	c int
	m map[interface{}]interface{}
	k []interface{}
}

func (t *HotMap) Push(k, v interface{}) {
	if t.Has(k) {
		t.m[k] = v
		return
	}
	if len(t.k) >= t.c {
		old := t.k[0]
		delete(t.m, old)
		t.k = append(t.k[1:], k)
		t.m[k] = v
	} else {
		t.k = append(t.k, k)
		t.m[k] = v
	}
}
func (t *HotMap) Has(k interface{}) bool {
	_, b := t.m[k]
	return b
}
func (t *HotMap) Get(k interface{}) interface{} {
	return t.m[k]
}
func (t *HotMap) Remove(k interface{}) {
	delete(t.m, k)
}
