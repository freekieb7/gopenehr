package oauth

type Scope string

const (
	ScopeAuditRead     Scope = "audit:read"
	ScopeWebhookManage Scope = "webhook:manage"
)

func (s Scope) String() string {
	return string(s)
}
