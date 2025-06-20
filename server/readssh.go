package server

import (
	"log"
	"time"
)

func (s *Server) readSsh() {
	for {
		if s.sshconn == nil {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		// 读取connssh的数据
		buf := make([]byte, 8192)
		n, err := s.sshconn.Read(buf)
		if err != nil {
			log.Printf("读取ssh连接失败: %v", err)
			s.sshconn.Close()
			s.sshconn = nil
			continue
		}
		s.pool.Write(buf[:n])
	}
}
