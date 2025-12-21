package audit

type Context struct {
	Event Event
}

func (c *Context) Fail(code string, msg string) {
	c.Event.Success = false
	c.Event.Details["error_code"] = code
	c.Event.Details["error"] = msg
}

func (c *Context) Success() {
	c.Event.Success = true
}
