package main

//create routes
func (s *Server) routes() {
	s.router.HandleFunc("/rate", s.handleratebuyer()).Methods("POST")
	s.router.HandleFunc("/rate", s.handlerateseller()).Methods("PUT")
}
