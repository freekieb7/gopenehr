package definition

import "encoding/xml"

// Template represents an OpenEHR operational template
type Template struct {
	XMLName     xml.Name       `xml:"template"`
	Language    CodePhrase     `xml:"language"`
	Description Description    `xml:"description"`
	UID         UID            `xml:"uid"`
	TemplateID  TemplateID     `xml:"template_id"`
	Concept     string         `xml:"concept"`
	Definition  CArchetypeRoot `xml:"definition"`
	Annotations []Annotation   `xml:"annotations"`
	View        *View          `xml:"view,omitempty"`
}

type View struct {
}

type Description struct {
	OriginalAuthor []OriginalAuthorItem `xml:"original_author"`
	LifecycleState string               `xml:"lifecycle_state"`
	OtherDetails   []OtherDetailsItem   `xml:"other_details"`
	Details        DescriptionDetail    `xml:"details"`
}

type OriginalAuthorItem struct {
	ID    string `xml:"id,attr"`
	Value string `xml:",chardata"`
}

type OtherDetailsItem struct {
	ID    string `xml:"id,attr"`
	Value string `xml:",chardata"`
}

type DescriptionDetail struct {
	Language CodePhrase `xml:"language"`
	Purpose  string     `xml:"purpose"`
	Keywords string     `xml:"keywords"`
	Use      string     `xml:"use"`
}

type CodePhrase struct {
	TerminologyID TerminologyID `xml:"terminology_id"`
	CodeString    string        `xml:"code_string"`
}

type TerminologyID struct {
	Value string `xml:"value"`
}

type UID struct {
	Value string `xml:"value"`
}

type TemplateID struct {
	Value string `xml:"value"`
}

type ArchetypeID struct {
	Value string `xml:"value"`
}

// CArchetypeRoot represents the root archetype definition
type CArchetypeRoot struct {
	RMTypeName   string           `xml:"rm_type_name"`
	Occurrences  Interval         `xml:"occurrences"`
	NodeID       string           `xml:"node_id"`
	Attributes   []CAttribute     `xml:"attributes"`
	ArchetypeID  ArchetypeID      `xml:"archetype_id"`
	TemplateID   TemplateID       `xml:"template_id"`
	TermDefs     []TermDefinition `xml:"term_definitions"`
	TermBindings []TermBinding    `xml:"term_bindings"`
	Cardinality  *Cardinality     `xml:"cardinality,omitempty"`
}

// CAttribute represents either a single or multiple attribute constraint
type CAttribute struct {
	Type            string       `xml:"http://www.w3.org/2001/XMLSchema-instance type,attr"`
	RMAttributeName string       `xml:"rm_attribute_name"`
	Existence       Interval     `xml:"existence"`
	Children        []CObject    `xml:"children"`
	Cardinality     *Cardinality `xml:"cardinality,omitempty"` // Only for C_MULTIPLE_ATTRIBUTE
}

// CObject represents a constraint on an object - can be various types
type CObject struct {
	Type        string       `xml:"http://www.w3.org/2001/XMLSchema-instance type,attr"`
	RMTypeName  string       `xml:"rm_type_name"`
	Occurrences Interval     `xml:"occurrences"`
	NodeID      string       `xml:"node_id"`
	Attributes  []CAttribute `xml:"attributes"`

	// For ARCHETYPE_SLOT
	Includes []Include `xml:"includes,omitempty"`
	Excludes []Exclude `xml:"excludes,omitempty"`

	// For C_CODE_PHRASE
	TerminologyID *TerminologyID `xml:"terminology_id,omitempty"`
	CodeList      []string       `xml:"code_list,omitempty"`

	// For C_DV_QUANTITY
	Property *CodePhrase        `xml:"property,omitempty"`
	List     []QuantityListItem `xml:"list,omitempty"`

	// For C_STRING
	Pattern *string `xml:"pattern,omitempty"`
}

type Interval struct {
	LowerIncluded  bool `xml:"lower_included"`
	UpperIncluded  bool `xml:"upper_included"`
	LowerUnbounded bool `xml:"lower_unbounded"`
	UpperUnbounded bool `xml:"upper_unbounded"`
	Lower          int  `xml:"lower"`
	Upper          int  `xml:"upper"`
}

type Cardinality struct {
	IsOrdered bool     `xml:"is_ordered"`
	IsUnique  bool     `xml:"is_unique"`
	Interval  Interval `xml:"interval"`
}

type TermDefinition struct {
	Code  string               `xml:"code,attr"`
	Items []TermDefinitionItem `xml:"items"`
}

type TermDefinitionItem struct {
	ID    string `xml:"id,attr"`
	Value string `xml:",chardata"`
}

type TermBinding struct {
	Terminology string            `xml:"terminology,attr"`
	Items       []TermBindingItem `xml:"items"`
}

type TermBindingItem struct {
	Code  string     `xml:"code,attr"`
	Value CodePhrase `xml:"value"`
}

type Annotation struct {
	Path  string           `xml:"path,attr"`
	Items []AnnotationItem `xml:"items"`
}

type AnnotationItem struct {
	ID    string `xml:"id,attr"`
	Value string `xml:",chardata"`
}

// Include represents an ARCHETYPE_SLOT include expression
type Include struct {
	Expression Expression `xml:"expression"`
}

// Exclude represents an ARCHETYPE_SLOT exclude expression
type Exclude struct {
	Expression Expression `xml:"expression"`
}

// Expression represents constraint expressions in archetype slots
type Expression struct {
	Type                 string      `xml:"http://www.w3.org/2001/XMLSchema-instance type,attr"`
	TypeValue            string      `xml:"type"`
	Operator             *int        `xml:"operator,omitempty"`
	PrecedenceOverridden *bool       `xml:"precedence_overridden,omitempty"`
	LeftOperand          *ExprLeaf   `xml:"left_operand,omitempty"`
	RightOperand         *ExprLeaf   `xml:"right_operand,omitempty"`
	Item                 interface{} `xml:"item,omitempty"`
	ReferenceType        *string     `xml:"reference_type,omitempty"`
}

// ExprLeaf represents a leaf node in an expression tree
type ExprLeaf struct {
	Type          string `xml:"http://www.w3.org/2001/XMLSchema-instance type,attr"`
	TypeValue     string `xml:"type"`
	Item          string `xml:"item"`
	ReferenceType string `xml:"reference_type"`
}

// CString represents a string constraint
type CString struct {
	Type    string `xml:"http://www.w3.org/2001/XMLSchema-instance type,attr"`
	Pattern string `xml:"pattern"`
}

// QuantityListItem represents a unit constraint in C_DV_QUANTITY
type QuantityListItem struct {
	Magnitude *Interval `xml:"magnitude,omitempty"`
	Precision *Interval `xml:"precision,omitempty"`
	Units     string    `xml:"units"`
}

// IntervalFloat for float-based intervals (used in quantity constraints)
type IntervalFloat struct {
	LowerIncluded  bool     `xml:"lower_included"`
	UpperIncluded  bool     `xml:"upper_included"`
	LowerUnbounded bool     `xml:"lower_unbounded"`
	UpperUnbounded bool     `xml:"upper_unbounded"`
	Lower          *float64 `xml:"lower,omitempty"`
	Upper          *float64 `xml:"upper,omitempty"`
}
