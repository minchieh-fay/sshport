package server

import (
	"fmt"
	"log"
	"net"
	"sshport/help"
	"sshport/pool"
)

type Server struct {
	Config  *help.ConfigInfo
	sshconn net.Conn
	pool    *pool.Pool
}

func NewServer(config *help.ConfigInfo) *Server {
	return &Server{
		Config: config,
		pool:   pool.NewPool(config, "server"),
	}
}

func (s *Server) Start() {
	s.pool.SetCallback(s.callback)
	go s.readSsh()
	go s.writeSsh()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Config.Port))
	if err != nil {
		log.Fatalf("监听端口失败: %v", err)
	}
	defer listener.Close()

	log.Printf("服务端已启动，监听端口: %d", s.Config.Port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("接受连接失败: %v", err)
			continue
		}
		if !s.authConn(conn) {
			conn.Close()
			continue
		}
		s.pool.AddConn(conn)
	}
}
