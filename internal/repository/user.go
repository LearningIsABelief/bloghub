package repository

import (
	"gohub/init/mysql"
	"gohub/internal/model"
)

func User(userID uint) (user *model.User, err error) {
	user = &model.User{ID: userID}
	function := func() error {
		err := mysql.DB.First(user).Error
		return err
	}
	err = dbOperations(function, "User")
	if err != nil {
		return nil, err
	}
	return user, nil
}

func All() ([]*model.User, error) {
	var users []*model.User
	function := func() error {
		err := mysql.DB.Find(&users).Error
		return err
	}
	err := dbOperations(function, "All")
	if err != nil {
		return nil, err
	}
	return users, nil
}

func Update(userID uint, field, fieldVal string) error {
	function := func() error {
		err := mysql.DB.Model(&model.User{}).Where("id = ?", userID).Update(field, fieldVal).Error
		return err
	}
	err := dbOperations(function, "Update")
	if err != nil {
		return err
	}
	return nil
}
