package main

import "github.com/gorilla/mux"

//create structs for JSON objects recieved and responses
type StartRating struct {
	AdvertisementID string `json:"advertisementid"`
	SellerID        string `json:"sellerid"`
	BuyerID         string `json:"buyerid"`
	BuyerRating     string `json:"buyerating"`
	BuyerComments   string `json:"buyercomments"`
}

type StartRatingResult struct {
	BuyerRated      bool   `json:"buyerrated"`
	AdvertisementID string `json:"advertisementid"`
	Message         string `json:"message"`
}

//touter service struct
type Server struct {
	router *mux.Router
}

//config struct
type Config struct {
	CRUDHost   string
	CRUDPort   string
	RATINGPort string
}
