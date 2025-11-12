package terminology

// Countries - ISO 3166-1 "alpha-2" country codes as defined in openEHR Support Terminology
// Used in: COMPOSITION.territory
// Source: https://specifications.openehr.org/releases/TERM/latest/SupportTerminology.html#_countries

const (
	COUNTRY_TERMINOLOGY_ID_ISO string = "ISO_3166-1"
)

var CountryTerminologyIds = map[string]string{
	COUNTRY_TERMINOLOGY_ID_ISO: "ISO 3166-1",
}

func IsValidCountryTerminologyID(terminologyID string) bool {
	_, exists := CountryTerminologyIds[terminologyID]
	return exists
}

const (
	// Country codes
	COUNTRY_CODE_AFGHANISTAN                       string = "AF"
	COUNTRY_CODE_ALAND_ISLANDS                     string = "AX"
	COUNTRY_CODE_ALBANIA                           string = "AL"
	COUNTRY_CODE_ALGERIA                           string = "DZ"
	COUNTRY_CODE_AMERICAN_SAMOA                    string = "AS"
	COUNTRY_CODE_ANDORRA                           string = "AD"
	COUNTRY_CODE_ANGOLA                            string = "AO"
	COUNTRY_CODE_ANGUILLA                          string = "AI"
	COUNTRY_CODE_ANTARCTICA                        string = "AQ"
	COUNTRY_CODE_ANTIGUA_AND_BARBUDA               string = "AG"
	COUNTRY_CODE_ARGENTINA                         string = "AR"
	COUNTRY_CODE_ARMENIA                           string = "AM"
	COUNTRY_CODE_ARUBA                             string = "AW"
	COUNTRY_CODE_AUSTRALIA                         string = "AU"
	COUNTRY_CODE_AUSTRIA                           string = "AT"
	COUNTRY_CODE_AZERBAIJAN                        string = "AZ"
	COUNTRY_CODE_BAHAMAS                           string = "BS"
	COUNTRY_CODE_BAHRAIN                           string = "BH"
	COUNTRY_CODE_BANGLADESH                        string = "BD"
	COUNTRY_CODE_BARBADOS                          string = "BB"
	COUNTRY_CODE_BELARUS                           string = "BY"
	COUNTRY_CODE_BELGIUM                           string = "BE"
	COUNTRY_CODE_BELIZE                            string = "BZ"
	COUNTRY_CODE_BENIN                             string = "BJ"
	COUNTRY_CODE_BERMUDA                           string = "BM"
	COUNTRY_CODE_BHUTAN                            string = "BT"
	COUNTRY_CODE_BOLIVIA                           string = "BO"
	COUNTRY_CODE_BONAIRE_SINT_EUSTATIUS_AND_SABA   string = "BQ"
	COUNTRY_CODE_BOSNIA_AND_HERZEGOVINA            string = "BA"
	COUNTRY_CODE_BOTSWANA                          string = "BW"
	COUNTRY_CODE_BOUVET_ISLAND                     string = "BV"
	COUNTRY_CODE_BRAZIL                            string = "BR"
	COUNTRY_CODE_BRITISH_INDIAN_OCEAN_TERRITORY    string = "IO"
	COUNTRY_CODE_BRUNEI_DARUSSALAM                 string = "BN"
	COUNTRY_CODE_BULGARIA                          string = "BG"
	COUNTRY_CODE_BURKINA_FASO                      string = "BF"
	COUNTRY_CODE_BURUNDI                           string = "BI"
	COUNTRY_CODE_CAMBODIA                          string = "KH"
	COUNTRY_CODE_CAMEROON                          string = "CM"
	COUNTRY_CODE_CANADA                            string = "CA"
	COUNTRY_CODE_CAPE_VERDE                        string = "CV"
	COUNTRY_CODE_CAYMAN_ISLANDS                    string = "KY"
	COUNTRY_CODE_CENTRAL_AFRICAN_REPUBLIC          string = "CF"
	COUNTRY_CODE_CHAD                              string = "TD"
	COUNTRY_CODE_CHILE                             string = "CL"
	COUNTRY_CODE_CHINA                             string = "CN"
	COUNTRY_CODE_CHRISTMAS_ISLAND                  string = "CX"
	COUNTRY_CODE_COCOS_KEELING_ISLANDS             string = "CC"
	COUNTRY_CODE_COLOMBIA                          string = "CO"
	COUNTRY_CODE_COMOROS                           string = "KM"
	COUNTRY_CODE_CONGO                             string = "CG"
	COUNTRY_CODE_CONGO_DEMOCRATIC_REPUBLIC         string = "CD"
	COUNTRY_CODE_COOK_ISLANDS                      string = "CK"
	COUNTRY_CODE_COSTA_RICA                        string = "CR"
	COUNTRY_CODE_COTE_DIVOIRE                      string = "CI"
	COUNTRY_CODE_CROATIA                           string = "HR"
	COUNTRY_CODE_CUBA                              string = "CU"
	COUNTRY_CODE_CURACAO                           string = "CW"
	COUNTRY_CODE_CYPRUS                            string = "CY"
	COUNTRY_CODE_CZECH_REPUBLIC                    string = "CZ"
	COUNTRY_CODE_DENMARK                           string = "DK"
	COUNTRY_CODE_DJIBOUTI                          string = "DJ"
	COUNTRY_CODE_DOMINICA                          string = "DM"
	COUNTRY_CODE_DOMINICAN_REPUBLIC                string = "DO"
	COUNTRY_CODE_ECUADOR                           string = "EC"
	COUNTRY_CODE_EGYPT                             string = "EG"
	COUNTRY_CODE_EL_SALVADOR                       string = "SV"
	COUNTRY_CODE_EQUATORIAL_GUINEA                 string = "GQ"
	COUNTRY_CODE_ERITREA                           string = "ER"
	COUNTRY_CODE_ESTONIA                           string = "EE"
	COUNTRY_CODE_ESWATINI                          string = "SZ"
	COUNTRY_CODE_ETHIOPIA                          string = "ET"
	COUNTRY_CODE_FALKLAND_ISLANDS                  string = "FK"
	COUNTRY_CODE_FAROE_ISLANDS                     string = "FO"
	COUNTRY_CODE_FIJI                              string = "FJ"
	COUNTRY_CODE_FINLAND                           string = "FI"
	COUNTRY_CODE_FRANCE                            string = "FR"
	COUNTRY_CODE_FRENCH_GUIANA                     string = "GF"
	COUNTRY_CODE_FRENCH_POLYNESIA                  string = "PF"
	COUNTRY_CODE_FRENCH_SOUTHERN_TERRITORIES       string = "TF"
	COUNTRY_CODE_GABON                             string = "GA"
	COUNTRY_CODE_GAMBIA                            string = "GM"
	COUNTRY_CODE_GEORGIA                           string = "GE"
	COUNTRY_CODE_GERMANY                           string = "DE"
	COUNTRY_CODE_GHANA                             string = "GH"
	COUNTRY_CODE_GIBRALTAR                         string = "GI"
	COUNTRY_CODE_GREECE                            string = "GR"
	COUNTRY_CODE_GREENLAND                         string = "GL"
	COUNTRY_CODE_GRENADA                           string = "GD"
	COUNTRY_CODE_GUADELOUPE                        string = "GP"
	COUNTRY_CODE_GUAM                              string = "GU"
	COUNTRY_CODE_GUATEMALA                         string = "GT"
	COUNTRY_CODE_GUERNSEY                          string = "GG"
	COUNTRY_CODE_GUINEA                            string = "GN"
	COUNTRY_CODE_GUINEA_BISSAU                     string = "GW"
	COUNTRY_CODE_GUYANA                            string = "GY"
	COUNTRY_CODE_HAITI                             string = "HT"
	COUNTRY_CODE_HEARD_ISLAND_AND_MCDONALD_ISLANDS string = "HM"
	COUNTRY_CODE_HOLY_SEE_VATICAN_CITY_STATE       string = "VA"
	COUNTRY_CODE_HONDURAS                          string = "HN"
	COUNTRY_CODE_HONG_KONG                         string = "HK"
	COUNTRY_CODE_HUNGARY                           string = "HU"
	COUNTRY_CODE_ICELAND                           string = "IS"
	COUNTRY_CODE_INDIA                             string = "IN"
	COUNTRY_CODE_INDONESIA                         string = "ID"
	COUNTRY_CODE_IRAN                              string = "IR"
	COUNTRY_CODE_IRAQ                              string = "IQ"
	COUNTRY_CODE_IRELAND                           string = "IE"
	COUNTRY_CODE_ISLE_OF_MAN                       string = "IM"
	COUNTRY_CODE_ISRAEL                            string = "IL"
	COUNTRY_CODE_ITALY                             string = "IT"
	COUNTRY_CODE_JAMAICA                           string = "JM"
	COUNTRY_CODE_JAPAN                             string = "JP"
	COUNTRY_CODE_JERSEY                            string = "JE"
	COUNTRY_CODE_JORDAN                            string = "JO"
	COUNTRY_CODE_KAZAKHSTAN                        string = "KZ"
	COUNTRY_CODE_KENYA                             string = "KE"
	COUNTRY_CODE_KIRIBATI                          string = "KI"
	COUNTRY_CODE_KOREA_DEMOCRATIC_PEOPLES_REPUBLIC string = "KP"
	COUNTRY_CODE_KOREA_REPUBLIC                    string = "KR"
	COUNTRY_CODE_KUWAIT                            string = "KW"
	COUNTRY_CODE_KYRGYZSTAN                        string = "KG"
	COUNTRY_CODE_LAO_PEOPLES_DEMOCRATIC_REPUBLIC   string = "LA"
	COUNTRY_CODE_LATVIA                            string = "LV"
	COUNTRY_CODE_LEBANON                           string = "LB"
	COUNTRY_CODE_LESOTHO                           string = "LS"
	COUNTRY_CODE_LIBERIA                           string = "LR"
	COUNTRY_CODE_LIBYA                             string = "LY"
	COUNTRY_CODE_LIECHTENSTEIN                     string = "LI"
	COUNTRY_CODE_LITHUANIA                         string = "LT"
	COUNTRY_CODE_LUXEMBOURG                        string = "LU"
	COUNTRY_CODE_MACAO                             string = "MO"
	COUNTRY_CODE_MADAGASCAR                        string = "MG"
	COUNTRY_CODE_MALAWI                            string = "MW"
	COUNTRY_CODE_MALAYSIA                          string = "MY"
	COUNTRY_CODE_MALDIVES                          string = "MV"
	COUNTRY_CODE_MALI                              string = "ML"
	COUNTRY_CODE_MALTA                             string = "MT"
	COUNTRY_CODE_MARSHALL_ISLANDS                  string = "MH"
	COUNTRY_CODE_MARTINIQUE                        string = "MQ"
	COUNTRY_CODE_MAURITANIA                        string = "MR"
	COUNTRY_CODE_MAURITIUS                         string = "MU"
	COUNTRY_CODE_MAYOTTE                           string = "YT"
	COUNTRY_CODE_MEXICO                            string = "MX"
	COUNTRY_CODE_MICRONESIA                        string = "FM"
	COUNTRY_CODE_MOLDOVA                           string = "MD"
	COUNTRY_CODE_MONACO                            string = "MC"
	COUNTRY_CODE_MONGOLIA                          string = "MN"
	COUNTRY_CODE_MONTENEGRO                        string = "ME"
	COUNTRY_CODE_MONTSERRAT                        string = "MS"
	COUNTRY_CODE_MOROCCO                           string = "MA"
	COUNTRY_CODE_MOZAMBIQUE                        string = "MZ"
	COUNTRY_CODE_MYANMAR                           string = "MM"
	COUNTRY_CODE_NAMIBIA                           string = "NA"
	COUNTRY_CODE_NAURU                             string = "NR"
	COUNTRY_CODE_NEPAL                             string = "NP"
	COUNTRY_CODE_NETHERLANDS                       string = "NL"
	// COUNTRY_CODE_NETHERLANDS_ANTILLES                     string = "AN" // DEPRECATED
	COUNTRY_CODE_NEW_CALEDONIA                            string = "NC"
	COUNTRY_CODE_NEW_ZEALAND                              string = "NZ"
	COUNTRY_CODE_NICARAGUA                                string = "NI"
	COUNTRY_CODE_NIGER                                    string = "NE"
	COUNTRY_CODE_NIGERIA                                  string = "NG"
	COUNTRY_CODE_NIUE                                     string = "NU"
	COUNTRY_CODE_NORFOLK_ISLAND                           string = "NF"
	COUNTRY_CODE_NORTH_MACEDONIA                          string = "MK"
	COUNTRY_CODE_NORTHERN_MARIANA_ISLANDS                 string = "MP"
	COUNTRY_CODE_NORWAY                                   string = "NO"
	COUNTRY_CODE_OMAN                                     string = "OM"
	COUNTRY_CODE_PAKISTAN                                 string = "PK"
	COUNTRY_CODE_PALAU                                    string = "PW"
	COUNTRY_CODE_PALESTINIAN_STATE                        string = "PS"
	COUNTRY_CODE_PANAMA                                   string = "PA"
	COUNTRY_CODE_PAPUA_NEW_GUINEA                         string = "PG"
	COUNTRY_CODE_PARAGUAY                                 string = "PY"
	COUNTRY_CODE_PERU                                     string = "PE"
	COUNTRY_CODE_PHILIPPINES                              string = "PH"
	COUNTRY_CODE_PITCAIRN                                 string = "PN"
	COUNTRY_CODE_POLAND                                   string = "PL"
	COUNTRY_CODE_PORTUGAL                                 string = "PT"
	COUNTRY_CODE_PUERTO_RICO                              string = "PR"
	COUNTRY_CODE_QATAR                                    string = "QA"
	COUNTRY_CODE_REUNION                                  string = "RE"
	COUNTRY_CODE_ROMANIA                                  string = "RO"
	COUNTRY_CODE_RUSSIAN_FEDERATION                       string = "RU"
	COUNTRY_CODE_RWANDA                                   string = "RW"
	COUNTRY_CODE_SAINT_BARTHELEMY                         string = "BL"
	COUNTRY_CODE_SAINT_HELENA_ASCENSION_TRISTAN_DA_CUNHA  string = "SH"
	COUNTRY_CODE_SAINT_KITTS_AND_NEVIS                    string = "KN"
	COUNTRY_CODE_SAINT_LUCIA                              string = "LC"
	COUNTRY_CODE_SAINT_MARTIN_FRENCH                      string = "MF"
	COUNTRY_CODE_SAINT_PIERRE_AND_MIQUELON                string = "PM"
	COUNTRY_CODE_SAINT_VINCENT_AND_THE_GRENADINES         string = "VC"
	COUNTRY_CODE_SAMOA                                    string = "WS"
	COUNTRY_CODE_SAN_MARINO                               string = "SM"
	COUNTRY_CODE_SAO_TOME_AND_PRINCIPE                    string = "ST"
	COUNTRY_CODE_SAUDI_ARABIA                             string = "SA"
	COUNTRY_CODE_SENEGAL                                  string = "SN"
	COUNTRY_CODE_SERBIA                                   string = "RS"
	COUNTRY_CODE_SEYCHELLES                               string = "SC"
	COUNTRY_CODE_SIERRA_LEONE                             string = "SL"
	COUNTRY_CODE_SINGAPORE                                string = "SG"
	COUNTRY_CODE_SINT_MAARTEN_DUTCH                       string = "SX"
	COUNTRY_CODE_SLOVAKIA                                 string = "SK"
	COUNTRY_CODE_SLOVENIA                                 string = "SI"
	COUNTRY_CODE_SOLOMON_ISLANDS                          string = "SB"
	COUNTRY_CODE_SOMALIA                                  string = "SO"
	COUNTRY_CODE_SOUTH_AFRICA                             string = "ZA"
	COUNTRY_CODE_SOUTH_GEORGIA_AND_SOUTH_SANDWICH_ISLANDS string = "GS"
	COUNTRY_CODE_SOUTH_SUDAN                              string = "SS"
	COUNTRY_CODE_SPAIN                                    string = "ES"
	COUNTRY_CODE_SRI_LANKA                                string = "LK"
	COUNTRY_CODE_SUDAN                                    string = "SD"
	COUNTRY_CODE_SURINAME                                 string = "SR"
	COUNTRY_CODE_SVALBARD_AND_JAN_MAYEN                   string = "SJ"
	COUNTRY_CODE_SWEDEN                                   string = "SE"
	COUNTRY_CODE_SWITZERLAND                              string = "CH"
	COUNTRY_CODE_SYRIAN_ARAB_REPUBLIC                     string = "SY"
	COUNTRY_CODE_TAIWAN                                   string = "TW"
	COUNTRY_CODE_TAJIKISTAN                               string = "TJ"
	COUNTRY_CODE_TANZANIA                                 string = "TZ"
	COUNTRY_CODE_THAILAND                                 string = "TH"
	COUNTRY_CODE_TIMOR_LESTE                              string = "TL"
	COUNTRY_CODE_TOGO                                     string = "TG"
	COUNTRY_CODE_TOKELAU                                  string = "TK"
	COUNTRY_CODE_TONGA                                    string = "TO"
	COUNTRY_CODE_TRINIDAD_AND_TOBAGO                      string = "TT"
	COUNTRY_CODE_TUNISIA                                  string = "TN"
	COUNTRY_CODE_TURKIYE                                  string = "TR"
	COUNTRY_CODE_TURKMENISTAN                             string = "TM"
	COUNTRY_CODE_TURKS_AND_CAICOS_ISLANDS                 string = "TC"
	COUNTRY_CODE_TUVALU                                   string = "TV"
	COUNTRY_CODE_UGANDA                                   string = "UG"
	COUNTRY_CODE_UKRAINE                                  string = "UA"
	COUNTRY_CODE_UNITED_ARAB_EMIRATES                     string = "AE"
	COUNTRY_CODE_UNITED_KINGDOM                           string = "GB"
	COUNTRY_CODE_UNITED_STATES                            string = "US"
	COUNTRY_CODE_UNITED_STATES_MINOR_OUTLYING_ISLANDS     string = "UM"
	COUNTRY_CODE_URUGUAY                                  string = "UY"
	COUNTRY_CODE_UZBEKISTAN                               string = "UZ"
	COUNTRY_CODE_VANUATU                                  string = "VU"
	COUNTRY_CODE_VENEZUELA                                string = "VE"
	COUNTRY_CODE_VIET_NAM                                 string = "VN"
	COUNTRY_CODE_VIRGIN_ISLANDS_BRITISH                   string = "VG"
	COUNTRY_CODE_VIRGIN_ISLANDS_US                        string = "VI"
	COUNTRY_CODE_WALLIS_AND_FUTUNA                        string = "WF"
	COUNTRY_CODE_WESTERN_SAHARA                           string = "EH"
	COUNTRY_CODE_YEMEN                                    string = "YE"
	COUNTRY_CODE_ZAMBIA                                   string = "ZM"
	COUNTRY_CODE_ZIMBABWE                                 string = "ZW"
)

// CountryNames provides a mapping from country codes to their full names
var CountryNames = map[string]string{
	COUNTRY_CODE_AFGHANISTAN:                       "AFGHANISTAN",
	COUNTRY_CODE_ALAND_ISLANDS:                     "ÅLAND ISLANDS",
	COUNTRY_CODE_ALBANIA:                           "ALBANIA",
	COUNTRY_CODE_ALGERIA:                           "ALGERIA",
	COUNTRY_CODE_AMERICAN_SAMOA:                    "AMERICAN SAMOA",
	COUNTRY_CODE_ANDORRA:                           "ANDORRA",
	COUNTRY_CODE_ANGOLA:                            "ANGOLA",
	COUNTRY_CODE_ANGUILLA:                          "ANGUILLA",
	COUNTRY_CODE_ANTARCTICA:                        "ANTARCTICA",
	COUNTRY_CODE_ANTIGUA_AND_BARBUDA:               "ANTIGUA AND BARBUDA",
	COUNTRY_CODE_ARGENTINA:                         "ARGENTINA",
	COUNTRY_CODE_ARMENIA:                           "ARMENIA",
	COUNTRY_CODE_ARUBA:                             "ARUBA",
	COUNTRY_CODE_AUSTRALIA:                         "AUSTRALIA",
	COUNTRY_CODE_AUSTRIA:                           "AUSTRIA",
	COUNTRY_CODE_AZERBAIJAN:                        "AZERBAIJAN",
	COUNTRY_CODE_BAHAMAS:                           "BAHAMAS",
	COUNTRY_CODE_BAHRAIN:                           "BAHRAIN",
	COUNTRY_CODE_BANGLADESH:                        "BANGLADESH",
	COUNTRY_CODE_BARBADOS:                          "BARBADOS",
	COUNTRY_CODE_BELARUS:                           "BELARUS",
	COUNTRY_CODE_BELGIUM:                           "BELGIUM",
	COUNTRY_CODE_BELIZE:                            "BELIZE",
	COUNTRY_CODE_BENIN:                             "BENIN",
	COUNTRY_CODE_BERMUDA:                           "BERMUDA",
	COUNTRY_CODE_BHUTAN:                            "BHUTAN",
	COUNTRY_CODE_BOLIVIA:                           "BOLIVIA",
	COUNTRY_CODE_BONAIRE_SINT_EUSTATIUS_AND_SABA:   "BONAIRE, SINT EUSTATIUS AND SABA",
	COUNTRY_CODE_BOSNIA_AND_HERZEGOVINA:            "BOSNIA AND HERZEGOVINA",
	COUNTRY_CODE_BOTSWANA:                          "BOTSWANA",
	COUNTRY_CODE_BOUVET_ISLAND:                     "BOUVET ISLAND",
	COUNTRY_CODE_BRAZIL:                            "BRAZIL",
	COUNTRY_CODE_BRITISH_INDIAN_OCEAN_TERRITORY:    "BRITISH INDIAN OCEAN TERRITORY",
	COUNTRY_CODE_BRUNEI_DARUSSALAM:                 "BRUNEI DARUSSALAM",
	COUNTRY_CODE_BULGARIA:                          "BULGARIA",
	COUNTRY_CODE_BURKINA_FASO:                      "BURKINA FASO",
	COUNTRY_CODE_BURUNDI:                           "BURUNDI",
	COUNTRY_CODE_CAMBODIA:                          "CAMBODIA",
	COUNTRY_CODE_CAMEROON:                          "CAMEROON",
	COUNTRY_CODE_CANADA:                            "CANADA",
	COUNTRY_CODE_CAPE_VERDE:                        "CAPE VERDE",
	COUNTRY_CODE_CAYMAN_ISLANDS:                    "CAYMAN ISLANDS",
	COUNTRY_CODE_CENTRAL_AFRICAN_REPUBLIC:          "CENTRAL AFRICAN REPUBLIC",
	COUNTRY_CODE_CHAD:                              "CHAD",
	COUNTRY_CODE_CHILE:                             "CHILE",
	COUNTRY_CODE_CHINA:                             "CHINA",
	COUNTRY_CODE_CHRISTMAS_ISLAND:                  "CHRISTMAS ISLAND",
	COUNTRY_CODE_COCOS_KEELING_ISLANDS:             "COCOS (KEELING) ISLANDS",
	COUNTRY_CODE_COLOMBIA:                          "COLOMBIA",
	COUNTRY_CODE_COMOROS:                           "COMOROS",
	COUNTRY_CODE_CONGO:                             "CONGO",
	COUNTRY_CODE_CONGO_DEMOCRATIC_REPUBLIC:         "CONGO, THE DEMOCRATIC REPUBLIC OF THE",
	COUNTRY_CODE_COOK_ISLANDS:                      "COOK ISLANDS",
	COUNTRY_CODE_COSTA_RICA:                        "COSTA RICA",
	COUNTRY_CODE_COTE_DIVOIRE:                      "CÔTE D'IVOIRE",
	COUNTRY_CODE_CROATIA:                           "CROATIA",
	COUNTRY_CODE_CUBA:                              "CUBA",
	COUNTRY_CODE_CURACAO:                           "CURAÇAO",
	COUNTRY_CODE_CYPRUS:                            "CYPRUS",
	COUNTRY_CODE_CZECH_REPUBLIC:                    "CZECH REPUBLIC",
	COUNTRY_CODE_DENMARK:                           "DENMARK",
	COUNTRY_CODE_DJIBOUTI:                          "DJIBOUTI",
	COUNTRY_CODE_DOMINICA:                          "DOMINICA",
	COUNTRY_CODE_DOMINICAN_REPUBLIC:                "DOMINICAN REPUBLIC",
	COUNTRY_CODE_ECUADOR:                           "ECUADOR",
	COUNTRY_CODE_EGYPT:                             "EGYPT",
	COUNTRY_CODE_EL_SALVADOR:                       "EL SALVADOR",
	COUNTRY_CODE_EQUATORIAL_GUINEA:                 "EQUATORIAL GUINEA",
	COUNTRY_CODE_ERITREA:                           "ERITREA",
	COUNTRY_CODE_ESTONIA:                           "ESTONIA",
	COUNTRY_CODE_ESWATINI:                          "ESWATINI",
	COUNTRY_CODE_ETHIOPIA:                          "ETHIOPIA",
	COUNTRY_CODE_FALKLAND_ISLANDS:                  "FALKLAND ISLANDS (MALVINAS)",
	COUNTRY_CODE_FAROE_ISLANDS:                     "FAROE ISLANDS",
	COUNTRY_CODE_FIJI:                              "FIJI",
	COUNTRY_CODE_FINLAND:                           "FINLAND",
	COUNTRY_CODE_FRANCE:                            "FRANCE",
	COUNTRY_CODE_FRENCH_GUIANA:                     "FRENCH GUIANA",
	COUNTRY_CODE_FRENCH_POLYNESIA:                  "FRENCH POLYNESIA",
	COUNTRY_CODE_FRENCH_SOUTHERN_TERRITORIES:       "FRENCH SOUTHERN TERRITORIES",
	COUNTRY_CODE_GABON:                             "GABON",
	COUNTRY_CODE_GAMBIA:                            "GAMBIA",
	COUNTRY_CODE_GEORGIA:                           "GEORGIA",
	COUNTRY_CODE_GERMANY:                           "GERMANY",
	COUNTRY_CODE_GHANA:                             "GHANA",
	COUNTRY_CODE_GIBRALTAR:                         "GIBRALTAR",
	COUNTRY_CODE_GREECE:                            "GREECE",
	COUNTRY_CODE_GREENLAND:                         "GREENLAND",
	COUNTRY_CODE_GRENADA:                           "GRENADA",
	COUNTRY_CODE_GUADELOUPE:                        "GUADELOUPE",
	COUNTRY_CODE_GUAM:                              "GUAM",
	COUNTRY_CODE_GUATEMALA:                         "GUATEMALA",
	COUNTRY_CODE_GUERNSEY:                          "GUERNSEY",
	COUNTRY_CODE_GUINEA:                            "GUINEA",
	COUNTRY_CODE_GUINEA_BISSAU:                     "GUINEA-BISSAU",
	COUNTRY_CODE_GUYANA:                            "GUYANA",
	COUNTRY_CODE_HAITI:                             "HAITI",
	COUNTRY_CODE_HEARD_ISLAND_AND_MCDONALD_ISLANDS: "HEARD ISLAND AND MCDONALD ISLANDS",
	COUNTRY_CODE_HOLY_SEE_VATICAN_CITY_STATE:       "HOLY SEE (VATICAN CITY STATE)",
	COUNTRY_CODE_HONDURAS:                          "HONDURAS",
	COUNTRY_CODE_HONG_KONG:                         "HONG KONG",
	COUNTRY_CODE_HUNGARY:                           "HUNGARY",
	COUNTRY_CODE_ICELAND:                           "ICELAND",
	COUNTRY_CODE_INDIA:                             "INDIA",
	COUNTRY_CODE_INDONESIA:                         "INDONESIA",
	COUNTRY_CODE_IRAN:                              "IRAN, ISLAMIC REPUBLIC OF",
	COUNTRY_CODE_IRAQ:                              "IRAQ",
	COUNTRY_CODE_IRELAND:                           "IRELAND",
	COUNTRY_CODE_ISLE_OF_MAN:                       "ISLE OF MAN",
	COUNTRY_CODE_ISRAEL:                            "ISRAEL",
	COUNTRY_CODE_ITALY:                             "ITALY",
	COUNTRY_CODE_JAMAICA:                           "JAMAICA",
	COUNTRY_CODE_JAPAN:                             "JAPAN",
	COUNTRY_CODE_JERSEY:                            "JERSEY",
	COUNTRY_CODE_JORDAN:                            "JORDAN",
	COUNTRY_CODE_KAZAKHSTAN:                        "KAZAKHSTAN",
	COUNTRY_CODE_KENYA:                             "KENYA",
	COUNTRY_CODE_KIRIBATI:                          "KIRIBATI",
	COUNTRY_CODE_KOREA_DEMOCRATIC_PEOPLES_REPUBLIC: "KOREA, DEMOCRATIC PEOPLE'S REPUBLIC OF",
	COUNTRY_CODE_KOREA_REPUBLIC:                    "KOREA, REPUBLIC OF",
	COUNTRY_CODE_KUWAIT:                            "KUWAIT",
	COUNTRY_CODE_KYRGYZSTAN:                        "KYRGYZSTAN",
	COUNTRY_CODE_LAO_PEOPLES_DEMOCRATIC_REPUBLIC:   "LAO PEOPLE'S DEMOCRATIC REPUBLIC",
	COUNTRY_CODE_LATVIA:                            "LATVIA",
	COUNTRY_CODE_LEBANON:                           "LEBANON",
	COUNTRY_CODE_LESOTHO:                           "LESOTHO",
	COUNTRY_CODE_LIBERIA:                           "LIBERIA",
	COUNTRY_CODE_LIBYA:                             "LIBYA",
	COUNTRY_CODE_LIECHTENSTEIN:                     "LIECHTENSTEIN",
	COUNTRY_CODE_LITHUANIA:                         "LITHUANIA",
	COUNTRY_CODE_LUXEMBOURG:                        "LUXEMBOURG",
	COUNTRY_CODE_MACAO:                             "MACAO",
	COUNTRY_CODE_MADAGASCAR:                        "MADAGASCAR",
	COUNTRY_CODE_MALAWI:                            "MALAWI",
	COUNTRY_CODE_MALAYSIA:                          "MALAYSIA",
	COUNTRY_CODE_MALDIVES:                          "MALDIVES",
	COUNTRY_CODE_MALI:                              "MALI",
	COUNTRY_CODE_MALTA:                             "MALTA",
	COUNTRY_CODE_MARSHALL_ISLANDS:                  "MARSHALL ISLANDS",
	COUNTRY_CODE_MARTINIQUE:                        "MARTINIQUE",
	COUNTRY_CODE_MAURITANIA:                        "MAURITANIA",
	COUNTRY_CODE_MAURITIUS:                         "MAURITIUS",
	COUNTRY_CODE_MAYOTTE:                           "MAYOTTE",
	COUNTRY_CODE_MEXICO:                            "MEXICO",
	COUNTRY_CODE_MICRONESIA:                        "MICRONESIA, FEDERATED STATES OF",
	COUNTRY_CODE_MOLDOVA:                           "MOLDOVA, REPUBLIC OF",
	COUNTRY_CODE_MONACO:                            "MONACO",
	COUNTRY_CODE_MONGOLIA:                          "MONGOLIA",
	COUNTRY_CODE_MONTENEGRO:                        "MONTENEGRO",
	COUNTRY_CODE_MONTSERRAT:                        "MONTSERRAT",
	COUNTRY_CODE_MOROCCO:                           "MOROCCO",
	COUNTRY_CODE_MOZAMBIQUE:                        "MOZAMBIQUE",
	COUNTRY_CODE_MYANMAR:                           "MYANMAR",
	COUNTRY_CODE_NAMIBIA:                           "NAMIBIA",
	COUNTRY_CODE_NAURU:                             "NAURU",
	COUNTRY_CODE_NEPAL:                             "NEPAL",
	COUNTRY_CODE_NETHERLANDS:                       "NETHERLANDS",
	// COUNTRY_CODE_NETHERLANDS_ANTILLES:              "NETHERLANDS ANTILLES - DEPRECATED",
	COUNTRY_CODE_NEW_CALEDONIA:                            "NEW CALEDONIA",
	COUNTRY_CODE_NEW_ZEALAND:                              "NEW ZEALAND",
	COUNTRY_CODE_NICARAGUA:                                "NICARAGUA",
	COUNTRY_CODE_NIGER:                                    "NIGER",
	COUNTRY_CODE_NIGERIA:                                  "NIGERIA",
	COUNTRY_CODE_NIUE:                                     "NIUE",
	COUNTRY_CODE_NORFOLK_ISLAND:                           "NORFOLK ISLAND",
	COUNTRY_CODE_NORTH_MACEDONIA:                          "NORTH MACEDONIA",
	COUNTRY_CODE_NORTHERN_MARIANA_ISLANDS:                 "NORTHERN MARIANA ISLANDS",
	COUNTRY_CODE_NORWAY:                                   "NORWAY",
	COUNTRY_CODE_OMAN:                                     "OMAN",
	COUNTRY_CODE_PAKISTAN:                                 "PAKISTAN",
	COUNTRY_CODE_PALAU:                                    "PALAU",
	COUNTRY_CODE_PALESTINIAN_STATE:                        "PALESTINIAN, STATE OF",
	COUNTRY_CODE_PANAMA:                                   "PANAMA",
	COUNTRY_CODE_PAPUA_NEW_GUINEA:                         "PAPUA NEW GUINEA",
	COUNTRY_CODE_PARAGUAY:                                 "PARAGUAY",
	COUNTRY_CODE_PERU:                                     "PERU",
	COUNTRY_CODE_PHILIPPINES:                              "PHILIPPINES",
	COUNTRY_CODE_PITCAIRN:                                 "PITCAIRN",
	COUNTRY_CODE_POLAND:                                   "POLAND",
	COUNTRY_CODE_PORTUGAL:                                 "PORTUGAL",
	COUNTRY_CODE_PUERTO_RICO:                              "PUERTO RICO",
	COUNTRY_CODE_QATAR:                                    "QATAR",
	COUNTRY_CODE_REUNION:                                  "RÉUNION",
	COUNTRY_CODE_ROMANIA:                                  "ROMANIA",
	COUNTRY_CODE_RUSSIAN_FEDERATION:                       "RUSSIAN FEDERATION",
	COUNTRY_CODE_RWANDA:                                   "RWANDA",
	COUNTRY_CODE_SAINT_BARTHELEMY:                         "SAINT BARTHÉLEMY",
	COUNTRY_CODE_SAINT_HELENA_ASCENSION_TRISTAN_DA_CUNHA:  "SAINT HELENA, ASCENSION AND TRISTAN DA CUNHA",
	COUNTRY_CODE_SAINT_KITTS_AND_NEVIS:                    "SAINT KITTS AND NEVIS",
	COUNTRY_CODE_SAINT_LUCIA:                              "SAINT LUCIA",
	COUNTRY_CODE_SAINT_MARTIN_FRENCH:                      "SAINT MARTIN (FRENCH PART)",
	COUNTRY_CODE_SAINT_PIERRE_AND_MIQUELON:                "SAINT PIERRE AND MIQUELON",
	COUNTRY_CODE_SAINT_VINCENT_AND_THE_GRENADINES:         "SAINT VINCENT AND THE GRENADINES",
	COUNTRY_CODE_SAMOA:                                    "SAMOA",
	COUNTRY_CODE_SAN_MARINO:                               "SAN MARINO",
	COUNTRY_CODE_SAO_TOME_AND_PRINCIPE:                    "SAO TOME AND PRINCIPE",
	COUNTRY_CODE_SAUDI_ARABIA:                             "SAUDI ARABIA",
	COUNTRY_CODE_SENEGAL:                                  "SENEGAL",
	COUNTRY_CODE_SERBIA:                                   "SERBIA",
	COUNTRY_CODE_SEYCHELLES:                               "SEYCHELLES",
	COUNTRY_CODE_SIERRA_LEONE:                             "SIERRA LEONE",
	COUNTRY_CODE_SINGAPORE:                                "SINGAPORE",
	COUNTRY_CODE_SINT_MAARTEN_DUTCH:                       "SINT MAARTEN (DUTCH PART)",
	COUNTRY_CODE_SLOVAKIA:                                 "SLOVAKIA",
	COUNTRY_CODE_SLOVENIA:                                 "SLOVENIA",
	COUNTRY_CODE_SOLOMON_ISLANDS:                          "SOLOMON ISLANDS",
	COUNTRY_CODE_SOMALIA:                                  "SOMALIA",
	COUNTRY_CODE_SOUTH_AFRICA:                             "SOUTH AFRICA",
	COUNTRY_CODE_SOUTH_GEORGIA_AND_SOUTH_SANDWICH_ISLANDS: "SOUTH GEORGIA AND THE SOUTH SANDWICH ISLANDS",
	COUNTRY_CODE_SOUTH_SUDAN:                              "SOUTH SUDAN",
	COUNTRY_CODE_SPAIN:                                    "SPAIN",
	COUNTRY_CODE_SRI_LANKA:                                "SRI LANKA",
	COUNTRY_CODE_SUDAN:                                    "SUDAN",
	COUNTRY_CODE_SURINAME:                                 "SURINAME",
	COUNTRY_CODE_SVALBARD_AND_JAN_MAYEN:                   "SVALBARD AND JAN MAYEN",
	COUNTRY_CODE_SWEDEN:                                   "SWEDEN",
	COUNTRY_CODE_SWITZERLAND:                              "SWITZERLAND",
	COUNTRY_CODE_SYRIAN_ARAB_REPUBLIC:                     "SYRIAN ARAB REPUBLIC",
	COUNTRY_CODE_TAIWAN:                                   "TAIWAN, PROVINCE OF CHINA",
	COUNTRY_CODE_TAJIKISTAN:                               "TAJIKISTAN",
	COUNTRY_CODE_TANZANIA:                                 "TANZANIA, UNITED REPUBLIC OF",
	COUNTRY_CODE_THAILAND:                                 "THAILAND",
	COUNTRY_CODE_TIMOR_LESTE:                              "TIMOR-LESTE",
	COUNTRY_CODE_TOGO:                                     "TOGO",
	COUNTRY_CODE_TOKELAU:                                  "TOKELAU",
	COUNTRY_CODE_TONGA:                                    "TONGA",
	COUNTRY_CODE_TRINIDAD_AND_TOBAGO:                      "TRINIDAD AND TOBAGO",
	COUNTRY_CODE_TUNISIA:                                  "TUNISIA",
	COUNTRY_CODE_TURKIYE:                                  "TÜRKIYE",
	COUNTRY_CODE_TURKMENISTAN:                             "TURKMENISTAN",
	COUNTRY_CODE_TURKS_AND_CAICOS_ISLANDS:                 "TURKS AND CAICOS ISLANDS",
	COUNTRY_CODE_TUVALU:                                   "TUVALU",
	COUNTRY_CODE_UGANDA:                                   "UGANDA",
	COUNTRY_CODE_UKRAINE:                                  "UKRAINE",
	COUNTRY_CODE_UNITED_ARAB_EMIRATES:                     "UNITED ARAB EMIRATES",
	COUNTRY_CODE_UNITED_KINGDOM:                           "UNITED KINGDOM OF GREAT BRITAIN AND NORTHERN IRELAND",
	COUNTRY_CODE_UNITED_STATES:                            "UNITED STATES OF AMERICA",
	COUNTRY_CODE_UNITED_STATES_MINOR_OUTLYING_ISLANDS:     "UNITED STATES MINOR OUTLYING ISLANDS",
	COUNTRY_CODE_URUGUAY:                                  "URUGUAY",
	COUNTRY_CODE_UZBEKISTAN:                               "UZBEKISTAN",
	COUNTRY_CODE_VANUATU:                                  "VANUATU",
	COUNTRY_CODE_VENEZUELA:                                "VENEZUELA, BOLIVARIAN REPUBLIC OF",
	COUNTRY_CODE_VIET_NAM:                                 "VIET NAM",
	COUNTRY_CODE_VIRGIN_ISLANDS_BRITISH:                   "VIRGIN ISLANDS, BRITISH",
	COUNTRY_CODE_VIRGIN_ISLANDS_US:                        "VIRGIN ISLANDS, U.S.",
	COUNTRY_CODE_WALLIS_AND_FUTUNA:                        "WALLIS AND FUTUNA",
	COUNTRY_CODE_WESTERN_SAHARA:                           "WESTERN SAHARA",
	COUNTRY_CODE_YEMEN:                                    "YEMEN",
	COUNTRY_CODE_ZAMBIA:                                   "ZAMBIA",
	COUNTRY_CODE_ZIMBABWE:                                 "ZIMBABWE",
}

// IsValidCountryCode checks if a given country code is valid according to ISO 3166-1
func IsValidCountryCode(code string) bool {
	_, exists := CountryNames[code]
	return exists
}

// GetCountryName returns the full name for a given country code
func GetCountryName(code string) (string, bool) {
	name, exists := CountryNames[code]
	return name, exists
}
