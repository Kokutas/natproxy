/**
 * @Author Kokutas
 * @Description 本机网络相关--本机局域网IP 测试
 * @Date 2020/11/6 22:46
 **/
package lan

import (
	"fmt"
	"log"
	"testing"
)
// 获取本机所在的局域网IP 测试
func TestLanIP(t *testing.T) {
	ip,err:=LanIP()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ip)
}
