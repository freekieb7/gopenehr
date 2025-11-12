package terminology

// Instruction Transitions vocabulary codes
// This vocabulary codifies the names of the transitions in the standard Instruction state machine.
// Used in: ISM_TRANSITION.transition

const (
	INSTRUCTION_TRANSITION_CODE_INITIATE         string = "535" // initiate
	INSTRUCTION_TRANSITION_CODE_PLAN_STEP        string = "536" // plan step
	INSTRUCTION_TRANSITION_CODE_POSTPONE         string = "537" // postpone
	INSTRUCTION_TRANSITION_CODE_RESTORE          string = "538" // restore
	INSTRUCTION_TRANSITION_CODE_CANCEL           string = "166" // cancel
	INSTRUCTION_TRANSITION_CODE_POSTPONED_STEP   string = "542" // postponed step
	INSTRUCTION_TRANSITION_CODE_SCHEDULE         string = "539" // schedule
	INSTRUCTION_TRANSITION_CODE_SCHEDULED_STEP   string = "534" // scheduled step
	INSTRUCTION_TRANSITION_CODE_START            string = "540" // start
	INSTRUCTION_TRANSITION_CODE_DO               string = "541" // do
	INSTRUCTION_TRANSITION_CODE_ACTIVE_STEP      string = "543" // active step
	INSTRUCTION_TRANSITION_CODE_SUSPEND          string = "544" // suspend
	INSTRUCTION_TRANSITION_CODE_SUSPENDED_STEP   string = "545" // suspended step
	INSTRUCTION_TRANSITION_CODE_RESUME           string = "546" // resume
	INSTRUCTION_TRANSITION_CODE_ABORT            string = "547" // abort
	INSTRUCTION_TRANSITION_CODE_FINISH           string = "548" // finish
	INSTRUCTION_TRANSITION_CODE_TIME_OUT         string = "549" // time out
	INSTRUCTION_TRANSITION_CODE_NOTIFY_ABORTED   string = "550" // notify aborted
	INSTRUCTION_TRANSITION_CODE_NOTIFY_COMPLETED string = "551" // notify completed
	INSTRUCTION_TRANSITION_CODE_NOTIFY_CANCELLED string = "552" // notify cancelled
)

// InstructionTransitionNames maps instruction transition codes to their display names
var InstructionTransitionNames = map[string]string{
	INSTRUCTION_TRANSITION_CODE_INITIATE:         "initiate",
	INSTRUCTION_TRANSITION_CODE_PLAN_STEP:        "plan step",
	INSTRUCTION_TRANSITION_CODE_POSTPONE:         "postpone",
	INSTRUCTION_TRANSITION_CODE_RESTORE:          "restore",
	INSTRUCTION_TRANSITION_CODE_CANCEL:           "cancel",
	INSTRUCTION_TRANSITION_CODE_POSTPONED_STEP:   "postponed step",
	INSTRUCTION_TRANSITION_CODE_SCHEDULE:         "schedule",
	INSTRUCTION_TRANSITION_CODE_SCHEDULED_STEP:   "scheduled step",
	INSTRUCTION_TRANSITION_CODE_START:            "start",
	INSTRUCTION_TRANSITION_CODE_DO:               "do",
	INSTRUCTION_TRANSITION_CODE_ACTIVE_STEP:      "active step",
	INSTRUCTION_TRANSITION_CODE_SUSPEND:          "suspend",
	INSTRUCTION_TRANSITION_CODE_SUSPENDED_STEP:   "suspended step",
	INSTRUCTION_TRANSITION_CODE_RESUME:           "resume",
	INSTRUCTION_TRANSITION_CODE_ABORT:            "abort",
	INSTRUCTION_TRANSITION_CODE_FINISH:           "finish",
	INSTRUCTION_TRANSITION_CODE_TIME_OUT:         "time out",
	INSTRUCTION_TRANSITION_CODE_NOTIFY_ABORTED:   "notify aborted",
	INSTRUCTION_TRANSITION_CODE_NOTIFY_COMPLETED: "notify completed",
	INSTRUCTION_TRANSITION_CODE_NOTIFY_CANCELLED: "notify cancelled",
}

// IsValidInstructionTransitionCode checks if the provided code is a valid instruction transition
func IsValidInstructionTransitionCode(code string) bool {
	_, exists := InstructionTransitionNames[code]
	return exists
}

// GetInstructionTransitionName returns the display name for an instruction transition code
func GetInstructionTransitionName(code string) string {
	if name, exists := InstructionTransitionNames[code]; exists {
		return name
	}
	return ""
}
