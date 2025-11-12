package terminology

// Media Type constants from IANA Media Types
// Based on OpenEHR Support Terminology specification
// Used in: DV_MULTIMEDIA.media_type
const (
	// Audio Media Types
	MEDIA_AUDIO_DVI4            string = "audio/DVI4"
	MEDIA_AUDIO_G722            string = "audio/G722"
	MEDIA_AUDIO_G723            string = "audio/G723"
	MEDIA_AUDIO_G726_16         string = "audio/G726-16"
	MEDIA_AUDIO_G726_24         string = "audio/G726-24"
	MEDIA_AUDIO_G726_32         string = "audio/G726-32"
	MEDIA_AUDIO_G726_40         string = "audio/G726-40"
	MEDIA_AUDIO_G728            string = "audio/G728"
	MEDIA_AUDIO_L8              string = "audio/L8"
	MEDIA_AUDIO_L16             string = "audio/L16"
	MEDIA_AUDIO_LPC             string = "audio/LPC"
	MEDIA_AUDIO_G729            string = "audio/G729"
	MEDIA_AUDIO_G729D           string = "audio/G729D"
	MEDIA_AUDIO_G729E           string = "audio/G729E"
	MEDIA_AUDIO_BASIC           string = "audio/basic"
	MEDIA_AUDIO_MPEG            string = "audio/mpeg"
	MEDIA_AUDIO_MPEG3           string = "audio/mpeg3"
	MEDIA_AUDIO_MPEG4_GENERIC   string = "audio/mpeg4-generic"
	MEDIA_AUDIO_MP4             string = "audio/mp4"
	MEDIA_AUDIO_L20             string = "audio/L20"
	MEDIA_AUDIO_L24             string = "audio/L24"
	MEDIA_AUDIO_TELEPHONE_EVENT string = "audio/telephone-event"
	MEDIA_AUDIO_OGG             string = "audio/ogg"
	MEDIA_AUDIO_VORBIS          string = "audio/vorbis"

	// Video Media Types
	MEDIA_VIDEO_BT656     string = "video/BT656"
	MEDIA_VIDEO_CELB      string = "video/CelB"
	MEDIA_VIDEO_JPEG      string = "video/JPEG"
	MEDIA_VIDEO_H261      string = "video/H261"
	MEDIA_VIDEO_H263      string = "video/H263"
	MEDIA_VIDEO_H263_1998 string = "video/H263-1998"
	MEDIA_VIDEO_H263_2000 string = "video/H263-2000"
	MEDIA_VIDEO_H264      string = "video/H264"
	MEDIA_VIDEO_MPV       string = "video/MPV"
	MEDIA_VIDEO_MP4       string = "video/mp4"
	MEDIA_VIDEO_OGG       string = "video/ogg"
	MEDIA_VIDEO_MPEG      string = "video/mpeg"
	MEDIA_VIDEO_QUICKTIME string = "video/quicktime"

	// Text Media Types
	MEDIA_TEXT_CALENDAR                   string = "text/calendar"
	MEDIA_TEXT_DIRECTORY                  string = "text/directory"
	MEDIA_TEXT_HTML                       string = "text/html"
	MEDIA_TEXT_PLAIN                      string = "text/plain"
	MEDIA_TEXT_RICHTEXT                   string = "text/richtext"
	MEDIA_TEXT_RTF                        string = "text/rtf"
	MEDIA_TEXT_RFC822_HEADERS             string = "text/rfc822-headers"
	MEDIA_TEXT_SGML                       string = "text/sgml"
	MEDIA_TEXT_TAB_SEPARATED_VALUES       string = "text/tab-separated-values"
	MEDIA_TEXT_URI_LIST                   string = "text/uri-list"
	MEDIA_TEXT_XML                        string = "text/xml"
	MEDIA_TEXT_XML_EXTERNAL_PARSED_ENTITY string = "text/xml-external-parsed-entity"

	// Image Media Types
	MEDIA_IMAGE_AVIF    string = "image/avif"
	MEDIA_IMAGE_BMP     string = "image/bmp"
	MEDIA_IMAGE_CGM     string = "image/cgm"
	MEDIA_IMAGE_GIF     string = "image/gif"
	MEDIA_IMAGE_PNG     string = "image/png"
	MEDIA_IMAGE_TIFF    string = "image/tiff"
	MEDIA_IMAGE_JPEG    string = "image/jpeg"
	MEDIA_IMAGE_JP2     string = "image/jp2"
	MEDIA_IMAGE_SVG_XML string = "image/svg+xml"

	// Application Media Types
	MEDIA_APPLICATION_CDA_XML      string = "application/cda+xml"
	MEDIA_APPLICATION_EDIFACT      string = "application/EDIFACT"
	MEDIA_APPLICATION_FHIR_JSON    string = "application/fhir+json"
	MEDIA_APPLICATION_FHIR_XML     string = "application/fhir+xml"
	MEDIA_APPLICATION_HL7V2_XML    string = "application/hl7v2+xml"
	MEDIA_APPLICATION_GZIP         string = "application/gzip"
	MEDIA_APPLICATION_JSON         string = "application/json"
	MEDIA_APPLICATION_MSWORD       string = "application/msword"
	MEDIA_APPLICATION_PDF          string = "application/pdf"
	MEDIA_APPLICATION_RTF          string = "application/rtf"
	MEDIA_APPLICATION_DICOM        string = "application/dicom"
	MEDIA_APPLICATION_DICOM_JSON   string = "application/dicom+json"
	MEDIA_APPLICATION_DICOM_XML    string = "application/dicom+xml"
	MEDIA_APPLICATION_OCTET_STREAM string = "application/octet-stream"
	MEDIA_APPLICATION_OGG          string = "application/ogg"

	// OpenDocument Media Types
	MEDIA_APPLICATION_ODF_BASE                  string = "application/vnd.oasis.opendocument.base"
	MEDIA_APPLICATION_ODF_CHART                 string = "application/vnd.oasis.opendocument.chart"
	MEDIA_APPLICATION_ODF_CHART_TEMPLATE        string = "application/vnd.oasis.opendocument.chart-template"
	MEDIA_APPLICATION_ODF_FORMULA               string = "application/vnd.oasis.opendocument.formula"
	MEDIA_APPLICATION_ODF_FORMULA_TEMPLATE      string = "application/vnd.oasis.opendocument.formula-template"
	MEDIA_APPLICATION_ODF_GRAPHICS              string = "application/vnd.oasis.opendocument.graphics"
	MEDIA_APPLICATION_ODF_GRAPHICS_TEMPLATE     string = "application/vnd.oasis.opendocument.graphics-template"
	MEDIA_APPLICATION_ODF_IMAGE                 string = "application/vnd.oasis.opendocument.image"
	MEDIA_APPLICATION_ODF_IMAGE_TEMPLATE        string = "application/vnd.oasis.opendocument.image-template"
	MEDIA_APPLICATION_ODF_PRESENTATION          string = "application/vnd.oasis.opendocument.presentation"
	MEDIA_APPLICATION_ODF_PRESENTATION_TEMPLATE string = "application/vnd.oasis.opendocument.presentation-template"
	MEDIA_APPLICATION_ODF_SPREADSHEET           string = "application/vnd.oasis.opendocument.spreadsheet"
	MEDIA_APPLICATION_ODF_SPREADSHEET_TEMPLATE  string = "application/vnd.oasis.opendocument.spreadsheet-template"
	MEDIA_APPLICATION_ODF_TEXT                  string = "application/vnd.oasis.opendocument.text"
	MEDIA_APPLICATION_ODF_TEXT_MASTER           string = "application/vnd.oasis.opendocument.text-master"
	MEDIA_APPLICATION_ODF_TEXT_TEMPLATE         string = "application/vnd.oasis.opendocument.text-template"
	MEDIA_APPLICATION_ODF_TEXT_WEB              string = "application/vnd.oasis.opendocument.text-web"

	// Microsoft Office Media Types
	MEDIA_APPLICATION_MS_WORD_MACRO_ENABLED         string = "application/vnd.ms-word.document.macroEnabled.12"
	MEDIA_APPLICATION_MS_WORD_DOCX                  string = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	MEDIA_APPLICATION_MS_WORD_TEMPLATE_MACRO        string = "application/vnd.ms-word.template.macroEnabled.12"
	MEDIA_APPLICATION_MS_WORD_TEMPLATE_DOCX         string = "application/vnd.openxmlformats-officedocument.wordprocessingml.template"
	MEDIA_APPLICATION_MS_POWERPOINT_SLIDESHOW_MACRO string = "application/vnd.ms-powerpoint.slideshow.macroEnabled.12"
	MEDIA_APPLICATION_MS_POWERPOINT_SLIDESHOW       string = "application/vnd.openxmlformats-officedocument.presentationml.slideshow"
	MEDIA_APPLICATION_MS_POWERPOINT_MACRO           string = "application/vnd.ms-powerpoint.presentation.macroEnabled.12"
	MEDIA_APPLICATION_MS_POWERPOINT_PPTX            string = "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	MEDIA_APPLICATION_MS_EXCEL_BINARY_MACRO         string = "application/vnd.ms-excel.sheet.binary.macroEnabled.12"
	MEDIA_APPLICATION_MS_EXCEL_MACRO                string = "application/vnd.ms-excel.sheet.macroEnabled.12"
	MEDIA_APPLICATION_MS_EXCEL_XLSX                 string = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	MEDIA_APPLICATION_MS_XPS                        string = "application/vnd.ms-xpsdocument"
	MEDIA_APPLICATION_MS_EXCEL                      string = "application/vnd.ms-excel"
	MEDIA_APPLICATION_MS_OUTLOOK                    string = "application/vnd.ms-outlook"
	MEDIA_APPLICATION_MS_POWERPOINT                 string = "application/vnd.ms-powerpoint"

	// Archive Media Types
	MEDIA_APPLICATION_RAR string = "application/vnd.rar"
	MEDIA_APPLICATION_ZIP string = "application/zip"
)

// MediaTypeNames provides a mapping of media type codes to their descriptions
var MediaTypeNames = map[string]string{
	// Audio Media Types
	MEDIA_AUDIO_DVI4:            "Audio DVI4",
	MEDIA_AUDIO_G722:            "Audio G.722",
	MEDIA_AUDIO_G723:            "Audio G.723",
	MEDIA_AUDIO_G726_16:         "Audio G.726-16",
	MEDIA_AUDIO_G726_24:         "Audio G.726-24",
	MEDIA_AUDIO_G726_32:         "Audio G.726-32",
	MEDIA_AUDIO_G726_40:         "Audio G.726-40",
	MEDIA_AUDIO_G728:            "Audio G.728",
	MEDIA_AUDIO_L8:              "Audio L8",
	MEDIA_AUDIO_L16:             "Audio L16",
	MEDIA_AUDIO_LPC:             "Audio LPC",
	MEDIA_AUDIO_G729:            "Audio G.729",
	MEDIA_AUDIO_G729D:           "Audio G.729D",
	MEDIA_AUDIO_G729E:           "Audio G.729E",
	MEDIA_AUDIO_BASIC:           "Basic Audio",
	MEDIA_AUDIO_MPEG:            "MPEG Audio",
	MEDIA_AUDIO_MPEG3:           "MP3 Audio",
	MEDIA_AUDIO_MPEG4_GENERIC:   "MPEG-4 Generic Audio",
	MEDIA_AUDIO_MP4:             "MP4 Audio",
	MEDIA_AUDIO_L20:             "Audio L20",
	MEDIA_AUDIO_L24:             "Audio L24",
	MEDIA_AUDIO_TELEPHONE_EVENT: "Telephone Event Audio",
	MEDIA_AUDIO_OGG:             "Ogg Audio",
	MEDIA_AUDIO_VORBIS:          "Vorbis Audio",

	// Video Media Types
	MEDIA_VIDEO_BT656:     "BT.656 Video",
	MEDIA_VIDEO_CELB:      "CelB Video",
	MEDIA_VIDEO_JPEG:      "JPEG Video",
	MEDIA_VIDEO_H261:      "H.261 Video",
	MEDIA_VIDEO_H263:      "H.263 Video",
	MEDIA_VIDEO_H263_1998: "H.263-1998 Video",
	MEDIA_VIDEO_H263_2000: "H.263-2000 Video",
	MEDIA_VIDEO_H264:      "H.264 Video",
	MEDIA_VIDEO_MPV:       "MPV Video",
	MEDIA_VIDEO_MP4:       "MP4 Video",
	MEDIA_VIDEO_OGG:       "Ogg Video",
	MEDIA_VIDEO_MPEG:      "MPEG Video",
	MEDIA_VIDEO_QUICKTIME: "QuickTime Video",

	// Text Media Types
	MEDIA_TEXT_CALENDAR:                   "iCalendar",
	MEDIA_TEXT_DIRECTORY:                  "vCard Directory",
	MEDIA_TEXT_HTML:                       "HTML",
	MEDIA_TEXT_PLAIN:                      "Plain Text",
	MEDIA_TEXT_RICHTEXT:                   "Rich Text",
	MEDIA_TEXT_RTF:                        "Rich Text Format",
	MEDIA_TEXT_RFC822_HEADERS:             "RFC 822 Headers",
	MEDIA_TEXT_SGML:                       "SGML",
	MEDIA_TEXT_TAB_SEPARATED_VALUES:       "Tab Separated Values",
	MEDIA_TEXT_URI_LIST:                   "URI List",
	MEDIA_TEXT_XML:                        "XML",
	MEDIA_TEXT_XML_EXTERNAL_PARSED_ENTITY: "XML External Parsed Entity",

	// Image Media Types
	MEDIA_IMAGE_AVIF:    "AVIF Image",
	MEDIA_IMAGE_BMP:     "Bitmap Image",
	MEDIA_IMAGE_CGM:     "Computer Graphics Metafile",
	MEDIA_IMAGE_GIF:     "GIF Image",
	MEDIA_IMAGE_PNG:     "PNG Image",
	MEDIA_IMAGE_TIFF:    "TIFF Image",
	MEDIA_IMAGE_JPEG:    "JPEG Image",
	MEDIA_IMAGE_JP2:     "JPEG 2000 Image",
	MEDIA_IMAGE_SVG_XML: "SVG Image",

	// Application Media Types
	MEDIA_APPLICATION_CDA_XML:      "Clinical Document Architecture XML",
	MEDIA_APPLICATION_EDIFACT:      "EDIFACT",
	MEDIA_APPLICATION_FHIR_JSON:    "FHIR JSON",
	MEDIA_APPLICATION_FHIR_XML:     "FHIR XML",
	MEDIA_APPLICATION_HL7V2_XML:    "HL7 v2 XML",
	MEDIA_APPLICATION_GZIP:         "Gzip Archive",
	MEDIA_APPLICATION_JSON:         "JSON",
	MEDIA_APPLICATION_MSWORD:       "Microsoft Word",
	MEDIA_APPLICATION_PDF:          "PDF Document",
	MEDIA_APPLICATION_RTF:          "Rich Text Format",
	MEDIA_APPLICATION_DICOM:        "DICOM",
	MEDIA_APPLICATION_DICOM_JSON:   "DICOM JSON",
	MEDIA_APPLICATION_DICOM_XML:    "DICOM XML",
	MEDIA_APPLICATION_OCTET_STREAM: "Binary Data",
	MEDIA_APPLICATION_OGG:          "Ogg Application",

	// OpenDocument Media Types
	MEDIA_APPLICATION_ODF_BASE:                  "OpenDocument Base",
	MEDIA_APPLICATION_ODF_CHART:                 "OpenDocument Chart",
	MEDIA_APPLICATION_ODF_CHART_TEMPLATE:        "OpenDocument Chart Template",
	MEDIA_APPLICATION_ODF_FORMULA:               "OpenDocument Formula",
	MEDIA_APPLICATION_ODF_FORMULA_TEMPLATE:      "OpenDocument Formula Template",
	MEDIA_APPLICATION_ODF_GRAPHICS:              "OpenDocument Graphics",
	MEDIA_APPLICATION_ODF_GRAPHICS_TEMPLATE:     "OpenDocument Graphics Template",
	MEDIA_APPLICATION_ODF_IMAGE:                 "OpenDocument Image",
	MEDIA_APPLICATION_ODF_IMAGE_TEMPLATE:        "OpenDocument Image Template",
	MEDIA_APPLICATION_ODF_PRESENTATION:          "OpenDocument Presentation",
	MEDIA_APPLICATION_ODF_PRESENTATION_TEMPLATE: "OpenDocument Presentation Template",
	MEDIA_APPLICATION_ODF_SPREADSHEET:           "OpenDocument Spreadsheet",
	MEDIA_APPLICATION_ODF_SPREADSHEET_TEMPLATE:  "OpenDocument Spreadsheet Template",
	MEDIA_APPLICATION_ODF_TEXT:                  "OpenDocument Text",
	MEDIA_APPLICATION_ODF_TEXT_MASTER:           "OpenDocument Text Master",
	MEDIA_APPLICATION_ODF_TEXT_TEMPLATE:         "OpenDocument Text Template",
	MEDIA_APPLICATION_ODF_TEXT_WEB:              "OpenDocument Text Web",

	// Microsoft Office Media Types
	MEDIA_APPLICATION_MS_WORD_MACRO_ENABLED:         "Microsoft Word (Macro-enabled)",
	MEDIA_APPLICATION_MS_WORD_DOCX:                  "Microsoft Word Document",
	MEDIA_APPLICATION_MS_WORD_TEMPLATE_MACRO:        "Microsoft Word Template (Macro-enabled)",
	MEDIA_APPLICATION_MS_WORD_TEMPLATE_DOCX:         "Microsoft Word Template",
	MEDIA_APPLICATION_MS_POWERPOINT_SLIDESHOW_MACRO: "Microsoft PowerPoint Slideshow (Macro-enabled)",
	MEDIA_APPLICATION_MS_POWERPOINT_SLIDESHOW:       "Microsoft PowerPoint Slideshow",
	MEDIA_APPLICATION_MS_POWERPOINT_MACRO:           "Microsoft PowerPoint (Macro-enabled)",
	MEDIA_APPLICATION_MS_POWERPOINT_PPTX:            "Microsoft PowerPoint Presentation",
	MEDIA_APPLICATION_MS_EXCEL_BINARY_MACRO:         "Microsoft Excel Binary (Macro-enabled)",
	MEDIA_APPLICATION_MS_EXCEL_MACRO:                "Microsoft Excel (Macro-enabled)",
	MEDIA_APPLICATION_MS_EXCEL_XLSX:                 "Microsoft Excel Spreadsheet",
	MEDIA_APPLICATION_MS_XPS:                        "Microsoft XPS Document",
	MEDIA_APPLICATION_MS_EXCEL:                      "Microsoft Excel",
	MEDIA_APPLICATION_MS_OUTLOOK:                    "Microsoft Outlook",
	MEDIA_APPLICATION_MS_POWERPOINT:                 "Microsoft PowerPoint",

	// Archive Media Types
	MEDIA_APPLICATION_RAR: "RAR Archive",
	MEDIA_APPLICATION_ZIP: "ZIP Archive",
}

// IsValidMediaType checks if the given string is a valid IANA media type
func IsValidMediaType(mediaType string) bool {
	_, exists := MediaTypeNames[mediaType]
	return exists
}

// GetMediaTypeName returns the description for the given media type
func GetMediaTypeName(mediaType string) (string, bool) {
	name, exists := MediaTypeNames[mediaType]
	return name, exists
}
