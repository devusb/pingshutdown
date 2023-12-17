package pushover

import (
	"github.com/hashicorp/go-retryablehttp"
	"net/http"
	"net/url"
	"time"
)

type Notification struct {
	Token string
	User  string
}

func (n *Notification) Send(message string) (*http.Response, error) {
	notificationURL := "https://api.pushover.net/1/messages.json"

	client := retryablehttp.NewClient()
	client.RetryWaitMax = time.Minute * 10
	client.RetryMax = 10

	data := url.Values{
		"token":   {n.Token},
		"user":    {n.User},
		"message": {message},
	}

	resp, err := client.PostForm(notificationURL, data)

	return resp, err
}
