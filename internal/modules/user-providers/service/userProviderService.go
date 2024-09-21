package service

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"context"
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/medium-messenger/messenger-backend/internal/config"
	"github.com/medium-messenger/messenger-backend/internal/modules/user-providers/dto"
	"github.com/medium-messenger/messenger-backend/internal/modules/user-providers/model"
	"github.com/medium-messenger/messenger-backend/internal/modules/user-providers/repo"
	auth "github.com/medium-messenger/messenger-backend/internal/modules/users/models"
	"github.com/medium-messenger/messenger-backend/utils/enums"
	"github.com/medium-messenger/messenger-backend/utils/exceptions"
	"github.com/medium-messenger/messenger-backend/utils/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"log"
)

type UserProviderService struct {
	repository          *repo.UserProviderRepository
	secretManagerClient *secretmanager.Client
	cnf                 *config.Schema
	projectId           string
	secretNameTemplate  string
}

func NewUserProviderService(
	repository *repo.UserProviderRepository,
	client *secretmanager.Client,
	cnf *config.Schema,
) *UserProviderService {
	projectId, err := util.GetProjectIdFromGCred(cnf.SecretManagerCredentials)
	if err != nil {
		log.Fatalf("cannot read project_id from credentials: %v", err.Error())
	}
	const template = "projects/%s/secrets/medium-messenger-%s-%s"
	return &UserProviderService{
		repository,
		client,
		cnf,
		projectId,
		template,
	}
}

func (s *UserProviderService) GetUserProviders(userId uuid.UUID) ([]dto.ResponseProviderDto, error) {
	providers, err := s.repository.GetUserProviders(userId)
	if err != nil {
		return nil, err
	}
	return util.Map(
		providers, func(provider model.UserProvider) dto.ResponseProviderDto {
			return *provider.ToResponseDto(s.cnf)
		},
	), nil
}

func (s *UserProviderService) GetAllProviders() ([]dto.ResponseProviderDto, error) {
	providers, err := s.repository.GetAllProviders()
	if err != nil {
		return nil, err
	}
	return util.Map(
		providers, func(provider model.UserProvider) dto.ResponseProviderDto {
			return *provider.ToResponseDto(s.cnf)
		},
	), nil
}

func (s *UserProviderService) CreateProvider(
	userId uuid.UUID,
	providerDto dto.UserProviderDto,
	credentials any,
) (*dto.ResponseProviderDto, error) {
	cred := credentials.(*dto.TwilioCredDto)

	provider := model.UserProvider{
		UserID:          userId,
		Name:            providerDto.Name,
		Type:            providerDto.Type,
		FromPhoneNumber: cred.TwilioFromPhoneNumber,
		Status:          enums.Approved,
	}

	twilioCred := model.TwilioCred{
		TwilioAccountSid:          cred.TwilioAccountSid,
		TwilioAuthToken:           cred.TwilioAuthToken,
		TwilioMessagingServiceSid: cred.TwilioMessagingServiceSid,
	}

	secret, err := s.saveCredentialsInGS(twilioCred)
	if err != nil {
		return nil, err
	}
	provider.ProviderCredentials = secret
	userProvider, err := s.repository.AddProvider(provider)
	if err != nil {
		return nil, err
	}
	return userProvider.ToResponseDto(s.cnf), nil
}

func (s *UserProviderService) DeleteProvider(user auth.UserDetail, providerId uuid.UUID) error {
	provider, err := s.checkAccess(user, providerId)
	if err != nil {
		return err
	}
	err = s.removeCredentialsFromGS(provider.ProviderCredentials)
	if err != nil {
		return err
	}
	return s.repository.DeleteProvider(providerId)
}

func (s *UserProviderService) checkAccess(user auth.UserDetail, provId uuid.UUID) (*model.UserProvider, error) {
	provider, err := s.repository.GetDetail(provId)
	if err != nil {
		return nil, err
	}
	if user.Role != enums.Admin && provider.UserID != user.ID {
		return nil, &exceptions.AccessDenied{}
	}
	return provider, nil
}

func (s *UserProviderService) GetProviderDetail(user auth.UserDetail, providerId uuid.UUID) (
	*dto.ResponseProviderDto,
	error,
) {
	provider, err := s.checkAccess(user, providerId)
	if err != nil {
		return nil, err
	}
	return provider.ToResponseDto(s.cnf), nil
}

func (s *UserProviderService) GetProviderWithCredentials(user auth.UserDetail, providerId uuid.UUID) (
	*model.UserProvider,
	error,
) {
	provider, err := s.checkAccess(user, providerId)
	if err != nil {
		return nil, err
	}
	return provider, nil
}

func (s *UserProviderService) saveCredentialsInGS(cred any) (string, error) {
	hash := util.NewSHA1Hash()
	secretName := fmt.Sprintf(s.secretNameTemplate, s.projectId, s.cnf.BranchName, hash)
	createSecretReq := &secretmanagerpb.CreateSecretRequest{
		Parent:   fmt.Sprintf("projects/%s", s.projectId),
		SecretId: fmt.Sprintf("medium-messenger-%s-%s", s.cnf.BranchName, hash),
		Secret: &secretmanagerpb.Secret{
			Replication: &secretmanagerpb.Replication{
				Replication: &secretmanagerpb.Replication_Automatic_{
					Automatic: &secretmanagerpb.Replication_Automatic{},
				},
			},
		},
	}

	secret, err := s.secretManagerClient.CreateSecret(context.Background(), createSecretReq)
	if err != nil {
		return "", err
	}

	// Declare the payload to store.
	payload, err := json.Marshal(cred)
	if err != nil {
		return "", err
	}

	// Build the request.
	addSecretVersionReq := &secretmanagerpb.AddSecretVersionRequest{
		Parent: secret.Name,
		Payload: &secretmanagerpb.SecretPayload{
			Data: payload,
		},
	}

	// Call the API.
	_, err = s.secretManagerClient.AddSecretVersion(context.Background(), addSecretVersionReq)
	if err != nil {
		return "", err
	}

	return secretName, nil
}

func (s *UserProviderService) removeCredentialsFromGS(cred string) error {
	req := &secretmanagerpb.DeleteSecretRequest{
		Name: cred,
	}
	if err := s.secretManagerClient.DeleteSecret(context.Background(), req); err != nil {
		return fmt.Errorf("failed to delete secret: %w", err)
	}
	return nil
}

func (s *UserProviderService) GetCredentials(cred string) ([]byte, error) {

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: cred,
	}

	result, err := s.secretManagerClient.AccessSecretVersion(context.Background(), req)
	if err != nil {
		if s, ok := status.FromError(err); ok && s.Code() == codes.NotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to access secret version: %w", err)
	}
	return result.Payload.Data, nil
}

func (s *UserProviderService) ValidateProviderDto(c echo.Context, providerDto *dto.UserProviderDto) (any, error) {
	if err := c.Validate(providerDto); err != nil {
		return nil,
			&exceptions.BadRequestError{
				Message: err.Error(),
			}
	}
	detail, err := dto.GetCredFromDto[dto.TwilioCredDto](*providerDto)
	if err != nil {
		return nil, &exceptions.BadRequestError{
			Message: err.Error(),
		}
	}
	if err := c.Validate(detail); err != nil {
		return nil, &exceptions.BadRequestError{
			Message: err.Error(),
		}
	}
	return detail, nil
}

func GetProviderWithCred[T any](
	db *gorm.DB,
	secretManagerClient *secretmanager.Client,
	user auth.UserDetail,
	providerId uuid.UUID,
) (
	*model.UserProvider,
	*T,
	error,
) {
	var provider model.UserProvider
	if err := db.Model(&model.UserProvider{}).Select("*").Where(
		"id = ?",
		providerId,
	).First(&provider).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, &exceptions.NotFoundError{}
		}
		return nil, nil, err
	}
	if user.Role != enums.Admin && provider.UserID != user.ID {
		return nil, nil, &exceptions.AccessDenied{}
	}
	result, err := secretManagerClient.AccessSecretVersion(
		context.Background(), &secretmanagerpb.AccessSecretVersionRequest{
			Name: provider.ProviderCredentials + "/versions/latest",
		},
	)
	if err != nil {
		if s, ok := status.FromError(err); ok && s.Code() == codes.NotFound {
			return nil, nil, &exceptions.NotFoundError{}
		}
		return nil, nil, fmt.Errorf("failed to access secret version: %w", err)
	}
	cred, err := dto.GetCredFromBytes[T](result.Payload.Data)
	if err != nil {
		return nil, nil, err
	}
	return &provider, cred, nil
}

func GetProviderWithCredWithoutCheck[T any](
	db *gorm.DB,
	secretManagerClient *secretmanager.Client,
	providerId uuid.UUID,
) (
	*model.UserProvider,
	*T,
	error,
) {
	var provider model.UserProvider
	if err := db.Model(&model.UserProvider{}).Select("*").Where(
		"id = ?",
		providerId,
	).First(&provider).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, &exceptions.NotFoundError{}
		}
		return nil, nil, err
	}
	result, err := secretManagerClient.AccessSecretVersion(
		context.Background(), &secretmanagerpb.AccessSecretVersionRequest{
			Name: provider.ProviderCredentials + "/versions/latest",
		},
	)
	if err != nil {
		if s, ok := status.FromError(err); ok && s.Code() == codes.NotFound {
			return nil, nil, &exceptions.NotFoundError{}
		}
		return nil, nil, fmt.Errorf("failed to access secret version: %w", err)
	}
	cred, err := dto.GetCredFromBytes[T](result.Payload.Data)
	if err != nil {
		return nil, nil, err
	}
	return &provider, cred, nil
}
