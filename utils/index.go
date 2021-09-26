package utils

import "sync/atomic"

// IndexI64 ...
type IndexI64 struct {
	i int64
}

// NewIndexI64 ...
func NewIndexI64() *IndexI64 {
	return &IndexI64{i: -1}
}

// NewIndex ...
func (i *IndexI64) NewIndex() int64 {
	return atomic.AddInt64(&i.i, 1)
}
