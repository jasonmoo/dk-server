package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"code.google.com/p/go.net/websocket"
)

func websocket_upgrade(w http.ResponseWriter, req *http.Request) {
	websocket.Server{Handler: sub_handler}.ServeHTTP(w, req)
}

func sub_handler(ws *websocket.Conn) {

	defer pubsub.Unsubscribe(ws)

	for {
		var list []string
		if err := websocket.JSON.Receive(ws, &list); err != nil {
			// log.Println("socket receive error", err)
			return
		}

		pubsub.Subscribe(ws, list)
	}

}

func add_handler(w http.ResponseWriter, r *http.Request) {

	start, h := time.Now(), w.Header()
	h.Set("Server", ServerName)
	h.Set("Access-Control-Allow-Origin", "*")
	h.Set("Access-Control-Allow-Methods", "GET, POST")

	g, k, v := r.FormValue("g"), r.FormValue("k"), r.FormValue("v")
	switch {
	case len(g) == 0, len(k) == 0:
		fmt.Fprintf(w, `Missing required param g (group) or k (key)

Column Count: %d
Key Count: %d

Column List:
%s
`, table.ColumnCount(), table.KeyCount(), column_list(r.Host))
		return
	}

	inc := 1.0
	if len(v) > 0 {
		inc, _ = strconv.ParseFloat(v, 64)
	}

	table.Add(g, k, inc)

	h.Set("X-Render-Time", time.Since(start).String())

}

func top_n_handler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("me panic: %v\n", r)
			http.Error(w, "probably a json marshalling error", http.StatusInternalServerError)
		}
	}()

	start, h := time.Now(), w.Header()
	h.Set("Server", ServerName)
	h.Set("Access-Control-Allow-Origin", "*")
	h.Set("Access-Control-Allow-Methods", "GET")

	g := r.FormValue("g")
	if len(g) == 0 {
		fmt.Fprintf(w, "Missing required field g (column name)\n\n%s", column_list(r.Host))
		return
	}

	n, _ := strconv.Atoi(r.FormValue("n"))
	if n < 1 {
		n = 10
	}

	if data, err := json.MarshalIndent(table.Report([]string{g}, n), "", "\t"); err == nil {
		h.Set("Content-Type", "application/json; charset=utf-8")
		h.Set("Cache-Control", "no-cache, no-store, must-revalidate")
		h.Set("X-Render-Time", time.Since(start).String())
		w.Write(data)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func column_list(host string) string {
	if columns := table.Columns(); len(columns) > 0 {
		glue := "\nhttp://" + host + "/top?g="
		return glue + strings.Join(columns, glue)
	}
	return "<no columns>"
}
