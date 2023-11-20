package dao

import (
	"worker-sample/config"
	"worker-sample/server/model"
)

type UserDao interface {
	AddUser(user *model.User) error
	GetUserByName(name string) (*model.User, error)
	GetUserById(id int64) (*model.User, error)
}

type UserRepo struct {
	ServiceContext *config.ServiceContext
}

func (d *UserRepo) AddUser(user *model.User) error {
	return d.ServiceContext.DB.Create(user).Error
}

func (d *UserRepo) GetUserByName(name string) (*model.User, error) {
	user := &model.User{}
	tx := d.ServiceContext.DB.Where("name = ?", name).Find(user)
	return user, tx.Error
}

func (d *UserRepo) GetUserById(id int64) (*model.User, error) {
	user := &model.User{}
	tx := d.ServiceContext.DB.Where("id = ?", id).Find(user)
	return user, tx.Error
}
