package user

import (
	"fmt"
)

type User struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Friends []int  `json:"friends"`
	id      int
}

func (u *User) toString() string {
	return fmt.Sprintf("name is %s and age is %d \n", u.Name, u.Age)
}

func (u *User) RemoveFriend(id int) error {
	for i, friend := range u.Friends {
		if friend == id {
			u.Friends = append(u.Friends[:i], u.Friends[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("пользователя с ID %v нет в друзьях у пользователя c ID %v", id, u.id)
}

func (u *User) AddFriend(id int) error {
	if u.id == id {
		return fmt.Errorf("ID пользователей совпадают")
	}
	for _, f := range u.Friends {
		if f == id {
			return fmt.Errorf("пользователь с ID %v уже есть в друзьях у пользователя c ID %v", id, u.id)
		}
	}
	u.Friends = append(u.Friends, id)
	return nil
}

func (u *User) AddID(id int) error {
	if u.id == 0 {
		u.id = id
		return nil
	} else {
		return fmt.Errorf("у пользователя уже есть ID")
	}
}
