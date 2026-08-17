package region

import (
	"fmt"
	"hash/fnv"
)

func CalcRegionCacheKey(values ...any) uint32 {
	h := fnv.New32a()

	for _, v := range values {
		fmt.Fprint(h, v)
		h.Write([]byte{0})
	}

	return h.Sum32()
}
