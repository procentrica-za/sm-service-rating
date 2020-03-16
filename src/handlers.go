package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func (s *Server) handleratebuyer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("handleaddchat Has Been Called!")
		//get JSON payload

		startrating := StartRating{}
		err := json.NewDecoder(r.Body).Decode(&startrating)
		//handle for bad JSON provided

		if err != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, err.Error())
			fmt.Println("Could not read body of request into proper JSON format for rate buyer.\n Please check that your data is in the correct format.")
			return
		}

		//create byte array from JSON payload
		requestByte, _ := json.Marshal(startrating)

		//post to crud service
		req, respErr := http.Post("http://"+config.CRUDHost+":"+config.CRUDPort+"/rating", "application/json", bytes.NewBuffer(requestByte))

		//check for response error of 500
		if respErr != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, respErr.Error())
			fmt.Println("Error in communication with CRUD service endpoint for request to rate buyer")
			return
		}
		if req.StatusCode != 200 {
			w.WriteHeader(req.StatusCode)
			fmt.Fprint(w, "Request to DB can't be completed...")
			fmt.Println("Unable to request to rate buyer to the CRUD service")
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
		var startratingResponse StartRatingResult

		//decode request into decoder which converts to the struct
		decoder := json.NewDecoder(req.Body)
		err = decoder.Decode(&startratingResponse)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, err.Error())
			fmt.Println("Error occured in decoding rate buyer response ")
			return
		}
		//convert struct back to JSON
		js, jserr := json.Marshal(startratingResponse)
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
