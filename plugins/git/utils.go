package git

import (
	"os"
	"os/exec"
	"strings"
	"fmt"
	"time"
)

// Checkout executes a git checkout command.
func checkoutHead() *exec.Cmd {
	return exec.Command(
		"git",
		"checkout",
		"-qf",
		"FETCH_HEAD",
	)
}

// Checkout executes a git checkout command.
func checkoutSha(commit string) *exec.Cmd {
	return exec.Command(
		"git",
		"reset",
		"--hard",
		"-q",
		commit,
	)
}

// fetch retuns git command that fetches from origin. If tags is true
// then tags will be fetched.
func fetch(ref string, tags bool, depth int) *exec.Cmd {
	tags_option := "--no-tags"
	if tags {
		tags_option = "--tags"
	}
	cmd := exec.Command(
		"git",
		"fetch",
		tags_option,
	)
	if depth != 0 {
		cmd.Args = append(cmd.Args, fmt.Sprintf("--depth=%d", depth))
	}
	cmd.Args = append(cmd.Args, "origin")
	cmd.Args = append(cmd.Args, fmt.Sprintf("+%s:", ref))
	return cmd
}

// updateSubmodules recursively initializes and updates submodules.
func updateSubmodules(remote bool) *exec.Cmd {
	cmd := exec.Command(
		"git",
		"submodule",
		"update",
		"--init",
		"--recursive",
	)

	if remote {
		cmd.Args = append(cmd.Args, "--remote")
	}

	return cmd
}

// skipVerify returns a git command that, when executed configures git to skip
// ssl verification. This should may be used with self-signed certificates.
func skipVerify() *exec.Cmd {
	return exec.Command(
		"git",
		"config",
		"--global",
		"http.sslVerify",
		"false",
	)
}

// helper function returns true if the commit is a pull_request.
func issPullRequest(event string) bool {
	return event == "pull_request"
}

// helper function returns true if the commit is a tag.
func isTag(event, ref string) bool {
	return event == "tag" ||
		strings.HasPrefix(ref, "refs/tags/")
}

// trace writes the command in the programs stdout for debug purposes.
// the command is wrapped in xml tags for easy parsing.
func trace(cmd *exec.Cmd) {
	fmt.Printf("+ %s\n", strings.Join(cmd.Args, " "))
}

// retryExec is a helper function that retries a command.
func retryExec(cmd *exec.Cmd, backoff time.Duration, retries int) (err error) {
	for i := 0; i < retries; i++ {
		// signal intent to retry
		fmt.Printf("retry in %v\n", backoff)

		// wait 5 seconds before retry
		<-time.After(backoff)

		// copy the original command
		retry := exec.Command(cmd.Args[0], cmd.Args[1:]...)
		retry.Dir = cmd.Dir
		retry.Stdout = os.Stdout
		retry.Stderr = os.Stderr
		trace(retry)
		err = retry.Run()
		if err == nil {
			return
		}
	}
	return
}

// shouldRetry returns true if the command should be re-executed. Currently
// this only returns true if the remote ref does not exist.
func shouldRetry(s string) bool {
	return strings.Contains(s, "find remote ref")
}
