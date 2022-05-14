package controller

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io/ioutil"
	"lesson30/internal/model"
	"lesson30/internal/storage"
	"net/http"
	"strconv"
)

func Build(r *chi.Mux, s *storage.Storage) {

	r.Use(middleware.Logger)

	r.Post("/users", func(w http.ResponseWriter, r *http.Request) {
		createUser(w, r, s)
	})

	r.Get("/users/{sourceID}", func(w http.ResponseWriter, r *http.Request) {
		getUser(w, r, s)
	})

	r.Patch("/users/{sourceID}", func(w http.ResponseWriter, r *http.Request) {
		putNewAge(w, r, s)
	})

	r.Delete("/users/{sourceID}", func(w http.ResponseWriter, r *http.Request) {
		deleteUser(w, r, s)
	})

	r.Put("/users/{sourceID}/friends", func(w http.ResponseWriter, r *http.Request) {
		makeFriends(w, r, s)
	})

	r.Get("/users/{sourceID}/friends", func(w http.ResponseWriter, r *http.Request) {
		getFriends(w, r, s)
	})

	r.Delete("/users/{sourceID}/friends/{targetID}", func(w http.ResponseWriter, r *http.Request) {
		deleteFriend(w, r, s)
	})
}

type inputParams struct {
	TargetID int `json:"target_id"`
	Age      int `json:"age"`
}

type message struct {
	Message string `json:"message"`
}

func readParams(w http.ResponseWriter, r *http.Request) (inputParams, error) {
	var p inputParams

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
func createUser(w http.ResponseWriter, r *http.Request, s *storage.Storage) {
	w.Header().Add("Content-Type", "application/json")
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(message{err.Error()})
		w.Write(b)
		return
	}
	defer r.Body.Close()

	var u model.User

	if err := json.Unmarshal(content, &u); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(message{err.Error()})
		w.Write(b)
		return
	}

	userID, err := s.PutUser(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(message{err.Error()})
		w.Write(b)
		return
	}

	result := message{"Username_" + strconv.Itoa(userID) + " was created"}
	b, _ := json.Marshal(result)
	w.Write(b)
	w.WriteHeader(http.StatusCreated)
}

func getUser(w http.ResponseWriter, r *http.Request, s *storage.Storage) {
	w.Header().Add("Content-Type", "application/json")

	sourceID, err := strconv.Atoi(chi.URLParam(r, "sourceID"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(message{err.Error()})
		w.Write(b)
		return
	}

	u, err := s.GetUser(sourceID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		b, _ := json.Marshal(message{err.Error()})
		w.Write(b)
		return
	}

	b, _ := json.Marshal(u)
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func makeFriends(w http.ResponseWriter, r *http.Request, s *storage.Storage) {
	w.Header().Add("Content-Type", "application/json")
	sourceID, err := strconv.Atoi(chi.URLParam(r, "sourceID"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(message{err.Error()})
		w.Write(b)
		return
	}

	p, err := readParams(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(message{err.Error()})
		w.Write(b)
		return
	}

	if err := s.MakeFriends(sourceID, p.TargetID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		b, _ := json.Marshal(message{err.Error()})
		w.Write(b)
		return
	}

	result := message{"Username_" + strconv.Itoa(p.TargetID) + " и username_" + strconv.Itoa(sourceID) + " теперь друзья"}
	b, _ := json.Marshal(result)
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func deleteUser(w http.ResponseWriter, r *http.Request, s *storage.Storage) {
	w.Header().Add("Content-Type", "application/json")
	sourceID, err := strconv.Atoi(chi.URLParam(r, "sourceID"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(message{err.Error()})
		w.Write(b)
		return
	}

	if err := s.DeleteUser(sourceID); err != nil {
		w.WriteHeader(http.StatusNotFound)
		b, _ := json.Marshal(message{err.Error()})
		w.Write(b)
		return
	}

	result := message{"Username_" + strconv.Itoa(sourceID) + " удален"}
	b, _ := json.Marshal(result)
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func getFriends(w http.ResponseWriter, r *http.Request, s *storage.Storage) {
	w.Header().Add("Content-Type", "application/json")

	sourceID, err := strconv.Atoi(chi.URLParam(r, "sourceID"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(message{err.Error()})
		w.Write(b)
		return
	}

	friends, err := s.GetFriends(sourceID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		b, _ := json.Marshal(message{err.Error()})
		w.Write(b)
		return
	}

	b, _ := json.Marshal(friends)
	b, _ = json.Marshal(message{string(b)})
	w.Write(b)
}

func putNewAge(w http.ResponseWriter, r *http.Request, s *storage.Storage) {
	w.Header().Add("Content-Type", "application/json")
	p, err := readParams(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(message{err.Error()})
		w.Write(b)
		return
	}

	sourceID, err := strconv.Atoi(chi.URLParam(r, "sourceID"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(message{err.Error()})
		w.Write(b)
		return
	}

	if err := s.PutNewAge(sourceID, p.Age); err != nil {
		w.WriteHeader(http.StatusNotFound)
		b, _ := json.Marshal(message{err.Error()})
		w.Write(b)
		return
	}

	result := message{"Возраст пользователя Username_" + strconv.Itoa(sourceID) + " успешно обновлен"}
	b, _ := json.Marshal(result)
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func deleteFriend(w http.ResponseWriter, r *http.Request, s *storage.Storage) {
	w.Header().Add("Content-Type", "application/json")
	sourceID, err := strconv.Atoi(chi.URLParam(r, "sourceID"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(message{err.Error()})
		w.Write(b)
		return
	}

	targetID, err := strconv.Atoi(chi.URLParam(r, "targetID"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(message{err.Error()})
		w.Write(b)
		return
	}

	if err := s.DeleteFriend(sourceID, targetID); err != nil {
		w.WriteHeader(http.StatusNotFound)
		b, _ := json.Marshal(message{err.Error()})
		w.Write(b)
		return
	}

	result := message{"Username_" + strconv.Itoa(targetID) + " и username_" + strconv.Itoa(sourceID) + " больше не друзья"}
	b, _ := json.Marshal(result)
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
