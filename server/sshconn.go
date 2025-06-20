package server

import (
	"log"
	"net"
	"os"
)

func (s *Server) keepsshconn() {
	conn, err := net.Dial("tcp", s.Config.SshAddress)
	if err != nil || conn == nil {
		log.Fatalf("连接ssh失败: %v", err)
		os.Exit(1)
	}
	s.sshconn = conn
}
