package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"koriebruh/restful/api/model/domain"
	"log"
)

type AuthRepositoryImpl struct {
	tx *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &AuthRepositoryImpl{
		tx: db,
	}
}

/// OMIT(clause.association) : agar data yg foreignKey tidak terpengaruh

func (repository AuthRepositoryImpl) Register(ctx context.Context, user *domain.User) error {
	//#VALIDATE
	var existUser domain.User
	if err := repository.tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("email = ?", user.Email).
		First(&existUser).Error; err == nil {
		return errors.New("username already exists")
	}

	//#CREATE
	result := repository.tx.Create(&user).Error
	if result != nil {
		log.Fatal("failed create")
	}

	log.Println("success create new user")
	return nil
}

func (repository AuthRepositoryImpl) UpdateAcc(ctx context.Context, id string, user *domain.User) error {
	//#LOCKING
	err := repository.tx.Select("*").Where("id = ?", id).
		Clauses(clause.Locking{Strength: "UPDATE"}).First(&user).Error
	if err != nil {
		return err
	}

	//#UPDATE,KECUALI YG TERKAID DENGAN ASSSOSIASI
	result := repository.tx.Omit(clause.Associations).Where("id=?", id).Updates(&user)
	if result.Error != nil {
		log.Fatal("failed to update data")
		return result.Error
	}

	log.Println("success updated user")
	return nil
}

func (repository AuthRepositoryImpl) DeleteAcc(ctx context.Context, id string) error {
	//#SOFT DELETE AND LOCKING
	var user domain.User

	if err := repository.tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", id).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	if err := repository.tx.Omit(clause.Associations).Delete(&user).Error; err != nil {
		log.Fatal("failed delete data")
		return err
	}

	log.Println("success delete user")
	return nil
}

func (repository AuthRepositoryImpl) FindById(ctx context.Context, id string) (domain.User, error) {
	var user domain.User
	result := repository.tx.Omit(clause.Associations).Take(&user, "id=?", id)
	if result.Error != nil {
		log.Fatal("failed find user by Id")
		return domain.User{}, result.Error
	}

	log.Println("success find by id")
	return user, nil
}

func (repository AuthRepositoryImpl) FindByUserName(ctx context.Context, username string) (domain.User, error) {
	var user domain.User
	result := repository.tx.Omit(clause.Associations).Take(&user, "user_name=?", username)
	if result.Error != nil {
		log.Fatal("failed find user by UserName")
		return domain.User{}, result.Error
	}

	log.Println("success find by UserName")
	return user, nil
}
