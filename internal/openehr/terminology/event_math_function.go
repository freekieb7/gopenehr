package terminology

// Event Math Function vocabulary codes
// This vocabulary codifies mathematical functions applied to non-instantaneous time series events.
// Used in: INTERVAL_EVENT.math_function

const (
	EVENT_MATH_FUNCTION_CODE_MINIMUM   string = "145" // minimum
	EVENT_MATH_FUNCTION_CODE_MAXIMUM   string = "144" // maximum
	EVENT_MATH_FUNCTION_CODE_MODE      string = "267" // mode
	EVENT_MATH_FUNCTION_CODE_MEDIAN    string = "268" // median
	EVENT_MATH_FUNCTION_CODE_MEAN      string = "146" // mean
	EVENT_MATH_FUNCTION_CODE_CHANGE    string = "147" // change
	EVENT_MATH_FUNCTION_CODE_TOTAL     string = "148" // total
	EVENT_MATH_FUNCTION_CODE_VARIATION string = "149" // variation
	EVENT_MATH_FUNCTION_CODE_DECREASE  string = "521" // decrease
	EVENT_MATH_FUNCTION_CODE_INCREASE  string = "522" // increase
	EVENT_MATH_FUNCTION_CODE_ACTUAL    string = "640" // actual
)

// EventMathFunctionNames maps event math function codes to their display names
var EventMathFunctionNames = map[string]string{
	EVENT_MATH_FUNCTION_CODE_MINIMUM:   "minimum",
	EVENT_MATH_FUNCTION_CODE_MAXIMUM:   "maximum",
	EVENT_MATH_FUNCTION_CODE_MODE:      "mode",
	EVENT_MATH_FUNCTION_CODE_MEDIAN:    "median",
	EVENT_MATH_FUNCTION_CODE_MEAN:      "mean",
	EVENT_MATH_FUNCTION_CODE_CHANGE:    "change",
	EVENT_MATH_FUNCTION_CODE_TOTAL:     "total",
	EVENT_MATH_FUNCTION_CODE_VARIATION: "variation",
	EVENT_MATH_FUNCTION_CODE_DECREASE:  "decrease",
	EVENT_MATH_FUNCTION_CODE_INCREASE:  "increase",
	EVENT_MATH_FUNCTION_CODE_ACTUAL:    "actual",
}

// IsValidEventMathFunctionCode checks if the provided code is a valid event math function
func IsValidEventMathFunctionCode(code string) bool {
	_, exists := EventMathFunctionNames[code]
	return exists
}

// GetEventMathFunctionName returns the display name for an event math function code
func GetEventMathFunctionName(code string) string {
	if name, exists := EventMathFunctionNames[code]; exists {
		return name
	}
	return ""
}
