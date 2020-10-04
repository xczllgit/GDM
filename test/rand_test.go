package test

import (
	"crypto/md5"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestRandFileName(t *testing.T) {
	rand.Seed(time.Now().Unix())
	str := "abcdefghijklmnopqrstuvwxyz"
	s1 := rand.Int63n(26)
	s2 := rand.Int63n(26)
	s3 := rand.Int63n(26)
	s4 := rand.Int63n(26)
	s5 := rand.Int63n(26)
	tempString := str[s1:s1+1] + str[s2:s2+1] + str[s3:s3+1] + str[s4:s4+1] + str[s5:s5+1]
	hash := md5.New()
	io.WriteString(hash, tempString)
	encrypt, _ := fmt.Printf("%x", hash.Sum(nil))
	fmt.Println(strconv.FormatInt(int64(encrypt), 10))
}

func TestRandString(t *testing.T) {
	str := "abcdef"
	//start := rand.Int63n(6)
	temp := str[5:6]
	fmt.Println(temp)
}
