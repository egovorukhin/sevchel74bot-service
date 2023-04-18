package telegram

import (
	"encoding/json"
	"fmt"
	"testing"
)

func initClient() {
	Init(Config{
		Url:     "https://api.telegram.org/bot6238493597:AAGsz--zRI2Z_SEtzLF60x1Ok3oHVSeLBhM",
		Timeout: 30,
	})
}

func TestGetMe(t *testing.T) {
	initClient()
	me, err := GetMe()
	if err != nil {
		t.Log(err)
	}
	fmt.Printf("%v+", me)
}

func TestGetUpdates(t *testing.T) {
	initClient()
	me, err := GetUpdates()
	if err != nil {
		t.Log(err)
	}
	data, _ := json.Marshal(me)
	fmt.Printf("%s", data)
}

func TestWatch(t *testing.T) {
	initClient()
	err := Watch()
	if err != nil {
		t.Log(err)
	}
}
