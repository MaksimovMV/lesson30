package storage

import (
	"bytes"
	"fmt"
	"lesson30/internal/user"
	"strconv"
)

type Storage struct {
	Repository map[int]*user.User
	UserID     int
}

func NewStorage() Storage {
	return Storage{make(map[int]*user.User), 0}
}

func (s *Storage) GetUser(targetID int) (*user.User, error) {
	u, ok := s.Repository[targetID]
	if !ok {
		return nil, fmt.Errorf("пользователя с ID %v не существует", targetID)
	}
	return u, nil
}

func (s *Storage) PutUser(u user.User) (int, error) {
	s.UserID++
	s.Repository[s.UserID] = &u
	return s.UserID, nil
}

func (s *Storage) MakeFriends(sourceID int, targetID int) error {
	if sourceID == targetID {
		return fmt.Errorf("ID пользователей совпадают")
	}

	sourceUser, err := s.GetUser(sourceID)
	if err != nil {
		return err
	}

	targetUser, err := s.GetUser(targetID)
	if err != nil {
		return err
	}

	for _, f := range sourceUser.Friends {
		if f == targetID {
			return fmt.Errorf("пользователь с ID %v уже есть в друзьях у пользователя c ID %v", sourceID, targetID)
		}
	}
	sourceUser.Friends = append(sourceUser.Friends, targetID)

	for _, f := range targetUser.Friends {
		if f == sourceID {
			return fmt.Errorf("пользователь с ID %v уже есть в друзьях у пользователя c ID %v", targetID, sourceID)
		}
	}
	targetUser.Friends = append(targetUser.Friends, sourceID)
	return nil
}

func (s *Storage) DeleteUser(targetID int) error {

	u, err := s.GetUser(targetID)

	if err != nil {
		return err
	}

	for _, friendID := range u.Friends {
		f := s.Repository[friendID]
		for i, fID := range f.Friends {
			if fID == targetID {
				f.Friends = append(f.Friends[:i], f.Friends[i+1:]...)
				break
			}
		}
	}

	delete(s.Repository, targetID)
	return nil
}

func (s *Storage) GetFriends(targetID int) (bytes.Buffer, error) {
	b := bytes.Buffer{}
	u, err := s.GetUser(targetID)
	if err != nil {
		return b, err
	}
	for _, f := range u.Friends {
		friend := s.Repository[f]
		b.Write([]byte("ID:" + strconv.Itoa(f) + " Имя: " + friend.Name + "\n"))
	}
	return b, nil
}

func (s *Storage) PutNewAge(targetID int, age int) error {
	u, err := s.GetUser(targetID)
	if err != nil {
		return err
	}
	u.Age = age
	return nil
}
