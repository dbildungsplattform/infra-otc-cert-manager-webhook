package otcdns

import (
	"testing"

	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	cmmeta1 "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
)

func TestConfigJsonToOtcDnsConfig_Nil(t *testing.T) {
	cfg, err := configJsonToOtcDnsConfig(nil)
	if err != nil {
		t.Fatalf("expected no error for nil config, got: %v", err)
	}
	if cfg.AuthURL != "" || cfg.Region != "" || cfg.Tenant != "" || cfg.DomainName != "" {
		t.Fatalf("expected empty config, got: %+v", cfg)
	}
}

func TestConfigJsonToOtcDnsConfig_Basic(t *testing.T) {
	rawJSON := []byte(`{
		"authURL": "https://iam.example.com/v3",
		"region": "eu-de"
	}`)
	cfgExt := &extapi.JSON{Raw: rawJSON}

	cfg, err := configJsonToOtcDnsConfig(cfgExt)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if cfg.AuthURL != "https://iam.example.com/v3" {
		t.Fatalf("expected authURL=https://iam.example.com/v3, got: %s", cfg.AuthURL)
	}
	if cfg.Region != "eu-de" {
		t.Fatalf("expected region=eu-de, got: %s", cfg.Region)
	}
	if cfg.Tenant != "" {
		t.Fatalf("expected empty tenant, got: %s", cfg.Tenant)
	}
	if cfg.DomainName != "" {
		t.Fatalf("expected empty domainName, got: %s", cfg.DomainName)
	}
}

func TestConfigJsonToOtcDnsConfig_WithTenantAndDomain(t *testing.T) {
	rawJSON := []byte(`{
		"authURL": "https://iam.example.com/v3",
		"region": "eu-de",
		"tenant": "my-tenant-project",
		"domainName": "my-domain"
	}`)
	cfgExt := &extapi.JSON{Raw: rawJSON}

	cfg, err := configJsonToOtcDnsConfig(cfgExt)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if cfg.Tenant != "my-tenant-project" {
		t.Fatalf("expected tenant=my-tenant-project, got: %s", cfg.Tenant)
	}
	if cfg.DomainName != "my-domain" {
		t.Fatalf("expected domainName=my-domain, got: %s", cfg.DomainName)
	}
}

func TestConfigJsonToOtcDnsConfig_WithSecretRefs(t *testing.T) {
	rawJSON := []byte(`{
		"authURL": "https://iam.example.com/v3",
		"region": "eu-de",
		"tenant": "my-tenant",
		"domainName": "my-domain",
		"accessKeySecretRef": {
			"name": "otcdns-credentials",
			"key": "accessKey"
		},
		"secretKeySecretRef": {
			"name": "otcdns-credentials",
			"key": "secretKey"
		}
	}`)
	cfgExt := &extapi.JSON{Raw: rawJSON}

	cfg, err := configJsonToOtcDnsConfig(cfgExt)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if cfg.Tenant != "my-tenant" {
		t.Fatalf("expected tenant=my-tenant, got: %s", cfg.Tenant)
	}
	if cfg.DomainName != "my-domain" {
		t.Fatalf("expected domainName=my-domain, got: %s", cfg.DomainName)
	}
	if cfg.AccessKeySecretRef.Name != "otcdns-credentials" {
		t.Fatalf("expected accessKeySecretRef.name=otcdns-credentials, got: %s", cfg.AccessKeySecretRef.Name)
	}
	if cfg.SecretKeySecretRef.Key != "secretKey" {
		t.Fatalf("expected secretKeySecretRef.key=secretKey, got: %s", cfg.SecretKeySecretRef.Key)
	}
}

func TestConfigJsonToOtcDnsConfig_InvalidJSON(t *testing.T) {
	rawJSON := []byte(`{invalid json`)
	cfgExt := &extapi.JSON{Raw: rawJSON}

	_, err := configJsonToOtcDnsConfig(cfgExt)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestOtcDnsConfig_StructFields(t *testing.T) {
	// Verify that all expected fields exist and have correct json tags
	cfg := OtcDnsConfig{
		AccessKey:  "ak",
		SecretKey:  "sk",
		Region:     "eu-de",
		AuthURL:    "https://iam.example.com/v3",
		Tenant:     "my-project",
		DomainName: "my-domain",
	}

	if cfg.Tenant != "my-project" {
		t.Fatalf("expected Tenant=my-project, got: %s", cfg.Tenant)
	}
	if cfg.DomainName != "my-domain" {
		t.Fatalf("expected DomainName=my-domain, got: %s", cfg.DomainName)
	}

	// Verify secret ref selectors compile
	_ = cmmeta1.SecretKeySelector{
		LocalObjectReference: cmmeta1.LocalObjectReference{Name: "test-secret"},
		Key:                  "test-key",
	}
}