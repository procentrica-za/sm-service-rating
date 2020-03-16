package main

//create routes
func (s *Server) routes() {
	s.router.HandleFunc("/rate", s.handleratebuyer()).Methods("POST")
}
