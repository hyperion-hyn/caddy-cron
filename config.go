package caddy_cron

import (
	"errors"
	"fmt"
	"github.com/mholt/caddy"
	"reflect"
	"strings"
	"time"
)

func parseConfig(c *caddy.Controller) ([]*Task, error) {
	tasks := make([]*Task, 0)

	for c.Next() {
		if !c.NextArg() {
			return tasks, c.ArgErr()
		}
		moment := c.Val()
		command := c.RemainingArgs()

		task, err := parseTask(moment, command)
		if err != nil {
			return tasks, err
		}

		tasks = append(tasks, &task)
	}

	return tasks, nil
}

/*
parse task
- Moment:
Just support `Predefined schedules` and `time Intervals`, For more detail, just view https://godoc.org/github.com/robfig/cron
support moment:
	Entry                  | Description                                | Equivalent To
	-----                  | -----------                                | -------------
	@yearly (or @annually) | Run once a year, midnight, Jan. 1st        | 0 0 0 1 1 *
	@monthly               | Run once a month, midnight, first of month | 0 0 0 1 * *
	@weekly                | Run once a week, midnight between Sat/Sun  | 0 0 0 * * 0
	@daily (or @midnight)  | Run once a day, midnight                   | 0 0 0 * * *
 	@hourly                | Run once an hour, beginning of hour        | 0 0 * * * *
	----------------------------------------------------------------------------------
 	You may also schedule a job to execute at fixed intervals, starting at the time it's added or cron is run.
 	This is supported by formatting the cron spec like this:
 	`@every:<duration>`
 	where "duration" is a string accepted by time.ParseDuration (http://golang.org/pkg/time/#ParseDuration).
- Command: is the command to execute; it may be followed by arguments, also look for https://caddyserver.com/docs/on
*/
func parseTask(moment string, command []string) (Task, error) {
	task := Task{}

	cronM, err := parseMoment(moment)
	if err != nil {
		return task, err
	}

	task.Moment = cronM
	task.Command = strings.Join(command, " ")

	return task, nil
}

func parseMoment(moment string) (string, error) {
	allows := []string{"@yearly", "@annually", "@monthly", "@weekly", "@daily", "@midnight", "@hourly", "@every"}
	moments := strings.Split(moment, ":")

	realM := ""
	cronM := ""
	if len(moments) == 1 {
		realM = moment
	} else if len(moments) == 2 {
		realM = moments[0]

		// check duration
		d, err := time.ParseDuration(moments[1])
		if err != nil || d < 1 {
			return cronM, errors.New(fmt.Sprintf("not support `moment` config, %v", err))
		}
	} else {
		return cronM, errors.New("not support `moment` config")
	}

	if e, _ := in_array(realM, allows); !e {
		return cronM, nil
	}

	return strings.Join(moments, " "), nil
}

func in_array(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}
