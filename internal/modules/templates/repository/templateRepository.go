package repository

import (
	"errors"
	"github.com/google/uuid"
	. "github.com/medium-messenger/messenger-backend/internal/modules/templates/models"
	"github.com/medium-messenger/messenger-backend/utils/exceptions"
	"gorm.io/gorm"
	"time"
)

type TemplateRepository struct {
	db *gorm.DB
}

func NewTemplateRepository(db *gorm.DB) *TemplateRepository {
	return &TemplateRepository{
		db,
	}
}

func (r *TemplateRepository) GetUserTemplates(userId uuid.UUID) ([]Template, error) {
	var list []Template
	if err := r.db.Model(&Template{}).Select("*").Where(
		"user_id = ?",
		userId,
	).Scan(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *TemplateRepository) GetTemplatesByProvider(providerId uuid.UUID) ([]Template, error) {
	var list []Template
	if err := r.db.Model(&Template{}).Select("*").Where(
		"provider_id = ?",
		providerId,
	).Scan(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *TemplateRepository) GetAllTemplates() ([]Template, error) {
	var list []Template
	if err := r.db.Model(&Template{}).Select("*").Scan(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *TemplateRepository) GetTemplatesBeforeTime(currentTime time.Time) ([]Template, error) {
	var list []Template
	if err := r.db.Model(&Template{}).Select("*").Where("next_check < ?", currentTime).Scan(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *TemplateRepository) GetDetail(templateId uuid.UUID) (*Template, error) {
	var templateDetail Template
	if err := r.db.Model(&Template{}).Select("*").Where(
		"id = ?",
		templateId,
	).First(&templateDetail).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &exceptions.NotFoundError{}
		}
		return nil, err
	}
	return &templateDetail, nil
}

func (r *TemplateRepository) AddTemplate(template Template) (*Template, error) {
	if err := r.db.Create(&template).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *TemplateRepository) UpdateTemplate(template Template) (*Template, error) {
	if err := r.db.Model(&Template{}).Where("id = ?", template.Id).Updates(template).Error; err != nil {
		return nil, err
	}
	return &template, nil
}
func (r *TemplateRepository) UpdateTemplateWithUpdates(tempId uuid.UUID, updates map[string]any) error {
	if err := r.db.Model(&Template{}).Where("id = ?", tempId).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}
func (r *TemplateRepository) DeleteTemplate(templateId uuid.UUID) error {
	if err := r.db.Delete(&Template{}, templateId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &exceptions.NotFoundError{}
		}
		return err
	}
	return nil
}
