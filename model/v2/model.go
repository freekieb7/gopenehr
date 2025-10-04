package v2

import "github.com/freekieb7/gopenehr/util"

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

func (f *Folder) Parent() util.Optional[Folder] {
	return util.None[Folder]()
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

func (d *Document) Parent() util.Optional[Folder] {
	return util.None[Folder]()
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

func (dr DocumentRevision) Content() util.Optional[string] {
	return util.None[string]()
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
	ID       string                `json:"id"`
	NewValue util.Optional[string] `json:"new_value"`
	OldValue util.Optional[string] `json:"old_value"`
}
