package main

import (
	"log"
	"sync"
	"time"

	"code.google.com/p/go.net/websocket"
	"github.com/jasonmoo/dk"
)

type PubSub struct {
	sync.Mutex
	sockets  map[*websocket.Conn]map[string]bool
	subs     map[string]int
	running  bool
	interval time.Duration
}

func NewPubSub(interval time.Duration) *PubSub {
	return &PubSub{
		sockets:  make(map[*websocket.Conn]map[string]bool),
		subs:     make(map[string]int),
		interval: interval,
	}
}

func (p *PubSub) Start() {
	p.running = true
	go func() {
		for p.running {
			p.Publish()
			time.Sleep(p.interval)
		}
	}()
}

func (p *PubSub) Stop() {
	p.running = false
}

func (p *PubSub) Subscribe(ws *websocket.Conn, list []string) {
	p.Lock()
	p.set(ws, list)
	p.Unlock()
}

func (p *PubSub) Unsubscribe(ws *websocket.Conn) {
	p.Lock()
	p.remove(ws)
	p.Unlock()
}

func (p *PubSub) set(ws *websocket.Conn, list []string) {

	if _, exists := p.sockets[ws]; !exists {
		p.sockets[ws] = make(map[string]bool)
	}

	for sub, _ := range p.sockets[ws] {
		p.remove_sub(sub)
		delete(p.sockets[ws], sub)
	}

	for _, sub := range list {
		p.add_sub(sub)
		p.sockets[ws][sub] = true
	}

}

func (p *PubSub) remove(ws *websocket.Conn) {

	if _, exists := p.sockets[ws]; exists {

		for sub, _ := range p.sockets[ws] {
			p.remove_sub(sub)
		}

		_ = ws.Close()

		delete(p.sockets, ws)

	}

}

func (p *PubSub) add_sub(sub string) {
	p.subs[sub]++
}

func (p *PubSub) remove_sub(sub string) {
	if ct, exists := p.subs[sub]; exists {
		if ct == 1 {
			delete(p.subs, sub)
		} else {
			p.subs[sub]--
		}
	}
}

func (p *PubSub) Publish() {

	p.Lock()
	defer p.Unlock()

	if len(p.subs) == 0 {
		return
	}

	// build a list of subs to report on
	var subs []string
	for group, _ := range p.subs {
		subs = append(subs, group)
	}

	report := table.Report(subs, 20)

	full_set := report.ResultSet

	// push grouped result sets out to each socket
	for ws, socket_subs := range p.sockets {

		if len(socket_subs) == 0 {
			continue
		}

		new_set := make(map[string]dk.Result)

		for group, _ := range socket_subs {
			if _, exists := full_set[group]; exists {
				new_set[group] = full_set[group]
			}
		}

		report.ResultSet = new_set

		if err := websocket.JSON.Send(ws, report); err != nil {
			log.Println("big boom", err)
			p.remove(ws)
		}

	}

}
