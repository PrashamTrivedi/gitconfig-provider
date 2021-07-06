package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mitchellh/go-homedir"
)

type GitProvider struct {
	Name       string            `json:"name,omitempty"`
	Url        string            `json:"url,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}

type GitRemote struct {
	Name    string
	Url     string
	Command string
}

var defaultProviders []GitProvider = []GitProvider{
	{Name: "github", Url: "https://github.com/"},
	{Name: "bitbucket", Url: "https://bitbucket.org/"},
	{Name: "gitlab", Url: "https://gitlab.com/"},
}

var currentProviders []GitProvider

func AddProviderPropertyFromName(providerName string, propertyKey string, propertyValue string) {
	provider, index := GetProviderByName(providerName)
	if provider.Properties == nil || len(provider.Properties) == 0 {
		provider.Properties = make(map[string]string)
	}
	provider.Properties[propertyKey] = propertyValue
	updateProvider(provider, index)
	runUpdateCommand(propertyKey, propertyValue)
}

func getRemotes() ([]GitRemote, error) {
	mainCommand := "bash"
	commandFlag := "-c"
	if runtime.GOOS == "windows" {
		mainCommand = "cmd"
		commandFlag = "/c"
	}

	configCmd := "git remote -v"

	command := exec.Command(mainCommand, commandFlag, configCmd)
	var stderr bytes.Buffer
	command.Stderr = &stderr
	output, commandError := command.Output()
	if commandError != nil {
		return nil, commandError
	}
	remoteData := string(output)
	remotes := strings.Split(remoteData, "\n")
	var gitRemotes []GitRemote
	for _, remote := range remotes {
		if remote != "" {

			name := ""
			url := ""
			command := ""
			fmt.Sscanf(remote, "%s %s (%s)", &name, &url, &command)

			gitRemotes = append(gitRemotes, GitRemote{Name: name, Url: url, Command: command})
		}
	}
	return gitRemotes, nil

}
func runUpdateCommand(key string, value string) {
	mainCommand := "bash"
	commandFlag := "-c"
	if runtime.GOOS == "windows" {
		mainCommand = "cmd"
		commandFlag = "/c"
	}

	configCmd := fmt.Sprintf("git config \"%s\" \"%s\"", key, value)

	command := exec.Command(mainCommand, commandFlag, configCmd)
	var stderr bytes.Buffer
	command.Stderr = &stderr
	_, commandError := command.Output()
	if commandError != nil {
		fmt.Println(command)
		fmt.Println(fmt.Sprint(commandError) + ": " + stderr.String())
	}

}
func AddProviderPropertyFromUrl(providerUrl string, propertyKey string, propertyValue string) {
	provider, index := GetProviderByUrl(providerUrl)
	if provider.Properties == nil || len(provider.Properties) == 0 {
		provider.Properties = make(map[string]string)
	}
	provider.Properties[propertyKey] = propertyValue
	updateProvider(provider, index)
	runUpdateCommand(propertyKey, propertyValue)
}

func RemoveProviderProperty(providerName string, propertyKey string) {
	provider, index := GetProviderByName(providerName)
	if provider.Properties == nil || len(provider.Properties) == 0 {
		provider.Properties = make(map[string]string)
	}
	delete(provider.Properties, propertyKey)
	updateProvider(provider, index)
	runUpdateCommand(propertyKey, "")
}
func GetProviderByName(providerName string) (GitProvider, int) {
	indexToReturn := -1
	var providerToReturn GitProvider
	if len(currentProviders) == 0 {
		readGitProviders()
	}
	for index, provider := range currentProviders {
		if provider.Name == providerName {
			providerToReturn = provider
			indexToReturn = index
			break
		}
	}
	return providerToReturn, indexToReturn

}

func GetProviderByUrl(providerUrl string) (GitProvider, int) {
	indexToReturn := -1
	var providerToReturn GitProvider
	if len(currentProviders) == 0 {
		readGitProviders()
	}
	for index, provider := range currentProviders {
		if strings.HasPrefix(provider.Url, providerUrl) {
			providerToReturn = provider
			indexToReturn = index
			break
		}
	}

	return providerToReturn, indexToReturn

}

func GetProviders() ([]GitProvider, error) {
	if len(currentProviders) > 0 {
		return currentProviders, nil
	}
	return readGitProviders()
}

func ApplyPropertiesForRemote() {
	remotes, err := getRemotes()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	firstRemote := remotes[0]
	for _, remote := range remotes {
		if remote.Name == "origin" {
			firstRemote = remote
			break
		}
	}

	provider, _ := GetProviderByUrl(firstRemote.Url)
	applyCommandFromProvider(provider)

}
func applyCommandFromProvider(provider GitProvider) {
	for key, value := range provider.Properties {
		runUpdateCommand(key, value)
	}
}
func updateProvider(provider GitProvider, index int) {
	currentProviders[index] = provider
	writeGitProviders(currentProviders)
}
func readGitProviders() ([]GitProvider, error) {
	var gitProviders []GitProvider
	home, _ := homedir.Dir()
	dbPath := filepath.Join(home, "gitProviders.json")
	fileBytes, err := ioutil.ReadFile(dbPath)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(fileBytes, &gitProviders); err != nil {
		return nil, err
	}
	currentProviders = gitProviders
	return gitProviders, nil
}

func ResetProviders() {
	home, _ := os.UserHomeDir()

	dbPath := filepath.Join(home, "gitProviders.json")
	err := os.Remove(dbPath)
	if err != nil {
		fmt.Println("Error in processing file:", err.Error())
		os.Exit(1)
	}
}
func init() {
	if len(currentProviders) == 0 {
		currentProviders, err := readGitProviders()
		if err != nil || len(currentProviders) == 0 {

			writeGitProviders(defaultProviders)
			currentProviders = defaultProviders
		}
	}
}

func writeGitProviders(providers []GitProvider) {
	home, _ := os.UserHomeDir()

	dbPath := filepath.Join(home, "gitProviders.json")

	defaultFilesData, err := json.Marshal(providers)
	if err != nil {
		fmt.Println("Error in writing to file:", err.Error())
		os.Exit(1)
	}

	dbFile, errorData := os.OpenFile(dbPath, os.O_RDWR|os.O_CREATE, 0600)
	if errorData != nil {
		fmt.Println("Error in processing file:", errorData.Error())
		os.Exit(1)
	}

	err = ioutil.WriteFile(dbFile.Name(), defaultFilesData, 0600)
	if err != nil {
		fmt.Println("Error in writing to file:", err.Error())
		os.Exit(1)
	}

}
