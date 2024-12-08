package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type CommonErrorResponse struct {
	Type     string `json:"type,omitempty"`
	Title    string `json:"title,omitempty"`
	Status   int    `json:"status,omitempty"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}

type RecordingOutputService struct {
	Service         string `json:"service"`
	Bucket          string `json:"bucket"`
	AccessKeyId     string `json:"accessKeyId,omitempty"`
	SecretAccessKey string `json:"secretAccessKey,omitempty"`
	Region          string `json:"region,omitempty"`
	Credential      string `json:"credential,omitempty"`
}

type CreateRecordingSessionParams struct {
	Input struct {
		Kind         string `json:"kind"`
		Publications []struct {
			Id          string `json:"id"`
			ContentType string `json:"contentType,omitempty"`
		} `json:"publications"`
	} `json:"input"`
	Output struct {
		Service         string `json:"service"`
		Bucket          string `json:"bucket"`
		Credential      string `json:"credential,omitempty"`
		AccessKeyId     string `json:"accessKeyId,omitempty"`
		SecretAccessKey string `json:"secretAccessKey,omitempty"`
		Region          string `json:"region,omitempty"`
	} `json:"output"`
}

type CreateRecordingSessionResponse struct {
	Id string `json:"id"`
	CommonErrorResponse
}

type GetRecordingSessionResponse struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	Input     struct {
		Kind         string `json:"kind"`
		Publications []struct {
			Id          string `json:"id"`
			ContentType string `json:"contentType,omitempty"`
		} `json:"publications"`
	} `json:"input"`
	Output struct {
		Service         string `json:"service"`
		Bucket          string `json:"bucket"`
		Credential      string `json:"credential,omitempty"`
		AccessKeyId     string `json:"accessKeyId,omitempty"`
		SecretAccessKey string `json:"secretAccessKey,omitempty"`
		Region          string `json:"region,omitempty"`
	} `json:"output"`
	Files []struct {
		Name   string `json:"name"`
		Path   string `json:"path"`
		Status string `json:"status"`
		Errors []struct {
			Detail     string `json:"detail"`
			Level      string `json:"level"`
			OccurredAt string `json:"occurredAt"`
		} `json:"errors"`
		Type            string   `json:"type"`
		MimeType        string   `json:"mimeType"`
		CreatedAt       string   `json:"createdAt"`
		Codecs          []string `json:"codecs"`
		Duration        int      `json:"duration"`
		PublisherId     string   `json:"publisherId"`
		OriginPublisher struct {
			Id   string `json:"id"`
			Name string `json:"name,omitempty"`
		} `json:"originPublisher"`
		Publications []struct {
			Id              string `json:"id"`
			PublisherId     string `json:"publisherId"`
			OriginPublisher struct {
				Id   string `json:"id"`
				Name string `json:"name,omitempty"`
			} `json:"originPublisher"`
			ContentType string `json:"contentType"`
		} `json:"publications"`
	} `json:"files"`
	CommonErrorResponse
}

type DeleteRecordingSessionResponse struct {
	GetRecordingSessionResponse
}

func LoadRecordingOutputServiceConfig(config map[string]interface{}) (RecordingOutputService, error) {
	var outputService RecordingOutputService

	if v, ok := config["bucket"]; ok {
		outputService.Bucket = v.(string)
	}
	if v, ok := config["access_key_id"]; ok {
		outputService.AccessKeyId = v.(string)
	}
	if v, ok := config["secret_access_key"]; ok {
		outputService.SecretAccessKey = v.(string)
	}
	if v, ok := config["region"]; ok {
		outputService.Region = v.(string)
	}
	if v, ok := config["credential_file"]; ok {
		data, err := os.ReadFile(v.(string))
		if err != nil {
			return outputService, err
		}

		outputService.Credential = string(data)
	}

	return outputService, nil
}

func CreateSession(channelId string, publicationId string, contentType string, output RecordingOutputService, token string, url string) (CreateRecordingSessionResponse, error) {
	var response CreateRecordingSessionResponse

	var request CreateRecordingSessionParams
	request.Input.Kind = "SFU"

	request.Input.Publications = []struct {
		Id          string `json:"id"`
		ContentType string `json:"contentType,omitempty"`
	}{
		{
			Id:          publicationId,
			ContentType: contentType,
		},
	}

	request.Output.Service = output.Service
	request.Output.Bucket = output.Bucket
	request.Output.AccessKeyId = output.AccessKeyId
	request.Output.SecretAccessKey = output.SecretAccessKey
	request.Output.Region = output.Region
	request.Output.Credential = output.Credential

	requestBody, err := json.Marshal(request)
	if err != nil {
		return response, err
	}

	bodyReader := bytes.NewReader(requestBody)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/channels/%s/sessions", url, channelId), bodyReader)
	if err != nil {
		return response, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return response, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, err
	}

	if res.StatusCode != 201 {
		return response, fmt.Errorf("failed to create recording session status: %d body: %v", res.StatusCode, response)
	}

	// TODO: 生のrequestとresponseも返して、デバッグログで表示させたい
	return response, nil
}

func GetSession(channelId string, sessionId string, token string, url string) (GetRecordingSessionResponse, error) {
	var response GetRecordingSessionResponse
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/channels/%s/sessions/%s", url, channelId, sessionId), nil)
	if err != nil {
		return response, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return response, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, err
	}

	if res.StatusCode != 200 {
		return response, fmt.Errorf("failed to get recording session status: %d body: %v", res.StatusCode, response)
	}

	// TODO: 生のresponseも返して、デバッグログで表示させたい
	return response, nil
}

func DeleteSession(channelId string, sessionId string, token string, url string) (DeleteRecordingSessionResponse, error) {
	var response DeleteRecordingSessionResponse
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/channels/%s/sessions/%s", url, channelId, sessionId), nil)
	if err != nil {
		return response, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("User-Agent", userAgent)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return response, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, err
	}

	if res.StatusCode != 200 {
		return response, fmt.Errorf("failed to delete recording session status: %d body: %v", res.StatusCode, response)
	}

	// TODO: 生のresponseも返して、デバッグログで表示させたい
	return response, nil
}
