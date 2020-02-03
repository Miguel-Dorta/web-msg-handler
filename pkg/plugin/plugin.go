package plugin

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"rogchap.com/v8go"
	"strings"
)

type plugin struct {
	name, data string
}

const (
	Directory = "/etc/web-msg-handler/plugins"
	ext = ".js"
)

var (
	jsVM *v8go.Isolate
	plugins map[string]*plugin
)

func init() {
	var err error
	if jsVM, err = v8go.NewIsolate(); err != nil {
		panic("error creating isolate: " + err.Error())
	}
}

func LoadPlugins() error {
	scripts, err := ioutil.ReadDir(Directory)
	if err != nil {
		return fmt.Errorf("error listing plugin directory: %w", err)
	}

	for _, script := range scripts {
		if !script.Mode().IsRegular() || !strings.HasSuffix(script.Name(), ext) {
			continue
		}

		data, err := ioutil.ReadFile(filepath.Join(Directory, script.Name()))
		if err != nil {
			return fmt.Errorf("error reading plugin \"%s\": %w", script.Name(), err)
		}

		plugins[strings.TrimSuffix(script.Name(), ext)] = &plugin{
			name: script.Name(),
			data: string(data),
		}
	}
	return nil
}

func Exec(pluginName string, args map[string]interface{}) error {
	p := plugins[pluginName]
	if p == nil {
		return fmt.Errorf("plugin \"%s\" doesn't exists", pluginName)
	}

	jsArgs, err := getJSArgsString(args)
	if err != nil {
		return fmt.Errorf("error parsing args for plugin \"%s\": %w", pluginName, err)
	}

	ctx, err := v8go.NewContext(jsVM)
	if err != nil {
		return fmt.Errorf("error creating new JS context: %w", err)
	}

	if _, jserr := ctx.RunScript(jsArgs, p.name); jserr != nil {
		return fmt.Errorf("error declaring args vars in plugin \"%s\": %w", pluginName, jserr)
	}

	if _, jserr := ctx.RunScript(p.data, p.name); jserr != nil {
		return fmt.Errorf("plugin \"%s\" returned an error: %w", p.name, jserr)
	}
	return nil
}

func getJSArgsString(args map[string]interface{}) (string, error) {
	sb := new(strings.Builder)
	for k, v := range args {
		s, err := getJSDeclaration(k, v)
		if err != nil {
			return "", err
		}
		sb.WriteString(s)
	}
	return sb.String(), nil
}

func getJSDeclaration(name string, i interface{}) (string, error) {
	switch i.(type) {
	case string:
		return fmt.Sprintf("const %s = \"%s\"", name, i), nil
	case int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64, float32, float64, bool:
		return fmt.Sprintf("const %s = %v", name, i), nil
	default:
		return "", fmt.Errorf("type of variable %s is not supported", name)
	}
}
