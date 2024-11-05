package report

type Reporter struct {
}

func New() *Reporter {
	return &Reporter{}
}

func (r *Reporter) ToOSCAL() error {
	return nil
}

func (r *Reporter) ToMarkdown() error {
	return nil
}
