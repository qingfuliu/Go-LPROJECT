package limiter

import (
	"MFile/db/mredis"
	"fmt"
	"testing"
	"time"
)

func TestRedis_TokenLimiter(t *testing.T) {
	limiter:=NewTokenLimiter(mredis.RedisDb,"testLimiter",1,10)
	ok:=limiter.AllowN(10)
	if !ok{
		t.Fatal("请求token失败")
	}else{
		fmt.Println("请求成功")
	}

	//ok=limiter.AllowN(10)
	//if !ok{
	//	t.Fatal("请求token失败")
	//}else{
	//	t.Log("请求成功")
	//}
	//ok=limiter.AllowN(10)
	//if ok{
	//	t.Fatal("逻辑错误")
	//}else{
	//	t.Log("请求成功")
	//}
	ch:=make(chan struct{})
	go func(){
		triker:=time.NewTicker(time.Second)
		defer  triker.Stop()
		Loop:
		for{
			select{
			case <-triker.C :
				ok:=limiter.AllowN(2)
				fmt.Println(ok)
			case <-ch:
				break Loop
			}
		}
		fmt.Println("break")
	}()

	time.Sleep(10*time.Second)
	close(ch)

	ok=limiter.AllowN(10)
	if !ok{
		t.Fatal("逻辑错误")
	}else{
		t.Log("请求成功")
	}

}
