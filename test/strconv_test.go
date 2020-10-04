package test

import (
	"fmt"
	"strconv"
	"testing"
)

func TestStrConv(t *testing.T) {
	string := strconv.FormatFloat(12.32424, 'f', 2, 64)
	fmt.Println(string)
}
