package server

func (s *Server) RegisterRoutes() {
	s.mux.HandleFunc("/health", s.handler.Health)
}