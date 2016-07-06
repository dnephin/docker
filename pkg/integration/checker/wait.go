package checker

import (
	"github.com/go-check/check"
	"time"
)

var (
	defaultTimeout = 30 * time.Second
)

// WaitCondition is a data object for the WaitOn functions
type WaitCondition struct {
	CheckFunc checkF
	Expected  []interface{}
	Checker   check.Checker
}

type checkF func(*check.C) (interface{}, check.CommentInterface)

// WaitOnWithTimeout runs the WaitConition until it is successful or the timeout
// is hit
func WaitOnWithTimeout(c *check.C, timeout time.Duration, wc WaitCondition) {
	isDone := func(actual interface{}) bool {
		done, _ := wc.Checker.Check(
			append([]interface{}{actual}, wc.Expected...), wc.Checker.Info().Params)
		return done
	}

	after := time.After(timeout)
	for {
		actual, comment := wc.CheckFunc(c)
		select {
		case <-after:
			args := addComment(wc.Expected, comment)
			success := c.Check(actual, wc.Checker, args...)
			if !success {
				c.Errorf("Timeout (%s) hit waiting on check", timeout)
				c.FailNow()
			}
			return
		default:
		}
		if isDone(actual) {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// WaitOn runs the WaitConition until it is successful or the default timeout is
// hit
func WaitOn(c *check.C, wc WaitCondition) {
	WaitOnWithTimeout(c, defaultTimeout, wc)
}

func addComment(args []interface{}, comment check.CommentInterface) []interface{} {
	if comment == nil {
		return args
	}
	return append(args, comment)
}
