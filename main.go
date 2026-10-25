package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type request struct {
	Kind     string             `json:"kind"`
	Donation *Donation          `json:"donation,omitempty"`
	Reminder *VolunteerReminder `json:"reminder,omitempty"`
	Report   *CampaignReport    `json:"report,omitempty"`
}

func main() {
	client, err := NewEmailClient()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var in request
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		var id string
		switch in.Kind {
		case "receipt":
			if in.Donation == nil {
				http.Error(w, "donation is required", 400)
				return
			}
			id, err = sendReceipt(r.Context(), client, *in.Donation)
		case "volunteer_reminder":
			if in.Reminder == nil {
				http.Error(w, "reminder is required", 400)
				return
			}
			id, err = sendReminder(r.Context(), client, *in.Reminder)
		default:
			http.Error(w, "kind must be receipt or volunteer_reminder", 400)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"message_id": id})
	})
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
