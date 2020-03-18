package main

import "github.com/gorilla/mux"

//create structs for JSON objects recieved and responses
type StartRating struct {
	AdvertisementID string `json:"advertisementid"`
	SellerID        string `json:"sellerid"`
	BuyerID         string `json:"buyerid"`
	BuyerRating     string `json:"buyerrating"`
	BuyerComments   string `json:"buyercomments"`
}

type StartRatingResult struct {
	BuyerRated bool   `json:"buyerrated"`
	RatingID   string `json:"ratingid"`
	Message    string `json:"message"`
}

type RateSeller struct {
	RatingID       string `json:"ratingid"`
	SellerRating   string `json:"sellerrating"`
	SellerComments string `json:"sellercomments"`
}

type RateSellerResult struct {
	SellerRated bool   `json:"sellerrated"`
	Message     string `json:"message"`
}

type GetOutstandingResult struct {
	RatingID    string `json:"ratingid"`
	UserName    string `json:"username"`
	Price       string `json:"price"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type OutstandingRatingList struct {
	Oustandingratings []GetOutstandingResult `json:"outstandingratings"`
}

type GetPreviousResult struct {
	RatingID string `json:"ratingid"`
	UserName string `json:"username"`
	Rating   string `json:"rating"`
	Comment  string `json:"comment"`
}

type PreviousRatingList struct {
	Previousratings []GetPreviousResult `json:"previousratings"`
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
