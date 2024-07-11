package spec

import (
	"context"
	"fmt"

	"github.com/google/go-github/v53/github"
)

type PullRequest struct {
	client *github.Client
}

func NewPullRequest() *PullRequest {
	return &PullRequest{
		client: CredGithubRepository(),
	}
}

func (pr *PullRequest) Check(repo string) error {
	// var err error
	var prs []*github.PullRequest

	for i := 1; i <= 5; i++ {
		temp, _, err := pr.client.PullRequests.List(context.Background(), Owner, repo, &github.PullRequestListOptions{
			State: "open",
			// Base: "main",
			// Sort: "desc"
			ListOptions: github.ListOptions{
				Page:    i,
				PerPage: 1000,
			},
		})
		if err != nil {
			// return nil, err
			return err
		}
		prs = append(prs, temp...)
	}

	//prs, _, err := pr.client.PullRequests.List(context.Background(), owner, repo, &github.PullRequestListOptions{
	//	State: "open",
	//	//Base: "main",
	//	//Sort: "desc"
	//	ListOptions: github.ListOptions{
	//		PerPage: 1000,
	//	},
	//})
	//if err != nil {
	//	//return nil, err
	//	return err
	//}
	//
	//newprs, _, err := pr.client.PullRequests.List(context.Background(), owner, repo, &github.PullRequestListOptions{
	//	State: "open",
	//	//Base: "main",
	//	//Sort: "desc"
	//	ListOptions: github.ListOptions{
	//		Page:    2,
	//		PerPage: 1000,
	//	},
	//})
	//if err != nil {
	//	//return nil, err
	//	return err
	//}

	// by CI-BreakingChange filter
	newPullRequests := make([]*github.PullRequest, 0, 20)
	for _, l := range prs {
		if *l.Draft {
			continue
		}
		if isNoRecentActivity(l.Labels) {
			continue
		}
		if havaLabel(l.Labels, GoBreakingChange) || havaLabel(l.Labels, BreakingChange_GO_SDK) || havaLabel(l.Labels, BreakingChange_GO_SDK_Suppression) {
			newPullRequests = append(newPullRequests, l)
		}
	}

	approves := make([]Approve, 0, len(newPullRequests))
	fmt.Printf("\nnon-compliant: need %s or %s or %s\n", ApprovedBreakingChange, ARMSignedOff, ArcSignedOff)
	for _, l := range newPullRequests {
		reviewRequire, approved := isApprovedBreakingChange(l.Labels)
		armReview, armSignedOff := isARMSignedOff(l.Labels)
		arcReview, arcSignedOff := isArcSignedOff(l.Labels)

		// true true pass false
		// true false fail true
		// false false pass false
		b1 := false
		b2 := false
		b3 := false
		if !((!reviewRequire && approved) || (reviewRequire && !approved)) {
			b1 = true
		}
		if !((!armReview && armSignedOff) || (armReview && !armSignedOff)) {
			b2 = true
		}
		if !((!arcReview && arcSignedOff) || (arcReview && !arcSignedOff)) {
			b3 = true
		}

		if b1 && b2 && b3 {
			flag := false
			if havaLabel(l.Labels, GoApprovedBreakingChange) || havaLabel(l.Labels, GoPrivateApproveBreakingChange) ||
				havaLabel(l.Labels, BreakingChange_GO_SDK_Approved) || havaLabel(l.Labels, BreakingChange_GO_SDK_Suppression_Approved) {
				flag = true
			}
			approves = append(approves, Approve{
				url:           *l.HTMLURL,
				reviewRequire: reviewRequire,
				armReview:     armReview,
				arcReview:     arcReview,
				approved:      flag,
			})
		} else {
			fmt.Printf("%s reviewRequire: %t, armReview: %t, arcReview: %t\n", *l.HTMLURL, reviewRequire, armReview, arcReview)
		}
	}

	fmt.Println("checks:", len(newPullRequests))

	// print need to approve
	fmt.Printf("\napproved breaking change:%s(%d)\n", GoApprovedBreakingChange, len(approves))
	for _, a := range approves {
		if a.approved {
			fmt.Print(a)
		}
	}
	fmt.Printf("\ncan to add %s label.\n", GoApprovedBreakingChange)
	for _, a := range approves {
		if !a.approved {
			fmt.Print(a)
		}
	}
	fmt.Println()
	return nil
}

type Approve struct {
	url           string
	approved      bool
	reviewRequire bool
	armReview     bool
	arcReview     bool
}

func (a Approve) String() string {
	return fmt.Sprintf("%s reviewRequire: %t, armReview: %t, arcReview: %t\n", a.url, a.reviewRequire, a.armReview, a.arcReview)
	// return s
}
