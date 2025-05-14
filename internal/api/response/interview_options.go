package response

type SelectOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type SelectOptionsResponse struct {
	Countries            []SelectOption `json:"countries"`
	Levels               []SelectOption `json:"levels"`
	JobTypes             []SelectOption `json:"jobTypes"`
	Languages            []SelectOption `json:"languages"`
	ProgrammingLanguages []SelectOption `json:"programmingLanguages"`
	InterviewTypes       []SelectOption `json:"interviewTypes"`
	Stages               []SelectOption `json:"stages"`
}
