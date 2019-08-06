package onlygreen

import (
	"github.com/dianelooney/gvd/filters"
)

func New() filters.Interface {
	return &filter{}
}

type filter struct{}

func (*filter) Apply(img []uint8) {
	for i := range img {
		if i%4 == 1 {
			img[i] = 0
		}
	}
}
