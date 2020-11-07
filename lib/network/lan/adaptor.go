/**
 * @Author Kokutas
 * @Description 本机网络相关--网卡/网卡适配器
 * @Date 2020/11/6 22:29
 **/
package lan

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// 网卡/网络适配器信息结构体
type Adaptor struct {
	// 网卡/网络适配器序号
	Index int `json:"index"`
	// 网卡/适配器名称
	Name string `json:"name"`
	// 网卡/网络适配器最大传输单元MTU
	Mtu int `json:"mtu"`
	// 网卡/网络适配器MAC地址
	Mac string `json:"mac"`
	// 网卡/网络适配器Flags标记
	Flags string `json:"flags"`
	// 网卡/网络适配器IPv4
	IPv4 net.IP `json:"ipv4"`
	// 网卡/网络适配器IPv4所在网络
	IPv4Network net.IP `json:"ipv4_network"`
	// 网卡/网络适配器IPv4掩码位/子网掩码
	IPv4SubMask int `json:"ipv4_sub_mask"`
	// 网卡/网络适配器IPv4掩码总位数
	IPv4MaskBits int `json:"ipv4_mask_bits"`
	// ======通过CIDR计算-开始======
	// 网卡/网络适配器IPv4的可用地址个数
	IPv4Available uint `json:"ipv4_available"`
	// 网卡/网络适配器IPv4的掩码地址
	IPv4Mask net.IP `json:"ipv4_mask"`
	// 网卡/网络适配器IPv4第一个可用地址
	IPv4FirstAvailable net.IP `json:"ipv4_first_available"`
	// 网卡/网络适配器IPv4最后一个可用地址
	IPv4LastAvailable net.IP `json:"ipv4_last_available"`
	// 网卡/网络适配器IPv4广播地址
	IPv4Broadcast net.IP `json:"ipv4_broadcast"`
	// ======通过CIDR计算-结束======
	// 网卡/网络适配器IPv6
	IPv6 net.IP `json:"ipv6"`
	// 网卡/网络适配器IPv6所在网络
	IPv6Network net.IP `json:"ipv6_network"`
	// 网卡/网络适配器IPv6掩码位/子网掩码
	IPv6SubMask int `json:"ipv6_sub_mask"`
	// 网卡/网络适配器IPv6掩码总位数
	IPv6MaskBits int `json:"ipv6_mask_bits"`
	// 网卡/网络适配器加入的组播地址
	MulticastAddress []net.Addr `json:"multicast_address"`
}

// 根据IP获取启用的网卡/网络适配器信息
func Adaptors(ip net.IP) (*Adaptor, error) {
	addresss, err := net.ResolveIPAddr("", ip.String())
	if err != nil {
		return nil, err
	}
	IP := net.ParseIP(addresss.String())
	var adaptor *Adaptor
	adaptors, err := adaptors()
	if err != nil {
		return nil, err
	}
	for _, adp := range adaptors {
		if adp.IPv4.Equal(IP) {
			adaptor = adp
			return adaptor, nil
		} else if adp.IPv6.Equal(IP) {
			adaptor = adp
			return adaptor, nil
		}
	}
	return nil, fmt.Errorf("Not fond adaptor of ip = \"%v\".\n", ip)
}

// 获取所有启用的网卡/网络适配器信息
func adaptors() ([]*Adaptor, error) {
	// 获取所有的网卡/网络适配器信息
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	adaptors := make([]*Adaptor, 0)
LAB:
	// 遍历所有的网卡
	for _, iface := range ifaces {
		// 如果状态标记不是UP的直接执行下一个
		if iface.Flags&net.FlagUp == 0 {
			//log.Printf("The adaptor \"%v\" Flags is not up.\n", iface.Name)
			continue LAB
		}
		addresses, err := iface.Addrs()
		if err != nil {
			//log.Printf("The adaptor \"%v\" get address failed.\n", iface.Name)
			continue LAB
		}
		// 如果地址切片的长度小于2执行下一个
		if len(addresses) < 2 {
			//log.Printf("The adaptor \"%v\" address slice < 2.\n", iface.Name)
			continue LAB
		}

		ipv4 := net.IP{}
		ipv6 := net.IP{}
		ipv4Network := net.IP{}
		ipv6Network := net.IP{}
		ipv4SubMask := 0
		ipv6SubMask := 0
		ipv4MaskBits := 0
		ipv6MaskBits := 0
		for _, address := range addresses {
			ip, ipnet, err := net.ParseCIDR(address.String())
			if err != nil {
				//log.Printf("The adaptor \"%v\" address \"%v\" parse cidr failed.\n", iface.Name, address.String())
				continue
			}
			// 如果是本地回环
			if ip.IsLoopback() {
				//log.Printf("The adaptor \"%v\" address \"%v\" is loopbak address.\n", iface.Name, ip)
				continue LAB
			}
			if ip.To4() != nil {
				ipv4SubMask, ipv4MaskBits = ipnet.Mask.Size()
				ipv4 = ip.To4()
				ipv4Network = ipnet.IP
			} else if ip.To16() != nil {
				ipv6SubMask, ipv6MaskBits = ipnet.Mask.Size()
				ipv6 = ip.To16()
				ipv6Network = ipnet.IP
			}
		}
		if ipv4 == nil || ipv6 == nil {
			continue LAB
		}
		adaptor := &Adaptor{
			Index:        iface.Index,
			Name:         iface.Name,
			Mtu:          iface.MTU,
			Mac:          iface.HardwareAddr.String(),
			Flags:        iface.Flags.String(),
			IPv4:         ipv4,
			IPv4Network:  ipv4Network,
			IPv4SubMask:  ipv4SubMask,
			IPv4MaskBits: ipv4MaskBits,
			IPv6:         ipv6,
			IPv6Network:  ipv6Network,
			IPv6SubMask:  ipv6SubMask,
			IPv6MaskBits: ipv6MaskBits,
		}
		if addrs,err:= iface.MulticastAddrs();err==nil{
			adaptor.MulticastAddress = addrs
		}
		adaptor.IPv4Available = cidrCalculateAvailable(adaptor.IPv4SubMask)
		adaptor.IPv4Mask = cidrCalculateIPMask(adaptor.IPv4SubMask)
		adaptor.IPv4FirstAvailable, adaptor.IPv4LastAvailable, adaptor.IPv4Broadcast = cidrCalculateIPRange(adaptor.IPv4, adaptor.IPv4SubMask)
		adaptors = append(adaptors, adaptor)
	}
	return adaptors, nil
}

// CIDR计算可用地址个数
func cidrCalculateAvailable(mask int) uint {
	num := uint(0)
	for i := uint(32 - mask - 1); i >= 1; i-- {
		num += 1 << i
	}
	return num
}

// CIDR计算掩码IP地址
func cidrCalculateIPMask(mask int) net.IP {
	// ^uint32(0)二进制为32个比特1，通过向左位移，得到CIDR掩码的二进制
	cidrMask := ^uint32(0) << uint(32-mask)
	//fmt.Println(fmt.Sprintf("%b \n", cidrMask))
	// 计算CIDR掩码的四个片段，将想要得到的片段移动到内存最低8位后，将其强转为8位整型，从而得到
	cidrIPSegment1 := uint8(cidrMask >> 24)
	cidrIPSegment2 := uint8(cidrMask >> 16)
	cidrIPSegment3 := uint8(cidrMask >> 8)
	cidrIPSegment4 := uint8(cidrMask & uint32(255))
	return net.IPv4(cidrIPSegment1, cidrIPSegment2, cidrIPSegment3, cidrIPSegment4)
}

// CIDR根据IP和CIDR掩码位计算一个IP的可用范围和广播地址
func cidrCalculateIPRange(ip net.IP, mask int) (net.IP, net.IP, net.IP) {
	ipSegs := strings.Split(ip.String(), ".")
	ipSegments := make([]uint8, 0)
	a, _ := strconv.Atoi(ipSegs[0])
	b, _ := strconv.Atoi(ipSegs[1])
	c, _ := strconv.Atoi(ipSegs[2])
	d, _ := strconv.Atoi(ipSegs[3])
	ipSegments = append(ipSegments, uint8(a), uint8(b), uint8(c), uint8(d))
	seg3MinIp, seg3MaxIp := cidrCalculateIPSegment3(ipSegments, mask)
	seg4MinIp, seg4MaxIp := cidrCalculateIPSegment4(ipSegments, mask)
	// 广播地址是范围的最后一段的最大可用+1
	return net.IPv4(uint8(a), uint8(b), seg3MinIp, seg4MinIp), net.IPv4(uint8(a), uint8(b), seg3MaxIp, seg4MaxIp), net.IPv4(uint8(a), uint8(b), seg3MaxIp, seg4MaxIp+1)
}

// CIDR根据IP地址的某一段和CIDR掩码计算一个IP片段的区间
func ipSegmentRange(ipSegment, mask uint8) (uint8, uint8) {
	var ipSegMax uint8 = 255
	netSegIp := ipSegMax << mask
	segMinIp := netSegIp & ipSegment
	segMaxIp := ipSegment&(255<<mask) | ^(255 << mask)
	return segMinIp, segMaxIp
}

// CIDR计算IP的第三段
func cidrCalculateIPSegment3(ipSegments []uint8, mask int) (uint8, uint8) {
	if mask > 24 {
		return ipSegments[2], ipSegments[2]
	}
	return ipSegmentRange(ipSegments[2], uint8(24-mask))
}

// CIDR计算IP的第四段
func cidrCalculateIPSegment4(ipSegments []uint8, mask int) (uint8, uint8) {
	segMinIP, segMaxIP := ipSegmentRange(ipSegments[3], uint8(32-mask))
	// 因为第一个是类似：192.168.6.0这种所在网络，所以第四段的最小值要从1开始，而最大值是广播地址，所以要减去1
	return segMinIP + 1, segMaxIP - 1
}
