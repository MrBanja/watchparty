package ssr

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Service) GetStatusBadge(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomName := vars["room_id"]
	pID := vars["participant_id"]

	room := s.party.GetRoom(roomName)
	if room == nil {
		if err := s.tmpl.ExecuteTemplate(w, statusBadgeTmpl, m{"Address": s.publicAddress, "IsLost": true, "Text": "Empty room", "ID": roomName, "PartID": pID}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	p := room.GetParticipantByID(pID)
	if p == nil {
		if err := s.tmpl.ExecuteTemplate(w, statusBadgeTmpl, m{"Address": s.publicAddress, "IsLost": true, "Text": "Lost Connection", "ID": roomName}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if err := s.tmpl.ExecuteTemplate(w, statusBadgeTmpl, m{"Address": s.publicAddress, "IsLost": false, "Text": "Connected", "ID": roomName, "PartID": p.ID}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return
}
