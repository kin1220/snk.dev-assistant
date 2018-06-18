//http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang

package main

import (
	"fmt"
	"math/rand"
	"time"
)

var num = []rune("0123456789")
var lenNum = len(num)

var alpha = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var lenAlpha = len(alpha)

func randSeq(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = num[rand.Intn(lenNum)]
	}
	return string(b)
}

func randAlpha(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = alpha[rand.Intn(lenAlpha)]
	}
	return string(b)
}

func main() {
	fmt.Println(randSeq(7))
	fmt.Println(randAlpha(100))
}
