package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dimfeld/httptreemux"
	"golang.org/x/net/websocket"
)

type ApiVisualizer struct {
	Name       string                        `json:"name"`
	ID         uint64                        `json:"id"`
	Properties map[string]VisualizerProperty `json:"properties"`
}

func (l *LedManager) StartApi() {
	router := httptreemux.New()
	router.GET("/buffer", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		websocket.Handler(func(ws *websocket.Conn) {
			for range time.Tick(100 * time.Millisecond) {
				j, _ := json.Marshal(l.buffer)
				_, err := ws.Write(j)
				if err != nil {
					break
				}
			}
		}).ServeHTTP(w, r)
	})
	router.GET("/api/visualizer", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		m := make([]ApiVisualizer, 0, len(l.visualizers))
		for _, v := range l.visualizers {
			m = append(m, ApiVisualizer{
				Name:       v.Name(),
				ID:         v.ID(),
				Properties: GetVisualizerProperties(v),
			})
		}
		j, _ := json.Marshal(m)
		w.Write(j)
	})
	router.NotFoundHandler = http.FileServer(http.Dir("./web")).ServeHTTP
	http.ListenAndServe(":8080", router)
}
