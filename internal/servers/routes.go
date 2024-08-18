package servers

func (s *Server) RegisterRoutes() {
	s.mux.HandleFunc("/ping", s.handler.Helath)
}