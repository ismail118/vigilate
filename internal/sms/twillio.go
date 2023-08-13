package sms

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/ismail118/vigilate/internal/config"
)

func SendTextTwillio(to, msg string, app *config.AppConfig) error {
	secret := app.PreferenceMap["twillio_auth_token"]
	key := app.PreferenceMap["twillio_sid"]

	urlString := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", key)

	msgData := url.Values{}
	msgData.Set("To", to)
	msgData.Set("From", app.PreferenceMap["twilio_phone_number"])
	msgData.Set("Body", msg)

	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlString, &msgDataReader)
	
	// auth
	req.SetBasicAuth(key, secret)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(req)
	if resp.StatusCode >= 200  && resp.StatusCode < 300 {
		var data map[string]interface{} 
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err != nil {
			log.Panicln(err)
			return err
		}
	} else {
		log.Panicln("Error sending sms!")
		return errors.New("error sending sms!")
	}

	return nil
}