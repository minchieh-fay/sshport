package server

import "sshport/help"

func (s *Server) callback(t int) {
	switch t {
	case help.CALLBACK_TYPE_SSHFIN:
		s.sshconn.Close()
		s.sshconn = nil
	}
}
