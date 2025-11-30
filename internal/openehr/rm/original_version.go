package rm

import (
	"encoding/json"
	"fmt"

	"github.com/freekieb7/gopenehr/internal/openehr/util"
	"github.com/freekieb7/gopenehr/pkg/utils"
)

const ORIGINAL_VERSION_TYPE = "ORIGINAL_VERSION"

type ORIGINAL_VERSION struct {
	Type_                 utils.Optional[string]              `json:"_type,omitzero"`
	UID                   OBJECT_VERSION_ID                   `json:"uid"`
	PrecedingVersionUID   utils.Optional[OBJECT_VERSION_ID]   `json:"preceding_version_uid,omitzero"`
	OtherInputVersionUIDs utils.Optional[[]OBJECT_VERSION_ID] `json:"other_input_version_uids,omitzero"`
	LifecycleState        DV_CODED_TEXT                       `json:"lifecycle_state"`
	Attestations          utils.Optional[[]ATTESTATION]       `json:"attestations,omitzero"`
	Data                  OriginalVersionDataUnion            `json:"data"`
}

func (ov *ORIGINAL_VERSION) SetModelName() {
	ov.Type_ = utils.Some(ORIGINAL_VERSION_TYPE)
	ov.UID.SetModelName()
	if ov.PrecedingVersionUID.E {
		ov.PrecedingVersionUID.V.SetModelName()
	}
	if ov.OtherInputVersionUIDs.E {
		for i := range ov.OtherInputVersionUIDs.V {
			ov.OtherInputVersionUIDs.V[i].SetModelName()
		}
	}
	ov.LifecycleState.SetModelName()
	if ov.Attestations.E {
		for i := range ov.Attestations.V {
			ov.Attestations.V[i].SetModelName()
		}
	}

}

func (ov *ORIGINAL_VERSION) Validate(path string) util.ValidateError {
	var validateErr util.ValidateError
	var attrPath string

	// Validate _type
	if ov.Type_.E && ov.Type_.V != "ORIGINAL_VERSION" {
		attrPath = path + "._type"
		validateErr.Errs = append(validateErr.Errs, util.ValidationError{
			Model:          "ORIGINAL_VERSION",
			Path:           attrPath,
			Message:        "_type must be ORIGINAL_VERSION",
			Recommendation: "Set _type to ORIGINAL_VERSION",
		})
	}

	// Validate uid
	attrPath = path + ".uid"
	validateErr.Errs = append(validateErr.Errs, ov.UID.Validate(attrPath).Errs...)

	// Validate preceding_version_uid
	if ov.PrecedingVersionUID.E {
		attrPath = path + ".preceding_version_uid"
		validateErr.Errs = append(validateErr.Errs, ov.PrecedingVersionUID.V.Validate(attrPath).Errs...)
	}

	// Validate other_input_version_uids
	if ov.OtherInputVersionUIDs.E {
		attrPath = path + ".other_input_version_uids"
		for i := range ov.OtherInputVersionUIDs.V {
			itemPath := fmt.Sprintf("%s[%d]", attrPath, i)
			validateErr.Errs = append(validateErr.Errs, ov.OtherInputVersionUIDs.V[i].Validate(itemPath).Errs...)
		}
	}

	// Validate lifecycle_state
	attrPath = path + ".lifecycle_state"
	validateErr.Errs = append(validateErr.Errs, ov.LifecycleState.Validate(attrPath).Errs...)

	// Validate attestations
	if ov.Attestations.E {
		attrPath = path + ".attestations"
		for i := range ov.Attestations.V {
			itemPath := fmt.Sprintf("%s[%d]", attrPath, i)
			validateErr.Errs = append(validateErr.Errs, ov.Attestations.V[i].Validate(itemPath).Errs...)
		}
	}

	// Validate data
	attrPath = path + ".data"
	validateErr.Errs = append(validateErr.Errs, ov.Data.Validate(attrPath).Errs...)

	return validateErr
}

type OriginalVersionDataKind int

const (
	OriginalVersionDataKind_Unknown OriginalVersionDataKind = iota
	OriginalVersionDataKind_EHR_STATUS
	OriginalVersionDataKind_EHR_ACCESS
	OriginalVersionDataKind_COMPOSITION
	OriginalVersionDataKind_FOLDER
	OriginalVersionDataKind_ROLE
	OriginalVersionDataKind_PERSON
	OriginalVersionDataKind_AGENT
	OriginalVersionDataKind_GROUP
	OriginalVersionDataKind_ORGANISATION
)

type OriginalVersionDataUnion struct {
	Kind  OriginalVersionDataKind
	Value any
}

func (ovd *OriginalVersionDataUnion) SetModelName() {
	switch ovd.Kind {
	case OriginalVersionDataKind_EHR_STATUS:
		ovd.Value.(*EHR_STATUS).SetModelName()
	case OriginalVersionDataKind_EHR_ACCESS:
		ovd.Value.(*EHR_ACCESS).SetModelName()
	case OriginalVersionDataKind_COMPOSITION:
		ovd.Value.(*COMPOSITION).SetModelName()
	case OriginalVersionDataKind_FOLDER:
		ovd.Value.(*FOLDER).SetModelName()
	case OriginalVersionDataKind_ROLE:
		ovd.Value.(*ROLE).SetModelName()
	case OriginalVersionDataKind_PERSON:
		ovd.Value.(*PERSON).SetModelName()
	case OriginalVersionDataKind_AGENT:
		ovd.Value.(*AGENT).SetModelName()
	case OriginalVersionDataKind_GROUP:
		ovd.Value.(*GROUP).SetModelName()
	case OriginalVersionDataKind_ORGANISATION:
		ovd.Value.(*ORGANISATION).SetModelName()
	}
}

func (ovd *OriginalVersionDataUnion) Validate(path string) util.ValidateError {
	switch ovd.Kind {
	case OriginalVersionDataKind_EHR_STATUS:
		return ovd.Value.(*EHR_STATUS).Validate(path)
	case OriginalVersionDataKind_EHR_ACCESS:
		return ovd.Value.(*EHR_ACCESS).Validate(path)
	case OriginalVersionDataKind_COMPOSITION:
		return ovd.Value.(*COMPOSITION).Validate(path)
	case OriginalVersionDataKind_FOLDER:
		return ovd.Value.(*FOLDER).Validate(path)
	case OriginalVersionDataKind_ROLE:
		return ovd.Value.(*ROLE).Validate(path)
	case OriginalVersionDataKind_PERSON:
		return ovd.Value.(*PERSON).Validate(path)
	case OriginalVersionDataKind_AGENT:
		return ovd.Value.(*AGENT).Validate(path)
	case OriginalVersionDataKind_GROUP:
		return ovd.Value.(*GROUP).Validate(path)
	case OriginalVersionDataKind_ORGANISATION:
		return ovd.Value.(*ORGANISATION).Validate(path)
	default:
		return util.ValidateError{
			Errs: []util.ValidationError{
				{
					Model:          ORIGINAL_VERSION_TYPE,
					Path:           path + ".data",
					Message:        "value is not known ORIGINAL_VERSION data subtype",
					Recommendation: "Ensure value is properly set",
				},
			},
		}
	}
}

func (ovd OriginalVersionDataUnion) Marshal() ([]byte, error) {
	return json.Marshal(ovd.Value)
}

func (ovd *OriginalVersionDataUnion) UnmarshalJSON(data []byte) error {
	t := util.UnsafeTypeFieldExtraction(data)
	switch t {
	case EHR_STATUS_TYPE:
		ovd.Kind = OriginalVersionDataKind_EHR_STATUS
		ovd.Value = &EHR_STATUS{}
		return json.Unmarshal(data, ovd.Value)
	case EHR_ACCESS_TYPE:
		ovd.Kind = OriginalVersionDataKind_EHR_ACCESS
		ovd.Value = &EHR_ACCESS{}
		return json.Unmarshal(data, ovd.Value)
	case COMPOSITION_TYPE:
		ovd.Kind = OriginalVersionDataKind_COMPOSITION
		ovd.Value = &COMPOSITION{}
		return json.Unmarshal(data, ovd.Value)
	case FOLDER_TYPE:
		ovd.Kind = OriginalVersionDataKind_FOLDER
		ovd.Value = &FOLDER{}
		return json.Unmarshal(data, ovd.Value)
	case ROLE_TYPE:
		ovd.Kind = OriginalVersionDataKind_ROLE
		ovd.Value = &ROLE{}
		return json.Unmarshal(data, ovd.Value)
	case PERSON_TYPE:
		ovd.Kind = OriginalVersionDataKind_PERSON
		ovd.Value = &PERSON{}
		return json.Unmarshal(data, ovd.Value)
	case AGENT_TYPE:
		ovd.Kind = OriginalVersionDataKind_AGENT
		ovd.Value = &AGENT{}
		return json.Unmarshal(data, ovd.Value)
	case GROUP_TYPE:
		ovd.Kind = OriginalVersionDataKind_GROUP
		ovd.Value = &GROUP{}
		return json.Unmarshal(data, ovd.Value)
	case ORGANISATION_TYPE:
		ovd.Kind = OriginalVersionDataKind_ORGANISATION
		ovd.Value = &ORGANISATION{}
		return json.Unmarshal(data, ovd.Value)
	default:
		ovd.Kind = OriginalVersionDataKind_Unknown
		return nil
	}
}

func (o *OriginalVersionDataUnion) EHRStatus() *EHR_STATUS {
	if o.Kind == OriginalVersionDataKind_EHR_STATUS {
		return o.Value.(*EHR_STATUS)
	}
	return nil
}

func (o *OriginalVersionDataUnion) EHRAccess() *EHR_ACCESS {
	if o.Kind == OriginalVersionDataKind_EHR_ACCESS {
		return o.Value.(*EHR_ACCESS)
	}
	return nil
}

func (o *OriginalVersionDataUnion) Composition() *COMPOSITION {
	if o.Kind == OriginalVersionDataKind_COMPOSITION {
		return o.Value.(*COMPOSITION)
	}
	return nil
}

func (o *OriginalVersionDataUnion) Folder() *FOLDER {
	if o.Kind == OriginalVersionDataKind_FOLDER {
		return o.Value.(*FOLDER)
	}
	return nil
}

func (o *OriginalVersionDataUnion) Role() *ROLE {
	if o.Kind == OriginalVersionDataKind_ROLE {
		return o.Value.(*ROLE)
	}
	return nil
}

func (o *OriginalVersionDataUnion) Person() *PERSON {
	if o.Kind == OriginalVersionDataKind_PERSON {
		return o.Value.(*PERSON)
	}
	return nil
}

func (o *OriginalVersionDataUnion) Agent() *AGENT {
	if o.Kind == OriginalVersionDataKind_AGENT {
		return o.Value.(*AGENT)
	}
	return nil
}

func (o *OriginalVersionDataUnion) Group() *GROUP {
	if o.Kind == OriginalVersionDataKind_GROUP {
		return o.Value.(*GROUP)
	}
	return nil
}

func (o *OriginalVersionDataUnion) Organisation() *ORGANISATION {
	if o.Kind == OriginalVersionDataKind_ORGANISATION {
		return o.Value.(*ORGANISATION)
	}
	return nil
}

func OriginalVersionDataFromEHRStatus(ehrStatus EHR_STATUS) OriginalVersionDataUnion {
	ehrStatus.Type_ = utils.Some(EHR_STATUS_TYPE)
	return OriginalVersionDataUnion{
		Kind:  OriginalVersionDataKind_EHR_STATUS,
		Value: &ehrStatus,
	}
}

func OriginalVersionDataFromEHRAccess(ehrAccess EHR_ACCESS) OriginalVersionDataUnion {
	ehrAccess.Type_ = utils.Some(EHR_ACCESS_TYPE)
	return OriginalVersionDataUnion{
		Kind:  OriginalVersionDataKind_EHR_ACCESS,
		Value: &ehrAccess,
	}
}

func OriginalVersionDataFromComposition(composition COMPOSITION) OriginalVersionDataUnion {
	composition.Type_ = utils.Some(COMPOSITION_TYPE)
	return OriginalVersionDataUnion{
		Kind:  OriginalVersionDataKind_COMPOSITION,
		Value: &composition,
	}
}

func OriginalVersionDataFromFolder(folder FOLDER) OriginalVersionDataUnion {
	folder.Type_ = utils.Some(FOLDER_TYPE)
	return OriginalVersionDataUnion{
		Kind:  OriginalVersionDataKind_FOLDER,
		Value: &folder,
	}
}

func OriginalVersionDataFromRole(role ROLE) OriginalVersionDataUnion {
	role.Type_ = utils.Some(ROLE_TYPE)
	return OriginalVersionDataUnion{
		Kind:  OriginalVersionDataKind_ROLE,
		Value: &role,
	}
}

func OriginalVersionDataFromPerson(person PERSON) OriginalVersionDataUnion {
	person.Type_ = utils.Some(PERSON_TYPE)
	return OriginalVersionDataUnion{
		Kind:  OriginalVersionDataKind_PERSON,
		Value: &person,
	}
}

func OriginalVersionDataFromAgent(agent AGENT) OriginalVersionDataUnion {
	agent.Type_ = utils.Some(AGENT_TYPE)
	return OriginalVersionDataUnion{
		Kind:  OriginalVersionDataKind_AGENT,
		Value: &agent,
	}
}

func OriginalVersionDataFromGroup(group GROUP) OriginalVersionDataUnion {
	group.Type_ = utils.Some(GROUP_TYPE)
	return OriginalVersionDataUnion{
		Kind:  OriginalVersionDataKind_GROUP,
		Value: &group,
	}
}

func OriginalVersionDataFromOrganisation(organisation ORGANISATION) OriginalVersionDataUnion {
	organisation.Type_ = utils.Some(ORGANISATION_TYPE)
	return OriginalVersionDataUnion{
		Kind:  OriginalVersionDataKind_ORGANISATION,
		Value: &organisation,
	}
}
