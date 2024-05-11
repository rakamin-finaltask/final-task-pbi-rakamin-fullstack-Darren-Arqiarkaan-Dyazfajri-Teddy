package usecase

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"user-personalize/internal/model/dto"
	"user-personalize/internal/model/entity"
	"user-personalize/internal/repository"
	"user-personalize/pkg/util/exception"
)

type PhotosUC interface {
	SavePhotos(photos entity.Photos) (dto.PhotosResponse, error)
	UpdatePhotos(payload entity.Photos) (dto.PhotosResponse, error)
	DeletePhotos(photos entity.Photos) error
	GetPhotosByUserId(userId string) (dto.PhotosResponse, error)
}

type photosUCImpl struct {
	photosRepository repository.PhotosRepository
}

func (p *photosUCImpl) GetPhotosByUserId(userId string) (dto.PhotosResponse, error) {
	photoByUserId, err := p.photosRepository.FindByUserId(userId)
	if err != nil {
		return dto.PhotosResponse{}, exception.NotFoundErr
	}
	return dto.PhotosResponse(photoByUserId), nil
}

func (p *photosUCImpl) SavePhotos(photos entity.Photos) (dto.PhotosResponse, error) {
	id := uuid.NewString()
	photos.Id = id
	photosInserted, err := p.photosRepository.Insert(photos)
	if err != nil {
		return dto.PhotosResponse{}, err
	}

	return dto.PhotosResponse{
		Id:       photosInserted.Id,
		Title:    photosInserted.Title,
		Caption:  photosInserted.Caption,
		PhotoUrl: photosInserted.PhotoUrl,
		UserId:   photosInserted.UserId,
	}, nil
}

func (p *photosUCImpl) UpdatePhotos(photos entity.Photos) (dto.PhotosResponse, error) {
	photosByUserId, err := p.photosRepository.FindByUserId(photos.UserId)
	if err != nil {
		return dto.PhotosResponse{}, exception.NotFoundErr
	}

	if photos.Id != photosByUserId.Id {
		return dto.PhotosResponse{}, exception.NotFoundErr
	}

	err = os.Remove(photosByUserId.PhotoUrl)
	if err != nil {
		return dto.PhotosResponse{}, fmt.Errorf("UpdatePhotosUC : %w", err)
	}

	photosUpdated, err := p.photosRepository.Update(photos)
	if err != nil {
		return dto.PhotosResponse{}, fmt.Errorf("UpdatePhotosUC : %w", err)
	}

	return dto.PhotosResponse(photosUpdated), nil
}

func (p *photosUCImpl) DeletePhotos(photos entity.Photos) error {
	photosByUserId, err := p.photosRepository.FindByUserId(photos.UserId)
	if err != nil {
		return exception.NotFoundErr
	}

	if photos.Id != photosByUserId.Id {
		return exception.NotFoundErr
	}

	err = p.photosRepository.Delete(photosByUserId.Id)
	if err != nil {
		return err
	}

	err = os.Remove(photosByUserId.PhotoUrl)
	if err != nil {
		return err

	}

	return nil
}

func NewPhotosUC(photosRepository repository.PhotosRepository) PhotosUC {
	return &photosUCImpl{photosRepository: photosRepository}
}
