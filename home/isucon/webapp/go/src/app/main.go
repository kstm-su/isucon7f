package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"github.com/parnurzeal/gorequest"

	"net"
	_ "net/http/pprof"
)

var (
	db     *sqlx.DB
	MItems map[int]mItem
	prec   = 70
)

func initDB() {
	db_host := os.Getenv("ISU_DB_HOST")
	if db_host == "" {
		db_host = "127.0.0.1"
	}
	db_port := os.Getenv("ISU_DB_PORT")
	if db_port == "" {
		db_port = "3306"
	}
	db_sock := os.Getenv("ISU_DB_SOCK")
	db_user := os.Getenv("ISU_DB_USER")
	if db_user == "" {
		db_user = "root"
	}
	db_password := os.Getenv("ISU_DB_PASSWORD")
	if db_password != "" {
		db_password = ":" + db_password
	}

	var dsn string
	if db_sock == "" {
		dsn = fmt.Sprintf("%s%s@tcp(%s:%s)/isudb?parseTime=true&loc=Local&charset=utf8mb4",
			db_user, db_password, db_host, db_port)
	} else {
		dsn = fmt.Sprintf("%s%s@unix(%s)/isudb?parseTime=true&loc=Local&charset=utf8mb4",
			db_user, db_password, db_sock)
	}

	log.Printf("Connecting to db: %q", dsn)
	db, _ = sqlx.Connect("mysql", dsn)
	for {
		err := db.Ping()
		if err == nil {
			break
		}
		log.Println(err)
		time.Sleep(time.Second * 3)
	}

	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(5 * time.Minute)
	log.Printf("Succeeded to connect db.")

	MItems = make(map[int]mItem)
	var items []mItem
	db.Select(&items, "SELECT * FROM m_item")
	//if err != nil {
	//	tx.Rollback()
	//	return nil, err
	//}
	for _, item := range items {
		MItems[item.ItemID] = item
	}

}

func getAllInitializeHandler(w http.ResponseWriter, r *http.Request) {
	gorequest.New().Get("http://app0251.isu7f.k0y.org:5000/sync/initialize").End()
	gorequest.New().Get("http://app0252.isu7f.k0y.org:5000/sync/initialize").End()
	gorequest.New().Get("http://app0253.isu7f.k0y.org:5000/sync/initialize").End()
	gorequest.New().Get("http://app0254.isu7f.k0y.org:5000/sync/initialize").End()
}

func getInitializeHandler(w http.ResponseWriter, r *http.Request) {
	db.MustExec("TRUNCATE TABLE adding")
	db.MustExec("TRUNCATE TABLE buying")
	db.MustExec("TRUNCATE TABLE room_time")
	w.WriteHeader(204)
}

func getRoomHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	roomName := vars["room_name"]
	path := "/ws/" + url.PathEscape(roomName)

	w.Header().Set("Content-Type", "application/json")

	hostNum := 0
	for _, r := range roomName {
		hostNum ^= int(r)
	}
	wsHostName := fmt.Sprintf("app025%d.isu7f.k0y.org", hostNum%4+1)
	log.Printf("websocket host name: %s\n", wsHostName)

	json.NewEncoder(w).Encode(struct {
		Host string `json:"host"`
		Path string `json:"path"`
	}{
		//Host: "",
		Host: wsHostName,
		Path: path,
	})
}

func wsGameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	roomName := vars["room_name"]

	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		log.Println("Failed to upgrade", err)
		return
	}
	go serveGameConn(ws, roomName)
}

func main() {
	l, _ := net.Listen("tcp", ":5001")
	go http.Serve(l, nil)

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	initDB()

	r := mux.NewRouter()
	r.HandleFunc("/sync/initialize", getInitializeHandler)
	r.HandleFunc("/initialize", getAllInitializeHandler)
	r.HandleFunc("/room/", getRoomHandler)
	r.HandleFunc("/room/{room_name}", getRoomHandler)
	r.HandleFunc("/ws/", wsGameHandler)
	r.HandleFunc("/ws/{room_name}", wsGameHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("../public/")))

	log.Fatal(http.ListenAndServe(":5000", handlers.LoggingHandler(os.Stderr, r)))
}
