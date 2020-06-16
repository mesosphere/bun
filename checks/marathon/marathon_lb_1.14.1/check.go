package marathon_lb_1_14_1

import (
	"strings"

	"github.com/mesosphere/bun/v2/bundle"
	"github.com/mesosphere/bun/v2/checks"
)

func init() {
	check := checks.Check{
		Name: "marathon-lb-1.14.1",
		Description: "Detects if Marathon-LB v1.14.1 is installed. " +
			"This Marathon-LB version is affected by a HAProxy bug: https://github.com/haproxy/haproxy/issues/278",
		Cure: "Upgrade to the latest Marathon-LB version. Please see " +
			"https://support.d2iq.com/s/article/Known-Issue-L4TOUT-on-Marathon-LB-1-14-1-Shipped-with-HAProxy-2-0-6 " +
			"for more details.",
		OKSummary:      "No Marathon-LB v1.14.1 installed",
		ProblemSummary: "Marathon-LB v1.14.1 is installed",
		Run:            checkFunc,
	}
	checks.RegisterCheck(check)
}

type marathonApps struct {
	Apps []struct {
		Id        string `json:"id"`
		Container struct {
			Docker struct {
				Image string `json:"image"`
			} `json:"docker"`
		} `json:"container"`
	} `json:"apps"`
}

func checkFunc(b bundle.Bundle) checks.Results {
	if len(b.Masters()) == 0 {
		return checks.Results{checks.Result{
			Status: checks.SUndefined,
			Value:  "Cannot find DC/OS masters to perform the check",
		}}
	}
	apps := marathonApps{}
	var err error
	for _, master := range b.Masters() {
		err = master.ReadJSON("marathon-apps", &apps)
		if err == nil {
			break
		}
	}
	if err != nil {
		return checks.Results{
			{
				Status: checks.SUndefined,
				Value:  err,
			},
		}
	}
	ids := make([]string, 0, len(apps.Apps))
	for _, a := range apps.Apps {
		if strings.HasSuffix(a.Container.Docker.Image, "marathon-lb:v1.14.1") {
			ids = append(ids, a.Id)
		}
	}
	if len(ids) > 0 {
		return checks.Results{
			checks.Result{
				Status: checks.SProblem,
				Value:  "Following Marathon-LB v1.14.1 services found: " + strings.Join(ids, ", "),
			},
		}
	}
	return checks.Results{
		checks.Result{
			Status: checks.SOK,
		},
	}
}
