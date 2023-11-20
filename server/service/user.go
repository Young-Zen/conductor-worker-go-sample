package service

import (
	"worker-sample/config"
	"worker-sample/server/dao"
	"worker-sample/server/form"
	"worker-sample/server/model"
)

type UserService interface {
	AddUser(req form.AddUserReq) (*model.User, error)
	GetUserByName(name string) (*model.User, error)
	GetUserById(id int64) (*model.User, error)
}

type UserServiceImpl struct {
	ServiceContext *config.ServiceContext
	UserDao        dao.UserDao
}

func (s *UserServiceImpl) AddUser(req form.AddUserReq) (*model.User, error) {
	err := s.UserDao.AddUser(&model.User{
		Name: req.Name,
		Age:  req.Age,
	})
	if err != nil {
		return nil, err
	}

	return s.GetUserByName(req.Name)
}

func (s *UserServiceImpl) GetUserByName(name string) (*model.User, error) {
	user, err := s.UserDao.GetUserByName(name)
	if err != nil {
		return nil, err
	}

	if user.Id == 0 {
		return nil, err
	}
	return user, nil
}

func (s *UserServiceImpl) GetUserById(id int64) (*model.User, error) {
	user, err := s.UserDao.GetUserById(id)
	if err != nil {
		return nil, err
	}

	if user.Id == 0 {
		return nil, err
	}
	return user, nil
}
