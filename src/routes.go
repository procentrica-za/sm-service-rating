package main

//create routes
func (s *Server) routes() {
	s.router.HandleFunc("/chat", s.handleaddchat()).Methods("POST")
	s.router.HandleFunc("/chat", s.handledeletechat()).Methods("DELETE")
	s.router.HandleFunc("/chats", s.handlegetactivechats()).Methods("GET")
	s.router.HandleFunc("/message", s.handlegetmessages()).Methods("GET")
	s.router.HandleFunc("/message", s.handleaddmessage()).Methods("POST")
}
