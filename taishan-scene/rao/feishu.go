package rao

type FsReqTemplate struct {
	Type string     `json:"type"`
	Data DataEntity `json:"data"`
}

type DataEntity struct {
	TemplateId          string `json:"template_id"`
	TemplateVersionName string `json:"template_version_name"`
	TemplateVariable    any    `json:"template_variable"`
}
