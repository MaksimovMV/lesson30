package storage

import (
	"fmt"
	"lesson30/internal/model"
	"sync"
)

type Storage struct {
	Repository map[int]*model.User
	UserID     int
	sync.Mutex
}

func NewStorage() Storage {
	return Storage{make(map[int]*model.User), 0, sync.Mutex{}}
}

func (s *Storage) GetUser(targetID int) (*model.User, error) {
	u, ok := s.Repository[targetID]
	if !ok {
		return nil, fmt.Errorf("пользователя с ID %v не существует", targetID)
	}
	return u, nil
}

func (s *Storage) PutUser(u model.User) (int, error) {
	s.Lock()
	defer s.Unlock()

	s.UserID++
	s.Repository[s.UserID] = &u
	return s.UserID, nil
}

func (s *Storage) MakeFriends(sourceID int, targetID int) error {
	s.Lock()
	defer s.Unlock()

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

	for _, f := range targetUser.Friends {
		if f == sourceID {
			return fmt.Errorf("пользователь с ID %v уже есть в друзьях у пользователя c ID %v", targetID, sourceID)
		}
	}

	sourceUser.Friends = append(sourceUser.Friends, targetID)
	targetUser.Friends = append(targetUser.Friends, sourceID)
	return nil
}

func (s *Storage) DeleteUser(targetID int) error {
	s.Lock()
	defer s.Unlock()

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

func (s *Storage) GetFriends(targetID int) ([]int, error) {
	s.Lock()
	defer s.Unlock()

	u, err := s.GetUser(targetID)
	if err != nil {
		return nil, err
	}
	return u.Friends, nil
}

func (s *Storage) PutNewAge(targetID int, age int) error {
	s.Lock()
	defer s.Unlock()

	u, err := s.GetUser(targetID)
	if err != nil {
		return err
	}
	u.Age = age
	return nil
}

func (s *Storage) DeleteFriend(sourceID int, targetID int) error {
	s.Lock()
	defer s.Unlock()

	sourceUser, err := s.GetUser(sourceID)
	if err != nil {
		return err
	}

	targetUser, err := s.GetUser(targetID)
	if err != nil {
		return err
	}

	for i, fID := range sourceUser.Friends {
		if fID == targetID {
			sourceUser.Friends = append(sourceUser.Friends[:i], sourceUser.Friends[i+1:]...)
			break
		}
	}

	for i, fID := range targetUser.Friends {
		if fID == sourceID {
			targetUser.Friends = append(targetUser.Friends[:i], targetUser.Friends[i+1:]...)
			break
		}
	}

	return nil
}
