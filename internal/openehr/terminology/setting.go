package terminology

// Setting vocabulary codes
// This vocabulary codifies broad types of settings in which clinical care is delivered.
// Used in: EVENT_CONTEXT.setting
const (
	SETTING_TERMINOLOGY_ID_OPENEHR string = "openehr"
)

var SettingTerminologyIds = map[string]string{
	SETTING_TERMINOLOGY_ID_OPENEHR: "openEHR",
}

// IsValidSettingTerminologyID checks if the provided terminology ID is valid for settings
func IsValidSettingTerminologyID(terminologyID string) bool {
	_, exists := SettingTerminologyIds[terminologyID]
	return exists
}

const (
	SETTING_CODE_HOME                         string = "225" // home
	SETTING_CODE_EMERGENCY_CARE               string = "227" // emergency care
	SETTING_CODE_PRIMARY_MEDICAL_CARE         string = "228" // primary medical care
	SETTING_CODE_PRIMARY_NURSING_CARE         string = "229" // primary nursing care
	SETTING_CODE_PRIMARY_ALLIED_HEALTH_CARE   string = "230" // primary allied health care
	SETTING_CODE_MIDWIFERY_CARE               string = "231" // midwifery care
	SETTING_CODE_SECONDARY_MEDICAL_CARE       string = "232" // secondary medical care
	SETTING_CODE_SECONDARY_NURSING_CARE       string = "233" // secondary nursing care
	SETTING_CODE_SECONDARY_ALLIED_HEALTH_CARE string = "234" // secondary allied health care
	SETTING_CODE_COMPLEMENTARY_HEALTH_CARE    string = "235" // complementary health care
	SETTING_CODE_DENTAL_CARE                  string = "236" // dental care
	SETTING_CODE_NURSING_HOME_CARE            string = "237" // nursing home care
	SETTING_CODE_MENTAL_HEALTHCARE            string = "802" // mental healthcare
	SETTING_CODE_OTHER_CARE                   string = "238" // other care
)

// SettingNames maps setting codes to their display names
var SettingNames = map[string]string{
	SETTING_CODE_HOME:                         "home",
	SETTING_CODE_EMERGENCY_CARE:               "emergency care",
	SETTING_CODE_PRIMARY_MEDICAL_CARE:         "primary medical care",
	SETTING_CODE_PRIMARY_NURSING_CARE:         "primary nursing care",
	SETTING_CODE_PRIMARY_ALLIED_HEALTH_CARE:   "primary allied health care",
	SETTING_CODE_MIDWIFERY_CARE:               "midwifery care",
	SETTING_CODE_SECONDARY_MEDICAL_CARE:       "secondary medical care",
	SETTING_CODE_SECONDARY_NURSING_CARE:       "secondary nursing care",
	SETTING_CODE_SECONDARY_ALLIED_HEALTH_CARE: "secondary allied health care",
	SETTING_CODE_COMPLEMENTARY_HEALTH_CARE:    "complementary health care",
	SETTING_CODE_DENTAL_CARE:                  "dental care",
	SETTING_CODE_NURSING_HOME_CARE:            "nursing home care",
	SETTING_CODE_MENTAL_HEALTHCARE:            "mental healthcare",
	SETTING_CODE_OTHER_CARE:                   "other care",
}

// IsValidSetting checks if the provided code is a valid setting
func IsValidSettingCode(code string) bool {
	_, exists := SettingNames[code]
	return exists
}

// GetSettingName returns the display name for a setting code
func GetSettingName(code string) string {
	if name, exists := SettingNames[code]; exists {
		return name
	}
	return ""
}
