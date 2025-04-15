package handlers

import (
	"errors"
	"inventory/api-gateway/config"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func ProxyRequestToInventory(c *gin.Context) {
	proxy := createProxy(config.InventoryServiceURL)
	err := proxy.ServeHTTP(c.Writer, c.Request)
	if err != nil {
		HandleError(c, errors.New("Error while proxying to Inventory Service"), http.StatusInternalServerError)
		return
	}
}

func ProxyRequestToOrder(c *gin.Context) {
	proxy := createProxy(config.OrderServiceURL)
	err := proxy.ServeHTTP(c.Writer, c.Request)
	if err != nil {
		HandleError(c, errors.New("Error while proxying to Order Service"), http.StatusInternalServerError)
		return
	}
}

func createProxy(target string) *httputil.ReverseProxy {
	parsedURL, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Error parsing URL: %v", err)
	}
	return httputil.NewSingleHostReverseProxy(parsedURL)
}
