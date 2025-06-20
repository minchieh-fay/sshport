package server

import (
	"log"
	"os"
	"time"
)

func (s *Server) writeSsh() {
	for {
		d := s.pool.Read()
		if d == nil {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		if s.sshconn == nil {
			s.keepsshconn()
		}
		_, err := s.sshconn.Write(d)
		if err == nil {
			continue
		}
		s.keepsshconn()
		_, err = s.sshconn.Write(d)
		if err != nil {
			log.Printf("ssh连接失败: %v", err)
			os.Exit(1)
		}
	}
}
