package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func (s *Server) handleaddchat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("handleaddchat Has Been Called!")
		//get JSON payload

		startchat := StartChat{}
		err := json.NewDecoder(r.Body).Decode(&startchat)
		//handle for bad JSON provided

		if err != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, err.Error())
			fmt.Println("Could not read body of request into proper JSON format for starting a chat.\n Please check that your data is in the correct format.")
			return
		}

		//create byte array from JSON payload
		requestByte, _ := json.Marshal(startchat)

		//post to crud service
		req, respErr := http.Post("http://"+config.CRUDHost+":"+config.CRUDPort+"/chat", "application/json", bytes.NewBuffer(requestByte))

		//check for response error of 500
		if respErr != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, respErr.Error())
			fmt.Println("Error in communication with CRUD service endpoint for request to start a chat")
			return
		}
		if req.StatusCode != 200 {
			w.WriteHeader(req.StatusCode)
			fmt.Fprint(w, "Request to DB can't be completed...")
			fmt.Println("Unable to request to start chat to the CRUD service")
		}
		if req.StatusCode == 500 {
			w.WriteHeader(500)
			bodyBytes, err := ioutil.ReadAll(req.Body)
			if err != nil {
				log.Fatal(err)
			}
			bodyString := string(bodyBytes)
			fmt.Fprintf(w, "Request to DB can't be completed with request: "+bodyString)
			fmt.Println("Request to DB can't be completed with request: " + bodyString)
			return
		}

		//close the request
		defer req.Body.Close()

		//create new response struct
		var startchatResponse StartChatResult

		//decode request into decoder which converts to the struct
		decoder := json.NewDecoder(req.Body)
		err = decoder.Decode(&startchatResponse)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, err.Error())
			fmt.Println("Error occured in decoding start new chat response ")
			return
		}
		//convert struct back to JSON
		js, jserr := json.Marshal(startchatResponse)
		if jserr != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, jserr.Error())
			fmt.Println("Error occured when trying to marshal the decoded response into specified JSON format! ")
			return
		}

		//return success back to Front-End user
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(js)
	}
}

func (s *Server) handledeletechat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//Get Advertisement ID from URL
		chatid := r.URL.Query().Get("id")

		//Check if Advertisement ID is null
		if chatid == "" {
			w.WriteHeader(500)
			fmt.Fprint(w, "AdvertisementID not properly provided in URL")
			fmt.Println("AdvertisementID not properly provided in URL")
			return
		}
		client := &http.Client{}

		//post to crud service
		req, respErr := http.NewRequest("DELETE", "http://"+config.CRUDHost+":"+config.CRUDPort+"/chat?id="+chatid, nil)
		if respErr != nil {

			//check for response error of 500
			w.WriteHeader(500)
			fmt.Fprint(w, respErr.Error())
			fmt.Println("Error in communication with CRUD service endpoint for request to delete a chat")
			return
		}
		// Fetch Request
		resp, err := client.Do(req)
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}

		//close the request
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			w.WriteHeader(resp.StatusCode)
			fmt.Fprint(w, "Request to DB can't be completed...")
			fmt.Println("Unable to request delete chat to the CRUD service")
		}
		if resp.StatusCode == 500 {
			w.WriteHeader(500)
			bodyBytes, err := ioutil.ReadAll(req.Body)
			if err != nil {
				log.Fatal(err)
			}
			bodyString := string(bodyBytes)
			fmt.Fprintf(w, "Request to DB can't be completed with request: "+bodyString)
			fmt.Println("Request to DB can't be completed with request: " + bodyString)
			return
		}

		//create new response struct
		var deletechatResponse DeleteChatResult

		//decode request into decoder which converts to the struct
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&deletechatResponse)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, err.Error())
			fmt.Println("Error occured in decoding delete Chat response ")
			return
		}
		//convert struct back to JSON
		js, jserr := json.Marshal(deletechatResponse)
		if jserr != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, jserr.Error())
			fmt.Println("Error occured when trying to marshal the decoded response into specified JSON format!")
			return
		}

		//return success back to Front-End user
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(js)
	}
}

func (s *Server) handlegetactivechats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//Get User ID from URL

		userid := r.URL.Query().Get("userid")

		//Check if User ID provided is null
		if userid == "" {
			w.WriteHeader(500)
			fmt.Fprint(w, "User ID not properly provided in URL")
			fmt.Println("User ID not proplery provided in URL")
			return
		}
		req, respErr := http.Get("http://" + config.CRUDHost + ":" + config.CRUDPort + "/chats?userid=" + userid)

		//check for response error of 500
		if respErr != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, respErr.Error())
			fmt.Println("Error in communication with CRUD service endpoint for request to retrieve user active chat information")
			return
		}
		if req.StatusCode != 200 {
			w.WriteHeader(req.StatusCode)
			fmt.Fprint(w, "Request to DB can't be completed...")
			fmt.Println("Request to DB can't be completed...")
		}
		if req.StatusCode == 500 {
			w.WriteHeader(500)
			bodyBytes, err := ioutil.ReadAll(req.Body)
			if err != nil {
				log.Fatal(err)
			}
			bodyString := string(bodyBytes)
			fmt.Fprintf(w, "An internal error has occured whilst trying to get a users active chat data"+bodyString)
			fmt.Println("An internal error has occured whilst trying to get a users active chat data" + bodyString)
			return
		}

		//close the request
		defer req.Body.Close()

		//create new response struct for JSON list
		activeChatList := ActiveChatList{}
		activeChatList.ActiveChats = []GetActiveChatResult{}

		//decode request into decoder which converts to the struct
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&activeChatList)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, err.Error())
			fmt.Println("Error occured in decoding get User active chats response ")
			return
		}
		//convert struct back to JSON
		js, jserr := json.Marshal(activeChatList)
		if jserr != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, jserr.Error())
			fmt.Println("Error occured when trying to marshal the decoded response into specified JSON format!")
			return
		}

		//return success back to Front-End user
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(js)
	}
}

func (s *Server) handlegetmessages() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//Get Ad post type from URL
		chatid := r.URL.Query().Get("chatid")

		//Check if Advertisement Post Type is not provided in URL
		if chatid == "" {
			w.WriteHeader(500)
			fmt.Fprint(w, "Post type not properly provided in URL")
			fmt.Println("Post type not properly provided in URL")
			return
		}

		//post to crud service
		req, respErr := http.Get("http://" + config.CRUDHost + ":" + config.CRUDPort + "/message?chatid=" + chatid)

		//check for response error of 500
		if respErr != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, respErr.Error())
			fmt.Println("Error in communication with CRUD service endpoint for request to retrieve advertisement information")
			return
		}
		if req.StatusCode != 200 {
			w.WriteHeader(req.StatusCode)
			fmt.Fprint(w, "Request to DB can't be completed...")
			fmt.Println("Request to DB can't be completed...")
		}
		if req.StatusCode == 500 {
			w.WriteHeader(500)
			bodyBytes, err := ioutil.ReadAll(req.Body)
			if err != nil {
				log.Fatal(err)
			}
			bodyString := string(bodyBytes)
			fmt.Fprintf(w, "An internal error has occured whilst trying to get message data"+bodyString)
			fmt.Println("An internal error has occured whilst trying to get message data" + bodyString)
			return
		}

		//close the request
		defer req.Body.Close()

		//create new response struct for JSON list
		messagesList := MessageList{}
		messagesList.Messages = []GetMessageResult{}

		//decode request into decoder which converts to the struct
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&messagesList)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, err.Error())
			fmt.Println("Error occured in decoding get Messages response ")
			return
		}
		//convert struct back to JSON
		js, jserr := json.Marshal(messagesList)
		if jserr != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, jserr.Error())
			fmt.Println("Error occured when trying to marshal the decoded response into specified JSON format!")
			return
		}

		//return success back to Front-End user
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(js)
	}
}

func (s *Server) handleaddmessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addmessage := SendMessage{}
		err := json.NewDecoder(r.Body).Decode(&addmessage)
		//handle for bad JSON provided

		if err != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, err.Error())
			fmt.Println("Could not read body of request into proper JSON format for adding a message.\n Please check that your data is in the correct format.")
			return
		}

		//create byte array from JSON payload
		requestByte, _ := json.Marshal(addmessage)

		//post to crud service
		req, respErr := http.Post("http://"+config.CRUDHost+":"+config.CRUDPort+"/message", "application/json", bytes.NewBuffer(requestByte))

		//check for response error of 500
		if respErr != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, respErr.Error())
			fmt.Println("Error in communication with CRUD service endpoint for request to retrieve advertisement information")
			return
		}
		if req.StatusCode != 200 {
			w.WriteHeader(req.StatusCode)
			fmt.Fprint(w, "Request to DB can't be completed...")
			fmt.Println("Request to DB can't be completed...")
		}
		if req.StatusCode == 500 {
			w.WriteHeader(500)
			bodyBytes, err := ioutil.ReadAll(req.Body)
			if err != nil {
				log.Fatal(err)
			}
			bodyString := string(bodyBytes)
			fmt.Fprintf(w, "An internal error has occured whilst trying to get message data"+bodyString)
			fmt.Println("An internal error has occured whilst trying to get message data" + bodyString)
			return
		}

		//close the request
		defer req.Body.Close()

		//create new response struct for JSON list
		messagesList := MessageList{}
		messagesList.Messages = []GetMessageResult{}

		//decode request into decoder which converts to the struct
		decoder := json.NewDecoder(req.Body)
		err1 := decoder.Decode(&messagesList)
		if err1 != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, err1.Error())
			fmt.Println("Error occured in decoding get Messages response ")
			return
		}
		//convert struct back to JSON
		js, jserr := json.Marshal(messagesList)
		if jserr != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, jserr.Error())
			fmt.Println("Error occured when trying to marshal the decoded response into specified JSON format!")
			return
		}

		//return success back to Front-End user
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(js)
	}
}
