package audit

import (
	"github.com/gofiber/fiber/v2"
)

type KeyType string

const ContextKey KeyType = "audit_context"

type Context struct {
	Event  Event
	Failed bool
}

func (c *Context) Fail(code string, msg string) {
	c.Failed = true
	c.Event.Details["error_code"] = code
	c.Event.Details["error"] = msg
}

func (c *Context) Success() {
	c.Event.Details["outcome"] = "success"
}

// func (c *Context) Commit() {
// 	if !c.Failed {
// 		c.Event.Success = true
// 	} else {
// 		c.Event.Success = false
// 		if _, ok := c.Event.Details["outcome"]; !ok {
// 			c.Event.Details["outcome"] = "failure"
// 		}
// 	}
// 	c.Sink.Enqueue(c.Event)
// }

func From(c *fiber.Ctx) *Context {
	raw := c.Locals(string(ContextKey))
	if raw == nil {
		panic("audit context not found in fiber context")
	}
	_, ok := raw.(*Context)
	if !ok {
		panic("audit context has wrong type")
	}
	return raw.(*Context)
}
