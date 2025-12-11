package rm

import (
	"fmt"

	"github.com/bytedance/sonic"
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
	ORIGINAL_VERSION_data_kind_unknown OriginalVersionDataKind = iota
	ORIGINAL_VERSION_data_kind_EHR_STATUS
	ORIGINAL_VERSION_data_kind_EHR_ACCESS
	ORIGINAL_VERSION_data_kind_COMPOSITION
	ORIGINAL_VERSION_data_kind_FOLDER
	ORIGINAL_VERSION_data_kind_ROLE
	ORIGINAL_VERSION_data_kind_PERSON
	ORIGINAL_VERSION_data_kind_AGENT
	ORIGINAL_VERSION_data_kind_GROUP
	ORIGINAL_VERSION_data_kind_ORGANISATION
)

type OriginalVersionDataUnion struct {
	Kind  OriginalVersionDataKind
	Value any
}

func (ovd *OriginalVersionDataUnion) SetModelName() {
	switch ovd.Kind {
	case ORIGINAL_VERSION_data_kind_EHR_STATUS:
		ovd.Value.(*EHR_STATUS).SetModelName()
	case ORIGINAL_VERSION_data_kind_EHR_ACCESS:
		ovd.Value.(*EHR_ACCESS).SetModelName()
	case ORIGINAL_VERSION_data_kind_COMPOSITION:
		ovd.Value.(*COMPOSITION).SetModelName()
	case ORIGINAL_VERSION_data_kind_FOLDER:
		ovd.Value.(*FOLDER).SetModelName()
	case ORIGINAL_VERSION_data_kind_ROLE:
		ovd.Value.(*ROLE).SetModelName()
	case ORIGINAL_VERSION_data_kind_PERSON:
		ovd.Value.(*PERSON).SetModelName()
	case ORIGINAL_VERSION_data_kind_AGENT:
		ovd.Value.(*AGENT).SetModelName()
	case ORIGINAL_VERSION_data_kind_GROUP:
		ovd.Value.(*GROUP).SetModelName()
	case ORIGINAL_VERSION_data_kind_ORGANISATION:
		ovd.Value.(*ORGANISATION).SetModelName()
	}
}

func (ovd *OriginalVersionDataUnion) Validate(path string) util.ValidateError {
	switch ovd.Kind {
	case ORIGINAL_VERSION_data_kind_EHR_STATUS:
		return ovd.Value.(*EHR_STATUS).Validate(path)
	case ORIGINAL_VERSION_data_kind_EHR_ACCESS:
		return ovd.Value.(*EHR_ACCESS).Validate(path)
	case ORIGINAL_VERSION_data_kind_COMPOSITION:
		return ovd.Value.(*COMPOSITION).Validate(path)
	case ORIGINAL_VERSION_data_kind_FOLDER:
		return ovd.Value.(*FOLDER).Validate(path)
	case ORIGINAL_VERSION_data_kind_ROLE:
		return ovd.Value.(*ROLE).Validate(path)
	case ORIGINAL_VERSION_data_kind_PERSON:
		return ovd.Value.(*PERSON).Validate(path)
	case ORIGINAL_VERSION_data_kind_AGENT:
		return ovd.Value.(*AGENT).Validate(path)
	case ORIGINAL_VERSION_data_kind_GROUP:
		return ovd.Value.(*GROUP).Validate(path)
	case ORIGINAL_VERSION_data_kind_ORGANISATION:
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

func (ovd OriginalVersionDataUnion) MarshalJSON() ([]byte, error) {
	return sonic.Marshal(ovd.Value)
}

func (ovd *OriginalVersionDataUnion) UnmarshalJSON(data []byte) error {
	t := util.UnsafeTypeFieldExtraction(data)
	switch t {
	case EHR_STATUS_TYPE:
		ovd.Kind = ORIGINAL_VERSION_data_kind_EHR_STATUS
		ovd.Value = &EHR_STATUS{}
	case EHR_ACCESS_TYPE:
		ovd.Kind = ORIGINAL_VERSION_data_kind_EHR_ACCESS
		ovd.Value = &EHR_ACCESS{}
	case COMPOSITION_TYPE:
		ovd.Kind = ORIGINAL_VERSION_data_kind_COMPOSITION
		ovd.Value = &COMPOSITION{}
	case FOLDER_TYPE:
		ovd.Kind = ORIGINAL_VERSION_data_kind_FOLDER
		ovd.Value = &FOLDER{}
	case ROLE_TYPE:
		ovd.Kind = ORIGINAL_VERSION_data_kind_ROLE
		ovd.Value = &ROLE{}
	case PERSON_TYPE:
		ovd.Kind = ORIGINAL_VERSION_data_kind_PERSON
		ovd.Value = &PERSON{}
	case AGENT_TYPE:
		ovd.Kind = ORIGINAL_VERSION_data_kind_AGENT
		ovd.Value = &AGENT{}
	case GROUP_TYPE:
		ovd.Kind = ORIGINAL_VERSION_data_kind_GROUP
		ovd.Value = &GROUP{}
	case ORGANISATION_TYPE:
		ovd.Kind = ORIGINAL_VERSION_data_kind_ORGANISATION
		ovd.Value = &ORGANISATION{}
	default:
		ovd.Kind = ORIGINAL_VERSION_data_kind_unknown
		return nil
	}

	return sonic.Unmarshal(data, ovd.Value)
}

func (o *OriginalVersionDataUnion) EHR_STATUS() *EHR_STATUS {
	if o.Kind == ORIGINAL_VERSION_data_kind_EHR_STATUS {
		return o.Value.(*EHR_STATUS)
	}
	return nil
}

func (o *OriginalVersionDataUnion) EHR_ACCESS() *EHR_ACCESS {
	if o.Kind == ORIGINAL_VERSION_data_kind_EHR_ACCESS {
		return o.Value.(*EHR_ACCESS)
	}
	return nil
}

func (o *OriginalVersionDataUnion) COMPOSITION() *COMPOSITION {
	if o.Kind == ORIGINAL_VERSION_data_kind_COMPOSITION {
		return o.Value.(*COMPOSITION)
	}
	return nil
}

func (o *OriginalVersionDataUnion) FOLDER() *FOLDER {
	if o.Kind == ORIGINAL_VERSION_data_kind_FOLDER {
		return o.Value.(*FOLDER)
	}
	return nil
}

func (o *OriginalVersionDataUnion) ROLE() *ROLE {
	if o.Kind == ORIGINAL_VERSION_data_kind_ROLE {
		return o.Value.(*ROLE)
	}
	return nil
}

func (o *OriginalVersionDataUnion) PERSON() *PERSON {
	if o.Kind == ORIGINAL_VERSION_data_kind_PERSON {
		return o.Value.(*PERSON)
	}
	return nil
}

func (o *OriginalVersionDataUnion) AGENT() *AGENT {
	if o.Kind == ORIGINAL_VERSION_data_kind_AGENT {
		return o.Value.(*AGENT)
	}
	return nil
}

func (o *OriginalVersionDataUnion) GROUP() *GROUP {
	if o.Kind == ORIGINAL_VERSION_data_kind_GROUP {
		return o.Value.(*GROUP)
	}
	return nil
}

func (o *OriginalVersionDataUnion) ORGANISATION() *ORGANISATION {
	if o.Kind == ORIGINAL_VERSION_data_kind_ORGANISATION {
		return o.Value.(*ORGANISATION)
	}
	return nil
}

func ORIGINAL_VERSION_DATA_from_EHR_STATUS(ehrStatus EHR_STATUS) OriginalVersionDataUnion {
	ehrStatus.Type_ = utils.Some(EHR_STATUS_TYPE)
	return OriginalVersionDataUnion{
		Kind:  ORIGINAL_VERSION_data_kind_EHR_STATUS,
		Value: &ehrStatus,
	}
}

func ORIGINAL_VERSION_DATA_from_EHR_ACCESS(ehrAccess EHR_ACCESS) OriginalVersionDataUnion {
	ehrAccess.Type_ = utils.Some(EHR_ACCESS_TYPE)
	return OriginalVersionDataUnion{
		Kind:  ORIGINAL_VERSION_data_kind_EHR_ACCESS,
		Value: &ehrAccess,
	}
}

func ORIGINAL_VERSION_DATA_from_COMPOSITION(composition COMPOSITION) OriginalVersionDataUnion {
	composition.Type_ = utils.Some(COMPOSITION_TYPE)
	return OriginalVersionDataUnion{
		Kind:  ORIGINAL_VERSION_data_kind_COMPOSITION,
		Value: &composition,
	}
}

func ORIGINAL_VERSION_DATA_from_FOLDER(folder FOLDER) OriginalVersionDataUnion {
	folder.Type_ = utils.Some(FOLDER_TYPE)
	return OriginalVersionDataUnion{
		Kind:  ORIGINAL_VERSION_data_kind_FOLDER,
		Value: &folder,
	}
}

func ORIGINAL_VERSION_DATA_from_ROLE(role ROLE) OriginalVersionDataUnion {
	role.Type_ = utils.Some(ROLE_TYPE)
	return OriginalVersionDataUnion{
		Kind:  ORIGINAL_VERSION_data_kind_ROLE,
		Value: &role,
	}
}

func ORIGINAL_VERSION_DATA_from_PERSON(person PERSON) OriginalVersionDataUnion {
	person.Type_ = utils.Some(PERSON_TYPE)
	return OriginalVersionDataUnion{
		Kind:  ORIGINAL_VERSION_data_kind_PERSON,
		Value: &person,
	}
}

func ORIGINAL_VERSION_DATA_from_AGENT(agent AGENT) OriginalVersionDataUnion {
	agent.Type_ = utils.Some(AGENT_TYPE)
	return OriginalVersionDataUnion{
		Kind:  ORIGINAL_VERSION_data_kind_AGENT,
		Value: &agent,
	}
}

func ORIGINAL_VERSION_DATA_from_GROUP(group GROUP) OriginalVersionDataUnion {
	group.Type_ = utils.Some(GROUP_TYPE)
	return OriginalVersionDataUnion{
		Kind:  ORIGINAL_VERSION_data_kind_GROUP,
		Value: &group,
	}
}

func ORIGINAL_VERSION_DATA_from_ORGANISATION(organisation ORGANISATION) OriginalVersionDataUnion {
	organisation.Type_ = utils.Some(ORGANISATION_TYPE)
	return OriginalVersionDataUnion{
		Kind:  ORIGINAL_VERSION_data_kind_ORGANISATION,
		Value: &organisation,
	}
}
