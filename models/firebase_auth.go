package models

import (
	"encoding/json"
	"fmt"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
)

type (
	FirebaseConfig struct {
		Type                    string `json:"type"`
		ProjectID               string `json:"project_id"`
		PrivateKeyID            string `json:"private_key_id"`
		PrivateKey              string `json:"private_key"`
		ClientEmail             string `json:"client_email"`
		ClientID                string `json:"client_id"`
		AuthURI                 string `json:"auth_uri"`
		TokenURI                string `json:"token_uri"`
		AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
		ClientX509CertURL       string `json:"client_x509_cert_url"`
		UniverseDomain          string `json:"universe_domain"`
	}

	FirebaseFile struct {
		Name         string `json:"name"`
		Size         int    `json:"size"`
		ObjectURL    string `json:"objectURL"`
		LastModified int    `json:"lastModified"`
		Type         string `json:"type"`
		FileBinary   string `json:"file_binary"`
	}

	DocumentList struct {
		File []FirebaseFile `json:"file"`
	}
)

func FirebaseAuth() ([]byte, error) {

	FirebaseType, _ := beego.AppConfig.String("firebase-storage::firebase_type")
	FirebaseProjectID, _ := beego.AppConfig.String("firebase-storage::firebase_project_id")
	FirebasePrivateKeyID, _ := beego.AppConfig.String("firebase-storage::firebase_private_key_id")
	FirebasePrivateKey, _ := beego.AppConfig.String("firebase-storage::firebase_private_key")
	FirebaseClientEmail, _ := beego.AppConfig.String("firebase-storage::firebase_client_email")
	FirebaseClientID, _ := beego.AppConfig.String("firebase-storage::firebase_client_id")
	FirebaseAuthUri, _ := beego.AppConfig.String("firebase-storage::firebase_auth_uri")
	FirebaseTokenUri, _ := beego.AppConfig.String("firebase-storage::firebase_token_uri")
	FirebaseAuthProvideX509CertUrl, _ := beego.AppConfig.String("firebase-storage::firebase_auth_provider_x509_cert_url")
	FirebaseClientX509CertUrl, _ := beego.AppConfig.String("firebase-storage::firebase_client_x509_cert_url")
	FirebaseUniverseDomain, _ := beego.AppConfig.String("firebase-storage::firebase_universe_domain")

	FirebasePrivateKey = strings.Replace(FirebasePrivateKey, "\\n", "\n", -1)
	fileContent := map[string]interface{}{
		"type":                        FirebaseType,
		"project_id":                  FirebaseProjectID,
		"private_key_id":              FirebasePrivateKeyID,
		"private_key":                 FirebasePrivateKey,
		"client_email":                FirebaseClientEmail,
		"client_id":                   FirebaseClientID,
		"auth_uri":                    FirebaseAuthUri,
		"token_uri":                   FirebaseTokenUri,
		"auth_provider_x509_cert_url": FirebaseAuthProvideX509CertUrl,
		"client_x509_cert_url":        FirebaseClientX509CertUrl,
		"universe_domain":             FirebaseUniverseDomain,
	}

	config := FirebaseConfig{
		Type:                    fileContent["type"].(string),
		ProjectID:               fileContent["project_id"].(string),
		PrivateKeyID:            fileContent["private_key_id"].(string),
		PrivateKey:              fileContent["private_key"].(string),
		ClientEmail:             fileContent["client_email"].(string),
		ClientID:                fileContent["client_id"].(string),
		AuthURI:                 fileContent["auth_uri"].(string),
		TokenURI:                fileContent["token_uri"].(string),
		AuthProviderX509CertURL: fileContent["auth_provider_x509_cert_url"].(string),
		ClientX509CertURL:       fileContent["client_x509_cert_url"].(string),
		UniverseDomain:          fileContent["universe_domain"].(string),
	}

	jsonData, err := json.Marshal(config)
	if err != nil {
		err = fmt.Errorf("error encoding json: %v", err.Error())
	}
	return jsonData, err
}

func ValidateFirebase() ([]byte, string, error) {
	filePath, errc := FirebaseAuth()

	storageBucket, _ := beego.AppConfig.String("firebase-storage::bucket_link")

	return filePath, storageBucket, errc
}
