package cmd

import (
	"strings"
	"testing"
)

func TestGetProviderByName(t *testing.T) {

	providerName, index := GetProviderByName("Github")

	want := GitProvider{
		Name: "github", Url: "https://github.com/",
	}

	if index == -1 {
		t.Errorf("Provider by name %s should be found", "Github")
	}

	if want.Name != providerName.Name || want.Url != providerName.Url {
		t.Errorf("Wanted %+v, got %+v", want, providerName)
	}
}

func TestGetProviderByUrl(t *testing.T) {
	// ResetProviders()

	// Test Case 1: Test with empty providerUrl, should return -1 and nil
	provider, index := GetProviderByUrl("")

	if index != -1 {
		t.Errorf("Expected index to be -1, but got %v", index)
	}
	if provider.Name != "" {
		t.Errorf("Expected provider to be nil, but got %v", provider)
	}

	// Test Case 2: Test with non-existent providerUrl, should return -1 and nil
	provider, index = GetProviderByUrl("https://example.com")
	if index != -1 {
		t.Errorf("Expected index to be -1, but got %v", index)
	}
	if provider.Name != "" {
		t.Errorf("Expected provider to be nil, but got %v", provider)
	}

	// Test Case 3: Test with existing providerUrl, should return a valid GitProvider and its index
	provider, index = GetProviderByUrl("https://github.com/")
	if index == -1 {
		t.Errorf("Expected index to be non-negative, but got %v", index)
	}
	if provider.Name == "" {
		t.Errorf("Expected provider to be non-nil, but got %v", provider)
	}
	if provider.Name != "github" {
		t.Errorf("Expected provider name to be Github, but got %v", provider.Name)
	}

	provider, index = GetProviderByUrl("https://github.com/myoffice/")
	if index != -1 {
		t.Errorf("Expected index to be -1, but got %v", index)
	}
	if provider.Name != "" {
		t.Errorf("Expected provider to be nil, but got %v", provider)
	}
}

func TestApply(t *testing.T) {
	failingScenario := []GitProvider{
		{
			Name: "Github", Url: "https://github.com/",
		},
		{
			Name: "GithubOffice", Url: "https://github.com/office/",
		},
	}

	writeGitProviders(failingScenario)

	provider, index := GetProviderByUrl("https://github.com/PrashamTrivedi/gitconfig-provider.git")

	if index == -1 {
		t.Errorf("Provider must be found")
	}
	if strings.ToLower(provider.Name) != "github" {
		t.Errorf("Provider Name should be github but got %s", provider.Name)
	}

	provider, index = GetProviderByUrl("https://github.com/office/my-repo.git")
	if index == -1 {
		t.Errorf("Provider must be found")
	}
	if strings.ToLower(provider.Name) != "githuboffice" {
		t.Errorf("Provider Name should be GithubOffice but got %s", provider.Name)
	}
	ResetProviders()
}
