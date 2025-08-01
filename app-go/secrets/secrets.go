package secrets

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// SecretsManager handles secrets for different environments
type SecretsManager struct {
	isProduction bool
	credentialsPath string
	examplePath     string
}

// MongoDBSecrets represents MongoDB connection credentials
type MongoDBSecrets struct {
	AdminUser     string `json:"admin_user"`
	AdminPassword string `json:"admin_password"`
	AppUser       string `json:"app_user"`
	AppPassword   string `json:"app_password"`
	Host          string `json:"host"`
	Port          string `json:"port"`
	Database      string `json:"database"`
	ReplicaSet    string `json:"replica_set"`
}

// Credentials represents the full credentials structure
type Credentials struct {
	MongoDB    MongoDBSecrets `json:"mongodb"`
	JWT        JWTSecrets     `json:"jwt"`
	Encryption EncryptionSecrets `json:"encryption"`
}

// JWTSecrets represents JWT configuration
type JWTSecrets struct {
	Secret string `json:"secret"`
}

// EncryptionSecrets represents encryption configuration
type EncryptionSecrets struct {
	Key string `json:"key"`
}

// NewSecretsManager creates a new secrets manager
func NewSecretsManager() *SecretsManager {
	wd, _ := os.Getwd()
	return &SecretsManager{
		isProduction:    os.Getenv("GO_ENV") == "production",
		credentialsPath: filepath.Join(wd, "credentials.local.json"),
		examplePath:     filepath.Join(wd, "env.credentials.example"),
	}
}

// GetSecrets retrieves secrets based on environment
func (sm *SecretsManager) GetSecrets() (*MongoDBSecrets, error) {
	if sm.isProduction {
		return sm.getProductionSecrets()
	}
	return sm.getLocalSecrets()
}

// getProductionSecrets retrieves secrets from AWS Secrets Manager
func (sm *SecretsManager) getProductionSecrets() (*MongoDBSecrets, error) {
	// In production, this would use AWS SDK
	// For now, we'll return an error to indicate this needs implementation
	return nil, fmt.Errorf("production secrets management not implemented. Please configure AWS Secrets Manager")
	
	// Example AWS Secrets Manager implementation:
	// sess := session.Must(session.NewSession())
	// svc := secretsmanager.New(sess)
	// input := &secretsmanager.GetSecretValueInput{
	//     SecretId: aws.String("mongodb-credentials"),
	// }
	// result, err := svc.GetSecretValue(input)
	// if err != nil {
	//     return nil, err
	// }
	// 
	// var secrets MongoDBSecrets
	// err = json.Unmarshal([]byte(*result.SecretString), &secrets)
	// if err != nil {
	//     return nil, err
	// }
	// return &secrets, nil
}

// getLocalSecrets retrieves secrets from local credentials file
func (sm *SecretsManager) getLocalSecrets() (*MongoDBSecrets, error) {
	if _, err := os.Stat(sm.credentialsPath); os.IsNotExist(err) {
		sm.createLocalCredentialsFile()
		return nil, fmt.Errorf("local credentials file not found. Please create %s based on env.credentials.example", sm.credentialsPath)
	}

	data, err := os.ReadFile(sm.credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("error reading credentials file: %v", err)
	}

	var credentials Credentials
	if err := json.Unmarshal(data, &credentials); err != nil {
		return nil, fmt.Errorf("error parsing credentials file: %v", err)
	}

	if err := sm.validateCredentials(&credentials); err != nil {
		return nil, err
	}

	return &credentials.MongoDB, nil
}

// createLocalCredentialsFile creates a local credentials file from example
func (sm *SecretsManager) createLocalCredentialsFile() error {
	if _, err := os.Stat(sm.examplePath); os.IsNotExist(err) {
		return fmt.Errorf("example credentials file not found: %s", sm.examplePath)
	}

	exampleData, err := os.ReadFile(sm.examplePath)
	if err != nil {
		return fmt.Errorf("error reading example file: %v", err)
	}

	err = os.WriteFile(sm.credentialsPath, exampleData, 0600)
	if err != nil {
		return fmt.Errorf("error creating credentials file: %v", err)
	}

	fmt.Printf("Created %s from example. Please update with your local credentials.\n", sm.credentialsPath)
	return nil
}

// validateCredentials validates the credentials structure
func (sm *SecretsManager) validateCredentials(creds *Credentials) error {
	if creds.MongoDB.AdminUser == "" {
		return fmt.Errorf("missing required field: admin_user")
	}
	if creds.MongoDB.AdminPassword == "" {
		return fmt.Errorf("missing required field: admin_password")
	}
	if creds.MongoDB.AppUser == "" {
		return fmt.Errorf("missing required field: app_user")
	}
	if creds.MongoDB.AppPassword == "" {
		return fmt.Errorf("missing required field: app_password")
	}
	if creds.MongoDB.Host == "" {
		return fmt.Errorf("missing required field: host")
	}
	if creds.MongoDB.Port == "" {
		return fmt.Errorf("missing required field: port")
	}
	if creds.MongoDB.Database == "" {
		return fmt.Errorf("missing required field: database")
	}

	return nil
}

// GetMongoURI constructs MongoDB connection URI from secrets
func (sm *SecretsManager) GetMongoURI(secrets *MongoDBSecrets, directConnection bool) string {
	if directConnection {
		return fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?directConnection=true",
			secrets.AppUser, secrets.AppPassword, secrets.Host, secrets.Port, secrets.Database)
	}
	return fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?replicaSet=%s",
		secrets.AppUser, secrets.AppPassword, secrets.Host, secrets.Port, secrets.Database, secrets.ReplicaSet)
}

// GetAdminMongoURI constructs admin MongoDB connection URI
func (sm *SecretsManager) GetAdminMongoURI(secrets *MongoDBSecrets) string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s/admin?replicaSet=%s&authSource=admin",
		secrets.AdminUser, secrets.AdminPassword, secrets.Host, secrets.Port, secrets.ReplicaSet)
}

// TestSecrets tests the secrets configuration
func (sm *SecretsManager) TestSecrets() error {
	secrets, err := sm.GetSecrets()
	if err != nil {
		return fmt.Errorf("secrets configuration error: %v", err)
	}

	if err := sm.validateCredentials(&Credentials{MongoDB: *secrets}); err != nil {
		return fmt.Errorf("secrets validation error: %v", err)
	}

	fmt.Println("✅ Secrets configuration is valid")
	return nil
}

// GetJWTSecret retrieves JWT secret
func (sm *SecretsManager) GetJWTSecret() string {
	if _, err := os.Stat(sm.credentialsPath); os.IsNotExist(err) {
		return "default-jwt-secret-change-in-production"
	}

	data, err := os.ReadFile(sm.credentialsPath)
	if err != nil {
		return "default-jwt-secret-change-in-production"
	}

	var credentials Credentials
	if err := json.Unmarshal(data, &credentials); err != nil {
		return "default-jwt-secret-change-in-production"
	}

	if credentials.JWT.Secret != "" {
		return credentials.JWT.Secret
	}

	return "default-jwt-secret-change-in-production"
}

// GetEncryptionKey retrieves encryption key
func (sm *SecretsManager) GetEncryptionKey() string {
	if _, err := os.Stat(sm.credentialsPath); os.IsNotExist(err) {
		return "default-encryption-key-32-chars"
	}

	data, err := os.ReadFile(sm.credentialsPath)
	if err != nil {
		return "default-encryption-key-32-chars"
	}

	var credentials Credentials
	if err := json.Unmarshal(data, &credentials); err != nil {
		return "default-encryption-key-32-chars"
	}

	if credentials.Encryption.Key != "" {
		return credentials.Encryption.Key
	}

	return "default-encryption-key-32-chars"
} 