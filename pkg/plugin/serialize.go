package plugin

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func ArgsToJS(args map[string]interface{}) (string, error) {
	data, err := json.Marshal(args)
	if err != nil {
		return "", fmt.Errorf("error serializing args JS: %w", err)
	}

	return jsonConstDeclaration("args", data), nil
}

func MsgToJS(name, mail, msg string) (string, error) {
	data, err := json.Marshal(map[string]string{
		"name": name,
		"mail": mail,
		"msg": msg,
	})
	if err != nil {
		return "", fmt.Errorf("error serializing msg JS: %w", err)
	}

	return jsonConstDeclaration("msg", data), nil
}

func jsonConstDeclaration(constName string, jsonData []byte) string {
	jsonData = bytes.ReplaceAll(jsonData, []byte{'\''}, []byte{'\\', '\''})
	return fmt.Sprintf("const %s = JSON.parse('%s');", constName, string(jsonData))
}
