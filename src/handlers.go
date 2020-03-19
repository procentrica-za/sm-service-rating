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
		fmt.Println("handleratebuyer Has Been Called!")
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
		req, respErr := http.Post("http://"+config.CRUDHost+":"+config.CRUDPort+"/rate", "application/json", bytes.NewBuffer(requestByte))

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

func (s *Server) handlerateseller() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//get JSON payload
		rateSeller := RateSeller{}
		err := json.NewDecoder(r.Body).Decode(&rateSeller)

		//handle for bad JSON provided
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, err.Error())
			return
		}

		client := &http.Client{}

		//create byte array from JSON payload
		requestByte, _ := json.Marshal(rateSeller)

		//put to crud service
		req, err := http.NewRequest("PUT", "http://"+config.CRUDHost+":"+config.CRUDPort+"/rate", bytes.NewBuffer(requestByte))
		if err != nil {
			fmt.Fprint(w, err.Error())
			fmt.Println("Error in communication with CRUD service endpoint for request to rate seller")
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

		//create new response struct
		var ratesellerResponse RateSellerResult
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&ratesellerResponse)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, err.Error())
			return
		}

		//convert struct back to JSON
		js, jserr := json.Marshal(ratesellerResponse)
		if jserr != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, jserr.Error())
			fmt.Println("Error occured when trying to marshal the response to rate seller")
			return
		}

		//return back to Front-End user
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(js)
	}
}

func (s *Server) handlegetoutstandingratings() http.HandlerFunc {
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
		req, respErr := http.Get("http://" + config.CRUDHost + ":" + config.CRUDPort + "/rate?userid=" + userid)

		//check for response error of 500
		if respErr != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, respErr.Error())
			fmt.Println("Error in communication with CRUD service endpoint for request to retrieve outstandingratings information")
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
			fmt.Fprintf(w, "An internal error has occured whilst trying to get a users outstanding ratings data"+bodyString)
			fmt.Println("An internal error has occured whilst trying to get a users outstanding rating data" + bodyString)
			return
		}

		//close the request
		defer req.Body.Close()

		//create new response struct for JSON list
		outstandingRatingList := OutstandingRatingList{}
		outstandingRatingList.Oustandingratings = []GetOutstandingResult{}

		//decode request into decoder which converts to the struct
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&outstandingRatingList)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, err.Error())
			fmt.Println("Error occured in decoding get User outstanding ratings response ")
			return
		}
		//convert struct back to JSON
		js, jserr := json.Marshal(outstandingRatingList)
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

func (s *Server) handlegetsellerratings() http.HandlerFunc {
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
		req, respErr := http.Get("http://" + config.CRUDHost + ":" + config.CRUDPort + "/sellerrating?userid=" + userid)

		//check for response error of 500
		if respErr != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, respErr.Error())
			fmt.Println("Error in communication with CRUD service endpoint for request to retrieve previousratings information")
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
			fmt.Fprintf(w, "An internal error has occured whilst trying to get a users previous ratings data"+bodyString)
			fmt.Println("An internal error has occured whilst trying to get a users previous rating data" + bodyString)
			return
		}

		//close the request
		defer req.Body.Close()

		//create new response struct for JSON list
		previousRatingList := PreviousRatingList{}
		previousRatingList.Previousratings = []GetPreviousResult{}

		//decode request into decoder which converts to the struct
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&previousRatingList)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, err.Error())
			fmt.Println("Error occured in decoding get User previous ratings response ")
			return
		}
		//convert struct back to JSON
		js, jserr := json.Marshal(previousRatingList)
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

func (s *Server) handlegetbuyerratings() http.HandlerFunc {
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
		req, respErr := http.Get("http://" + config.CRUDHost + ":" + config.CRUDPort + "/buyerrating?userid=" + userid)

		//check for response error of 500
		if respErr != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, respErr.Error())
			fmt.Println("Error in communication with CRUD service endpoint for request to retrieve previousratings information")
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
			fmt.Fprintf(w, "An internal error has occured whilst trying to get a users previous ratings data"+bodyString)
			fmt.Println("An internal error has occured whilst trying to get a users previous rating data" + bodyString)
			return
		}

		//close the request
		defer req.Body.Close()

		//create new response struct for JSON list
		previousRatingList := PreviousRatingList{}
		previousRatingList.Previousratings = []GetPreviousResult{}

		//decode request into decoder which converts to the struct
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&previousRatingList)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, err.Error())
			fmt.Println("Error occured in decoding get User previous ratings response ")
			return
		}
		//convert struct back to JSON
		js, jserr := json.Marshal(previousRatingList)
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

func (s *Server) handlegetinterestedbuyers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		interestedratings := InterestedBuyers{}
		err := json.NewDecoder(r.Body).Decode(&interestedratings)
		//handle for bad JSON provided

		if err != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, err.Error())
			fmt.Println("Could not read body of request into proper JSON format for getting interested buyers.\n Please check that your data is in the correct format.")
			return
		}

		//create byte array from JSON payload
		requestByte, _ := json.Marshal(interestedratings)

		//post to crud service
		req, respErr := http.Post("http://"+config.CRUDHost+":"+config.CRUDPort+"/interest", "application/json", bytes.NewBuffer(requestByte))

		//check for response error of 500
		if respErr != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, respErr.Error())
			fmt.Println("Error in communication with CRUD service endpoint for request to retrieve interested buyers")
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
			fmt.Fprintf(w, "An internal error has occured whilst trying to get interested buyers"+bodyString)
			fmt.Println("An internal error has occured whilst trying to get interested buyers" + bodyString)
			return
		}

		//close the request
		defer req.Body.Close()

		//create new response struct for JSON list
		interestedRatingList := InterestedRatingList{}
		interestedRatingList.Interestedratings = []GetInterestedResult{}

		//decode request into decoder which converts to the struct
		decoder := json.NewDecoder(req.Body)
		err1 := decoder.Decode(&interestedRatingList)
		if err1 != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, err1.Error())
			fmt.Println("Error occured in decoding get Messages response ")
			return
		}
		//convert struct back to JSON
		js, jserr := json.Marshal(interestedRatingList)
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
