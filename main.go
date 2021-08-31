package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Note struct {
	Time     time.Time
	Text     string
	Lifetime time.Time `json:"-"`
}

func main() {
	Database := make(map[string][]Note)

	// api/add-new-user?user=kek
	go AddNewUser(Database)

	// /api/add-note?user=kek&text=LOL&lifetime=3
	go AddNewNote(Database)

	// api/delete-note?id=123
	go DeleteNote(Database)

	// api/get-all-notes?user=kek
	go GetAllNotes(Database)

	// api/get-first-note?user=kek
	go GetFirstNote(Database)

	// api/get-last-note?user=kek
	GetLastNote(Database)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

//::TODO check user duplicate

func GetFirstNote(Database map[string][]Note) {
	http.HandleFunc("/api/get-first-note", func(w http.ResponseWriter, r *http.Request) {
		user := r.URL.Query().Get("user")
		if user == "" {
			SendMsg(w, "/api/add-note", "Empty user")
			return
		}
		if _, ok := Database[user]; !ok {
			SendMsg(w, "/api/get-first-note", "Not found user")
			return
		}

		// check lifetime
		for i, note := range Database[user] {
			if !note.Lifetime.IsZero() && time.Now().After(note.Lifetime) {
				Database[user] = append(Database[user][:i], Database[user][i+1:]...)
			}
		}

		if len(Database[user]) == 0 {
			SendMsg(w, "/api/get-first-note", "Empty notes")
			return
		}

		js, err := json.Marshal(Database[user][0])
		if err != nil {
			log.Fatal("get-first-note: " + err.Error())
		}
		SendMsg(w, "/api/get-first-note", string(js))
	})
}

func GetLastNote(Database map[string][]Note) {
	http.HandleFunc("/api/get-last-note", func(w http.ResponseWriter, r *http.Request) {
		user := r.URL.Query().Get("user")
		if user == "" {
			SendMsg(w, "/api/get-last-note", "Empty user")
			return
		}
		if _, ok := Database[user]; !ok {
			SendMsg(w, "/api/get-last-note", "Not found user")
			return
		}

		// check lifetime
		for i, note := range Database[user] {
			if !note.Lifetime.IsZero() && time.Now().After(note.Lifetime) {
				Database[user] = append(Database[user][:i], Database[user][i+1:]...)
			}
		}

		if len(Database[user]) == 0 {
			SendMsg(w, "/api/get-last-note", "Empty notes")
			return
		}

		js, err := json.Marshal(Database[user][len(Database[user])-1])
		if err != nil {
			log.Fatal("get-last-note: " + err.Error())
		}
		SendMsg(w, "/api/get-last-note", string(js))
	})
}

func GetAllNotes(Database map[string][]Note) {
	http.HandleFunc("/api/get-all-notes", func(w http.ResponseWriter, r *http.Request) {
		user := r.URL.Query().Get("user")
		if user == "" {
			SendMsg(w, "/api/get-all-notes", "Empty user")
			return
		}
		if _, ok := Database[user]; !ok {
			SendMsg(w, "/api/get-all-notes", "Not found user")
			return
		}

		// check lifetime
		for i, note := range Database[user] {
			if !note.Lifetime.IsZero() && time.Now().After(note.Lifetime) {
				Database[user] = append(Database[user][:i], Database[user][i+1:]...)
			}
		}

		if len(Database[user]) == 0 {
			SendMsg(w, "/api/get-all-notes", "Empty notes")
			return
		}

		js, err := json.Marshal(Database[user])
		if err != nil {
			log.Fatal("get-all-notes: " + err.Error())
		}
		SendMsg(w, "/api/get-all-notes", string(js))
	})
}

func DeleteNote(Database map[string][]Note) {
	http.HandleFunc("/api/delete-note", func(w http.ResponseWriter, r *http.Request) {
		user := r.URL.Query().Get("user")
		if user == "" {
			SendMsg(w, "/api/delete-note", "Empty user")
			return
		}
		if _, ok := Database[user]; !ok {
			SendMsg(w, "/api/delete-note", "Not found user")
			return
		}
		id := r.URL.Query().Get("id")
		if id == "" {
			SendMsg(w, "/api/delete-note", "Empty id note")
			return
		}

		res, err := strconv.Atoi(id)
		if err != nil {
			SendMsg(w, "/api/delete-note", "Invalid id")
			return
		}

		// ::Todo no effect
		for i, _ := range Database[user] {
			if i == res {
				Database[user] = append(Database[user][:i], Database[user][i+1:]...)
				SendMsg(w, "/api/delete-note", "Note deleted")
				return
			}
		}
		SendMsg(w, "/api/delete-note", "Not found Note")
	})
}

func AddNewNote(Database map[string][]Note) {
	http.HandleFunc("/api/add-note", func(w http.ResponseWriter, r *http.Request) {
		var reply string
		var timeHours int
		user := r.URL.Query().Get("user")
		if _, ok := Database[user]; !ok {
			SendMsg(w, "/api/add-note", "Not found user")
			return
		}
		text := r.URL.Query().Get("text")
		lifetime := r.URL.Query().Get("lifetime")
		switch {
		case user == "":
			reply = "Unidentified user"
			SendMsg(w, "/api/add-note", reply)
			return
		case text == "":
			reply = "Empty text"
			SendMsg(w, "/api/add-note", reply)
			return
		case lifetime != "":
			{
				t, err := strconv.Atoi(lifetime)
				if err != nil {
					reply = "Invalid lifetime"
					SendMsg(w, "/api/add-note", reply)
					return
				}
				timeHours = t
			}
		}
		var lt time.Time
		if timeHours != 0 {
			lt = time.Now()
			lt.Add(time.Duration(timeHours) * time.Hour)
		}
		currentTime := time.Now()
		Database[user] = append(Database[user],
			Note{
				Time:     currentTime,
				Text:     text,
				Lifetime: lt})
		reply = "Added new Note"
		SendMsg(w, "/api/add-note", reply)
	})
}

func SendMsg(w http.ResponseWriter, uri string, reply string) {
	_, err := fmt.Fprintln(w, reply)
	if err != nil {
		log.Fatal(uri, err)
	}
}

func AddNewUser(Database map[string][]Note) {
	http.HandleFunc("/api/add-new-user", func(w http.ResponseWriter, r *http.Request) {
		var reply string
		user := r.URL.Query().Get("user")
		if _, ok := Database[user]; !ok {
			Database[user] = []Note{}
			reply = "Added new User " + user
		} else {
			reply = "Such a user already exists"
		}
		_, err := fmt.Fprintln(w, reply)
		if err != nil {
			log.Fatal("/api/add-new-user: ", err)
		}
	})
}
