package tts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type TtsClient struct {
	token    string
	iamToken *IAMTokenResponse
}

type OAuthTokenRequest struct {
	YandexPassportOauthToken string `json:"yandexPassportOauthToken"`
}

type IAMTokenResponse struct {
	IAMToken  string `json:"iamToken"`
	ExpiresAt string `json:"expiresAt"`
}

func NewClient(token string) *TtsClient {
	client := &TtsClient{
		token: token,
	}
	go func() {
		for {
			tokenResp, err := client.getIAMToken()
			if err != nil {
				fmt.Println("Ошибка при обновлении IAM токена:", err)
			} else {
				client.iamToken = tokenResp
				fmt.Println("IAM токен обновлен")
			}
			time.Sleep(2 * time.Hour)
		}
	}()

	return client
}

func (c *TtsClient) getIAMToken() (*IAMTokenResponse, error) {
	const url = "https://iam.api.cloud.yandex.net/iam/v1/tokens"
	data := `{"yandexPassportOauthToken":"` + c.token + `"}`

	resp, err := http.Post(url, "application/json", bytes.NewBufferString(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result IAMTokenResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *TtsClient) CreateVoice(text string, pathToFiles string) error {
	if c.iamToken == nil {
		iamToken, err := c.getIAMToken()
		if err != nil {
			return err
		}
		c.iamToken = iamToken
	}

	_ = os.Remove(pathToFiles)
	const (
		folderID = "b1g337dj0b6b1dcrre63"
		endpoint = "https://tts.api.cloud.yandex.net/speech/v1/tts:synthesize"
		voice    = "alena"
		lang     = "ru-RU"
	)

	data := url.Values{}
	data.Set("text", text)
	data.Add("lang", lang)
	data.Add("voice", voice)
	data.Add("folderId", folderID)

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+c.iamToken.IAMToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Fatalf("Failed with status: %s, body: %s", resp.Status, string(bodyBytes))
	}

	file, err := os.Create(pathToFiles)
	if err != nil {
		log.Fatalf("Failed to create the output file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Fatalf("Failed to write to the file: %v", err)
		return err
	}

	log.Printf("Speech saved to %s", pathToFiles)

	return nil
}
