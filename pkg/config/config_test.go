package config_test

import (
	"fmt"
	"github.com/Miguel-Dorta/web-msg-handler/pkg/config"
	"github.com/Miguel-Dorta/web-msg-handler/pkg/sender"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	configPath, err := downloadConfig()
	if err != nil {
		t.Skipf("TestLoadConfig requires an Internet connection and a writable \"/tmp\" directory with a few KiB of free space. Skipping details: %s", err)
	}
	defer os.Remove(configPath)

	senders, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("execution error of LoadConfig: %s", err)
	}

	if len(senders) != 4 {
		t.Error("not all senders read")
	}

	sender1, exists := senders[5577006791947779410]
	if !exists {
		t.Error("Sender with ID 5577006791947779410 not found")
	} else {
		sender1telegram, ok := sender1.(*sender.Telegram)
		if !ok {
			t.Error("Sender with ID 5577006791947779410 cannot be parsed as Telegram Sender")
		} else {
			if sender1telegram.URL != "website1.com" {
				t.Errorf("URL don't match: found %s", sender1telegram.URL)
			}
			if sender1telegram.RecaptchaSecret != "Uv38ByGCZU8WP18PmmIdcpVmx00QA3xNe7sEB9Hi" {
				t.Errorf("RecaptchaSecret don't match: found %s", sender1telegram.RecaptchaSecret)
			}
			if sender1telegram.ChatId != "9167320" {
				t.Errorf("ChatId don't match: found %s", sender1telegram.ChatId)
			}
			if sender1telegram.BotToken != "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11" {
				t.Errorf("BotToken don't match: found %s", sender1telegram.BotToken)
			}
		}
	}

	sender2, exists := senders[15352856648520921629]
	if !exists {
		t.Error("Sender with ID 15352856648520921629 not found")
	} else {
		sender2telegram, ok := sender2.(*sender.Telegram)
		if !ok {
			t.Error("Sender with ID 15352856648520921629 cannot be parsed as Telegram Sender")
		} else {
			if sender2telegram.URL != "website2.org" {
				t.Errorf("URL don't match: found %s", sender2telegram.URL)
			}
			if sender2telegram.RecaptchaSecret != "xkmBhVrYaB0NhtHpHgAWeTnLZpTSxCKs0gigByk5" {
				t.Errorf("RecaptchaSecret don't match: found %s", sender2telegram.RecaptchaSecret)
			}
			if sender2telegram.ChatId != "87745566" {
				t.Errorf("ChatId don't match: found %s", sender2telegram.ChatId)
			}
			if sender2telegram.BotToken != "654321:ABC-DEF1234ghIkl-zyx57W2v1u123ew12" {
				t.Errorf("BotToken don't match: found %s", sender2telegram.BotToken)
			}
		}
	}

	sender3, exists := senders[8674665223082153551]
	if !exists {
		t.Error("Sender with ID 8674665223082153551 not found")
	} else {
		sender3mail, ok := sender3.(*sender.Mail)
		if !ok {
			t.Error("Sender with ID 8674665223082153551 cannot be parsed as Mail Sender")
		} else {
			if sender3mail.URL != "website3.com" {
				t.Errorf("URL don't match: found %s", sender3mail.URL)
			}
			if sender3mail.RecaptchaSecret != "SH9pmeudGKRHhARdh_PGfPInRumVr1olNnlRuqL_" {
				t.Errorf("RecaptchaSecret don't match: found %s", sender3mail.RecaptchaSecret)
			}
			if sender3mail.Mailto != "contact@website3.com" {
				t.Errorf("Mailto don't match: found %s", sender3mail.Mailto)
			}
			if sender3mail.Username != "no-reply@website3.com" {
				t.Errorf("Username don't match: found %s", sender3mail.Username)
			}
			if sender3mail.Password != "bNRxxIPxX7kLrbN8WCG22VUmpBqVBGgLTnyLdjob" {
				t.Errorf("Password don't match: found %s", sender3mail.Password)
			}
			if sender3mail.Hostname != "smtp.website3.com" {
				t.Errorf("Hostname don't match: found %s", sender3mail.Hostname)
			}
			if sender3mail.Port != "587" {
				t.Errorf("Port don't match: found %s", sender3mail.Port)
			}
		}
	}

	sender4, exists := senders[13260572831089785859]
	if !exists {
		t.Error("Sender with ID 13260572831089785859 not found")
	} else {
		sender4mail, ok := sender4.(*sender.Mail)
		if !ok {
			t.Error("Sender with ID 13260572831089785859 cannot be parsed as Mail Sender")
		} else {
			if sender4mail.URL != "website4.net" {
				t.Errorf("URL don't match: found %s", sender4mail.URL)
			}
			if sender4mail.RecaptchaSecret != "HUnUlVyEhiFjJSU_7HON16nii_khEZwWDwcCRIYV" {
				t.Errorf("RecaptchaSecret don't match: found %s", sender4mail.RecaptchaSecret)
			}
			if sender4mail.Mailto != "personal-mail@gmail.com" {
				t.Errorf("Mailto don't match: found %s", sender4mail.Mailto)
			}
			if sender4mail.Username != "contact-forms@website4.net" {
				t.Errorf("Username don't match: found %s", sender4mail.Username)
			}
			if sender4mail.Password != "u9oIMT9qjrZo0gv1BZh1kh5milvfLH_EhEWS0lcr" {
				t.Errorf("Password don't match: found %s", sender4mail.Password)
			}
			if sender4mail.Hostname != "mail.website4.net" {
				t.Errorf("Hostname don't match: found %s", sender4mail.Hostname)
			}
			if sender4mail.Port != "25" {
				t.Errorf("Port don't match: found %s", sender4mail.Port)
			}
		}
	}
}

// downloadConfig is a helper of TestLoadConfig.
// It will download the example config from the official repo and save it in "/tmp/web-msg-handler_config_<timestamp>.json".
// It will return the filepath of the config or an error if the operation failed.
func downloadConfig() (string, error) {
	resp, err := http.Get("https://raw.githubusercontent.com/Miguel-Dorta/web-msg-handler/master/examples/config.json")
	if err != nil {
		return "", fmt.Errorf("failed to get config: %s", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read the request body: %s", err)
	}

	if err = resp.Body.Close(); err != nil {
		return "", fmt.Errorf("error closing body request: %s", err)
	}

	filename := "/tmp/web-msg-handler_config_" + strconv.FormatInt(time.Now().UnixNano(), 10) + ".json"
	if err = ioutil.WriteFile(filename, data, 0777); err != nil {
		return "", fmt.Errorf("error writing file in \"%s\": %s", filename, err)
	}

	return filename, nil
}
