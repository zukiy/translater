package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	v1 "translator/api/v1"
	"translator/config"
	"translator/model"

	"translator/providers/storage"

	"github.com/stretchr/testify/suite"
)

type TranslateSuite struct {
	suite.Suite
	storage *storage.Storage
	server  *v1.Server
}

func (suite *TranslateSuite) SetupSuite() {
	cfg := &config.Config{
		Translators: config.Translators{
			Yandex: config.Yandex{
				ApiVer: "v1.5",
				ApiKey: "trnsl.1.1.20190205T140220Z.1c75b6a3a6d8311c.23fa815511980868b6b5d09eacdfccc1b87ccead",
			},
		},
		MySQL: config.MySQL{
			Address:  "127.0.0.1:3306",
			DBName:   "translator",
			User:     "root",
			Password: "pass",
		},
	}

	suite.storage = storage.Init(cfg)
	suite.server = v1.New(8008, suite.storage, model.New(suite.storage))
	go suite.server.Serve()
}

func (suite *TranslateSuite) TearDownSuite() {
	suite.server.Stop()
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(TranslateSuite))
}

func (suite *TranslateSuite) Test_NewWord_Translate_Choose_Save() {

	transReq := v1.TranslateRequest{
		Word: "copy",
		Lang: "ru",
	}

	raw := post(transReq, "v1/translate")

	var resp v1.TranslateResponse
	err := json.Unmarshal(raw, &resp)
	if err != nil {
		panic(err)
	}

	saveReq := v1.SaveRequest{
		Word: v1.Word{
			Text: transReq.Word,
			Lang: "en", // todo
		},
		Translate: v1.Word{
			Text: resp.Results[0],
			Lang: transReq.Lang,
		},
	}

	raw = post(saveReq, "v1/save")

	listReq := v1.ListRequest{
		PerPage: 50,
		Page:    1,
	}

	// list
	raw = post(listReq, "v1/")
	suite.T().Logf("%+v", string(raw))
}

func get(data interface{}, path string) []byte {

	u := fmt.Sprintf("http://127.0.0.1:8008/%s", path)

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		panic(err)
	}

	params := url.Values{}
	params.Add("page", "1")
	params.Add("perpage", "50")
	req.URL.RawQuery = params.Encode()

	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return raw
}

func post(data interface{}, path string) []byte {
	raw, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	u := fmt.Sprintf("http://127.0.0.1:8008/%s", path)
	req, err := http.NewRequest(http.MethodPost, u, bytes.NewReader(raw))
	if err != nil {
		panic(err)
	}

	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}

	raw, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return raw
}
