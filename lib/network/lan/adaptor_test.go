/**
 * @Author Kokutas
 * @Description 本机网络相关--网卡/网卡适配器 测试
 * @Date 2020/11/6 22:32
 **/
package lan

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"testing"
)

// 获取所有启用的网卡/网络适配器信息 测试
func TestAdaptors(t *testing.T) {
	adaptors, err := Adaptors()
	if err != nil {
		log.Fatal(err)
	}
	data, err := json.Marshal(adaptors)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("%s\n", data)
}

// 根据IP获取启用的网卡/网络适配器信息 测试
func TestAdaptorByIP(t *testing.T) {
	// 获取本机所在的局域网IP
	ip, err := LanIP()
	if err != nil {
		log.Fatal(err)
	}
	adaptor, err := AdaptorByIP(ip)
	if err != nil {
		log.Fatal(err)
	}
	data, err := json.Marshal(adaptor)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", data)
}

// CIDR计算可用地址个数 测试
func Test_cidrCalculateAvailable(t *testing.T) {
	num := cidrCalculateAvailable(24)
	fmt.Println(num)
}

// CIDR计算掩码 测试
func Test_cidrCalculateIPMask(t *testing.T) {
	ip := cidrCalculateIPMask(24)
	fmt.Println(ip)
}

func Test_cidrCalculateIPRange(t *testing.T) {
	fmt.Println(cidrCalculateIPRange(net.IPv4(192,168,18,1),24))
}
