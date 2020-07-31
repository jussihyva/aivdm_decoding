package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const (
	updateInterval = 1 * time.Second
	pushInterval   = 5 * time.Second
)

// Message ...
// [{'AIS': {'MMSI': 230985650, 'TIMESTAMP': '2020-07-30 16:35:25 UTC', 'LATITUDE': 60.46008, 'LONGITUDE': 21.91846, 'COURSE': 296.0, 'SPEED': 0.2, 'HEADING': 68, 'NAVSTAT': 5, 'IMO': 0, 'NAME': 'OSMERUS', 'CALLSIGN': 'OH5332', 'TYPE': 30, 'A': 5, 'B': 10, 'C': 4, 'D': 1, 'DRAUGHT': 3.0, 'DESTINATION': 'MERIMASKU', 'ETA_AIS': '04-28 15:30', 'ETA': '2020-04-28 15:30:00', 'SRC': 'TER', 'ZONE': 'Baltic Sea', 'ECA': True}},
type Message struct {
	MMSI string
	// TODO: Change this
	Timestamp string
	Name      string
}

type server struct {
	err          error
	db           *gorm.DB
	messages     []Message
	jsonResponse []byte
}

func (s *server) createResponse() {
	s.jsonResponse, s.err = json.Marshal(s.messages)
	if s.err != nil {
		fmt.Println(s.err)
		s.jsonResponse = []byte("{}")
	}
	fmt.Println("LEN:", len(s.messages))
}

// Update data every nth second
func (s *server) updateData() {
	ticker := time.NewTicker(updateInterval)
	for range ticker.C {
		s.db.Find(&s.messages)
		// Because response is same for all clients, we can create it here
		s.createResponse()
	}
}

func (s *server) wsHandle(w http.ResponseWriter, r *http.Request) {
	ws, err := Upgrade(w, r)
	if err != nil {
		fmt.Println(err)
	}
	s.Writer(ws)
}

func main() {

	s := &server{jsonResponse: []byte("{}")}
	db, err := gorm.Open("sqlite3", "../data.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	s.db = db
	defer s.db.Close()

	s.db.AutoMigrate(&Message{})
	go s.updateData()

	http.HandleFunc("/ws", s.wsHandle)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	if err := http.ListenAndServe(":8001", nil); err != nil {
		fmt.Println(err)
		return
	}
}
