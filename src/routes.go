package main

//create routes
func (s *Server) routes() {
	s.router.HandleFunc("/rate", s.handleratebuyer()).Methods("POST")
	s.router.HandleFunc("/rate", s.handlerateseller()).Methods("PUT")
	s.router.HandleFunc("/rate", s.handlegetoutstandingratings()).Methods("GET")
	s.router.HandleFunc("/sellerrating", s.handlegetsellerratings()).Methods("GET")
	s.router.HandleFunc("/buyerrating", s.handlegetbuyerratings()).Methods("GET")
	s.router.HandleFunc("/interest", s.handlegetinterestedbuyers()).Methods("POST")
	s.router.HandleFunc("/rating", s.handlegetratingstodo()).Methods("GET")
	s.router.HandleFunc("/buyer", s.handlegetbuyerdashboard()).Methods("GET")
	s.router.HandleFunc("/seller", s.handlegetsellerdashboard()).Methods("GET")
}
