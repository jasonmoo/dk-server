package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/jasonmoo/dk"
)

const ServerName = "dk server"

var (
	table  *dk.Table
	pubsub *PubSub

	// cli options
	http_host = flag.String("host", ":80", "addr:port to listen on for http")

	decay_rate      = flag.Float64("decay_rate", .05, "rate of decay per second")
	decay_floor     = flag.Float64("decay_floor", 1, "minimum value to keep")
	decay_interval  = flag.Duration("decay_interval", 2*time.Second, "maximum amount to time between decays")
	socket_interval = flag.Duration("socket_interval", time.Second, "interval to publish to listening websockets")
)

func main() {

	log.Println("dk starting up")

	flag.Parse()
	if *http_host == "" {
		fmt.Println(ServerName + "\n\nUsage:\n")
		flag.PrintDefaults()
		fmt.Println()
		os.Exit(0)
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	table = dk.NewTable(*decay_rate, *decay_floor, *decay_interval)
	table.Start()

	pubsub = NewPubSub(*socket_interval)
	pubsub.Start()

	http.HandleFunc("/", add_handler)
	http.HandleFunc("/top", top_n_handler)
	http.HandleFunc("/sub", websocket_upgrade)

	log.Fatal(http.ListenAndServe(*http_host, nil))

}
