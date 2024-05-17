package handlers

import (
	"errors"
	"github.com/EupravaProjekat/prosecution/Repo"
	"github.com/EupravaProjekat/prosecution/Models"
	protos "github.com/MihajloJankovic/profile-service/protos/main"
	"github.com/gorilla/mux"
	"log"
	"mime"
	"net/http"
	"encoding/json"
)

type ProsecutionHandler struct {
	l    *log.Logger
	repo *Repo.Repo
}

func NewProsecutionHandler(l *log.Logger, r *Repo.Repo) *ProsecutionHandler {
	return &ProsecutionHandler{l, r}

}

func (h *ProsecutionHandler) CheckIfUserExists(w http.ResponseWriter, r *http.Request) {
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
func (h *ProsecutionHandler) CheckIfPersonIsProsecuted(w http.ResponseWriter, r *http.Request) {
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

	requestBody := struct {
		JMBG string `json:"jmbg"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if requestBody.JMBG == "" {
		err := errors.New("empty JMBG")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	prosecutions, err := h.repo.GetProsecutionsByJmbg(requestBody.JMBG)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hasProsecutions := len(prosecutions) > 0

	response := struct {
		Prosecuted bool `json:"prosecuted"`
	}{
		Prosecuted: hasProsecutions,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}


func (h *ProsecutionHandler) ProsecuteHandler(w http.ResponseWriter, r *http.Request) {
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

    // Decode prosecution data from request body
    prosecutionData, err := DecodeProsecutionBody(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Create a new prosecution instance
    prosecution := &Models.Prosecution{
        JMBG:         prosecutionData.JMBG,
        TypeOfBreach: prosecutionData.TypeOfBreach,
    }

    // Call CreateProsecution method from Repo
    err = h.repo.CreateProsecution(prosecution)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Respond with success
    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("Prosecution created successfully"))
}






func (h *ProsecutionHandler) GetProfile(w http.ResponseWriter, r *http.Request) {

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

func (h *ProsecutionHandler) NewRequest(w http.ResponseWriter, r *http.Request) {

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
