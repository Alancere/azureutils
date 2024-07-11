package spec

import (
	"strings"

	"github.com/google/go-github/v53/github"
)

type label string

// enable
const (
	GoBreakingChange           label = "CI-BreakingChange-Go"
	JsBreakingChange           label = "CI-BreakingChange-JavaScript"
	PythonBreakingChange       label = "CI-BreakingChange-Python"
	PythonTrack2BreakingChange label = "CI-BreakingChange-Python-Track2"
	DotNetBreakingChange       label = "CI-BreakingChange-DotNet"
	JavaBreakingChange         label = "CI-BreakingChange-Java"

	GoApprovedBreakingChange     label = "Approved-SdkBreakingChange-Go"
	JsApprovedBreakingChange     label = "Approved-SdkBreakingChange-JavaScript"
	PythonApprovedBreakingChange label = "Approved-SdkBreakingChange-Python"

	// private specs
	GoPrivateApproveBreakingChange label = "Approved-BreakingChange-Go"
)

// when
const (
	BreakingChangeReviewRequired label = "BreakingChangeReviewRequired"
	NewApiVersionRequired        label = "NewApiVersionRequired"

	ApprovedBreakingChange label = "Approved-BreakingChange"
	BreakingChangeApproved label = "BreakingChange-Approved"
)

const (
	ARMReview           label = "ARMReview"
	WaitForARMFeedback  label = "WaitForARMFeedback"
	ARMChangesRequested label = "ARMChangesRequested"

	ARMSignedOff label = "ARMSignedOff"
)

const (
	ArcReview label = "ArcReview"

	ArcSignedOff label = "ArcSignedOff"
)

const (
	noRecentActivity label = "no-recent-activity"
)

// 新规则
const (
	BreakingChange_GO_SDK         label = "BreakingChange-Go-Sdk"
	BreakingChange_JavaScript_SDK label = "BreakingChange-JavaScript-Sdk"
	BreakingChange_Python_SDK     label = "BreakingChange-Python-Sdk"

	BreakingChange_GO_SDK_Suppression         label = "BreakingChange-Go-Sdk-Suppression"
	BreakingChange_JavaScript_SDK_Suppression label = "BreakingChange-JavaScript-Sdk-Suppression"
	BreakingChange_Python_SDK_Suppression     label = "BreakingChange-Python-Sdk-Suppression"

	// Changes are not breaking at the REST API level and have at most minor impact to generated SDKs
	BreakingChange_Approved_Benign label = "BreakingChange-Approved-Benign"
	// Changes are to correct the REST API definition to correctly describe service behavior
	BreakingChange_Approved_BigFix label = "BreakingChange-Approved-BugFix"
	// Changes were reviewed and approved in a previous PR
	BreakingChange_Approved_Previously label = "BreakingChange-Approved-Previously"
	// Changes are not backward compatible and may cause customer disruption.
	BreakingChange_Approved_UserImpact label = "BreakingChange-Approved-UserImpact"

	BreakingChange_GO_SDK_Approved         label = "BreakingChange-Go-Sdk-Approved"
	BreakingChange_JavaScript_SDK_Approved label = "BreakingChange-JavaScript-Sdk-Approved"
	BreakingChange_Python_SDK_Approved     label = "BreakingChange-Python-Sdk-Approved"

	BreakingChange_GO_SDK_Suppression_Approved         label = "BreakingChange-Go-Sdk-Suppression-Approved"
	BreakingChange_JavaScript_SDK_Suppression_Approved label = "BreakingChange-JavaScript-Sdk-Suppression-Approved"
	BreakingChange_Python_SDK_Suppression_Approved     label = "BreakingChange-Python-Sdk-Suppression-Approved"
)

func isApprovedBreakingChange(labels []*github.Label) (bool, bool) { // reviewRequired, approved

	approved := false
	reviewRequire := false

	for _, l := range labels {
		lName := *l.Name
		if !approved && (lName == string(ApprovedBreakingChange) || strings.Contains(lName, string(BreakingChangeApproved))) {
			approved = true
		}
		if !reviewRequire && (lName == string(BreakingChangeReviewRequired) || lName == string(NewApiVersionRequired)) {
			reviewRequire = true
		}

		if reviewRequire && approved {
			return true, true
		}
	}

	return reviewRequire, false
}

func isARMSignedOff(labels []*github.Label) (bool, bool) {
	armReview := false
	signedOff := false

	for _, l := range labels {
		lName := *l.Name
		if !signedOff && lName == string(ARMSignedOff) {
			signedOff = true
		}
		if !armReview && (lName == string(WaitForARMFeedback) || lName == string(ARMChangesRequested)) || lName == string(ARMReview) { // || lName == string(ARMReview)
			armReview = true
		}

		if armReview && signedOff {
			return true, true
		}
	}

	return armReview, false
}

func isArcSignedOff(labels []*github.Label) (bool, bool) {
	arcReview := false
	signedOff := false

	for _, l := range labels {
		lName := *l.Name
		if !signedOff && lName == string(ArcSignedOff) {
			signedOff = true
		}
		if !arcReview && lName == string(ArcReview) {
			arcReview = true
		}

		if arcReview && signedOff {
			return true, true
		}
	}

	return arcReview, false
}

func isNoRecentActivity(labels []*github.Label) bool {
	for _, l := range labels {
		lName := *l.Name
		if lName == string(noRecentActivity) {
			return true
		}
	}

	return false
}

func havaLabel(labels []*github.Label, languageBreaking label) bool {
	for _, v := range labels {
		if *v.Name == string(languageBreaking) {
			return true
		}
	}
	return false
}
