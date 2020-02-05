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
	ext       = ".js"
)

var (
	jsVM    *v8go.Isolate
	plugins map[string]*plugin
)

func Start() error {
	var err error
	if jsVM, err = v8go.NewIsolate(); err != nil {
		panic("error creating isolate: " + err.Error())
	}

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

func Exec(pluginName, args, msg string) error {
	p := plugins[pluginName]
	if p == nil {
		return fmt.Errorf("plugin \"%s\" doesn't exists", pluginName)
	}

	ctx, err := v8go.NewContext(jsVM)
	if err != nil {
		return fmt.Errorf("error creating new JS context: %w", err)
	}

	sb := new(strings.Builder)
	sb.WriteString(args)
	sb.WriteString(msg)
	sb.WriteString(p.data)

	if _, jsErr := ctx.RunScript(sb.String(), p.name); jsErr != nil {
		return fmt.Errorf("plugin \"%s\" returned an error: %w", p.name, jsErr)
	}
	return nil
}
