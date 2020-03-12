package main

import "github.com/gorilla/mux"

//create structs for JSON objects recieved and responses
type StartChat struct {
	SellerID        string `json:"sellerid"`
	BuyerID         string `json:"buyerid"`
	AdvertisementID string `json:"advertisementid"`
}

type StartChatResult struct {
	ChatPosted bool   `json:"chatposted"`
	ChatID     string `json:"chatid"`
	Message    string `json:"message"`
}

type ChatID struct {
	ChatID string `json:"id"`
}

type DeleteChatResult struct {
	ChatDeleted bool   `json:"chatposted"`
	Message     string `json:"message"`
}

type GetActiveChatResult struct {
	ChatID      string `json:"chatid"`
	UserName    string `json:"username"`
	Message     string `json:"message"`
	MessageDate string `json:"messagedate"`
}

type ActiveChatList struct {
	ActiveChats []GetActiveChatResult `json:"activechats"`
}

type GetMessageResult struct {
	MessageID   string `json:"messageid"`
	UserName    string `json:"username"`
	Message     string `json:"message"`
	MessageDate string `json:"messagedate"`
}

type MessageList struct {
	Messages []GetMessageResult `json:"messages"`
}

type SendMessage struct {
	ChatID   string `json:"chatid"`
	AuthorID string `json:"authorid"`
	Message  string `json:"message"`
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
