package terminology

// Participation Mode vocabulary codes
// This vocabulary codifies modes of participation of parties in an interaction.
// External reference: HL7 ParticipationMode domain
// Used in: PARTICIPATION.mode

const (
	PARTICIPATION_MODE_CODE_NOT_SPECIFIED            string = "193" // not specified
	PARTICIPATION_MODE_CODE_FACE_TO_FACE             string = "216" // face-to-face communication
	PARTICIPATION_MODE_CODE_INTERPRETED_FACE_TO_FACE string = "223" // interpreted face-to-face communication
	PARTICIPATION_MODE_CODE_SIGNING_FACE_TO_FACE     string = "217" // signing (face-to-face)
	PARTICIPATION_MODE_CODE_LIVE_AUDIOVISUAL         string = "195" // live audiovisual; videoconference; videophone
	PARTICIPATION_MODE_CODE_VIDEOCONFERENCING        string = "198" // videoconferencing
	PARTICIPATION_MODE_CODE_VIDEOPHONE               string = "197" // videophone
	PARTICIPATION_MODE_CODE_SIGNING_OVER_VIDEO       string = "218" // signing over video
	PARTICIPATION_MODE_CODE_INTERPRETED_VIDEO        string = "224" // interpreted video communication
	PARTICIPATION_MODE_CODE_ASYNCHRONOUS_AUDIOVISUAL string = "194" // asynchronous audiovisual; recorded video
	PARTICIPATION_MODE_CODE_RECORDED_VIDEO           string = "196" // recorded video
	PARTICIPATION_MODE_CODE_LIVE_AUDIO_ONLY          string = "202" // live audio-only; telephone; internet phone; teleconference
	PARTICIPATION_MODE_CODE_TELEPHONE                string = "204" // telephone
	PARTICIPATION_MODE_CODE_TELECONFERENCE           string = "203" // teleconference
	PARTICIPATION_MODE_CODE_INTERNET_TELEPHONE       string = "205" // internet telephone
	PARTICIPATION_MODE_CODE_INTERPRETED_AUDIO_ONLY   string = "222" // interpreted audio-only
	PARTICIPATION_MODE_CODE_ASYNCHRONOUS_AUDIO_ONLY  string = "199" // asynchronous audio-only; dictated; voice mail
	PARTICIPATION_MODE_CODE_DICTATED                 string = "200" // dictated
	PARTICIPATION_MODE_CODE_VOICE_MAIL               string = "201" // voice-mail
	PARTICIPATION_MODE_CODE_LIVE_TEXT_ONLY           string = "212" // live text-only; internet chat; SMS chat; interactive written note
	PARTICIPATION_MODE_CODE_INTERNET_CHAT            string = "213" // internet chat
	PARTICIPATION_MODE_CODE_SMS_CHAT                 string = "214" // SMS chat
	PARTICIPATION_MODE_CODE_INTERACTIVE_WRITTEN_NOTE string = "215" // interactive written note
	PARTICIPATION_MODE_CODE_ASYNCHRONOUS_TEXT        string = "206" // asynchronous text; email; fax; letter; handwritten note; SMS message
	PARTICIPATION_MODE_CODE_HANDWRITTEN_NOTE         string = "211" // handwritten note
	PARTICIPATION_MODE_CODE_PRINTED_TYPED_LETTER     string = "210" // printed/typed letter
	PARTICIPATION_MODE_CODE_EMAIL                    string = "207" // email
	PARTICIPATION_MODE_CODE_FACSIMILE_TELEFAX        string = "208" // facsimile/telefax
	PARTICIPATION_MODE_CODE_TRANSLATED_TEXT          string = "221" // translated text
	PARTICIPATION_MODE_CODE_SMS_MESSAGE              string = "209" // SMS message
	PARTICIPATION_MODE_CODE_PHYSICALLY_PRESENT       string = "219" // physically present
	PARTICIPATION_MODE_CODE_PHYSICALLY_REMOTE        string = "220" // physically remote
)

// ParticipationModeNames maps participation mode codes to their display names
var ParticipationModeNames = map[string]string{
	PARTICIPATION_MODE_CODE_NOT_SPECIFIED:            "not specified",
	PARTICIPATION_MODE_CODE_FACE_TO_FACE:             "face-to-face communication",
	PARTICIPATION_MODE_CODE_INTERPRETED_FACE_TO_FACE: "interpreted face-to-face communication",
	PARTICIPATION_MODE_CODE_SIGNING_FACE_TO_FACE:     "signing (face-to-face)",
	PARTICIPATION_MODE_CODE_LIVE_AUDIOVISUAL:         "live audiovisual; videoconference; videophone",
	PARTICIPATION_MODE_CODE_VIDEOCONFERENCING:        "videoconferencing",
	PARTICIPATION_MODE_CODE_VIDEOPHONE:               "videophone",
	PARTICIPATION_MODE_CODE_SIGNING_OVER_VIDEO:       "signing over video",
	PARTICIPATION_MODE_CODE_INTERPRETED_VIDEO:        "interpreted video communication",
	PARTICIPATION_MODE_CODE_ASYNCHRONOUS_AUDIOVISUAL: "asynchronous audiovisual; recorded video",
	PARTICIPATION_MODE_CODE_RECORDED_VIDEO:           "recorded video",
	PARTICIPATION_MODE_CODE_LIVE_AUDIO_ONLY:          "live audio-only; telephone; internet phone; teleconference",
	PARTICIPATION_MODE_CODE_TELEPHONE:                "telephone",
	PARTICIPATION_MODE_CODE_TELECONFERENCE:           "teleconference",
	PARTICIPATION_MODE_CODE_INTERNET_TELEPHONE:       "internet telephone",
	PARTICIPATION_MODE_CODE_INTERPRETED_AUDIO_ONLY:   "interpreted audio-only",
	PARTICIPATION_MODE_CODE_ASYNCHRONOUS_AUDIO_ONLY:  "asynchronous audio-only; dictated; voice mail",
	PARTICIPATION_MODE_CODE_DICTATED:                 "dictated",
	PARTICIPATION_MODE_CODE_VOICE_MAIL:               "voice-mail",
	PARTICIPATION_MODE_CODE_LIVE_TEXT_ONLY:           "live text-only; internet chat; SMS chat; interactive written note",
	PARTICIPATION_MODE_CODE_INTERNET_CHAT:            "internet chat",
	PARTICIPATION_MODE_CODE_SMS_CHAT:                 "SMS chat",
	PARTICIPATION_MODE_CODE_INTERACTIVE_WRITTEN_NOTE: "interactive written note",
	PARTICIPATION_MODE_CODE_ASYNCHRONOUS_TEXT:        "asynchronous text; email; fax; letter; handwritten note; SMS message",
	PARTICIPATION_MODE_CODE_HANDWRITTEN_NOTE:         "handwritten note",
	PARTICIPATION_MODE_CODE_PRINTED_TYPED_LETTER:     "printed/typed letter",
	PARTICIPATION_MODE_CODE_EMAIL:                    "email",
	PARTICIPATION_MODE_CODE_FACSIMILE_TELEFAX:        "facsimile/telefax",
	PARTICIPATION_MODE_CODE_TRANSLATED_TEXT:          "translated text",
	PARTICIPATION_MODE_CODE_SMS_MESSAGE:              "SMS message",
	PARTICIPATION_MODE_CODE_PHYSICALLY_PRESENT:       "physically present",
	PARTICIPATION_MODE_CODE_PHYSICALLY_REMOTE:        "physically remote",
}

// IsValidParticipationModeCode checks if the provided code is a valid participation mode
func IsValidParticipationModeCode(code string) bool {
	_, exists := ParticipationModeNames[code]
	return exists
}

// GetParticipationModeName returns the display name for a participation mode code
func GetParticipationModeName(code string) string {
	if name, exists := ParticipationModeNames[code]; exists {
		return name
	}
	return ""
}
