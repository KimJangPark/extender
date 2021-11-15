package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"flag"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	// get bandwidth
	// flag.Int(flag, default, explaination)
	pointer_bandwidth = flag.Int("b", 10, "minimum bandwidth")
	flag.Parse()
	minBandwidth = *pointer_bandwidth

	router := httprouter.New()
	router.GET("/", Index)
	router.POST("/filter", Filter)
	router.POST("/prioritize", Prioritize)

	log.Print("info: server starting on the port :8888")
	if err := http.ListenAndServe(":8888", router); err != nil {
		log.Fatal(err)
	}

}
