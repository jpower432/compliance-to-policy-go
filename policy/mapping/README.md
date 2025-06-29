# Notes

By having the plugins store the checks we have more flexibility for plugin that have pre-defined check profiles.
Core should only pass the in scope rule and the applicability (profile)

To support Gemara, the assessment methods information should be colocated with the check logic and registered against
the Layer 2 Catalog. To accomplish this, the plugins need awareness of the evaluation plans, not the core. The `CompletePolicy`
method is added to fill in the check information from various sources for the applied ruleSet.

# Options

1. OSCAL Validation Component
2. Gemara L4 Evaluation
3. Plugin store the mapping internally

The plugin will report the rule to check mapping back in the Observation
