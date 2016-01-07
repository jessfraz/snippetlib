package mailchimp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Sirupsen/logrus"
)

const (
	// APIURI is the host for Mailchimp API requests.
	APIURI = "https://us7.api.mailchimp.com"
)

var (
	// SubscribeURL is the url to send subscription requests to.
	SubscribeURL = APIURI + "/2.0/lists/subscribe"
)

// SubscribeRequest defines the parameters sent to the mailchimp API to
// subscribe someone to a newsletter.
type SubscribeRequest struct {
	APIKey      string `json:"apikey"`
	ID          string `json:"id"`
	Email       Email  `json:"email"`
	DoubleOptin bool   `json:"double_optin"`
	SendWelcome bool   `json:"send_welcome"`
}

// Email defines a user's email.
type Email struct {
	Email string `json:"email"`
}

// Subscribe takes a mailchimp apikey, list id, and email address and subscribes
// it to a mailchimp newsletter.
func Subscribe(apikey, listID, email string) error {
	req := SubscribeRequest{
		APIKey: apikey,
		ID:     listID,
		Email: Email{
			Email: email,
		},
		DoubleOptin: false,
		SendWelcome: false,
	}

	// marshal the data to send
	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal subscribe request failed: %v", err)
	}

	// create the request
	client := &http.Client{}
	r, err := http.NewRequest("POST", SubscribeURL, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("creating request to %s failed: %v", SubscribeURL, err)
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Content-Length", strconv.Itoa(len(data)))

	// make the request
	resp, err := client.Do(r)
	if err != nil {
		return fmt.Errorf("completing request to %s failed: %v", SubscribeURL, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading body from subscribe response failed: %v", err)
	}

	logrus.Debugf("subscribe response: %s", body)

	return nil
}
