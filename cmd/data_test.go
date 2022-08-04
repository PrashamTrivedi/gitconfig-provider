package cmd

import (
	"reflect"
	"testing"
)

func TestGetProviderByName(t *testing.T) {
	providerName,index := GetProviderByName("Github")

	want := GitProvider{
		Name: "github", Url: "https://github.com/",
	}

	if index==-1 {
		t.Errorf("Provider by name %s should be found","Github")
	}

	if !reflect.DeepEqual(want,providerName){
		t.Errorf("Wanted %+v, got %+v",want,providerName)
	}
}
