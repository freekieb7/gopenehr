package terminology

// Property vocabulary constants
// This vocabulary codifies purposes for physical properties corresponding to formal unit specifications,
// and allows comparison of Quantities with different units but which measure the same property.
// Used in: DV_QUANTITY
// External reference: Regenstrief Institute - Unified Codes for Units of Measure
const (
	PROPERTY_CODE_ACCELERATION               string = "339"
	PROPERTY_CODE_ACCELERATION_ANGULAR       string = "342"
	PROPERTY_CODE_AMOUNT_EQ                  string = "381"
	PROPERTY_CODE_AMOUNT_MOLE                string = "384"
	PROPERTY_CODE_ANGLE_PLANE                string = "497"
	PROPERTY_CODE_ANGLE_SOLID                string = "500"
	PROPERTY_CODE_AREA                       string = "335"
	PROPERTY_CODE_CONCENTRATION              string = "119"
	PROPERTY_CODE_DENSITY                    string = "350"
	PROPERTY_CODE_DIFFUSION_COEFFICIENT      string = "362"
	PROPERTY_CODE_ELECTRIC_CAPACITANCE       string = "501"
	PROPERTY_CODE_ELECTRIC_CHARGE            string = "498"
	PROPERTY_CODE_ELECTRIC_CONDUCTANCE       string = "502"
	PROPERTY_CODE_ELECTRIC_CURRENT           string = "334"
	PROPERTY_CODE_ELECTRIC_FIELD_STRENGTH    string = "377"
	PROPERTY_CODE_ELECTRIC_POTENTIAL_TIME    string = "655"
	PROPERTY_CODE_ENERGY                     string = "121"
	PROPERTY_CODE_ENERGY_DENSITY             string = "366"
	PROPERTY_CODE_ENERGY_DOSE                string = "508"
	PROPERTY_CODE_ENERGY_PER_AREA            string = "365"
	PROPERTY_CODE_ENERGY_LINEAR              string = "364"
	PROPERTY_CODE_FLOW_RATE_MASS             string = "347"
	PROPERTY_CODE_FLOW_RATE_MASS_FORCE       string = "352"
	PROPERTY_CODE_FLOW_RATE_MASS_VOLUME      string = "351"
	PROPERTY_CODE_FLOW_RATE_VOLUME           string = "126"
	PROPERTY_CODE_FLUX_MASS                  string = "348"
	PROPERTY_CODE_FORCE                      string = "355"
	PROPERTY_CODE_FORCE_PER_MASS             string = "358"
	PROPERTY_CODE_FORCE_BODY                 string = "357"
	PROPERTY_CODE_FREQUENCY                  string = "382"
	PROPERTY_CODE_GLOMERULAR_FILTRATION_RATE string = "586"
	PROPERTY_CODE_HEAT_TRANSFER_COEFFICIENT  string = "373"
	PROPERTY_CODE_ILLUMINANCE                string = "505"
	PROPERTY_CODE_INDUCTANCE                 string = "379"
	PROPERTY_CODE_LENGTH                     string = "122"
	PROPERTY_CODE_LIGHT_INTENSITY            string = "499"
	PROPERTY_CODE_LOUDNESS                   string = "123"
	PROPERTY_CODE_LUMINOUS_FLUX              string = "504"
	PROPERTY_CODE_MAGNETIC_FLUX              string = "378"
	PROPERTY_CODE_MAGNETIC_FLUX_DENSITY      string = "503"
	PROPERTY_CODE_MASS                       string = "124"
	PROPERTY_CODE_MASS_IU                    string = "385"
	PROPERTY_CODE_MASS_UNITS                 string = "445"
	PROPERTY_CODE_MASS_PER_AREA              string = "349"
	PROPERTY_CODE_MOMENT_INERTIA_AREA        string = "344"
	PROPERTY_CODE_MOMENT_INERTIA_MASS        string = "345"
	PROPERTY_CODE_MOMENTUM                   string = "340"
	PROPERTY_CODE_MOMENTUM_FLOW_RATE         string = "346"
	PROPERTY_CODE_MOMENTUM_ANGULAR           string = "343"
	PROPERTY_CODE_POWER                      string = "363"
	PROPERTY_CODE_POWER_DENSITY              string = "369"
	PROPERTY_CODE_POWER_FLUX                 string = "368"
	PROPERTY_CODE_POWER_LINEAR               string = "367"
	PROPERTY_CODE_PRESSURE                   string = "125"
	PROPERTY_CODE_PROPORTION                 string = "507"
	PROPERTY_CODE_QUALIFIED_REAL             string = "380"
	PROPERTY_CODE_RADIOACTIVITY              string = "506"
	PROPERTY_CODE_RESISTANCE                 string = "375"
	PROPERTY_CODE_SPECIFIC_ENERGY            string = "370"
	PROPERTY_CODE_SPECIFIC_HEAT_GAS_CONSTANT string = "371"
	PROPERTY_CODE_SPECIFIC_SURFACE           string = "337"
	PROPERTY_CODE_SPECIFIC_VOLUME            string = "336"
	PROPERTY_CODE_SPECIFIC_WEIGHT            string = "354"
	PROPERTY_CODE_SURFACE_TENSION            string = "356"
	PROPERTY_CODE_TEMPERATURE                string = "127"
	PROPERTY_CODE_THERMAL_CONDUCTIVITY       string = "372"
	PROPERTY_CODE_TIME                       string = "128"
	PROPERTY_CODE_TORQUE                     string = "359"
	PROPERTY_CODE_VELOCITY                   string = "338"
	PROPERTY_CODE_VELOCITY_ANGULAR           string = "341"
	PROPERTY_CODE_VISCOSITY_DYNAMIC          string = "360"
	PROPERTY_CODE_VISCOSITY_KINEMATIC        string = "361"
	PROPERTY_CODE_ELECTRIC_POTENTIAL         string = "374"
	PROPERTY_CODE_VOLUME                     string = "129"
	PROPERTY_CODE_WORK                       string = "130"
	PROPERTY_CODE_REFRACTIVE_POWER           string = "685"
	PROPERTY_CODE_NOT_SET                    string = "118"
	PROPERTY_CODE_TIME_FRACTION              string = "709"
	PROPERTY_CODE_RATE_OF_CHANGE_PRESSURE    string = "708"
	PROPERTY_CODE_RATE_OF_CHANGE_FREQUENCY   string = "754"
	PROPERTY_CODE_ARBITRARY                  string = "755"
	PROPERTY_CODE_MEDICATION_DOSE_RATE       string = "756"
	PROPERTY_CODE_SPECTRAL_POWER             string = "757"
	PROPERTY_CODE_SPECTRAL_POWER_DENSITY     string = "758"
	PROPERTY_CODE_PACE                       string = "759"
	PROPERTY_CODE_ENZYME_ACTIVITY            string = "760"
)

// PropertyNames maps property codes to their human-readable names
var PropertyNames = map[string]string{
	PROPERTY_CODE_ACCELERATION:               "Acceleration",
	PROPERTY_CODE_ACCELERATION_ANGULAR:       "Acceleration, angular",
	PROPERTY_CODE_AMOUNT_EQ:                  "Amount (Eq)",
	PROPERTY_CODE_AMOUNT_MOLE:                "Amount (mole)",
	PROPERTY_CODE_ANGLE_PLANE:                "Angle, plane",
	PROPERTY_CODE_ANGLE_SOLID:                "Angle, solid",
	PROPERTY_CODE_AREA:                       "Area",
	PROPERTY_CODE_CONCENTRATION:              "Concentration",
	PROPERTY_CODE_DENSITY:                    "Density",
	PROPERTY_CODE_DIFFUSION_COEFFICIENT:      "Diffusion coefficient",
	PROPERTY_CODE_ELECTRIC_CAPACITANCE:       "Electric capacitance",
	PROPERTY_CODE_ELECTRIC_CHARGE:            "Electric charge",
	PROPERTY_CODE_ELECTRIC_CONDUCTANCE:       "Electric conductance",
	PROPERTY_CODE_ELECTRIC_CURRENT:           "Electric current",
	PROPERTY_CODE_ELECTRIC_FIELD_STRENGTH:    "Electric field strength",
	PROPERTY_CODE_ELECTRIC_POTENTIAL_TIME:    "Electric potential time",
	PROPERTY_CODE_ENERGY:                     "Energy",
	PROPERTY_CODE_ENERGY_DENSITY:             "Energy density",
	PROPERTY_CODE_ENERGY_DOSE:                "Energy dose",
	PROPERTY_CODE_ENERGY_PER_AREA:            "Energy per area",
	PROPERTY_CODE_ENERGY_LINEAR:              "Energy, linear",
	PROPERTY_CODE_FLOW_RATE_MASS:             "Flow rate, mass",
	PROPERTY_CODE_FLOW_RATE_MASS_FORCE:       "Flow rate, mass/force",
	PROPERTY_CODE_FLOW_RATE_MASS_VOLUME:      "Flow rate, mass/volume",
	PROPERTY_CODE_FLOW_RATE_VOLUME:           "Flow rate, volume",
	PROPERTY_CODE_FLUX_MASS:                  "Flux, mass",
	PROPERTY_CODE_FORCE:                      "Force",
	PROPERTY_CODE_FORCE_PER_MASS:             "Force per mass",
	PROPERTY_CODE_FORCE_BODY:                 "Force, body",
	PROPERTY_CODE_FREQUENCY:                  "Frequency",
	PROPERTY_CODE_GLOMERULAR_FILTRATION_RATE: "Glomerular filtration rate",
	PROPERTY_CODE_HEAT_TRANSFER_COEFFICIENT:  "Heat transfer coefficient",
	PROPERTY_CODE_ILLUMINANCE:                "Illuminance",
	PROPERTY_CODE_INDUCTANCE:                 "Inductance",
	PROPERTY_CODE_LENGTH:                     "Length",
	PROPERTY_CODE_LIGHT_INTENSITY:            "Light intensity",
	PROPERTY_CODE_LOUDNESS:                   "Loudness",
	PROPERTY_CODE_LUMINOUS_FLUX:              "Luminous flux",
	PROPERTY_CODE_MAGNETIC_FLUX:              "Magnetic flux",
	PROPERTY_CODE_MAGNETIC_FLUX_DENSITY:      "Magnetic flux density",
	PROPERTY_CODE_MASS:                       "Mass",
	PROPERTY_CODE_MASS_IU:                    "Mass (IU)",
	PROPERTY_CODE_MASS_UNITS:                 "Mass (Units)",
	PROPERTY_CODE_MASS_PER_AREA:              "Mass per area",
	PROPERTY_CODE_MOMENT_INERTIA_AREA:        "Moment inertia, area",
	PROPERTY_CODE_MOMENT_INERTIA_MASS:        "Moment inertia, mass",
	PROPERTY_CODE_MOMENTUM:                   "Momentum",
	PROPERTY_CODE_MOMENTUM_FLOW_RATE:         "Momentum flow rate",
	PROPERTY_CODE_MOMENTUM_ANGULAR:           "Momentum, angular",
	PROPERTY_CODE_POWER:                      "Power",
	PROPERTY_CODE_POWER_DENSITY:              "Power density",
	PROPERTY_CODE_POWER_FLUX:                 "Power flux",
	PROPERTY_CODE_POWER_LINEAR:               "Power, linear",
	PROPERTY_CODE_PRESSURE:                   "Pressure",
	PROPERTY_CODE_PROPORTION:                 "Proportion",
	PROPERTY_CODE_QUALIFIED_REAL:             "Qualified real",
	PROPERTY_CODE_RADIOACTIVITY:              "Radioactivity",
	PROPERTY_CODE_RESISTANCE:                 "Resistance",
	PROPERTY_CODE_SPECIFIC_ENERGY:            "Specific energy",
	PROPERTY_CODE_SPECIFIC_HEAT_GAS_CONSTANT: "Specific heat, gas constant",
	PROPERTY_CODE_SPECIFIC_SURFACE:           "Specific surface",
	PROPERTY_CODE_SPECIFIC_VOLUME:            "Specific volume",
	PROPERTY_CODE_SPECIFIC_WEIGHT:            "Specific weight",
	PROPERTY_CODE_SURFACE_TENSION:            "Surface tension",
	PROPERTY_CODE_TEMPERATURE:                "Temperature",
	PROPERTY_CODE_THERMAL_CONDUCTIVITY:       "Thermal conductivity",
	PROPERTY_CODE_TIME:                       "Time",
	PROPERTY_CODE_TORQUE:                     "Torque",
	PROPERTY_CODE_VELOCITY:                   "Velocity",
	PROPERTY_CODE_VELOCITY_ANGULAR:           "Velocity, angular",
	PROPERTY_CODE_VISCOSITY_DYNAMIC:          "Viscosity, dynamic",
	PROPERTY_CODE_VISCOSITY_KINEMATIC:        "Viscosity, kinematic",
	PROPERTY_CODE_ELECTRIC_POTENTIAL:         "Electric potential",
	PROPERTY_CODE_VOLUME:                     "Volume",
	PROPERTY_CODE_WORK:                       "Work",
	PROPERTY_CODE_REFRACTIVE_POWER:           "Refractive power",
	PROPERTY_CODE_NOT_SET:                    "<not set>",
	PROPERTY_CODE_TIME_FRACTION:              "Time fraction",
	PROPERTY_CODE_RATE_OF_CHANGE_PRESSURE:    "Rate of change, pressure",
	PROPERTY_CODE_RATE_OF_CHANGE_FREQUENCY:   "Rate of change, frequency",
	PROPERTY_CODE_ARBITRARY:                  "Arbitrary",
	PROPERTY_CODE_MEDICATION_DOSE_RATE:       "Medication dose rate",
	PROPERTY_CODE_SPECTRAL_POWER:             "Spectral power",
	PROPERTY_CODE_SPECTRAL_POWER_DENSITY:     "Spectral power density",
	PROPERTY_CODE_PACE:                       "Pace",
	PROPERTY_CODE_ENZYME_ACTIVITY:            "Enzyme activity",
}

// IsValidPropertyCode checks if the given code is a valid property
func IsValidPropertyCode(code string) bool {
	_, exists := PropertyNames[code]
	return exists
}

// GetPropertyName returns the human-readable name for the given property code
func GetPropertyName(code string) string {
	if name, exists := PropertyNames[code]; exists {
		return name
	}
	return ""
}
