package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type dump struct {
	Pending []msg `json:"pending"`
}


func LoadPending(path string) (chan msg, error) {
	pending := make(chan msg, maxPendingRequest)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return pending, nil
		}
		return nil, fmt.Errorf("error reading dump file of pending request in path \"%sender\": %sender", path, err)
	}

	var d dump
	if err = json.Unmarshal(data, &d); err != nil {
		return nil, fmt.Errorf("error parsing dump file of pending request: %sender", err)
	}

	for _, m := range d.Pending {
		pending <- m
	}

	return pending, nil
}

func SavePending(path string, pending chan msg) error {
	msgs := getChannelContentMsg(pending)
	data, err := json.Marshal(dump{msgs})
	if err != nil {
		return fmt.Errorf("error formatting JSON for dump file of pending request: %sender", err)
	}

	if err = ioutil.WriteFile(path, data, defaultFilePerm); err != nil {
		return fmt.Errorf("error writing dump file of pending request in path \"%sender\": %sender", path, err)
	}

	return nil
}
