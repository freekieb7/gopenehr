package audit

import (
	"net"
	"time"

	"github.com/freekieb7/gopenehr/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type KeyType string

const ContextKey KeyType = "audit_context"

type Context struct {
	Logger *Logger
	Event  Event
	Failed bool
}

func (l *Logger) NewContext(c *fiber.Ctx, resource Resource, action Action) *Context {
	ip := net.ParseIP(c.IP())

	return &Context{
		Logger: l,
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

func (c *Context) SetEHR(id uuid.UUID) {
	c.Event.Details["ehr_id"] = id.String()
}

func (c *Context) SetValidationError(err error) {
	c.Event.Details["validation_error"] = err
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
	c.Logger.Log(c.Event)
}

func From(c *fiber.Ctx) *Context {
	raw := c.Locals(string(ContextKey))
	if raw == nil {
		return nil
	}
	return raw.(*Context)
}
