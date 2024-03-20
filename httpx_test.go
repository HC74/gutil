package gutil

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var httpClient *HttpClient

func init() {
	httpClient = NewHttpClient()
}

func TestHttpGet(t *testing.T) {
	core := 999999

	wg := sync.WaitGroup{}
	wg.Add(core)
	for i := 0; i < core; i++ {
		fmt.Printf("%v 发送 \n", i)
		time.Sleep(100)
		go func() {
			n, _ := httpClient.GetStringN("http://10.12.200.111:31001/usercenter/v1/user/test")
			fmt.Printf("%v 接受 \n", n)
		}()

	}
	wg.Wait()
	fmt.Println("ok")
}
