package usecase

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"lesson30/internal/storage"
	"lesson30/internal/user"
	"net/http"
	"strconv"
)

type UseCase struct {
	storage.Storage
}

func NewUseCase(s storage.Storage) UseCase {
	return UseCase{s}
}

type InputParams struct {
	SourceID int `json:"source_id"`
	TargetID int `json:"target_id"`
	NewAge   int `json:"new_age"`
}

type message struct {
	Message string `json:"message"`
}

func readParams(w http.ResponseWriter, r *http.Request) (InputParams, error) {
	var p InputParams

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return p, err
	}
	defer r.Body.Close()

	if err := json.Unmarshal(content, &p); err != nil {
		return p, err
	}

	return p, nil
}
func (uc *UseCase) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer r.Body.Close()

	var u user.User

	if err := json.Unmarshal(content, &u); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	userID, err := uc.PutUser(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	result := message{"Username_" + strconv.Itoa(userID) + " was created"}
	b, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(b)
	w.WriteHeader(http.StatusCreated)
}

func (uc *UseCase) Make(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	p, err := readParams(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if err := uc.MakeFriends(p.SourceID, p.TargetID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	result := &message{"Username_" + strconv.Itoa(p.SourceID) + " и username_" + strconv.Itoa(p.TargetID) + " теперь друзья"}
	b, err := json.Marshal(result)

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (uc *UseCase) User(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	p, err := readParams(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	uc.DeleteUser(p.TargetID)

	result := &message{"Username_" + strconv.Itoa(p.TargetID) + " удален"}
	b, err := json.Marshal(result)

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (uc *UseCase) Friends(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	targetID, err := strconv.Atoi(chi.URLParam(r, "targetID"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	b, err := uc.GetFriends(targetID)
	w.Write(b.Bytes())
}

func (uc *UseCase) NewAge(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	p, err := readParams(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	targetID, err := strconv.Atoi(chi.URLParam(r, "targetID"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if err := uc.PutNewAge(targetID, p.NewAge); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	result := &message{"Возраст пользователя Username_" + strconv.Itoa(targetID) + " успешно обновлен"}
	b, err := json.Marshal(result)

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
