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
	// pointer_bandwidth = flag.Int("b", 10, "minimum bandwidth")
	worker1 := flag.Int("1", 100, "worker node1 bandwidth")
	worker2 := flag.Int("2", 100, "worker node2 bandwidth")
	worker3 := flag.Int("3", 100, "worker node3 bandwidth")
	flag.Parse()
	// minBandwidth = *pointer_bandwidth
	worker1Bandwidth := *worker1
	worker2Bandwidth := *worker2
	worker3Bandwidth := *worker3

	router := httprouter.New()
	router.GET("/", Index)
	router.POST("/filter", Filter)
	router.POST("/prioritize", Prioritize)

	log.Print("info: server starting on the port :8888")
	if err := http.ListenAndServe(":8888", router); err != nil {
		log.Fatal(err)
	}

}
