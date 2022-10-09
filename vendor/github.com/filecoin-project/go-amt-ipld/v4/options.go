package amt

import (
	"fmt"
)

var defaultBitWidth = uint(3)

type config struct {
	bitWidth uint
}

type Option func(*config) error

func UseTreeBitWidth(bitWidth uint) Option {
	return func(c *config) error {
		if bitWidth < 1 {
			return fmt.Errorf("bit width must be at least 1 (i.e. 2 children per node), is %d", bitWidth)
		}
		c.bitWidth = bitWidth
		return nil
	}
}

func defaultConfig() *config {
	return &config{
		bitWidth: defaultBitWidth,
	}
}
