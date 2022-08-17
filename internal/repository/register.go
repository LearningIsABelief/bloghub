package repository

import (
	"fmt"
	"gohub/init/mysql"
	"gohub/internal/model"
	"gorm.io/gorm"
)

var funcNameMap = map[string]struct{}{}

func UserExist(queryFieldVal string, phoneOrEmailOrName int) (user *model.User, exist bool, err error) {
	user = &model.User{}
	queryField := "phone"
	switch phoneOrEmailOrName {
	case 2:
		queryField = "email"
	case 3:
		queryField = "name"
	}
	function := func() error {
		err := mysql.DB.Where(fmt.Sprintf("%v = ?", queryField), queryFieldVal).First(user).Error
		return err
	}
	err = dbOperations(function, "UserExist")
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}
	return user, true, nil
}

func CreateAUser(user *model.User) (err error) {
	function := func() error {
		err := mysql.DB.Create(user).Error
		return err
	}
	err = dbOperations(function, "CreateAUser")
	return
}
