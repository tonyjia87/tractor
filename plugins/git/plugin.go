package git

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

type Plugin struct {
	Repo    Repo
	Build   Build
	Config  Config
	Backoff Backoff
}


func (p Plugin) Exec() error {
	if p.Build.Path != "" {
		err := os.MkdirAll(p.Build.Path, 0777)
		if err != nil {
			return err
		}
	}
	var cmds []*exec.Cmd
	switch {
	case issPullRequest(p.Build.Event) || isTag(p.Build.Event, p.Build.Ref):
		cmds = append(cmds, fetch(p.Build.Ref, p.Config.Tags, p.Config.Depth))
		cmds = append(cmds, checkoutHead())
	default:
		cmds = append(cmds, fetch(p.Build.Ref, p.Config.Tags, p.Config.Depth))
		cmds = append(cmds, checkoutSha(p.Build.Commit))
	}

	for _, cmd := range cmds {
		buf := new(bytes.Buffer)
		cmd.Dir = p.Build.Path
		cmd.Stdout = io.MultiWriter(os.Stdout, buf)
		cmd.Stderr = io.MultiWriter(os.Stderr, buf)
		trace(cmd)
		err := cmd.Run()
		switch {
		case err != nil && shouldRetry(buf.String()):
			err = retryExec(cmd, p.Backoff.Duration, p.Backoff.Attempts)
			if err != nil {
				return err
			}
		case err != nil:
			return err
		}
	}
	return nil
}
