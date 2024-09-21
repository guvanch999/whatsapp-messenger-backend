package repository

import (
	"errors"
	"github.com/google/uuid"
	. "github.com/medium-messenger/messenger-backend/internal/modules/organization/models"
	"github.com/medium-messenger/messenger-backend/utils/exceptions"
	"gorm.io/gorm"
)

type OrganizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) *OrganizationRepository {
	return &OrganizationRepository{
		db,
	}
}

func (r *OrganizationRepository) GetUserOrganizations(userId uuid.UUID) ([]Organization, error) {
	var organizations []Organization
	if err := r.db.Model(&Organization{}).Where(
		"owner_id = ?",
		userId,
	).Select("*").Scan(&organizations).Error; err != nil {
		return nil, err
	}
	return organizations, nil
}
func (r *OrganizationRepository) GetAllOrganizations() ([]Organization, error) {
	var organizations []Organization
	if err := r.db.Model(&Organization{}).Select("*").Scan(&organizations).Error; err != nil {
		return nil, err
	}
	return organizations, nil
}
func (r *OrganizationRepository) GetOrganizationDetail(id uuid.UUID) (*Organization, error) {
	var organization Organization
	if err := r.db.Model(&Organization{}).Where("id = ?", id).First(&organization).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &exceptions.NotFoundError{}
		}
		return nil, err
	}
	return &organization, nil
}

func (r *OrganizationRepository) AddOrganization(orgModel Organization) (*Organization, error) {
	if err := r.db.Save(&orgModel).Error; err != nil {
		return nil, err
	}
	return &orgModel, nil
}

func (r *OrganizationRepository) UpdateOrganization(orgModel Organization) (*Organization, error) {
	if err := r.db.Model(&Organization{}).Where("id = ?", orgModel.Id).Updates(orgModel).Error; err != nil {
		return nil, err
	}
	return &orgModel, nil
}

func (r *OrganizationRepository) DeleteOrganization(id uuid.UUID) error {
	if err := r.db.Delete(&Organization{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &exceptions.NotFoundError{}
		}
		return err
	}
	return nil
}
