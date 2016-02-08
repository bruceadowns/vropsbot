package comm

import (
	"fmt"
	"strconv"

	"github.com/bruceadowns/vropsbot/common"
)

// Mention returns the default mention string for the given user
func Mention(id string) (res string) {
	res = fmt.Sprintf("<@%s>: ", id)

	user, err := GetUser(id)
	if err != nil {
		return
	}

	res = fmt.Sprintf("<@%s|%s>: ", id, user)
	return
}

// BuildwebBranchURL returns a md formatted buildweb url for a branch
func BuildwebBranchURL(branch string) string {
	config, _ := common.NewConfig()
	if config.URLTemplates.BuildwebBranch == "" {
		return branch
	}

	return fmt.Sprintf(config.URLTemplates.BuildwebBranch, branch, branch)
}

// BuildwebSbURL returns a md formatted buildweb url for a sandbox build
func BuildwebSbURL(sb int, friendly string) string {
	config, _ := common.NewConfig()
	if config.URLTemplates.BuildwebSb == "" {
		return strconv.Itoa(sb)
	}

	return fmt.Sprintf(config.URLTemplates.BuildwebSb, sb, friendly)
}

// BuildwebTargetURL returns a md formatted buildweb url for a target
func BuildwebTargetURL(target string) string {
	config, _ := common.NewConfig()
	if config.URLTemplates.BuildwebTarget == "" {
		return target
	}

	return fmt.Sprintf(config.URLTemplates.BuildwebTarget, target, target)
}

// CatAreaURL returns a md formatted cat area url
func CatAreaURL(area string) string {
	config, _ := common.NewConfig()
	if config.URLTemplates.CatArea == "" {
		return area
	}

	return fmt.Sprintf(config.URLTemplates.CatArea, area, area)
}

// CatTestrunURL returns a md formatted cat testrun url
func CatTestrunURL(testrun int, friendly string) string {
	config, _ := common.NewConfig()
	if config.URLTemplates.CatTestrun == "" {
		return friendly
	}

	return fmt.Sprintf(config.URLTemplates.CatTestrun, testrun, friendly)
}

// CatTestrunResultsURL returns a md formatted cat testrun results url
func CatTestrunResultsURL(dir string) string {
	config, _ := common.NewConfig()
	if config.URLs.BaseCatURL == "" {
		return dir
	}

	return fmt.Sprintf("<%s%s|results dir>", config.URLs.BaseCatURL, dir)
}

// P4WebBranchURL returns a md formatted p4web url
func P4WebBranchURL(branch, cln string) string {
	config, _ := common.NewConfig()
	if config.URLTemplates.P4WebBranch == "" {
		return branch
	}

	return fmt.Sprintf(config.URLTemplates.P4WebBranch, branch, cln, cln)
}
