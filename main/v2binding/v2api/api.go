package main

import (
	"time"

	"github.com/gzjjjfree/v5-result/main/v2binding"
)

func main() {
	v2binding.StartAPIInstance(".")
	for {
		time.Sleep(time.Minute)
	}
}
