package plugin
// Package plugin is the package that will execute the web-msg-handler plugins. It depends of having Node.js installed.

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Miguel-Dorta/web-msg-handler/pkg/config"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	// Directory is the subdirectory of config.Directory where plugins will be saved
	Directory = "plugins"

	// ext is the extension of the plugins
	ext       = ".js"
)

// Exec will execute the plugin with the name provided. It requires args and msg being JSON,
// the first should contain the plugin config (and therefore is up to the plugin creator to define it and check it) and
// the second will contain 3 fields: "name", "mail" and "msg", all of them strings.
func Exec(pluginName, args, msg string) error {
	pluginName += ext
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	stderr := bytes.NewBuffer(nil)

	cmd := exec.CommandContext(ctx, "node", filepath.Join(config.Directory, Directory, pluginName), args, msg)
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error executing plugin %s: %w", pluginName, err)
	}

	if stderr.Len() != 0 {
		return fmt.Errorf("%s: %s", pluginName, stderr.String())
	}
	return nil
}
