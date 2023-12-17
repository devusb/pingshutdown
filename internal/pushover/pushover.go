package pushover

import (
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

type Notification struct {
	Token 	     string
	User         string
	RetryWaitMax time.Duration
	RetryMax     int
}

func (n *Notification) Send(message string) (*http.Response, error) {
	notificationURL := "https://api.pushover.net/1/messages.json"

	client := retryablehttp.NewClient()
	client.RetryWaitMax = n.RetryWaitMax
	client.RetryMax = n.RetryMax

	data := url.Values{
		"token":   {n.Token},
		"user":    {n.User},
		"message": {message},
	}

	resp, err := client.PostForm(notificationURL, data)

	return resp, err
}
