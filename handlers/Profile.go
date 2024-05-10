package handlers

import (
	"errors"
	"github.com/MihajloJankovic/border-police/Repo"
	protos "github.com/MihajloJankovic/profile-service/protos/main"
	"github.com/gorilla/mux"
	"log"
	"mime"
	"net/http"
)

type Borderhendler struct {
	l    *log.Logger
	repo *Repo.Repo
}

func NewBorderhendler(l *log.Logger, r *Repo.Repo) *Borderhendler {
	return &Borderhendler{l, r}

}

func (h *Borderhendler) CheckIfUserExists(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}
	res := ValidateJwt(r, h.repo)
	if res == nil {
		err := errors.New("jwt error")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	re := res
	response, err := h.repo.GetByEmail(re.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if re.Email != response.Email {
		err := errors.New("authorization error")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

}

func (h *Borderhendler) GetProfile(w http.ResponseWriter, r *http.Request) {

	emaila := mux.Vars(r)["email"]
	ee := new(protos.ProfileRequest)
	ee.Email = emaila
	res := ValidateJwt(r, h.repo)
	if res == nil {
		err := errors.New("jwt error")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	re := res
	responsea, err := h.repo.GetByEmail(re.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if re.Email != responsea.Email {
		err := errors.New("authorization error")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	if re.Email != ee.GetEmail() {
		err := errors.New("authorization error")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	response, err := h.repo.GetByEmail(ee.Email)
	if err != nil || response == nil {
		log.Printf("Operation Failed: %v\n", err)
		w.WriteHeader(http.StatusNotAcceptable)
		_, err := w.Write([]byte("Profile not found"))
		if err != nil {
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	RenderJSON(w, response)
}

func (h *Borderhendler) NewRequest(w http.ResponseWriter, r *http.Request) {

	contentType := r.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}
	rt, err := DecodeBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusAccepted)
		return
	}
	res := ValidateJwt(r, h.repo)
	if res == nil {
		err := errors.New("jwt error")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	re := res
	responsea, err := h.repo.GetByEmail(re.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if re.Email != responsea.Email {
		err := errors.New("authorization error")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	err = h.repo.NewRequest(*rt,re.Email)
	if err != nil {
		log.Printf("Operation failed: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("couldn't add request"))
		if err != nil {
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("successfully added request"))
	if err != nil {
		return
	}
}
