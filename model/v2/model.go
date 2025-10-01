package v2

import "encoding/json"

type Option[T any] struct {
	value *T
}

func Some[T any](v T) Option[T] {
	return Option[T]{value: &v}
}

func None[T any]() Option[T] {
	return Option[T]{value: nil}
}

func (o Option[T]) IsSome() bool {
	return o.value != nil
}

func (o Option[T]) IsNone() bool {
	return o.value == nil
}

func (o Option[T]) Unwrap() T {
	if o.value == nil {
		panic("called Unwrap on a None value")
	}
	return *o.value
}

func (o Option[T]) UnwrapOr(defaultValue T) T {
	if o.value == nil {
		return defaultValue
	}
	return *o.value
}

func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.value == nil {
		return []byte("null"), nil
	}
	return json.Marshal(*o.value)
}

func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		o.value = nil
		return nil
	}
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	o.value = &v
	return nil
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Bot  bool   `json:"bot"`
}

type Group struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (g Group) Members() []User {
	return []User{}
}

type Folder struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (f *Folder) Parent() Option[Folder] {
	return None[Folder]()
}

func (f Folder) Children() []Folder {
	return []Folder{}
}

func (f Folder) Documents() []Document {
	return []Document{}
}

type Document struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (d *Document) Parent() Option[Folder] {
	return None[Folder]()
}

func (d Document) Revisions() []DocumentRevision {
	return []DocumentRevision{}
}

type DocumentRevision struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	AuthorID  string `json:"author_id"`
}

func (dr DocumentRevision) Author() User {
	return User{}
}

func (dr DocumentRevision) Content() Option[string] {
	return None[string]()
}

type Audit struct {
	ID string `json:"id"`
}

func (a Audit) Events() []AuditEvent {
	return []AuditEvent{}
}

type AuditEvent struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	UserID    string `json:"user_id"`
	Action    string `json:"action"`
	Reason    string `json:"reason"`
}

func (a AuditEvent) User() User {
	return User{}
}

func (a AuditEvent) Changes() []AuditChange {
	return []AuditChange{}
}

type AuditChange struct {
	ID       string         `json:"id"`
	NewValue Option[string] `json:"new_value"`
	OldValue Option[string] `json:"old_value"`
}

type Webhook struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedBy string `json:"created_by"`
	URL       string `json:"url"`
	Active    bool   `json:"active"`
}

type WebhookEvent struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	WebhookID string `json:"webhook_id"`
	EventType string `json:"event_type"`
	Payload   string `json:"payload"`
}

func (w WebhookEvent) Webhook() Webhook {
	return Webhook{}
}

type Calender struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c Calender) Events() []CalenderEvent {
	return []CalenderEvent{}
}

type CalenderEvent struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Location    string `json:"location"`
}

func (c CalenderEvent) Attendees() []User {
	return []User{}
}
