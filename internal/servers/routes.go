package servers

func (s *Server) RegisterRoutes() {
	s.mux.HandleFunc("/health", s.handler.Helath)
}