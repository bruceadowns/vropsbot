package common

// BiMap encapsulates a bidirectional map
type BiMap struct {
	kv map[string]string
	vk map[string]string
}

// NewBiMap returns an initialized BiMap
func NewBiMap() *BiMap {
	return &BiMap{
		kv: make(map[string]string),
		vk: make(map[string]string),
	}
}

// Put adds the k,v pair into each map
func (bm *BiMap) Put(k, v string) *BiMap {
	bm.kv[k] = v
	bm.vk[v] = k
	return bm
}

// GetByKey returns a value given a key
func (bm *BiMap) GetByKey(k string) (v string, e bool) {
	v, e = bm.kv[k]
	return
}

// GetByValue returns a key given a value
func (bm *BiMap) GetByValue(v string) (k string, e bool) {
	k, e = bm.vk[v]
	return
}
