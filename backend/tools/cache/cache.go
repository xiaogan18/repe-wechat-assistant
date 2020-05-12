package cache

type ICache interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}) error
	Remove(key string)
	Clear()
}

type CacherType uint8

const (
	CacherMemory CacherType = iota
)

func NewCacher() ICache {
	var cacher ICache
	switch cacherContext.t {
	case CacherMemory:
		cacher = &memoryCache{c: make(map[string]interface{})}
	default:
		panic("not implement yet")
	}
	cacherContext.s = append(cacherContext.s, cacher)
	return cacher
}

type cacheStorage struct {
	s []ICache
	t CacherType
}

var cacherContext cacheStorage

func init() {
	cacherContext = cacheStorage{
		s: make([]ICache, 0),
		t: CacherMemory,
	}
}
