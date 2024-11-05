package observations

import (
	proto "github.com/oscal-compass/compliance-to-policy-go/v2/api/proto/v1alpha1"
	"time"
)

// Result represents the kind of result statuses.
type Result int

const (
	ResultInvalid Result = iota
	ResultFail
	ResultError
	ResultPass
	ResultWarning
)

// String prints a string representation of the result
func (r Result) String() string {
	switch r {
	case ResultInvalid:
		return "INVALID"
	case ResultFail:
		return "fail"
	case ResultError:
		return "error"
	case ResultPass:
		return "pass"
	case ResultWarning:
		return "warning"
	default:
		panic("invalid result")
	}
}

var protoByResult = map[Result]proto.Result{
	ResultPass:    proto.Result_RESULT_PASS,
	ResultInvalid: proto.Result_RESULT_UNSPECIFIED,
	ResultError:   proto.Result_RESULT_ERROR,
	ResultWarning: proto.Result_RESULT_WARNING,
	ResultFail:    proto.Result_RESULT_FAILURE,
}

var resultByProto = map[proto.Result]Result{
	proto.Result_RESULT_UNSPECIFIED: ResultInvalid,
	proto.Result_RESULT_ERROR:       ResultError,
	proto.Result_RESULT_WARNING:     ResultWarning,
	proto.Result_RESULT_PASS:        ResultPass,
	proto.Result_RESULT_FAILURE:     ResultFail,
}

type Property struct {
	Name  string
	Value string
}

type Link struct {
	Description string
	Href        string
}

type Subject struct {
	Title       string
	Type        string
	ResourceID  string
	Result      Result
	EvaluatedOn time.Time
	Reason      string
	Props       []*Property
}

type ObservationByCheck struct {
	Title             string
	Description       string
	CheckID           string
	Methods           []string
	Subjects          []*Subject
	Collected         time.Time
	RelevantEvidences []*Link
	Props             []*Property
}
