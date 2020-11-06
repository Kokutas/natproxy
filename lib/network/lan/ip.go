/**
 * @Author Kokutas
 * @Description 本机网络相关--本机局域网IP
 * @Date 2020/11/6 22:39
 **/
package lan

import (
	"net"
)
// 获取本机所在的局域网IP
func LanIP() (net.IP,error) {
	conn,err:=net.Dial("udp","google.com:80")
	if err != nil {
		return nil, err
	}
	ipAddr,err:=net.ResolveUDPAddr("udp",conn.LocalAddr().String())
	if err!=nil{
		return nil,err
	}
	_ = conn.Close()
	return ipAddr.IP,nil
}