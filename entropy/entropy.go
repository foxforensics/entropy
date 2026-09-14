// Package entropy provides methods to calculate entropy and average of block.
//
// Source: https://gist.github.com/n2p5/4eda328b080c9f09eff928ad47228ab1
package entropy

import "math"

// Calculate returns the entropy and average of the given block.
func Calculate(block []byte) (byte, float64) {
	var n = float64(len(block))
	var m = uint64(len(block))
	var a [256]float64
	var e float64
	var i uint64

	if len(block) == 0 {
		return 0, 0.0
	}

	for _, b := range block {
		i += uint64(b)
		a[b]++
	}

	for i := range 256 {
		if a[i] != 0 {
			v := a[i] / n
			e -= v * math.Log2(v)
		}
	}

	return byte(i / m), e
}
