package service

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"lesson30/pkg/user"
	"net/http"
	"strconv"
)

type Service struct {
	Store  map[int]*user.User
	UserID int
}

type InputParams struct {
	SourceID int `json:"source_id"`
	TargetID int `json:"target_id"`
	NewAge   int `json:"new_age"`
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

func (s Service) GetUser(id int) (*user.User, error) {
	u, ok := s.Store[id]
	if !ok {
		return nil, fmt.Errorf("пользователя с ID %v не существует", id)
	}
	return u, nil
}

func (s *Service) Create(w http.ResponseWriter, r *http.Request) {
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

	s.UserID++
	s.Store[s.UserID] = &u
	if err := u.AddID(s.UserID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Username_" + strconv.Itoa(s.UserID) + " was created"))
}

func (s *Service) MakeFriends(w http.ResponseWriter, r *http.Request) {
	p, err := readParams(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	sourceUser, err := s.GetUser(p.SourceID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	targetUser, err := s.GetUser(p.TargetID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if err := sourceUser.AddFriend(p.TargetID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if err := targetUser.AddFriend(p.SourceID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Username_" + strconv.Itoa(p.SourceID) + " и username_" + strconv.Itoa(p.TargetID) + " теперь друзья"))
}

func (s *Service) DeleteUser(w http.ResponseWriter, r *http.Request) {
	p, err := readParams(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	u, err := s.GetUser(p.TargetID)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	for _, i := range u.Friends {
		fr := s.Store[i]
		fr.RemoveFriend(p.TargetID)
	}

	delete(s.Store, p.TargetID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Username_" + strconv.Itoa(p.TargetID) + " удален"))
}

func (s *Service) GetFriends(w http.ResponseWriter, r *http.Request) {
	targetID, err := strconv.Atoi(chi.URLParam(r, "targetID"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	u, err := s.GetUser(targetID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	for _, f := range u.Friends {
		friend := s.Store[f]
		w.Write([]byte("ID:" + strconv.Itoa(f) + " Имя: " + friend.Name + "\n"))
	}
}

func (s *Service) PutNewAge(w http.ResponseWriter, r *http.Request) {
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

	u, err := s.GetUser(targetID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	u.Age = p.NewAge

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Возраст пользователя Username_" + strconv.Itoa(targetID) + " успешно обновлен"))
}
