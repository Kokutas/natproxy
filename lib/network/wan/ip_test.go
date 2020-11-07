/**
 * @Author Kokutas
 * @Description //TODO
 * @Date 2020/11/7 13:25
 **/
package wan

import (
	"fmt"
	"log"
	"testing"
)

func TestWanIP(t *testing.T) {
	var qy *Query
	qy,err:=WanIP()
	if err != nil {
		log.Fatal(err)
	}
	if qy.Success == 1{
		fmt.Println(qy.Result)
	}
}
