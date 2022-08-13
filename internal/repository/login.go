package repository

import (
	"gohub/init/mysql"
	"gohub/internal/model"
)

func UpdateAUser(newUser *model.User) error {
	function := func() error {
		err := mysql.DB.Model(&model.User{}).Updates(newUser).Error
		return err
	}
	err := dbOperations(function, "UpdateAUser")
	if err != nil {
		return err
	}
	return nil
}
