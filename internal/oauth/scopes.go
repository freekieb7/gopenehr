package oauth

type Scope string

const (
	ScopeAuditRead         Scope = "audit:read"
	ScopeWebhookManage     Scope = "webhook:manage"
	ScopeTenantManage      Scope = "tenant:manage"
	ScopeEHRRead           Scope = "ehr:read"
	ScopeEHRWrite          Scope = "ehr:write"
	ScopeEHRDelete         Scope = "ehr:delete"
	ScopeDemographicsRead  Scope = "demographics:read"
	ScopeDemographicsWrite Scope = "demographics:write"
	ScopeQueryRead         Scope = "query:read"
	ScopeQueryWrite        Scope = "query:write"
	ScopeQueryExecute      Scope = "query:execute"
	ScopeTemplateRead      Scope = "template:read"
	ScopeTemplateWrite     Scope = "template:write"
)

func (s Scope) String() string {
	return string(s)
}
