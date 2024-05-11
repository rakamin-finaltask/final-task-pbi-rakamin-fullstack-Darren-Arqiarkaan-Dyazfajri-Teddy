package repository

import (
	"database/sql"
	"fmt"
	"user-personalize/internal/model/entity"
)

type UserRepository interface {
	Create(user entity.User) (entity.User, error)
	Update(user entity.User) (entity.User, error)
	Delete(id string) error
	GetAll() ([]entity.User, error)
	GetById(id string) (entity.User, error)
	UpdatePassword(id string, newPassword string) (entity.User, error)
	GetByEmail(email string) (entity.User, error)
}

type userRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (u *userRepositoryImpl) GetByEmail(email string) (entity.User, error) {
	query := "select id, username, email, password, created_at, updated_at from users where email = $1"

	var user entity.User
	err := u.db.QueryRow(query, email).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return entity.User{}, fmt.Errorf("GetUserByEmailRepository: %w", err)
	}

	return user, nil
}

func (u *userRepositoryImpl) Create(user entity.User) (entity.User, error) {
	query := "insert into users (id, username, email, password, created_at, updated_at) values ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) returning id, username, email, created_at, updated_at"

	var result entity.User
	err := u.db.QueryRow(query, user.Id, user.Username, user.Email, user.Password).Scan(&result.Id, &result.Username, &result.Email, &result.CreatedAt, &result.UpdatedAt)

	if err != nil {
		return result, fmt.Errorf("CreateRepository: %w", err)
	}

	return result, nil
}

func (u *userRepositoryImpl) Update(user entity.User) (entity.User, error) {
	query := "update users set username = $1, email = $2, updated_at = CURRENT_TIMESTAMP where id = $3 returning id, username, email, created_at, updated_at"

	var result entity.User
	err := u.db.QueryRow(query, user.Username, user.Email, user.Id).Scan(&result.Id, &result.Username, &result.Email, &result.CreatedAt, &result.UpdatedAt)

	if err != nil {
		return entity.User{}, fmt.Errorf("UpdateRepository: %w", err)
	}

	return result, nil
}

func (u *userRepositoryImpl) Delete(id string) error {
	query := "delete from users where id = $1"
	_, err := u.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("DeleteRepository: %w", err)
	}
	return nil
}

func (u *userRepositoryImpl) GetAll() ([]entity.User, error) {
	query := "select id, username, email, created_at, updated_at from users"
	rows, err := u.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("GetAllRepository: %w", err)
	}

	defer rows.Close()
	users := make([]entity.User, 0)
	for rows.Next() {
		var user entity.User

		err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("GetAllRepository: %w", err)
		}

		users = append(users, user)
	}

	return users, nil
}

func (u *userRepositoryImpl) GetById(id string) (entity.User, error) {
	query := "select id, username, email, created_at, updated_at from users where id = $1"

	var user entity.User
	err := u.db.QueryRow(query, id).Scan(&user.Id, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return entity.User{}, fmt.Errorf("GetByIdRepository: %w", err)
	}

	return user, nil
}

func (u *userRepositoryImpl) UpdatePassword(id string, newPassword string) (entity.User, error) {
	query := "update users set password = $1 where id = $2 returning id, username, email, created_at, updated_at"

	var result entity.User

	err := u.db.QueryRow(query, newPassword, id).Scan(&result.Id, &result.Username, &result.Email, &result.CreatedAt, &result.UpdatedAt)

	if err != nil {
		return entity.User{}, fmt.Errorf("UpdatePasswordRepository: %w", err)
	}

	return result, nil
}
