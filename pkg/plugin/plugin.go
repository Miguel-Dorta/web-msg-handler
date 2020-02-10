package plugin

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
	Directory = "plugins"
	ext       = ".js"
)

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
