# Notes

By having the plugins store the checks we have more flexibility for plugin that have pre-defined check profiles.
Core should only pass the in scope rule and the applicability (profile)

To support Gemara, the assessment methods information should be colocated with the check logic and registered against
the Layer 2 Catalog in the form of a Layer 4 Evaluation. The rule-to-check mappings is held within the validator.
To accomplish this, the plugins need awareness of the evaluation plans, not the core.

# Options

1. OSCAL Validation Component
2. Gemara L4 Evaluation
3. Plugin store the mapping internally

The component flow would be "Evaluation Plans" -> "Policy Engines" much like the receiver -> processor relationship