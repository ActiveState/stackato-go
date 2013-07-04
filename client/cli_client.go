// Client interface using the official stackato client binary

package client

import (
	"fmt"
	"github.com/ActiveState/log"
	"github.com/ActiveState/run"
	"os/exec"
)

type CliClient struct {
	TargetURL string
	Token     string
	Group     string
}

func NewCliClient(targetUrl, token, group string) (*CliClient, error) {
	if token == "" {
		return nil, fmt.Errorf("token string must not be empty")
	}
	c := &CliClient{targetUrl, token, group}
	return c, nil
}

// PushAppNoCreate emulates `s push --no-create ...` and sends the
// output in outputCh channel.
func (c *CliClient) PushAppNoCreate(name string, dir string, autoStart bool, outputCh chan string) (bool, error) {
	options := []string{
		name,
		"--no-tail", "--no-prompt",
		"--target", c.TargetURL,
		"--token", c.Token,
		"--path", dir}

	if !autoStart {
		options = append(options, "--no-start")
	}

	if c.Group != "" {
		options = append(options, "--group", c.Group)
	}

	pushOptions := append([]string{"push", "--no-create"}, options...)

	ret, err := run.Run(exec.Command("stackato", pushOptions...), outputCh)
	if err != nil {
		log.Error("cannot read line: ", err)
		return false, err
	}
	if r, ok := ret.(*exec.ExitError); ok {
		log.Errorf("Client exited abruptly: %v", r)
		return false, nil
	}else{
		return true, ret
	}
}
