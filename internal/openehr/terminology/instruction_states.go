package terminology

// Instruction States vocabulary codes
// This vocabulary codifies the names of the states in the standard Instruction state machine.
// Used in: ISM_TRANSITION.current_state

const (
	INSTRUCTION_STATE_CODE_INITIAL   string = "524" // initial
	INSTRUCTION_STATE_CODE_PLANNED   string = "526" // planned
	INSTRUCTION_STATE_CODE_POSTPONED string = "527" // postponed
	INSTRUCTION_STATE_CODE_CANCELLED string = "528" // cancelled
	INSTRUCTION_STATE_CODE_SCHEDULED string = "529" // scheduled
	INSTRUCTION_STATE_CODE_ACTIVE    string = "245" // active
	INSTRUCTION_STATE_CODE_SUSPENDED string = "530" // suspended
	INSTRUCTION_STATE_CODE_ABORTED   string = "531" // aborted
	INSTRUCTION_STATE_CODE_COMPLETED string = "532" // completed
	INSTRUCTION_STATE_CODE_EXPIRED   string = "533" // expired
)

// InstructionStateNames maps instruction state codes to their display names
var InstructionStateNames = map[string]string{
	INSTRUCTION_STATE_CODE_INITIAL:   "initial",
	INSTRUCTION_STATE_CODE_PLANNED:   "planned",
	INSTRUCTION_STATE_CODE_POSTPONED: "postponed",
	INSTRUCTION_STATE_CODE_CANCELLED: "cancelled",
	INSTRUCTION_STATE_CODE_SCHEDULED: "scheduled",
	INSTRUCTION_STATE_CODE_ACTIVE:    "active",
	INSTRUCTION_STATE_CODE_SUSPENDED: "suspended",
	INSTRUCTION_STATE_CODE_ABORTED:   "aborted",
	INSTRUCTION_STATE_CODE_COMPLETED: "completed",
	INSTRUCTION_STATE_CODE_EXPIRED:   "expired",
}

// IsValidInstructionStateCode checks if the provided code is a valid instruction state
func IsValidInstructionStateCode(code string) bool {
	_, exists := InstructionStateNames[code]
	return exists
}

// GetInstructionStateName returns the display name for an instruction state code
func GetInstructionStateName(code string) string {
	if name, exists := InstructionStateNames[code]; exists {
		return name
	}
	return ""
}
