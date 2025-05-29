package convert

import (
	"fmt"
	"time"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/models"
	"github.com/revanite-io/sci/layer4"

	"github.com/oscal-compass/compliance-to-policy-go/v2/framework/policy"
	"github.com/oscal-compass/compliance-to-policy-go/v2/pkg"
)

func SCI2AssessmentResults(plans []policy.PlanRef, catalogId string) (oscalTypes.AssessmentResults, error) {
	comps := map[string]oscalTypes.SystemComponent{}
	var findings []oscalTypes.Finding
	var observations []oscalTypes.Observation
	assessedControls := oscalTypes.AssessedControls{
		IncludeControls: &[]oscalTypes.AssessedControlsSelectControlById{},
	}
	var err error
	assessmentResults := oscalTypes.AssessmentResults{
		UUID: uuid.NewUUID(),
		ImportAp: oscalTypes.ImportAp{
			Href: "REPLACE_ME",
		},
		Results: []oscalTypes.Result{},
	}
	assessmentResults.Metadata = models.NewSampleMetadata()

	localDefinition := &oscalTypes.LocalDefinitions{
		Components:     &[]oscalTypes.SystemComponent{},
		InventoryItems: &[]oscalTypes.InventoryItem{},
	}
	for _, plan := range plans {
		if plan.Plan == nil {
			if err := plan.Load(); err != nil {
				return assessmentResults, err
			}
		}
		if plan.Plan.CatalogID != catalogId {
			continue
		}
		comp, ok := comps[plan.PluginID.String()]
		if !ok {
			comp = component(plan.PluginID.String())
			comps[plan.PluginID.String()] = comp
			*localDefinition.Components = append(*localDefinition.Components, comp)
		}
		ii := inventoryItem(comp, plan.Service)
		*localDefinition.InventoryItems = append(*localDefinition.InventoryItems, ii)

		for _, controlEval := range plan.Plan.ControlEvaluations {
			selection := oscalTypes.AssessedControlsSelectControlById{
				ControlId:    controlEval.ControlID,
				StatementIds: &[]string{},
			}
			for _, assessment := range controlEval.Assessments {
				*selection.StatementIds = append(*selection.StatementIds, assessment.RequirementID)
				for _, method := range assessment.Methods {
					obs := observation(method, plan.Plan.EndTime, ii, assessment.RequirementID)
					observations = append(observations, obs)
					if !method.Run || method.Result.Status != layer4.Status(layer4.Passed.String()) {
						findings, err = generateFindings(findings, obs, []string{controlEval.ControlID})
						if err != nil {
							return oscalTypes.AssessmentResults{}, err
						}
					}
				}
			}
			*assessedControls.IncludeControls = append(*assessedControls.IncludeControls)
		}
	}

	// Check for empty optional arrays
	localDefinition.Components = pkg.NilIfEmpty(localDefinition.Components)
	localDefinition.InventoryItems = pkg.NilIfEmpty(localDefinition.InventoryItems)
	assessedControls.IncludeControls = pkg.NilIfEmpty(assessedControls.IncludeControls)

	results := oscalTypes.Result{
		UUID:             uuid.NewUUID(),
		Title:            "Results",
		Start:            time.Now(),
		LocalDefinitions: localDefinition,
		ReviewedControls: oscalTypes.ReviewedControls{
			Description:       fmt.Sprintf("Review controls from %s", catalogId),
			ControlSelections: []oscalTypes.AssessedControls{assessedControls},
		},
		Observations: pkg.NilIfEmpty(&observations),
		Findings:     pkg.NilIfEmpty(&findings),
	}
	assessmentResults.Results = append(assessmentResults.Results, results)
	return assessmentResults, nil
}

func observation(method layer4.AssessmentMethod, end time.Time, inventoryItem oscalTypes.InventoryItem, req string) oscalTypes.Observation {
	return oscalTypes.Observation{
		Collected:   end,
		Description: method.Description,
		Title:       method.Name,
		UUID:        uuid.NewUUID(),
		Subjects: &[]oscalTypes.SubjectReference{
			{
				SubjectUuid: inventoryItem.UUID,
				Title:       inventoryItem.Description,
				Type:        "inventory-item",
				Props: &[]oscalTypes.Property{
					{
						Name:  "resource-id",
						Value: req,
						Ns:    extensions.TrestleNameSpace,
					},
					{
						Name:  "result",
						Value: string(method.Result.Status),
						Ns:    extensions.TrestleNameSpace,
					},
					{
						Name:  "reason",
						Value: method.RemediationGuide,
						Ns:    extensions.TrestleNameSpace,
					},
				},
			},
		},
	}
}

func component(validator string) oscalTypes.SystemComponent {
	return oscalTypes.SystemComponent{
		Title:       validator,
		UUID:        uuid.NewUUID(),
		Description: fmt.Sprintf("Validated by %s", validator),
		Status: oscalTypes.SystemComponentStatus{
			State: "operational",
		},
	}
}

func inventoryItem(comp oscalTypes.SystemComponent, service string) oscalTypes.InventoryItem {
	return oscalTypes.InventoryItem{
		UUID: uuid.NewUUID(),
		ImplementedComponents: &[]oscalTypes.ImplementedComponent{
			{
				ComponentUuid: comp.UUID,
			},
		},
		Description: service,
	}
}

func getFindingForTarget(findings []oscalTypes.Finding, targetId string) *oscalTypes.Finding {
	for i := range findings {
		if findings[i].Target.TargetId == targetId {
			return &findings[i] // if finding is found, return a pointer to that slice element
		}
	}
	return nil
}

// Generate OSCAL Findings for all non-passing controls in the OSCAL Observation
func generateFindings(findings []oscalTypes.Finding, observation oscalTypes.Observation, targets []string) ([]oscalTypes.Finding, error) {
	for _, targetId := range targets {
		finding := getFindingForTarget(findings, targetId)
		if finding == nil { // if an empty finding was returned, create a new one and append to findings
			newFinding := oscalTypes.Finding{
				UUID: uuid.NewUUID(),
				RelatedObservations: &[]oscalTypes.RelatedObservation{
					{
						ObservationUuid: observation.UUID,
					},
				},
				Target: oscalTypes.FindingTarget{
					TargetId: targetId,
					Type:     "statement-id",
					Status: oscalTypes.ObjectiveStatus{
						State: "not-satisfied",
					},
				},
			}
			findings = append(findings, newFinding)
		} else {
			relObs := oscalTypes.RelatedObservation{
				ObservationUuid: observation.UUID,
			}
			*finding.RelatedObservations = append(*finding.RelatedObservations, relObs) // add new related obs to existing finding for targetId
		}
	}
	return findings, nil
}
