package test

import (
	"fmt"
	"testing"
)

func TestSprintf(t *testing.T) {
	ranges := fmt.Sprintf("bytes=%d-%d", 0, 3)
	fmt.Println(ranges)
}
