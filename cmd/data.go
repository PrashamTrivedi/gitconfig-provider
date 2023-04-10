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

//Check if the provider is already in the list, if not, add into the list and save to the file.
func AddProvider(providerName, providerUrl string, doForcefully bool) (GitProvider, int) {
	if doForcefully {
		fmt.Println("CAUTION!! This is an experimental API under testing, Adding a provider forcefully may lead to a situation where a wrong provider can be recognised.\n If this is the situation, use one of these 3 remedies. \n1. Use apply with desired provider name e.g. `gitconfig-provider apply githubOffice` \n2. Remove the provider you have added forcefully and \n3. Reset all providers and add desired providers and properties accordingly ")
		return addNewProvider(providerName, providerUrl)
	} else {

		provider, index := GetProviderByName(providerName)
		//if index is -1 get provider by url
		if index == -1 {
			provider, index = GetProviderByUrl(providerUrl)
		}
		//return provider if index is not -1
		if index != -1 {
			currentProviders[index] = provider
			return provider, index
		} else {
			//create a git provider with providerName and providerUrl, keep properties blank
			return addNewProvider(providerName, providerUrl)
		}
	}
}

func addNewProvider(providerName string, providerUrl string) (GitProvider, int) {
	provider := GitProvider{Name: providerName, Url: providerUrl, Properties: make(map[string]string)}
	currentProviders = append(currentProviders, provider)
	writeGitProviders(currentProviders)
	return provider, -1
}

func AddProviderPropertyFromName(providerName string, propertyKey string, propertyValue string) {
	provider, index := GetProviderByName(providerName)
	//throw error if index is -1
	if index == -1 {
		fmt.Println("Error: Provider not found. Run gitconfig-provider listProviders to get available git providers")
		os.Exit(1)
	}
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

	//check if key has space
	if strings.Contains(key, " ") {
		key = "\"" + key + "\""
	}
	if strings.Contains(value, " ") {
		value = "\"" + value + "\""
	}
	configCmd := fmt.Sprintf("git config %s %s", key, value)

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

	//throw error if index is -1
	if index == -1 {
		fmt.Println("Error: Provider not found. Run gitconfig-provider listProviders to get available git providers")
		os.Exit(1)
	}

	if provider.Properties == nil || len(provider.Properties) == 0 {
		provider.Properties = make(map[string]string)
	}
	provider.Properties[propertyKey] = propertyValue
	updateProvider(provider, index)
	runUpdateCommand(propertyKey, propertyValue)
}

func RemoveProvider(providerName, providerUrl string) {
	provider, index := GetProviderByName(providerName)
	if index == -1 {
		provider, index = GetProviderByUrl(providerUrl)
	}
	if index == -1 {
		fmt.Printf("No provider with name %s or url %s is found\n", providerName, providerUrl)
	}

	currentProviders = append(currentProviders[:index], currentProviders[index+1:]...)
	writeGitProviders(currentProviders)
	fmt.Printf("Provider %s (with Url: %s) is removed along with the properties\n", provider.Name, provider.Url)
}

func RemoveProviderProperty(providerName string, propertyKey string) {
	provider, index := GetProviderByName(providerName)

	//throw error if index is -1
	if index == -1 {
		fmt.Println("Error: Provider not found. Run gitconfig-provider listProviders to get available git providers")
		os.Exit(1)
	}

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
		if strings.EqualFold(provider.Name, providerName) {
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
		if strings.HasPrefix(providerUrl, provider.Url) {
			providerToReturn = provider
			indexToReturn = index
			// break
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

func ApplyPropertiesForRemote(providerName string) {
	url := ""
	if providerName != "" {
		provider, index := GetProviderByName(providerName)
		if index == -1 {
			fmt.Println("Provider with name:", providerName, "Not found")
			os.Exit(1)
		}
		url = provider.Url
	} else {

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
		url = firstRemote.Url

	}
	provider, index := GetProviderByUrl(url)
	//if index is -1 then no provider found
	if index == -1 {
		fmt.Println("Error: Provider not found, Checking for ", url)
		os.Exit(1)
	}
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
