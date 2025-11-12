package terminology

// Language constants from ISO 639-1 and RFC-5646
// Based on OpenEHR Support Terminology specification
// Used in: ENTRY.language, COMPOSITION.language, DV_ENCAPSULATED.language
const (
	LANG_TERMINOLOGY_ID_ISO string = "ISO_639-1"
	LANG_TERMINOLOGY_ID_RFC string = "RFC_5646"
)

var LanguageTerminologyIDs = map[string]string{
	LANG_TERMINOLOGY_ID_ISO: "ISO 639-1",
	LANG_TERMINOLOGY_ID_RFC: "RFC 5646",
}

func IsValidLanguageTerminologyID(id string) bool {
	_, exists := LanguageTerminologyIDs[id]
	return exists
}

const (
	LANG_CODE_AA    string = "aa"    // Afar
	LANG_CODE_AF    string = "af"    // Afrikaans
	LANG_CODE_AK    string = "ak"    // Akan
	LANG_CODE_SQ    string = "sq"    // Albanian
	LANG_CODE_AM    string = "am"    // Amharic
	LANG_CODE_AR    string = "ar"    // Arabic
	LANG_CODE_AR_SA string = "ar-sa" // Arabic (Saudi Arabia)
	LANG_CODE_AR_IQ string = "ar-iq" // Arabic (Iraq)
	LANG_CODE_AR_EG string = "ar-eg" // Arabic (Egypt)
	LANG_CODE_AR_LY string = "ar-ly" // Arabic (Libya)
	LANG_CODE_AR_DZ string = "ar-dz" // Arabic (Algeria)
	LANG_CODE_AR_MA string = "ar-ma" // Arabic (Morocco)
	LANG_CODE_AR_TN string = "ar-tn" // Arabic (Tunisia)
	LANG_CODE_AR_OM string = "ar-om" // Arabic (Oman)
	LANG_CODE_AR_YE string = "ar-ye" // Arabic (Yemen)
	LANG_CODE_AR_SY string = "ar-sy" // Arabic (Syria)
	LANG_CODE_AR_JO string = "ar-jo" // Arabic (Jordan)
	LANG_CODE_AR_LB string = "ar-lb" // Arabic (Lebanon)
	LANG_CODE_AR_KW string = "ar-kw" // Arabic (Kuwait)
	LANG_CODE_AR_AE string = "ar-ae" // Arabic (U.A.E.)
	LANG_CODE_AR_BH string = "ar-bh" // Arabic (Bahrain)
	LANG_CODE_AR_QA string = "ar-qa" // Arabic (Qatar)
	LANG_CODE_AN    string = "an"    // Aragonese
	LANG_CODE_HY    string = "hy"    // Armenian
	LANG_CODE_AS    string = "as"    // Assamese
	LANG_CODE_AV    string = "av"    // Avaric, Avar
	LANG_CODE_AY    string = "ay"    // Aymara
	LANG_CODE_AZ    string = "az"    // Azerbaijani, Azeri
	LANG_CODE_BM    string = "bm"    // Bambara
	LANG_CODE_BA    string = "ba"    // Bashkir
	LANG_CODE_EU    string = "eu"    // Basque
	LANG_CODE_BE    string = "be"    // Belarusian
	LANG_CODE_BN    string = "bn"    // Bengali, Bangla
	LANG_CODE_BI    string = "bi"    // Bislama
	LANG_CODE_BS    string = "bs"    // Bosnian
	LANG_CODE_BR    string = "br"    // Breton
	LANG_CODE_BG    string = "bg"    // Bulgarian
	LANG_CODE_MY    string = "my"    // Burmese, Myanmar
	LANG_CODE_CA    string = "ca"    // Catalan, Valencian
	LANG_CODE_CH    string = "ch"    // Chamorro
	LANG_CODE_CE    string = "ce"    // Chechen
	LANG_CODE_NY    string = "ny"    // Chichewa, Chewa, Nyanja
	LANG_CODE_ZH    string = "zh"    // Chinese
	LANG_CODE_ZH_TW string = "zh-tw" // Chinese (Taiwan)
	LANG_CODE_ZH_CN string = "zh-cn" // Chinese (PRC)
	LANG_CODE_ZH_HK string = "zh-hk" // Chinese (Hong Kong SAR)
	LANG_CODE_ZH_SG string = "zh-sg" // Chinese (Singapore)
	LANG_CODE_ZH_MO string = "zh-mo" // Chinese (Macau)
	LANG_CODE_CV    string = "cv"    // Chuvash
	LANG_CODE_KW    string = "kw"    // Cornish
	LANG_CODE_CO    string = "co"    // Corsican
	LANG_CODE_CR    string = "cr"    // Cree
	LANG_CODE_HR    string = "hr"    // Croatian
	LANG_CODE_HR_BA string = "hr-ba" // Croatian (Bosnia and Herzegovina)
	LANG_CODE_CS    string = "cs"    // Czech
	LANG_CODE_DA    string = "da"    // Danish
	LANG_CODE_DV    string = "dv"    // Divehi, Dhivehi, Maldivian
	LANG_CODE_NL    string = "nl"    // Dutch
	LANG_CODE_NL_BE string = "nl-be" // Dutch (Belgium)
	LANG_CODE_DZ    string = "dz"    // Dzongkha
	LANG_CODE_EN    string = "en"    // English
	LANG_CODE_EN_US string = "en-us" // English (United States)
	LANG_CODE_EN_GB string = "en-gb" // English (United Kingdom)
	LANG_CODE_EN_AU string = "en-au" // English (Australia)
	LANG_CODE_EN_CA string = "en-ca" // English (Canada)
	LANG_CODE_EN_NZ string = "en-nz" // English (New Zealand)
	LANG_CODE_EN_IE string = "en-ie" // English (Ireland)
	LANG_CODE_EN_ZA string = "en-za" // English (South Africa)
	LANG_CODE_EN_JM string = "en-jm" // English (Jamaica)
	LANG_CODE_EN_CB string = "en-cb" // English (Caribbean)
	LANG_CODE_EN_BZ string = "en-bz" // English (Belize)
	LANG_CODE_EN_TT string = "en-tt" // English (Trinidad and Tobago)
	LANG_CODE_EN_PH string = "en-ph" // English (Republic of the Philippines)
	LANG_CODE_EN_ZW string = "en-zw" // English (Zimbabwe)
	LANG_CODE_EO    string = "eo"    // Esperanto
	LANG_CODE_ET    string = "et"    // Estonian
	LANG_CODE_EE    string = "ee"    // Ewe
	LANG_CODE_FO    string = "fo"    // Faroese
	LANG_CODE_FJ    string = "fj"    // Fijian
	LANG_CODE_FI    string = "fi"    // Finnish
	LANG_CODE_FR    string = "fr"    // French
	LANG_CODE_FR_BE string = "fr-be" // French (Belgium)
	LANG_CODE_FR_CA string = "fr-ca" // French (Canada)
	LANG_CODE_FR_CH string = "fr-ch" // French (Switzerland)
	LANG_CODE_FR_LU string = "fr-lu" // French (Luxembourg)
	LANG_CODE_FR_MC string = "fr-mc" // French (Principality of Monaco)
	LANG_CODE_FY    string = "fy"    // Frisian, Western Frisian
	LANG_CODE_FF    string = "ff"    // Fulah, Fulani
	LANG_CODE_GD    string = "gd"    // Gaelic, Scottish Gaelic
	LANG_CODE_GD_IE string = "gd-ie" // Gaelic (Ireland)
	LANG_CODE_GL    string = "gl"    // Galician
	LANG_CODE_LG    string = "lg"    // Ganda
	LANG_CODE_KA    string = "ka"    // Georgian
	LANG_CODE_DE    string = "de"    // German
	LANG_CODE_DE_CH string = "de-ch" // German (Switzerland)
	LANG_CODE_DE_AT string = "de-at" // German (Austria)
	LANG_CODE_DE_LU string = "de-lu" // German (Luxembourg)
	LANG_CODE_DE_LI string = "de-li" // German (Liechtenstein)
	LANG_CODE_EL    string = "el"    // Greek
	LANG_CODE_KL    string = "kl"    // Kalaallisut, Greenlandic
	LANG_CODE_GN    string = "gn"    // Guarani
	LANG_CODE_GU    string = "gu"    // Gujarati
	LANG_CODE_HT    string = "ht"    // Haitian, Haitian Creole
	LANG_CODE_HA    string = "ha"    // Hausa
	LANG_CODE_HE    string = "he"    // Hebrew
	LANG_CODE_HZ    string = "hz"    // Herero
	LANG_CODE_HI    string = "hi"    // Hindi
	LANG_CODE_HO    string = "ho"    // Hiri Motu, Pidgin Motu
	LANG_CODE_HU    string = "hu"    // Hungarian
	LANG_CODE_IS    string = "is"    // Icelandic
	LANG_CODE_IG    string = "ig"    // Igbo
	LANG_CODE_ID    string = "id"    // Indonesian
	LANG_CODE_IU    string = "iu"    // Inuktitut
	LANG_CODE_IK    string = "ik"    // Inupiaq
	LANG_CODE_GA    string = "ga"    // Irish
	LANG_CODE_IT    string = "it"    // Italian
	LANG_CODE_IT_CH string = "it-ch" // Italian (Switzerland)
	LANG_CODE_JA    string = "ja"    // Japanese
	LANG_CODE_JV    string = "jv"    // Javanese
	LANG_CODE_KN    string = "kn"    // Kannada
	LANG_CODE_KR    string = "kr"    // Kanuri
	LANG_CODE_KS    string = "ks"    // Kashmiri
	LANG_CODE_KK    string = "kk"    // Kazakh
	LANG_CODE_KM    string = "km"    // Central Khmer, Cambodian
	LANG_CODE_KI    string = "ki"    // Kikuyu, Gikuyu
	LANG_CODE_RW    string = "rw"    // Kinyarwanda
	LANG_CODE_KY    string = "ky"    // Kirghiz, Kyrgyz
	LANG_CODE_KV    string = "kv"    // Komi
	LANG_CODE_KG    string = "kg"    // Kongo
	LANG_CODE_KO    string = "ko"    // Korean
	LANG_CODE_KJ    string = "kj"    // Kuanyama, Kwanyama
	LANG_CODE_KU    string = "ku"    // Kurdish
	LANG_CODE_LO    string = "lo"    // Lao
	LANG_CODE_LA    string = "la"    // Latin
	LANG_CODE_LV    string = "lv"    // Latvian
	LANG_CODE_LI    string = "li"    // Limburgan, Limburger, Limburgish
	LANG_CODE_LN    string = "ln"    // Lingala
	LANG_CODE_LT    string = "lt"    // Lithuanian
	LANG_CODE_LU    string = "lu"    // Luba-Katanga, Luba-Shaba
	LANG_CODE_LB    string = "lb"    // Luxembourgish, Letzeburgesch
	LANG_CODE_MK    string = "mk"    // Macedonian
	LANG_CODE_MG    string = "mg"    // Malagasy
	LANG_CODE_MS    string = "ms"    // Malay
	LANG_CODE_ML    string = "ml"    // Malayalam
	LANG_CODE_MT    string = "mt"    // Maltese
	LANG_CODE_GV    string = "gv"    // Manx
	LANG_CODE_MI    string = "mi"    // Maori
	LANG_CODE_MR    string = "mr"    // Marathi
	LANG_CODE_MH    string = "mh"    // Marshallese
	LANG_CODE_MN    string = "mn"    // Mongolian
	LANG_CODE_NA    string = "na"    // Nauru, Nauruan
	LANG_CODE_NV    string = "nv"    // Navajo, Navaho
	LANG_CODE_ND    string = "nd"    // North Ndebele
	LANG_CODE_NR    string = "nr"    // South Ndebele
	LANG_CODE_NG    string = "ng"    // Ndonga
	LANG_CODE_NE    string = "ne"    // Nepali
	LANG_CODE_NB    string = "nb"    // Norwegian Bokmal
	LANG_CODE_NN    string = "nn"    // Norwegian Nynorsk
	LANG_CODE_II    string = "ii"    // Sichuan Yi, Nuosu, Northern Yi, Liangshan Yi
	LANG_CODE_OC    string = "oc"    // Occitan
	LANG_CODE_OJ    string = "oj"    // Ojibwa, Ojibwe
	LANG_CODE_OR    string = "or"    // Oriya, Odia
	LANG_CODE_OM    string = "om"    // Oromo
	LANG_CODE_OS    string = "os"    // Ossetian, Ossetic
	LANG_CODE_PS    string = "ps"    // Pashto, Pushto
	LANG_CODE_FA    string = "fa"    // Persian, Farsi
	LANG_CODE_PL    string = "pl"    // Polish
	LANG_CODE_PT    string = "pt"    // Portuguese
	LANG_CODE_PT_BR string = "pt-br" // Portuguese (Brazil)
	LANG_CODE_PT_PT string = "pt-pt" // Portuguese (Portugal) - DEPRECATED
	LANG_CODE_PA    string = "pa"    // Punjabi, Panjabi
	LANG_CODE_QU    string = "qu"    // Quechua
	LANG_CODE_QU_BO string = "qu-bo" // Quechua (Bolivia)
	LANG_CODE_QU_EC string = "qu-ec" // Quechua (Ecuador)
	LANG_CODE_QU_PE string = "qu-pe" // Quechua (Peru)
	LANG_CODE_RO    string = "ro"    // Romanian
	LANG_CODE_RO_MO string = "ro-mo" // Romanian (Moldavia)
	LANG_CODE_RM    string = "rm"    // Romansh, Rhaeto-Romanic
	LANG_CODE_RN    string = "rn"    // Rundi, Kirundi
	LANG_CODE_RU    string = "ru"    // Russian
	LANG_CODE_RU_MO string = "ru-mo" // Russian (Moldavia)
	LANG_CODE_SE    string = "se"    // Northern Sami
	LANG_CODE_SZ    string = "sz"    // Sami (Lappish) - DEPRECATED
	LANG_CODE_SM    string = "sm"    // Samoan
	LANG_CODE_SG    string = "sg"    // Sango
	LANG_CODE_SC    string = "sc"    // Sardinian
	LANG_CODE_SR    string = "sr"    // Serbian
	LANG_CODE_SR_BA string = "sr-ba" // Serbian (Bosnia and Herzegovina)
	LANG_CODE_SB    string = "sb"    // Serbian - DEPRECATED
	LANG_CODE_SN    string = "sn"    // Shona
	LANG_CODE_SD    string = "sd"    // Sindhi
	LANG_CODE_SI    string = "si"    // Sinhala, Sinhalese
	LANG_CODE_SK    string = "sk"    // Slovak
	LANG_CODE_SL    string = "sl"    // Slovenian, Slovene
	LANG_CODE_SO    string = "so"    // Somali
	LANG_CODE_ST    string = "st"    // Southern Sotho, Sesotho, Sutu
	LANG_CODE_ES    string = "es"    // Spanish, Castilian
	LANG_CODE_ES_MX string = "es-mx" // Spanish (Mexico)
	LANG_CODE_ES_GT string = "es-gt" // Spanish (Guatemala)
	LANG_CODE_ES_CR string = "es-cr" // Spanish (Costa Rica)
	LANG_CODE_ES_PA string = "es-pa" // Spanish (Panama)
	LANG_CODE_ES_DO string = "es-do" // Spanish (Dominican Republic)
	LANG_CODE_ES_VE string = "es-ve" // Spanish (Venezuela)
	LANG_CODE_ES_CO string = "es-co" // Spanish (Colombia)
	LANG_CODE_ES_PE string = "es-pe" // Spanish (Peru)
	LANG_CODE_ES_AR string = "es-ar" // Spanish (Argentina)
	LANG_CODE_ES_EC string = "es-ec" // Spanish (Ecuador)
	LANG_CODE_ES_CL string = "es-cl" // Spanish (Chile)
	LANG_CODE_ES_UY string = "es-uy" // Spanish (Uruguay)
	LANG_CODE_ES_PY string = "es-py" // Spanish (Paraguay)
	LANG_CODE_ES_BO string = "es-bo" // Spanish (Bolivia)
	LANG_CODE_ES_SV string = "es-sv" // Spanish (El Salvador)
	LANG_CODE_ES_HN string = "es-hn" // Spanish (Honduras)
	LANG_CODE_ES_NI string = "es-ni" // Spanish (Nicaragua)
	LANG_CODE_ES_PR string = "es-pr" // Spanish (Puerto Rico)
	LANG_CODE_SU    string = "su"    // Sundanese
	LANG_CODE_SX    string = "sx"    // Sutu - DEPRECATED
	LANG_CODE_SW    string = "sw"    // Swahili
	LANG_CODE_SS    string = "ss"    // Swati, Swazi
	LANG_CODE_SV    string = "sv"    // Swedish
	LANG_CODE_SV_FI string = "sv-fi" // Swedish (Finland)
	LANG_CODE_TL    string = "tl"    // Tagalog
	LANG_CODE_TY    string = "ty"    // Tahitian
	LANG_CODE_TG    string = "tg"    // Tajik
	LANG_CODE_TA    string = "ta"    // Tamil
	LANG_CODE_TT    string = "tt"    // Tatar
	LANG_CODE_TE    string = "te"    // Telugu
	LANG_CODE_TH    string = "th"    // Thai
	LANG_CODE_BO    string = "bo"    // Tibetan
	LANG_CODE_TI    string = "ti"    // Tigrinya
	LANG_CODE_TO    string = "to"    // Tonga, Tongan
	LANG_CODE_TS    string = "ts"    // Tsonga
	LANG_CODE_TN    string = "tn"    // Tswana
	LANG_CODE_TR    string = "tr"    // Turkish
	LANG_CODE_TK    string = "tk"    // Turkmen
	LANG_CODE_TW    string = "tw"    // Twi
	LANG_CODE_UG    string = "ug"    // Uighur, Uyghur
	LANG_CODE_UK    string = "uk"    // Ukrainian
	LANG_CODE_UR    string = "ur"    // Urdu
	LANG_CODE_UZ    string = "uz"    // Uzbek
	LANG_CODE_VE    string = "ve"    // Venda
	LANG_CODE_VI    string = "vi"    // Vietnamese
	LANG_CODE_WA    string = "wa"    // Walloon
	LANG_CODE_CY    string = "cy"    // Welsh
	LANG_CODE_CY_GB string = "cy-gb" // Welsh (United Kingdom)
	LANG_CODE_CY_AR string = "cy-ar" // Welsh (Argentina)
	LANG_CODE_WO    string = "wo"    // Wolof
	LANG_CODE_XH    string = "xh"    // Xhosa
	LANG_CODE_YI    string = "yi"    // Yiddish
	LANG_CODE_JI    string = "ji"    // Yiddish - DEPRECATED
	LANG_CODE_YO    string = "yo"    // Yoruba
	LANG_CODE_ZA    string = "za"    // Zhuang, Chuang
	LANG_CODE_ZU    string = "zu"    // Zulu
)

// LanguageNames provides a mapping of language codes to their names
var LanguageNames = map[string]string{
	LANG_CODE_AA:    "Afar",
	LANG_CODE_AF:    "Afrikaans",
	LANG_CODE_AK:    "Akan",
	LANG_CODE_SQ:    "Albanian",
	LANG_CODE_AM:    "Amharic",
	LANG_CODE_AR:    "Arabic",
	LANG_CODE_AR_SA: "Arabic (Saudi Arabia)",
	LANG_CODE_AR_IQ: "Arabic (Iraq)",
	LANG_CODE_AR_EG: "Arabic (Egypt)",
	LANG_CODE_AR_LY: "Arabic (Libya)",
	LANG_CODE_AR_DZ: "Arabic (Algeria)",
	LANG_CODE_AR_MA: "Arabic (Morocco)",
	LANG_CODE_AR_TN: "Arabic (Tunisia)",
	LANG_CODE_AR_OM: "Arabic (Oman)",
	LANG_CODE_AR_YE: "Arabic (Yemen)",
	LANG_CODE_AR_SY: "Arabic (Syria)",
	LANG_CODE_AR_JO: "Arabic (Jordan)",
	LANG_CODE_AR_LB: "Arabic (Lebanon)",
	LANG_CODE_AR_KW: "Arabic (Kuwait)",
	LANG_CODE_AR_AE: "Arabic (U.A.E.)",
	LANG_CODE_AR_BH: "Arabic (Bahrain)",
	LANG_CODE_AR_QA: "Arabic (Qatar)",
	LANG_CODE_AN:    "Aragonese",
	LANG_CODE_HY:    "Armenian",
	LANG_CODE_AS:    "Assamese",
	LANG_CODE_AV:    "Avaric, Avar",
	LANG_CODE_AY:    "Aymara",
	LANG_CODE_AZ:    "Azerbaijani, Azeri",
	LANG_CODE_BM:    "Bambara",
	LANG_CODE_BA:    "Bashkir",
	LANG_CODE_EU:    "Basque",
	LANG_CODE_BE:    "Belarusian",
	LANG_CODE_BN:    "Bengali, Bangla",
	LANG_CODE_BI:    "Bislama",
	LANG_CODE_BS:    "Bosnian",
	LANG_CODE_BR:    "Breton",
	LANG_CODE_BG:    "Bulgarian",
	LANG_CODE_MY:    "Burmese, Myanmar",
	LANG_CODE_CA:    "Catalan, Valencian",
	LANG_CODE_CH:    "Chamorro",
	LANG_CODE_CE:    "Chechen",
	LANG_CODE_NY:    "Chichewa, Chewa, Nyanja",
	LANG_CODE_ZH:    "Chinese",
	LANG_CODE_ZH_TW: "Chinese (Taiwan)",
	LANG_CODE_ZH_CN: "Chinese (PRC)",
	LANG_CODE_ZH_HK: "Chinese (Hong Kong SAR)",
	LANG_CODE_ZH_SG: "Chinese (Singapore)",
	LANG_CODE_ZH_MO: "Chinese (Macau)",
	LANG_CODE_CV:    "Chuvash",
	LANG_CODE_KW:    "Cornish",
	LANG_CODE_CO:    "Corsican",
	LANG_CODE_CR:    "Cree",
	LANG_CODE_HR:    "Croatian",
	LANG_CODE_HR_BA: "Croatian (Bosnia and Herzegovina)",
	LANG_CODE_CS:    "Czech",
	LANG_CODE_DA:    "Danish",
	LANG_CODE_DV:    "Divehi, Dhivehi, Maldivian",
	LANG_CODE_NL:    "Dutch",
	LANG_CODE_NL_BE: "Dutch (Belgium)",
	LANG_CODE_DZ:    "Dzongkha",
	LANG_CODE_EN:    "English",
	LANG_CODE_EN_US: "English (United States)",
	LANG_CODE_EN_GB: "English (United Kingdom)",
	LANG_CODE_EN_AU: "English (Australia)",
	LANG_CODE_EN_CA: "English (Canada)",
	LANG_CODE_EN_NZ: "English (New Zealand)",
	LANG_CODE_EN_IE: "English (Ireland)",
	LANG_CODE_EN_ZA: "English (South Africa)",
	LANG_CODE_EN_JM: "English (Jamaica)",
	LANG_CODE_EN_CB: "English (Caribbean)",
	LANG_CODE_EN_BZ: "English (Belize)",
	LANG_CODE_EN_TT: "English (Trinidad and Tobago)",
	LANG_CODE_EN_PH: "English (Republic of the Philippines)",
	LANG_CODE_EN_ZW: "English (Zimbabwe)",
	LANG_CODE_EO:    "Esperanto",
	LANG_CODE_ET:    "Estonian",
	LANG_CODE_EE:    "Ewe",
	LANG_CODE_FO:    "Faroese",
	LANG_CODE_FJ:    "Fijian",
	LANG_CODE_FI:    "Finnish",
	LANG_CODE_FR:    "French",
	LANG_CODE_FR_BE: "French (Belgium)",
	LANG_CODE_FR_CA: "French (Canada)",
	LANG_CODE_FR_CH: "French (Switzerland)",
	LANG_CODE_FR_LU: "French (Luxembourg)",
	LANG_CODE_FR_MC: "French (Principality of Monaco)",
	LANG_CODE_FY:    "Frisian, Western Frisian",
	LANG_CODE_FF:    "Fulah, Fulani",
	LANG_CODE_GD:    "Gaelic, Scottish Gaelic",
	LANG_CODE_GD_IE: "Gaelic (Ireland)",
	LANG_CODE_GL:    "Galician",
	LANG_CODE_LG:    "Ganda",
	LANG_CODE_KA:    "Georgian",
	LANG_CODE_DE:    "German",
	LANG_CODE_DE_CH: "German (Switzerland)",
	LANG_CODE_DE_AT: "German (Austria)",
	LANG_CODE_DE_LU: "German (Luxembourg)",
	LANG_CODE_DE_LI: "German (Liechtenstein)",
	LANG_CODE_EL:    "Greek",
	LANG_CODE_KL:    "Kalaallisut, Greenlandic",
	LANG_CODE_GN:    "Guarani",
	LANG_CODE_GU:    "Gujarati",
	LANG_CODE_HT:    "Haitian, Haitian Creole",
	LANG_CODE_HA:    "Hausa",
	LANG_CODE_HE:    "Hebrew",
	LANG_CODE_HZ:    "Herero",
	LANG_CODE_HI:    "Hindi",
	LANG_CODE_HO:    "Hiri Motu, Pidgin Motu",
	LANG_CODE_HU:    "Hungarian",
	LANG_CODE_IS:    "Icelandic",
	LANG_CODE_IG:    "Igbo",
	LANG_CODE_ID:    "Indonesian",
	LANG_CODE_IU:    "Inuktitut",
	LANG_CODE_IK:    "Inupiaq",
	LANG_CODE_GA:    "Irish",
	LANG_CODE_IT:    "Italian",
	LANG_CODE_IT_CH: "Italian (Switzerland)",
	LANG_CODE_JA:    "Japanese",
	LANG_CODE_JV:    "Javanese",
	LANG_CODE_KN:    "Kannada",
	LANG_CODE_KR:    "Kanuri",
	LANG_CODE_KS:    "Kashmiri",
	LANG_CODE_KK:    "Kazakh",
	LANG_CODE_KM:    "Central Khmer, Cambodian",
	LANG_CODE_KI:    "Kikuyu, Gikuyu",
	LANG_CODE_RW:    "Kinyarwanda",
	LANG_CODE_KY:    "Kirghiz, Kyrgyz",
	LANG_CODE_KV:    "Komi",
	LANG_CODE_KG:    "Kongo",
	LANG_CODE_KO:    "Korean",
	LANG_CODE_KJ:    "Kuanyama, Kwanyama",
	LANG_CODE_KU:    "Kurdish",
	LANG_CODE_LO:    "Lao",
	LANG_CODE_LA:    "Latin",
	LANG_CODE_LV:    "Latvian",
	LANG_CODE_LI:    "Limburgan, Limburger, Limburgish",
	LANG_CODE_LN:    "Lingala",
	LANG_CODE_LT:    "Lithuanian",
	LANG_CODE_LU:    "Luba-Katanga, Luba-Shaba",
	LANG_CODE_LB:    "Luxembourgish, Letzeburgesch",
	LANG_CODE_MK:    "Macedonian",
	LANG_CODE_MG:    "Malagasy",
	LANG_CODE_MS:    "Malay",
	LANG_CODE_ML:    "Malayalam",
	LANG_CODE_MT:    "Maltese",
	LANG_CODE_GV:    "Manx",
	LANG_CODE_MI:    "Maori",
	LANG_CODE_MR:    "Marathi",
	LANG_CODE_MH:    "Marshallese",
	LANG_CODE_MN:    "Mongolian",
	LANG_CODE_NA:    "Nauru, Nauruan",
	LANG_CODE_NV:    "Navajo, Navaho",
	LANG_CODE_ND:    "North Ndebele",
	LANG_CODE_NR:    "South Ndebele",
	LANG_CODE_NG:    "Ndonga",
	LANG_CODE_NE:    "Nepali",
	LANG_CODE_NB:    "Norwegian Bokmal",
	LANG_CODE_NN:    "Norwegian Nynorsk",
	LANG_CODE_II:    "Sichuan Yi, Nuosu, Northern Yi, Liangshan Yi",
	LANG_CODE_OC:    "Occitan",
	LANG_CODE_OJ:    "Ojibwa, Ojibwe",
	LANG_CODE_OR:    "Oriya, Odia",
	LANG_CODE_OM:    "Oromo",
	LANG_CODE_OS:    "Ossetian, Ossetic",
	LANG_CODE_PS:    "Pashto, Pushto",
	LANG_CODE_FA:    "Persian, Farsi",
	LANG_CODE_PL:    "Polish",
	LANG_CODE_PT:    "Portuguese",
	LANG_CODE_PT_BR: "Portuguese (Brazil)",
	LANG_CODE_PT_PT: "Portuguese (Portugal) - DEPRECATED",
	LANG_CODE_PA:    "Punjabi, Panjabi",
	LANG_CODE_QU:    "Quechua",
	LANG_CODE_QU_BO: "Quechua (Bolivia)",
	LANG_CODE_QU_EC: "Quechua (Ecuador)",
	LANG_CODE_QU_PE: "Quechua (Peru)",
	LANG_CODE_RO:    "Romanian",
	LANG_CODE_RO_MO: "Romanian (Moldavia)",
	LANG_CODE_RM:    "Romansh, Rhaeto-Romanic",
	LANG_CODE_RN:    "Rundi, Kirundi",
	LANG_CODE_RU:    "Russian",
	LANG_CODE_RU_MO: "Russian (Moldavia)",
	LANG_CODE_SE:    "Northern Sami",
	LANG_CODE_SZ:    "Sami (Lappish) - DEPRECATED",
	LANG_CODE_SM:    "Samoan",
	LANG_CODE_SG:    "Sango",
	LANG_CODE_SC:    "Sardinian",
	LANG_CODE_SR:    "Serbian",
	LANG_CODE_SR_BA: "Serbian (Bosnia and Herzegovina)",
	LANG_CODE_SB:    "Serbian - DEPRECATED",
	LANG_CODE_SN:    "Shona",
	LANG_CODE_SD:    "Sindhi",
	LANG_CODE_SI:    "Sinhala, Sinhalese",
	LANG_CODE_SK:    "Slovak",
	LANG_CODE_SL:    "Slovenian, Slovene",
	LANG_CODE_SO:    "Somali",
	LANG_CODE_ST:    "Southern Sotho, Sesotho, Sutu",
	LANG_CODE_ES:    "Spanish, Castilian",
	LANG_CODE_ES_MX: "Spanish (Mexico)",
	LANG_CODE_ES_GT: "Spanish (Guatemala)",
	LANG_CODE_ES_CR: "Spanish (Costa Rica)",
	LANG_CODE_ES_PA: "Spanish (Panama)",
	LANG_CODE_ES_DO: "Spanish (Dominican Republic)",
	LANG_CODE_ES_VE: "Spanish (Venezuela)",
	LANG_CODE_ES_CO: "Spanish (Colombia)",
	LANG_CODE_ES_PE: "Spanish (Peru)",
	LANG_CODE_ES_AR: "Spanish (Argentina)",
	LANG_CODE_ES_EC: "Spanish (Ecuador)",
	LANG_CODE_ES_CL: "Spanish (Chile)",
	LANG_CODE_ES_UY: "Spanish (Uruguay)",
	LANG_CODE_ES_PY: "Spanish (Paraguay)",
	LANG_CODE_ES_BO: "Spanish (Bolivia)",
	LANG_CODE_ES_SV: "Spanish (El Salvador)",
	LANG_CODE_ES_HN: "Spanish (Honduras)",
	LANG_CODE_ES_NI: "Spanish (Nicaragua)",
	LANG_CODE_ES_PR: "Spanish (Puerto Rico)",
	LANG_CODE_SU:    "Sundanese",
	LANG_CODE_SX:    "Sutu - DEPRECATED",
	LANG_CODE_SW:    "Swahili",
	LANG_CODE_SS:    "Swati, Swazi",
	LANG_CODE_SV:    "Swedish",
	LANG_CODE_SV_FI: "Swedish (Finland)",
	LANG_CODE_TL:    "Tagalog",
	LANG_CODE_TY:    "Tahitian",
	LANG_CODE_TG:    "Tajik",
	LANG_CODE_TA:    "Tamil",
	LANG_CODE_TT:    "Tatar",
	LANG_CODE_TE:    "Telugu",
	LANG_CODE_TH:    "Thai",
	LANG_CODE_BO:    "Tibetan",
	LANG_CODE_TI:    "Tigrinya",
	LANG_CODE_TO:    "Tonga, Tongan",
	LANG_CODE_TS:    "Tsonga",
	LANG_CODE_TN:    "Tswana",
	LANG_CODE_TR:    "Turkish",
	LANG_CODE_TK:    "Turkmen",
	LANG_CODE_TW:    "Twi",
	LANG_CODE_UG:    "Uighur, Uyghur",
	LANG_CODE_UK:    "Ukrainian",
	LANG_CODE_UR:    "Urdu",
	LANG_CODE_UZ:    "Uzbek",
	LANG_CODE_VE:    "Venda",
	LANG_CODE_VI:    "Vietnamese",
	LANG_CODE_WA:    "Walloon",
	LANG_CODE_CY:    "Welsh",
	LANG_CODE_CY_GB: "Welsh (United Kingdom)",
	LANG_CODE_CY_AR: "Welsh (Argentina)",
	LANG_CODE_WO:    "Wolof",
	LANG_CODE_XH:    "Xhosa",
	LANG_CODE_YI:    "Yiddish",
	LANG_CODE_JI:    "Yiddish - DEPRECATED",
	LANG_CODE_YO:    "Yoruba",
	LANG_CODE_ZA:    "Zhuang, Chuang",
	LANG_CODE_ZU:    "Zulu",
}

// IsValidLanguageCode checks if the given string is a valid ISO 639-1 language code
func IsValidLanguageCode(code string) bool {
	_, exists := LanguageNames[code]
	return exists
}

// GetLanguageName returns the name of the language for the given code
func GetLanguageName(code string) (string, bool) {
	name, exists := LanguageNames[code]
	return name, exists
}
