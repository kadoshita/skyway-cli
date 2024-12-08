package internal

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/ybbus/jsonrpc/v3"
)

type Member struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Subtype  string `json:"subtype"`
	Metadata string `json:"metadata"`
}

type Publication struct {
	Id                string `json:"id"`
	PublisherId       string `json:"publisherId"`
	ContentType       string `json:"contentType"`
	IsEnabled         bool   `json:"isEnabled"`
	OriginId          string `json:"originId"`
	OriginPublisherId string `json:"originPublisherId"`
	Metadata          string `json:"metadata"`
}

type Subscription struct {
	Id            string `json:"id"`
	PublicationId string `json:"publicationId"`
	SubscriberId  string `json:"subscriberId"`
}

type Channel struct {
	Id            string         `json:"id"`
	Name          string         `json:"name"`
	Metadata      string         `json:"metadata"`
	Members       []Member       `json:"members"`
	Publications  []Publication  `json:"publications"`
	Subscriptions []Subscription `json:"subscriptions"`
}

type FindChannelParams struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type FindChannelResult struct {
	Channel Channel `json:"channel"`
}

type CreateChannelParams struct {
	Name     string `json:"name,omitempty"`
	Metadata string `json:"metadata,omitempty"`
}

type CreateChannelResult struct {
	Channel Channel `json:"channel"`
}

const subscribeChannelEventsRequest = `{
	"id":"%s",
	"jsonrpc":"2.0",
	"method":"subscribeChannelEvents",
	"params":{
		"authToken":"%s",
		"appId":"%s",
		"channelId":"%s"
	}
}`

func GetChannel(id string, name string, token string, url string) (Channel, error) {
	rpcClient := jsonrpc.NewClientWithOpts(url, &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			"Authorization": "Bearer " + token,
			"User-Agent":    userAgent,
		},
	})
	var channel *FindChannelResult
	err := rpcClient.CallFor(context.Background(), &channel, "findChannel", &FindChannelParams{Id: id, Name: name})

	if err != nil || channel == nil {
		return Channel{}, err
	}

	return channel.Channel, nil
}

func CreateChannel(name string, metadata string, token string, url string) (Channel, error) {
	rpcClient := jsonrpc.NewClientWithOpts(url, &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			"Authorization": "Bearer " + token,
			"User-Agent":    userAgent,
		},
	})
	var result *CreateChannelResult
	err := rpcClient.CallFor(context.Background(), &result, "createChannel", &CreateChannelParams{Name: name, Metadata: metadata})

	if err != nil || result == nil {
		return Channel{}, err
	}

	return result.Channel, nil
}

func SubscribeEvents(id string, name string, token string, appId string, url string, handler chan string) error {
	client, _, err := websocket.DefaultDialer.Dial(url, http.Header{"Sec-WebSocket-Protocol": []string{token}})
	if err != nil {
		return err
	}
	defer client.Close()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		err := client.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(subscribeChannelEventsRequest, uuid.New().String(), token, appId, id)))
		if err != nil {
			return
		}

		for {
			_, message, err := client.ReadMessage()
			if err != nil {
				return
			}

			handler <- string(message)
		}
	}()

	<-interrupt
	fmt.Println("shutting down...")

	// ref: https://github.com/gorilla/websocket/issues/448
	err = client.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return err
	}
	return nil
}
