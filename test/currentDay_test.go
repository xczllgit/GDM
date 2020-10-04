package test

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestGetCurrentDay(t *testing.T) {
	now := time.Now()
	fmt.Println(now)
	year := time.Now().Year()
	month := time.Now().Format("01")
	day := time.Now().Day()
	fmt.Println("Year: ", year, " Month: ", month, " Day: ", day)
	fmt.Println(month)
}

func TestGetCurrentDay2(t *testing.T) {
	now := time.Now()
	fmt.Println(now)
	year, month, day := time.Now().Date()
	fmt.Println("Year: ", year, " Month: ", strconv.FormatInt(int64(month), 10), " Day: ", day)
}
