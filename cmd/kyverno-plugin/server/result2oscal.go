/*
Copyright 2023 IBM Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package server

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/complytime/complybeacon/proofwatch"
	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	typepolr "sigs.k8s.io/wg-policy-prototypes/policy-report/pkg/api/wgpolicyk8s.io/v1beta1"

	"github.com/oscal-compass/compliance-to-policy-go/v2/internal/utils"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"

	ocsf "github.com/Santiago-Labs/go-ocsf/ocsf/v1_5_0"
)

type ResultToOscal struct {
	policy                  policy.Policy
	policyResultsDir        string
	policyReportList        *typepolr.PolicyReportList
	clusterPolicyReportList *typepolr.ClusterPolicyReportList
	policyList              *kyvernov1.PolicyList
	clusterPolicyList       *kyvernov1.ClusterPolicyList
}

type PolicyReportContainer struct {
	PolicyReports        []*typepolr.PolicyReport
	ClusterPolicyReports []*typepolr.ClusterPolicyReport
}

type PolicyResourceIndexContainer struct {
	PolicyResourceIndex PolicyResourceIndex
	ControlIds          []string
}

func NewResultToOscal(pl policy.Policy, policyResultsDir string) *ResultToOscal {
	r := ResultToOscal{
		policy:                  pl,
		policyResultsDir:        policyResultsDir,
		policyReportList:        &typepolr.PolicyReportList{},
		clusterPolicyReportList: &typepolr.ClusterPolicyReportList{},
		policyList:              &kyvernov1.PolicyList{},
		clusterPolicyList:       &kyvernov1.ClusterPolicyList{},
	}
	return &r
}

func (r *ResultToOscal) retrievePolicyReportResults(name string) []*typepolr.PolicyReportResult {
	prrs := []*typepolr.PolicyReportResult{}
	for _, polr := range r.policyReportList.Items {
		for _, result := range polr.Results {
			policy := result.Policy
			if policy == name {
				prrs = append(prrs, &result)
			}
		}
	}
	return prrs
}

func (r *ResultToOscal) loadData(path string, out interface{}) error {
	if err := utils.LoadYamlFileToK8sTypedObject(r.policyResultsDir+"/"+path, &out); err != nil {
		return err
	}
	return nil
}

func makeProp(name string, value string) policy.Property {
	return policy.Property{
		Name:  name,
		Value: value,
	}
}

func (r *ResultToOscal) GenerateResults(ctx context.Context, watcher *proofwatch.ProofWatch) (policy.PVPResult, error) {
	var polList kyvernov1.PolicyList
	if err := r.loadData("/policies.kyverno.io.yaml", &polList); err != nil {
		return policy.PVPResult{}, err
	}

	var cpolList kyvernov1.ClusterPolicyList
	if err := r.loadData("/clusterpolicies.kyverno.io.yaml", &cpolList); err != nil {
		return policy.PVPResult{}, err
	}

	var polrList typepolr.PolicyReportList
	if err := r.loadData("/policyreports.wgpolicyk8s.io.yaml", &polrList); err != nil {
		return policy.PVPResult{}, err
	}
	r.policyReportList = &polrList

	var cpolrList typepolr.ClusterPolicyReportList
	if err := r.loadData("/clusterpolicyreports.wgpolicyk8s.io.yaml", &cpolrList); err != nil {
		return policy.PVPResult{}, err
	}

	var observations []policy.ObservationByCheck
	for _, rule := range r.policy {
		for _, check := range rule.Checks {
			name := check.ID
			prrs := r.retrievePolicyReportResults(name)
			observation := policy.ObservationByCheck{
				Title:       rule.Rule.ID,
				CheckID:     name,
				Description: fmt.Sprintf("Observation of check %s", name),
				Methods:     []string{"TEST-AUTOMATED"},
				Props: []policy.Property{
					makeProp("assessment-rule-id", rule.Rule.ID),
				},
				Collected: time.Now(),
				Subjects:  []policy.Subject{},
			}
			for _, prr := range prrs {
				for _, resource := range prr.Subjects {
					gvknsn := fmt.Sprintf("ApiVersion: %s, Kind: %s, Namespace: %s, Name: %s", resource.APIVersion, resource.Kind, resource.Namespace, resource.Name)
					subject := policy.Subject{
						Title:       gvknsn,
						ResourceID:  string(resource.UID),
						Type:        "resource",
						Result:      mapResults(prr.Result),
						EvaluatedOn: time.Now(),
						Reason:      prr.Description,
					}
					observation.Subjects = append(observation.Subjects, subject)
				}
				evidence, err := ToOCSF(name, prr)
				if err != nil {
					return policy.PVPResult{}, err
				}
				if err := watcher.Log(ctx, evidence); err != nil {
					return policy.PVPResult{}, err
				}
			}
			observations = append(observations, observation)
		}
	}
	result := policy.PVPResult{
		ObservationsByCheck: observations,
	}

	return result, nil
}

func mapResults(result typepolr.PolicyResult) policy.Result {
	switch result {
	case "pass":
		return policy.ResultPass
	case "fail", "warn":
		return policy.ResultFail
	case "error":
		return policy.ResultError
	default:
		return policy.ResultInvalid
	}
}

func ToOCSF(checkId string, prr *typepolr.PolicyReportResult) (proofwatch.Evidence, error) {
	classUID := 6007
	categoryUID := 6
	categoryName := "Application Activity"
	className := "Scan Activity"
	completedScan := 60070

	// Map operation to OCSF activity type
	var activityID int
	var activityName string
	var typeName string

	vendorName := "kyverno"
	productName := "kyverno"
	action := "observed"
	actionId := int32(3)
	status, statusID := mapResultsStatus(prr.Result)

	numFilesInt := len(prr.Subjects)
	if numFilesInt > math.MaxInt32 {
		return proofwatch.Evidence{}, fmt.Errorf("number of subjects (%d) exceeds the maximum value for an int32 (%d)", numFilesInt, math.MaxInt32)
	}
	numFiles := int32(numFilesInt)

	severity, id := mapSeverity(prr)

	uid := fmt.Sprintf("c2p-kyerno-%s", prr.Policy)
	activity := ocsf.ScanActivity{
		ActivityId:   int32(activityID),
		ActivityName: &activityName,
		CategoryName: &categoryName,
		CategoryUid:  int32(categoryUID),
		ClassName:    &className,
		ClassUid:     int32(classUID),
		Status:       &status,
		StatusId:     &statusID,
		Severity:     &severity,
		SeverityId:   id,
		NumFiles:     &numFiles,
		Metadata: ocsf.Metadata{
			Uid: &uid,
			Product: ocsf.Product{
				Name:       &productName,
				VendorName: &vendorName,
			},
			Version:     "v1beta1",
			LogProvider: &productName,
		},
		Time:     prr.Timestamp.Seconds,
		TypeName: &typeName,
		TypeUid:  int64(completedScan),
	}

	policyData, err := json.Marshal(prr)
	if err != nil {
		return proofwatch.Evidence{}, err
	}
	policyDataStr := string(policyData)

	policy := ocsf.Policy{
		Name: &prr.Policy,
		Uid:  &checkId,
		Data: &policyDataStr,
		Desc: &prr.Description,
	}

	files := "Other"
	for _, input := range prr.Subjects {
		observable := ocsf.Observable{
			Name:   &input.Name,
			Type:   &files,
			TypeId: int32(99),
			Value:  &input.Kind,
		}
		activity.Observables = append(activity.Observables, &observable)
	}

	evidenceEvent := proofwatch.Evidence{
		ScanActivity: activity,
		Policy:       policy,
		Action:       &action,
		ActionID:     &actionId,
	}

	return evidenceEvent, nil
}

func mapResultsStatus(result typepolr.PolicyResult) (string, int32) {
	switch result {
	case "pass":
		return "success", 1
	case "fail", "warn":
		return "failure", 2
	case "error":
		return "error", 6
	default:
		return "unknown", 0
	}
}

func mapSeverity(policy *typepolr.PolicyReportResult) (string, int32) {
	switch policy.Severity {
	case "critical":
		return "critical", 5
	case "high":
		return "high", 4
	case "low":
		return "low", 2
	case "medium":
		return "medium", 3
	case "info":
		return "information", 1
	default:
		return "unknown", 0
	}
}
