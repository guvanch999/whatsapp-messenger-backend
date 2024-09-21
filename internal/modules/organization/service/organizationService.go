package service

import (
	"github.com/google/uuid"
	"github.com/medium-messenger/messenger-backend/internal/modules/organization/dto"
	"github.com/medium-messenger/messenger-backend/internal/modules/organization/models"
	"github.com/medium-messenger/messenger-backend/internal/modules/organization/repository"
	auth "github.com/medium-messenger/messenger-backend/internal/modules/users/models"
	"github.com/medium-messenger/messenger-backend/utils/enums"
	"github.com/medium-messenger/messenger-backend/utils/exceptions"
	"github.com/medium-messenger/messenger-backend/utils/util"
)

type OrganizationService struct {
	repository *repository.OrganizationRepository
}

func NewOrganizationService(organizationRepository *repository.OrganizationRepository) *OrganizationService {
	return &OrganizationService{
		organizationRepository,
	}
}

func (s *OrganizationService) GetUserOrganizations(userId uuid.UUID) ([]dto.ResponseOrganizationDto, error) {
	organizations, err := s.repository.GetUserOrganizations(userId)
	if err != nil {
		return nil, err
	}
	return util.Map(
		organizations, func(organization models.Organization) dto.ResponseOrganizationDto {
			return *organization.ToResponseDto()
		},
	), nil
}

func (s *OrganizationService) GetAllOrganizations() ([]dto.ResponseOrganizationDto, error) {
	organizations, err := s.repository.GetAllOrganizations()
	if err != nil {
		return nil, err
	}
	return util.Map(
		organizations, func(organization models.Organization) dto.ResponseOrganizationDto {
			return *organization.ToResponseDto()
		},
	), nil
}

func (s *OrganizationService) CreateOrganization(
	userId uuid.UUID,
	organizationDto dto.OrganizationDto,
) (*dto.ResponseOrganizationDto, error) {
	organizationModel := models.Organization{
		OwnerId: userId,
		Name:    organizationDto.Name,
	}
	contact, err := s.repository.AddOrganization(organizationModel)
	if err != nil {
		return nil, err
	}
	return contact.ToResponseDto(), nil
}

func (s *OrganizationService) UpdateOrganization(
	user auth.UserDetail,
	updateOrgDto dto.UpdateOrganization,
) (*dto.ResponseOrganizationDto, error) {
	orgModel, err := s.checkAccess(user, updateOrgDto.Id)
	if err != nil {
		return nil, err
	}
	orgModel.Name = updateOrgDto.Name
	contact, err := s.repository.UpdateOrganization(*orgModel)
	if err != nil {
		return nil, err
	}
	return contact.ToResponseDto(), nil
}

func (s *OrganizationService) DeleteOrganization(user auth.UserDetail, organizationId uuid.UUID) error {
	_, err := s.checkAccess(user, organizationId)
	if err != nil {
		return err
	}
	return s.repository.DeleteOrganization(organizationId)
}
func (s *OrganizationService) checkAccess(user auth.UserDetail, orgId uuid.UUID) (*models.Organization, error) {
	organization, err := s.repository.GetOrganizationDetail(orgId)
	if err != nil {
		return nil, err
	}
	if user.Role != enums.Admin && organization.OwnerId != user.ID {
		return nil, &exceptions.AccessDenied{}
	}
	return organization, nil
}

func (s *OrganizationService) GetOrganizationDetail(user auth.UserDetail, organizationId uuid.UUID) (
	*dto.ResponseOrganizationDto,
	error,
) {
	contact, err := s.checkAccess(user, organizationId)
	if err != nil {
		return nil, err
	}
	return contact.ToResponseDto(), nil
}
