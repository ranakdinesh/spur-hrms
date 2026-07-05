package domain

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestNewTenantProfileRejectsBusinessSuffixSubdomain(t *testing.T) {
	_, err := NewTenantProfile(uuid.New(), "aanvi-infotech", "AANVI-1234", strPtr("Aanvi infotech"), nil)
	if !errors.Is(err, ErrConfusingSubdomain) {
		t.Fatalf("error = %v, want %v", err, ErrConfusingSubdomain)
	}
}

func TestTenantSubdomainCollisionKeyIgnoresBusinessSuffix(t *testing.T) {
	if got, want := TenantSubdomainCollisionKey("aanvi-infotech"), "aanvi"; got != want {
		t.Fatalf("collision key = %q, want %q", got, want)
	}
	if got, want := TenantSubdomainCollisionKey("mash-virtual"), "mashvirtual"; got != want {
		t.Fatalf("collision key = %q, want %q", got, want)
	}
}

func TestNewTenantProvisioningDerivesShortSubdomainFromCompanyName(t *testing.T) {
	provisioning, err := NewTenantProvisioning(TenantProvisioningInput{
		TenantID:             uuid.New(),
		CompanyName:          "Aanvi Infotech",
		MobileActivationCode: "AANVI-1234",
	})
	if err != nil {
		t.Fatalf("NewTenantProvisioning returned error: %v", err)
	}
	if got, want := provisioning.Subdomain, "aanvi"; got != want {
		t.Fatalf("subdomain = %q, want %q", got, want)
	}
}

func strPtr(value string) *string {
	return &value
}
