package pigcache

import (
	"fmt"
	"strings"
)

//var m sync.Mutex
//var set = make(map[int]bool, 0)
//
//func printOnce(num int) {
//	m.Lock()
//	defer m.Unlock()
//	if _, exist := set[num]; !exist {
//		println(num)
//	}
//	set[num] = true
//
//}
//
//func main()  {
//	for i :=0; i < 10; i++ {
//		go printOnce(100)
//	}
//	time.Sleep(time.Second)
//}

func fib(c, quit chan int)  {
	x, y := 1, 1
	for {
		select {
		case c <- x:
			x, y = y, x + y
		case <- quit:
			fmt.Println("quit")
			return
		}
	}
}
func main()  {
	//c := make(chan int)
	//o := make(chan bool)
	//go func() {
	//	for {
	//		select {
	//		case v := <-c:
	//			println(v)
	//		case <- time.After(5 * time.Second):
	//			println("Timeout")
	//			o <- true
	//			break
	//		}
	//	}
	//}()
	//<- o
	s := "hello world hello world"
	idx := strings.SplitN(s, " ", 2)
	println(idx[0])
	println(idx[1])
	println(idx[2])
	fmt.Println(len(idx))
}