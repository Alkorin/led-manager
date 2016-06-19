package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/dimfeld/httptreemux"
	"golang.org/x/net/websocket"
)

func (l *LedManager) StartApi() {
	router := httptreemux.New()
	router.GET("/events", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		websocket.Handler(func(ws *websocket.Conn) {
			eventListener, eventCloser := l.apiEvents.Listen()
			defer close(eventCloser)
			for {
				data := <-eventListener
				j, _ := json.Marshal(data)
				_, err := ws.Write(j)
				if err != nil {
					break
				}
			}
		}).ServeHTTP(w, r)
	})
	router.GET("/api/visualizer", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		m := make([]*ApiVisualizer, 0, len(l.visualizers))
		for _, v := range l.visualizers {
			m = append(m, NewApiVisualizer(v))
		}
		j, _ := json.Marshal(m)
		w.Write(j)
	})
	router.GET("/api/visualizer/:id", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		ID, err := strconv.ParseUint(params["id"], 10, 64)
		if err == nil {
			for _, v := range l.visualizers {
				if ID == v.ID() {
					j, _ := json.Marshal(NewApiVisualizer(v))
					w.Write(j)
					return
				}
			}
		}
		w.WriteHeader(http.StatusNotFound)
	})
	router.PUT("/api/visualizer/:id/properties", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		ID, err := strconv.ParseUint(params["id"], 10, 64)
		if err == nil {
			for _, v := range l.visualizers {
				if ID == v.ID() {
					body, _ := ioutil.ReadAll(r.Body)
					data := map[string]interface{}{}
					err := json.Unmarshal(body, &data)
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
						w.Write([]byte(err.Error()))
						return
					}
					if err := SetVisualizerProperties(v, data); err != nil {
						w.WriteHeader(http.StatusBadRequest)
						w.Write([]byte(err.Error()))
						return
					} else {
						// OK
						l.apiEvents.Write(NewApiVisualizerPropertiesChangedEvent(v.ID()))
						return
					}
				}
			}
		}
		w.WriteHeader(http.StatusNotFound)
	})

	// Generate bufferUpdate events
	go func() {
		for range time.Tick(100 * time.Millisecond) {
			l.apiEvents.Write(NewApiBufferEvent(l.buffer))
		}
	}()

	// Default: server static files
	router.NotFoundHandler = http.FileServer(http.Dir("./web")).ServeHTTP
	http.ListenAndServe(":8080", router)
}
