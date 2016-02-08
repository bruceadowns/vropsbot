package plugins

import (
	"gopkg.in/jmcvetta/napping.v3"
)

// DefaultCatParams returns the default set of params
// consumers will append their specifics
func DefaultCatParams() napping.Params {
	return napping.Params{
		"format":   "json",
		"order_by": "-endtime",
		"limit":    "1",
	}
}

// DefaultRecommendedParams returns the default set of params for a recommended build
func DefaultRecommendedParams() napping.Params {
	return napping.Params{
		"format":    "json",
		"order_by":  "-updated",
		"limit":     "1",
		"site":      "mbu",
		"sla__name": "VA_Bats",
	}
}

// DefaultCurrBuildParams returns the default set of params for a current build
func DefaultCurrBuildParams() napping.Params {
	return napping.Params{
		"format": "json",
	}
}
