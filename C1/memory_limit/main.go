package main

import (
	"fmt"
	"strings"
	"time"
)

func main() {
	var longStrs []string
	times := 50
	for i := 1; i <= times; i++ {
		fmt.Printf("===========%d\n===========", i)
		longStrs = append(longStrs, buildString(1000000, byte(i)))
	}
	time.Sleep(3600)
}

func buildString(n int, b byte) string {
	var builder strings.Builder
	builder.Grow(n)
	for i := 0; i < n; i++ {
		builder.WriteByte(b)
	}
	return builder.String()
}
