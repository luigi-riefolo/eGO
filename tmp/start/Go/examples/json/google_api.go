// This sample program demonstrates how to decode a JSON response
// using the json package and NewDecoder function.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type (
	// gResult maps to the result document received from the search.
	gResult struct {
		GsearchResultClass string `json:"GsearchResultClass"`
		UnescapedURL       string `json:"unescapedURL"`
		UR                 string `json:"url"`
		VisibleURL         string `json:"visibleURL"`
		CacheURL           string `json:"cacheURL"`
		Title              string `json:"title"`
		TitleNoFormatting  string `json:"titleNoFormatting"`
		Content            string `json:"content"`
	}
	// gResponse contains the top level document.
	gResponse struct {
		ResponseData struct {
			Results []gResult `json:"results"`
		} `json:"responseData"`
	}
)

// Contact ...
type Contact struct {
	Name    string `json:"name"`
	Title   string `json:"title"`
	Contact struct {
		Home string `json:"home"`
		Cell string `json:"cell"`
	} `json:"contact"`
}

func main() {

	uri := "http://ajax.googleapis.com/ajax/services/search/web?v=1.0&rsz=8&q=golang"

	// Issue the search against Google.
	resp, err := http.Get(uri)
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	defer resp.Body.Close()
	// Decode the JSON response into our struct type.
	var gr gResponse
	err = json.NewDecoder(resp.Body).Decode(&gr)

	if err != nil {
		fmt.Println("ERROR: ", err)
		return
	}

	fmt.Println(gr)

	var JSON = `{
		"name": "Gopher",
		"title": "programmer",
		"contact": {
	    		"home": "415.333.3333",
	    		"cell": "415.555.5555"
		}
	}`
	// Unmarshal the JSON string into our variable.
	var c Contact
	err = json.Unmarshal([]byte(JSON), &c)
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	fmt.Println(c)

	// Unmarshal the JSON string into our map variable.
	var m map[string]interface{}
	err = json.Unmarshal([]byte(JSON), &m)
	if err != nil {
		log.Println("ERROR:", err)
		return
	}

	// Marshall
	// Create a map of key/value pairs.
	p := make(map[string]interface{})
	p["name"] = "Gopher"
	p["title"] = "programmer"
	p["contact"] = map[string]interface{}{
		"home": "415.333.3333",
		"cell": "415.555.5555",
	}

	data, perr := json.MarshalIndent(p, "", " ")
	if perr != nil {
		log.Println("ERROR:", perr)
		return
	}
	fmt.Println(string(data))
	/*if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}*/
}
