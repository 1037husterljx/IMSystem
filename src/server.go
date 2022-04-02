package imsystem

import {
	"fmt",
}

type Server struct{
	Ip string
	Port int
}
func (this *Server) NewServer(ip string, port int) *Server{
	server := &Server{
		Ip : ip,
		Port : port
	}
	
	return server
}
