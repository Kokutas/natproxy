/**
 * @Author Kokutas
 * @Description 本机网络相关--本机互联网IP信息
 * @Date 2020/11/7 13:21
 **/
package wan

import (
	"encoding/json"
	"net"
	"net/http"
)

// 查询结果结构体
// 采用接口： 进行查询
type Query struct {
	Success int     `json:"success,string"`
	Result  *Result `json:"result"`
}
type Result struct {
	IP        net.IP `json:"ip,string"`
	Proxy     string `json:"proxy"`
	Att       string `json:"att"`
	Operators string `json:"operators"`
}

// 查询公网IP
func WanIP() (*Query, error) {
	rsp, err := http.Get("http://api.k780.com/?app=ip.local&format=json")
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	var query Query
	if err := json.NewDecoder(rsp.Body).Decode(&query); err != nil {
		return nil, err
	}
	return &query, nil
}
