package utility

import "sync"

type BytesPool struct {
	pool sync.Pool
}

func (p *BytesPool) Get(minCap int) []byte {
	v := p.pool.Get()
	if v == nil {
		return make([]byte, 0, minCap)
	}
	bs := v.([]byte)
	if cap(bs) < minCap {
		bs = make([]byte, 0, minCap)
	}
	return bs
}

func (p *BytesPool) Put(v []byte) {
	if v != nil {
		p.pool.Put(v[:0])
	}
}
