package twilio

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// From https://www.twilio.com/blog/2017/09/send-text-messages-golang.html
// curl -X POST https://api.twilio.com/2010-04-01/Accounts/ACdcd62d5409bfe350f8b2e482b4b9cce3/Messages.json \
// --data-urlencode "Body=This is the ship that made the Kessel Run in fourteen parsecs?" \
// --data-urlencode "From=+15017122661" \
// --data-urlencode "To=+15558675310" \
// -u ACdcd62d5409bfe350f8b2e482b4b9cce3:your_auth_token

func SendText(phone string, body string) bool {
	accountSid, defined := os.LookupEnv("TWILIO_SID")
	if !defined {
		log.Fatalln("TWILIO_SID environment value not provided")
	}
	authToken, defined := os.LookupEnv("TWILIO_TOKEN")
	if !defined {
		log.Fatalln("TWILIO_TOKEN environment value not provided")
	}
	twilioNumber, defined := os.LookupEnv("TWILIO_NUMBER")
	if !defined {
		log.Fatalln("TWILIO_NUMBER environment value not provided")
	}

	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

	msgData := url.Values{}
	msgData.Set("To", phone)
	msgData.Set("From", twilioNumber)
	msgData.Set("Body", body)
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			// https://support.twilio.com/hc/en-us/articles/223134387-What-is-a-Message-SID-
			// The Message SID is the unique ID for any message successfully created by Twilio’s API.
			// It is a 34 character string that starts with “SM…” for text messages and “MM…” for media messages.
			log.Print("sent text message (sid: " + data["sid"].(string) + ")")
			return true
		}
	} else {
		log.Print("failed to send text message (status: " + resp.Status + ")")
	}

	return false
}
