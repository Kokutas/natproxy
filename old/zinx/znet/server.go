package znet

import (
	"errors"
	"fmt"
	"log"
	"natproxy/old/zinx/utils"
	"natproxy/old/zinx/ziface"
	"net"
	"os"
	"os/signal"
	"syscall"
)

// 定义服务类，实现iserver接口
type Server struct {
	// 服务器名称
	Name string
	// 服务器绑定的IP
	IP string
	// 服务器绑定的端口
	Port uint
	//当前Server由用户绑定的回调router,也就是Server注册的链接对应的处理业务
	Router ziface.IRouter
}

//============== 定义当前客户端链接的handle api ===========
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	//回显业务
	fmt.Println("[Conn Handle] CallBackToClient ... ")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf err ", err)
		return errors.New("CallBackToClient error")
	}
	return nil
}

//============== 实现 ziface.IServer 里的全部接口方法 ========

// 开启网络服务
func (s *Server) Start() {
	log.Printf("[START] Server listenner at IP: %v, Port %v, is starting.\n", s.IP, s.Port)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d,  MaxPacketSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize)
	// 开启一个go去做服务端Linster业务
	go func() {
		// 获取tcp的addr
		addr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(s.IP, fmt.Sprintf("%v", s.Port)))
		if err != nil {
			log.Printf("Resolve tcp address %s error : %v.\n", net.JoinHostPort(s.IP, fmt.Sprintf("%v", s.Port)), err.Error())
			return
		}
		// 开启监听服务器地址
		listener, err := net.ListenTCP("tcp", addr)
		if err != nil {
			log.Printf("Listen tcp addr %v error : %v.\n", net.JoinHostPort(s.IP, fmt.Sprintf("%v", s.Port)), err.Error())
			return
		}
		log.Printf("Listen tcp addr %v success.\n", net.JoinHostPort(s.IP, fmt.Sprintf("%v", s.Port)))
		//TODO server.go 应该有一个自动生成ID的方法
		var cid uint32
		cid = 0
		// 启动网络连接服务
		for {
			//3.1 阻塞等待客户端建立连接请求
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err ", err)
				continue
			}
			//3.2 TODO Server.Start() 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接
			//3.3 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
			dealConn := NewConntion(conn, cid, s.Router)
			cid++
			//3.4 启动当前链接的处理业务
			go dealConn.Start()
		}

	}()
}
func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server , name ", s.Name)
	//TODO  Server.Stop() 将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
}
func (s *Server) Service() {
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	// start server

	s.Start()
	//TODO Server.Serve() 是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加

	<-sig
}

//路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	fmt.Println("Add Router succ! ")
}

/*
  创建一个服务器句柄
*/
func NewServer() ziface.IServer {
	s := &Server{
		Name:   utils.GlobalObject.Name,    //从全局参数获取
		IP:     utils.GlobalObject.Host,    //从全局参数获取
		Port:   utils.GlobalObject.TcpPort, //从全局参数获取
		Router: nil,
	}
	return s
}
