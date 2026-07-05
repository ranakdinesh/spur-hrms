package storage

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
	"github.com/ranakdinesh/spur-hrms/internal/logging"
	"github.com/rs/zerolog"
)

type ObjectStorage struct {
	log *zerolog.Logger
}

func NewObjectStorage(log *zerolog.Logger) *ObjectStorage {
	return &ObjectStorage{log: logging.Component(log, "object_storage")}
}

func (s *ObjectStorage) PutObject(ctx context.Context, settings *domain.StorageProviderSettings, input ports.StoreObjectInput) (string, error) {
	if err := validateObject(settings, input); err != nil {
		return "", err
	}
	client, err := s.client(settings)
	if err != nil {
		return "", err
	}
	key := domain.StorageObjectKey(settings.ObjectPrefix, input.TenantID, input.Category, input.OwnerID, input.EntityID, input.FileName)
	_, err = client.PutObject(ctx, settings.Bucket, key, bytes.NewReader(input.Content), int64(len(input.Content)), minio.PutObjectOptions{ContentType: strings.TrimSpace(input.ContentType)})
	if err != nil {
		if s.log != nil {
			s.log.Error().Err(err).Str("tenant_id", input.TenantID.String()).Str("provider", settings.Provider).Str("bucket", settings.Bucket).Str("object_key", key).Msg("hrms object storage put failed")
		}
		return "", err
	}
	return objectReference(settings, key), nil
}

func (s *ObjectStorage) TestStorage(ctx context.Context, settings *domain.StorageProviderSettings) error {
	if settings == nil {
		return domain.ErrStorageProviderSettingsNotFound
	}
	if !settings.IsEnabled {
		return domain.ErrStorageProviderDisabled
	}
	client, err := s.client(settings)
	if err != nil {
		return err
	}
	ok, err := client.BucketExists(ctx, settings.Bucket)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("storage bucket %q does not exist or is not accessible", settings.Bucket)
	}
	return nil
}

func (s *ObjectStorage) client(settings *domain.StorageProviderSettings) (*minio.Client, error) {
	endpoint := strings.TrimSpace(valueOrDefault(settings.Endpoint, ""))
	if endpoint == "" && settings.Provider == domain.StorageProviderS3 {
		region := strings.TrimSpace(valueOrDefault(settings.Region, ""))
		if region == "" {
			endpoint = "s3.amazonaws.com"
		} else {
			endpoint = "s3." + region + ".amazonaws.com"
		}
	}
	if endpoint == "" {
		return nil, domain.ErrInvalidStorageEndpoint
	}
	secure := settings.UseSSL
	endpoint = normalizeEndpoint(endpoint)
	lookup := minio.BucketLookupAuto
	if settings.ForcePathStyle || settings.Provider == domain.StorageProviderMinIO {
		lookup = minio.BucketLookupPath
	}
	return minio.New(endpoint, &minio.Options{Creds: credentials.NewStaticV4(valueOrDefault(settings.AccessKeyID, ""), valueOrDefault(settings.SecretAccessKey, ""), ""), Secure: secure, Region: valueOrDefault(settings.Region, ""), BucketLookup: lookup})
}

type TenantFileStorage struct {
	repo            ports.StorageProviderRepo
	objectStorage   ports.ObjectStorage
	defaultSettings *domain.StorageProviderSettings
	log             *zerolog.Logger
}

func NewTenantFileStorage(repo ports.StorageProviderRepo, objectStorage ports.ObjectStorage, defaultSettings *domain.StorageProviderSettings, log *zerolog.Logger) *TenantFileStorage {
	return &TenantFileStorage{repo: repo, objectStorage: objectStorage, defaultSettings: defaultSettings, log: logging.Component(log, "tenant_file_storage")}
}

func (s *TenantFileStorage) StorePolicyFile(ctx context.Context, input ports.StorePolicyFileInput) (string, error) {
	return s.store(ctx, ports.StoreObjectInput{TenantID: input.TenantID, Category: ports.StorageCategoryPolicyFile, EntityID: input.PolicyID, FileName: input.FileName, ContentType: input.ContentType, Content: input.Content})
}

func (s *TenantFileStorage) StoreEmployeeDocument(ctx context.Context, input ports.StoreEmployeeDocumentInput) (string, error) {
	return s.store(ctx, ports.StoreObjectInput{TenantID: input.TenantID, Category: ports.StorageCategoryEmployeeDoc, OwnerID: input.EmployeeID, EntityID: input.DocumentID, FileName: input.FileName, ContentType: input.ContentType, Content: input.Content})
}

func (s *TenantFileStorage) StoreSalarySlipPDF(ctx context.Context, input ports.StoreSalarySlipPDFInput) (string, error) {
	return s.store(ctx, ports.StoreObjectInput{TenantID: input.TenantID, Category: ports.StorageCategorySalarySlip, OwnerID: input.UserID, EntityID: input.SlipID, FileName: path.Join(fmt.Sprintf("%04d", input.Year), fmt.Sprintf("%02d", input.Month), input.FileName), ContentType: input.ContentType, Content: input.Content})
}

func (s *TenantFileStorage) StoreEmployeeLetterPDF(ctx context.Context, input ports.StoreEmployeeLetterPDFInput) (string, error) {
	return s.store(ctx, ports.StoreObjectInput{TenantID: input.TenantID, Category: ports.StorageCategoryEmployeeLetter, OwnerID: input.EmployeeID, EntityID: input.LetterID, FileName: path.Join(input.LetterType, input.FileName), ContentType: input.ContentType, Content: input.Content})
}

func (s *TenantFileStorage) StoreAgreementPDF(ctx context.Context, input ports.StoreAgreementPDFInput) (string, error) {
	return s.store(ctx, ports.StoreObjectInput{TenantID: input.TenantID, Category: ports.StorageCategoryAgreement, EntityID: input.AgreementID, FileName: path.Join(input.AgreementType, input.FileName), ContentType: input.ContentType, Content: input.Content})
}

func (s *TenantFileStorage) StoreHRCaseAttachment(ctx context.Context, input ports.StoreHRCaseAttachmentInput) (string, error) {
	return s.store(ctx, ports.StoreObjectInput{TenantID: input.TenantID, Category: ports.StorageCategoryHRCase, OwnerID: input.CaseID, EntityID: input.CommentID, FileName: input.FileName, ContentType: input.ContentType, Content: input.Content})
}

func (s *TenantFileStorage) StoreLearningCertificate(ctx context.Context, input ports.StoreLearningCertificateInput) (string, error) {
	return s.store(ctx, ports.StoreObjectInput{TenantID: input.TenantID, Category: ports.StorageCategoryLearningCert, OwnerID: input.WorkerProfileID, EntityID: input.EnrollmentID, FileName: input.FileName, ContentType: input.ContentType, Content: input.Content})
}

func (s *TenantFileStorage) store(ctx context.Context, input ports.StoreObjectInput) (string, error) {
	settings, err := s.resolveSettings(ctx, input.TenantID)
	if err != nil {
		return "", err
	}
	if s.objectStorage == nil {
		return "", domain.ErrStorageProviderSettingsNotFound
	}
	return s.objectStorage.PutObject(ctx, settings, input)
}

func (s *TenantFileStorage) resolveSettings(ctx context.Context, tenantID uuid.UUID) (*domain.StorageProviderSettings, error) {
	if s.repo != nil {
		item, err := s.repo.GetStorageProviderSettings(ctx, tenantID)
		if err == nil {
			return item, nil
		}
		if err != domain.ErrStorageProviderSettingsNotFound {
			return nil, err
		}
	}
	if s.defaultSettings == nil {
		return nil, domain.ErrStorageProviderSettingsNotFound
	}
	copy := *s.defaultSettings
	copy.TenantID = tenantID
	return &copy, nil
}

func validateObject(settings *domain.StorageProviderSettings, input ports.StoreObjectInput) error {
	if settings == nil {
		return domain.ErrStorageProviderSettingsNotFound
	}
	if !settings.IsEnabled {
		return domain.ErrStorageProviderDisabled
	}
	if settings.MaxFileSizeBytes > 0 && int64(len(input.Content)) > settings.MaxFileSizeBytes {
		return domain.ErrStorageFileTooLarge
	}
	if settings.AllowedContentTypes != nil && strings.TrimSpace(*settings.AllowedContentTypes) != "" {
		contentType := strings.ToLower(strings.TrimSpace(input.ContentType))
		allowed := false
		for _, item := range strings.Split(*settings.AllowedContentTypes, ",") {
			if strings.ToLower(strings.TrimSpace(item)) == contentType {
				allowed = true
				break
			}
		}
		if !allowed {
			return domain.ErrStorageContentTypeNotAllowed
		}
	}
	return nil
}

func objectReference(settings *domain.StorageProviderSettings, key string) string {
	if settings.PublicBaseURL != nil && strings.TrimSpace(*settings.PublicBaseURL) != "" {
		base := strings.TrimRight(strings.TrimSpace(*settings.PublicBaseURL), "/")
		return base + "/" + strings.TrimLeft(key, "/")
	}
	return settings.Provider + "://" + settings.Bucket + "/" + key
}

func normalizeEndpoint(value string) string {
	clean := strings.TrimSpace(value)
	if parsed, err := url.Parse(clean); err == nil && parsed.Host != "" {
		return parsed.Host
	}
	return strings.TrimPrefix(strings.TrimPrefix(clean, "https://"), "http://")
}

func valueOrDefault(value *string, fallback string) string {
	if value == nil || strings.TrimSpace(*value) == "" {
		return fallback
	}
	return strings.TrimSpace(*value)
}
