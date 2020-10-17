package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
)

func main() {
	filePath, _ := filepath.Abs("./src.yml")

	config := getConfig(filePath)
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

// Структура для конфига
type Service struct {
	Proxy struct {
		ListenPort int `yaml:"listenPort"`
		Servers    []struct {
			Name      string `yaml:"name"`
			Prefix    string `yaml:"prefix"`
			TargetUrl string `yaml:"target"`
		}
	}
}

// getConfig считывает yaml-файл с настройками обратного прокси
// Возвращает структуру Service
func getConfig(filePath string) Service {

	var service Service
	yamlFile, _ := ioutil.ReadFile(filePath)

	err := yaml.Unmarshal(yamlFile, &service)
	if err != nil {
		panic(err)
	}
	return service
}
