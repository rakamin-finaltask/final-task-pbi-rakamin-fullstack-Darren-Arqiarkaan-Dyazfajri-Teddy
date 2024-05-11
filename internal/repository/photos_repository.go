package repository

import (
	"database/sql"
	"fmt"
	"user-personalize/internal/model/entity"
)

type PhotosRepository interface {
	Insert(photos entity.Photos) (entity.Photos, error)
	Update(photos entity.Photos) (entity.Photos, error)
	Delete(userId string) error
	FindByUserId(userId string) (entity.Photos, error)
	FindById(id string) (entity.Photos, error)
}

type photosRepositoryImpl struct {
	db *sql.DB
}

func (p *photosRepositoryImpl) FindById(id string) (entity.Photos, error) {
	query := "select id, title, caption, photo_url, user_id from photos where id = $1"

	var result entity.Photos
	err := p.db.QueryRow(query, id).Scan(&result.Id, &result.Title, &result.Caption, &result.PhotoUrl, &result.UserId)
	if err != nil {
		return entity.Photos{}, fmt.Errorf("FindByIdRepository : %w", err)
	}

	return result, nil
}

func (p *photosRepositoryImpl) FindByUserId(userId string) (entity.Photos, error) {
	query := "select id, title, caption, photo_url, user_id from photos where user_id = $1"

	var photosEntity entity.Photos
	err := p.db.QueryRow(query, userId).Scan(&photosEntity.Id, &photosEntity.Title, &photosEntity.Caption, &photosEntity.PhotoUrl, &photosEntity.UserId)

	if err != nil {
		return entity.Photos{}, fmt.Errorf("PhotosFindByUserIdRepository : %v", err)
	}

	return photosEntity, nil
}

func (p *photosRepositoryImpl) Insert(photos entity.Photos) (entity.Photos, error) {
	query := "insert into photos (id, title, caption, photo_url, user_id) values($1, $2, $3, $4, $5) returning id, title, caption, photo_url, user_id"

	var photosEntity entity.Photos
	err := p.db.QueryRow(query, photos.Id, photos.Title, photos.Caption, photos.PhotoUrl, photos.UserId).Scan(&photosEntity.Id, &photosEntity.Title, &photosEntity.Caption, &photosEntity.PhotoUrl, &photosEntity.UserId)

	if err != nil {
		return entity.Photos{}, fmt.Errorf("insertPhotosRepository : %v", err)
	}

	return photosEntity, nil
}

func (p *photosRepositoryImpl) Update(photos entity.Photos) (entity.Photos, error) {
	query := "update photos set title = $1, caption = $2, photo_url = $3 where user_id = $4 returning id, title, caption, photo_url, user_id"

	var photosEntity entity.Photos
	err := p.db.QueryRow(query, photos.Title, photos.Caption, photos.PhotoUrl, photos.UserId).Scan(&photosEntity.Id, &photosEntity.Title, &photosEntity.Caption, &photosEntity.PhotoUrl, &photosEntity.UserId)

	if err != nil {
		return entity.Photos{}, fmt.Errorf("updatePhotosRepository : %v", err)
	}

	return photosEntity, nil
}

func (p *photosRepositoryImpl) Delete(userId string) error {
	query := "delete from photos where id = $1"

	_, err := p.db.Exec(query, userId)
	if err != nil {
		return fmt.Errorf("deletePhotosRepository : %v", err)
	}

	return nil
}

func NewPhotosRepository(db *sql.DB) PhotosRepository {
	return &photosRepositoryImpl{db: db}
}
