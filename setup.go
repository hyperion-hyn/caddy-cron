package caddy_cron

import (
	"fmt"
	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
	"github.com/robfig/cron"
	"os/exec"
)

type Task struct {
	Moment string
	Command string
}

func init() {
	pluginN := "ccron"
	caddy.RegisterPlugin(pluginN, caddy.Plugin{
		ServerType: "http",
		Action: setup,
	})
	httpserver.RegisterDevDirective(pluginN, "")
}

func setup(c *caddy.Controller) error {
	missions, err := parseConfig(c)
	if err != nil {
		return err
	}

	if len(missions) > 0 {
		go startTask(missions)
	} else {
		println("cron task is not set...")
	}

	return nil
}

func startTask(tasks []*Task) {
	c := cron.New()

	for _, task := range tasks {
		err := c.AddJob(task.Moment, &funcCommandJob{task.Command, execute})

		if err != nil {
			println(fmt.Sprintf("[error] start `%s` at `%s` get error, %v", task.Command, task.Moment, err))
		} else {
			println(fmt.Sprintf("[info] set `%s` at `%s` success", task.Moment, task.Command))
		}
	}

	c.Start()
	defer c.Stop()
	select{}
}

func execute(command string) {
	cmd := exec.Command("sh", "-c", command)
	out, err := cmd.Output()
	if err != nil {
		println(fmt.Sprintf("[error] execute command `%s` get error, %v", command, err))
	} else {
		println(fmt.Sprintf("[info] execute command `%s` success, output: %v", command, string(out)))
	}
}

type funcCommandJob struct {
	command   string
	function func(string)
}

func (f *funcCommandJob) Run() {
	if nil != f.function {
		f.function(f.command)
	}
}

