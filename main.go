package main

import (
	"cctv/helper"
	"time"
)

func main() {
	var i = 0

	helper.InitFolder()

	for {
		helper.ReadChunk()

		time.Sleep(2 * time.Minute)
		i++
	}
}
