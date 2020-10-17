package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

func main() {
	port := 777
	targetURL := "http://localhost:8008"

	u, err := url.Parse(targetURL)
	if err != nil {
		log.Fatalf("Could not parse downstream url: %s", targetURL)
	}

	proxy := httputil.NewSingleHostReverseProxy(u)

	//proxy.ModifyResponse = func(res *http.Response) error {
	//	responseContent := map[string]interface{}{}
	//	err := parseResponse(res, &responseContent)
	//	if err != nil {
	//		return err
	//	}
	//
	//	return captureMetrics(responseContent)
	//}

	director := proxy.Director
	proxy.Director = func(req *http.Request) {
		director(req)
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.Host = req.URL.Host
	}

	http.HandleFunc("/", proxy.ServeHTTP)
	log.Printf("Listening on port %d", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

//func parseResponse(res *http.Response, unmarshalStruct interface{}) error {
//	body, err := ioutil.ReadAll(res.Body)
//	if err != nil {
//		return err
//	}
//	res.Body.Close()
//
//	res.Body = ioutil.NopCloser(bytes.NewBuffer(body))
//	return json.Unmarshal(body, unmarshalStruct)
//}
