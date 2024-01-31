package main

import (
	"cctv/helper"
	"fmt"
	"time"
)

func main() {
	var i = 0
	for {
		helper.ReadChunk("https://mam.jogjaprov.go.id:1937/cctv-public/ViewParangtritis.stream/", fmt.Sprintf("chunk-%d.mp4", i))
		time.Sleep(4 * time.Second)
		i++
	}
}
