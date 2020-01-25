package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	startTime := time.Now()
	var slc []int
	for len(slc) <= 10010 {
		slc = append(slc, rand.Intn(100))
	}

	// var r int32
	// for index := 0; index < len(slc); index++ {
	// 	r += int32(slc[index])
	// }
	// fmt.Println(r)

	deli := 40
	resultChan := make(chan int32, deli)
	var result int32
	var wg = &sync.WaitGroup{}
	i := len(slc) / deli

	for index := 0; index < deli; index++ {
		wg.Add(1)
		// fmt.Printf("%d,%d:%d", i, i*index, (index+1)*i)
		// fmt.Println(slc[i*index : (index+1)*i])
		s := i * index
		e := i * (index + 1)

		var slice []int
		if e == len(slc)-len(slc)%deli {
			slice = slc[s:]
		} else {
			slice = slc[s:e]
		}
		go cal(slice, resultChan, wg)
	}

	wg.Wait()
	close(resultChan)
	for v := range resultChan {
		result += v
	}
	fmt.Println(result)
	endTime := time.Now()
	fmt.Println(deli, "gorotine, total time: ", endTime.Sub(startTime).Seconds())
}

func cal(slc []int, chn chan int32, wg *sync.WaitGroup) {
	var result int32
	for _, v := range slc {
		result += int32(v)
	}
	chn <- result
	wg.Done()
}
