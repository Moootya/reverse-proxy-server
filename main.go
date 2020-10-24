package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
)

var config ServiceConfig
var TargetUrl, prefixRedirect, name, fileName string
var listenPort int

func init() {
	flag.StringVar(&fileName, "file", "servers.yml", "File name with settings")
	flag.StringVar(&TargetUrl, "target", "http://localhost:8008/", "Target url")
	flag.StringVar(&prefixRedirect, "prefix", "/app/", "Target prefix")
	flag.StringVar(&name, "name", "Example", "Set name to proxy")
	flag.IntVar(&listenPort, "listen", 777, "Main listen port")
}

func main() {
	flag.Parse()
	filePath, _ := filepath.Abs(fileName)
	_, err := os.Stat(filePath)

	if err == nil {
		config = getConfig(filePath)
	} else {
		config.Proxy.ListenPort = listenPort
		config.Proxy.Servers = []Server{
			{Name: name, Prefix: prefixRedirect, TargetUrl: TargetUrl},
		}
	}

	listenPort := config.Proxy.ListenPort
	targets := config.Proxy.Servers

	for _, target := range targets {

		urlTarget, _ := url.Parse(target.TargetUrl)
		proxy := httputil.NewSingleHostReverseProxy(urlTarget)

		director := proxy.Director
		proxy.Director = func(request *http.Request) {
			director(request)
			request.Header.Set("X-Forwarded-Host", request.Header.Get("Host"))
			request.Host = request.URL.Host
		}

		http.HandleFunc(target.Prefix, proxy.ServeHTTP)
	}

	log.Printf("Listening on port %d", listenPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", listenPort), nil))
}

// Server struct for one server
type Server struct {
	Name      string `yaml:"name"`
	Prefix    string `yaml:"prefix"`
	TargetUrl string `yaml:"target"`
}

// ServiceConfig wrapper for Server - structs
type ServiceConfig struct {
	Proxy struct {
		ListenPort int `yaml:"listenPort"`
		Servers    []Server
	}
}

// getConfig read and parse Config
// return ServiceConfig - struct
func getConfig(filePath string) ServiceConfig {

	var service ServiceConfig
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(yamlFile, &service)
	if err != nil {
		log.Fatal(err)
	}
	return service
}
