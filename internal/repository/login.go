package repository

import (
	"gohub/init/mysql"
	"gohub/internal/model"
)

func UpdateAUser(newUser *model.User) (err error) {
	function := func() error {
		err := mysql.DB.Model(model.User{}).Where("id = ?", newUser.ID).Updates(newUser).Error
		return err
	}
	err = dbOperations(function, "UpdateAUser")
	return
}
