package audit

import (
	"net"
	"time"

	"github.com/freekieb7/gopenehr/internal/config"
	"github.com/gofiber/fiber/v2"
)

type KeyType string

const ContextKey KeyType = "audit_context"

type Context struct {
	Sink   *Sink
	Event  Event
	Failed bool
}

func (s *Sink) NewContext(c *fiber.Ctx, resource Resource, action Action) *Context {
	ip := net.ParseIP(c.IP())

	return &Context{
		Sink: s,
		Event: Event{
			ActorID:   config.SystemUserID,
			ActorType: "system",
			Resource:  string(resource),
			Action:    string(action),
			IPAddress: ip,
			UserAgent: c.Get("User-Agent"),
			Details:   map[string]any{},
			CreatedAt: time.Now(),
		},
	}
}

func (c *Context) Fail(code string, msg string) {
	c.Failed = true
	c.Event.Details["error_code"] = code
	c.Event.Details["error"] = msg
}

func (c *Context) Success() {
	c.Event.Details["outcome"] = "success"
}

func (c *Context) Commit() {
	if !c.Failed {
		c.Event.Success = true
	} else {
		c.Event.Success = false
		if _, ok := c.Event.Details["outcome"]; !ok {
			c.Event.Details["outcome"] = "failure"
		}
	}
	c.Sink.Enqueue(c.Event)
}

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
