package main

import (
	"fmt"
	"log"
	"net/http"
)

func digestHandler(w http.ResponseWriter, r *http.Request) {
	orders := []WorkOrder{{ID: "WO-1042", Site: "North plant", Status: "dispatched", PhotoURL: "https://photos.example/wo-1042.jpg", FollowUp: "confirm pump reading"}}
	text, send := DigestText(orders)
	if !send {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	client, err := NewInfrai()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if _, err = client.EmailSend("ops@example.com", "Field-service digest", text); err != nil {
		http.Error(w, err.Error(), 502)
		return
	}
	fmt.Fprintln(w, "digest sent")
}

func main() {
	client, err := NewInfrai()
	if err != nil {
		log.Fatal(err)
	}
	job, err := client.CronCreate("0 8 * * *", "http://localhost:8080/digest")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("scheduled job %s; serving digest endpoint", job)
	http.HandleFunc("/digest", digestHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
