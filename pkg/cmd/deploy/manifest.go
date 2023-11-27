package deploy

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	apiapp "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"

	"github.com/aziontech/azion-cli/utils"
)

type Manifest struct {
	Routes Routes `json:"routes"`
	Fs     []any  `json:"fs"`
}

type Routes struct {
	Deliver []Deliver `json:"deliver"`
	Compute []Compute `json:"compute"`
}

type Compute struct {
	Variable   string `json:"from"`
	InputValue string `json:"to"`
	Priority   int    `json:"priority"`
}

type Deliver struct {
	Variable   string `json:"from"`
	InputValue string `json:"to"`
	Priority   int    `json:"priority"`
}

var manifestFilePath = "/.edge/manifest.json"

func readManifest() (*Manifest, error) {
	path, err := utils.GetWorkingDir()
	if err != nil {
		return nil, err
	}

	b, err := os.ReadFile(utils.Concat(path, manifestFilePath))
	if err != nil {
		return nil, err
	}

	manifest := Manifest{}
	err = json.Unmarshal(b, &manifest)
	if err != nil {
		return nil, err
	}

	return &manifest, err
}

func prepareRequestDeliverRulesEngine(manifest Manifest) []apiapp.RequestsRulesEngine {
	deliver := manifest.Routes.Deliver
	cri := make([][]sdk.RulesEngineCriteria, len(deliver))

	for i := 0; i < len(deliver); i++ {
		cri[i] = make([]sdk.RulesEngineCriteria, len(deliver))
	}

	requestList := []apiapp.RequestsRulesEngine{}

	for i, v := range deliver {
		req := apiapp.CreateRulesEngineRequest{}
		req.SetName("deliver")

		var beh sdk.RulesEngineBehaviorString
		beh.SetName("deliver")
		beh.SetTarget("")

		req.SetBehaviors([]sdk.RulesEngineBehaviorEntry{
			{
				RulesEngineBehaviorString: &beh,
			},
		})

		var criteria sdk.RulesEngineCriteria

		criteria.SetConditional("if")
		criteria.SetOperator("starts_with")
		criteria.SetVariable(v.Variable)
		criteria.SetInputValue(utils.Concat(".edge/storage", v.InputValue))

		cri[i][i] = criteria
		req.SetCriteria(cri)

		requestList = append(requestList, apiapp.RequestsRulesEngine{
			Request: req.CreateRulesEngineRequest,
			Phase:   "response",
		})
	}

	return requestList
}

func prepareRequestComputeRulesEngine(manifest Manifest) []apiapp.RequestsRulesEngine {
	compute := manifest.Routes.Compute
	cri := make([][]sdk.RulesEngineCriteria, len(compute))

	for i := 0; i < len(compute); i++ {
		cri[i] = make([]sdk.RulesEngineCriteria, len(compute))
	}

	requestList := []apiapp.RequestsRulesEngine{}

	for i, v := range compute {
		req := apiapp.CreateRulesEngineRequest{}
		req.SetName("compute")

		var beh sdk.RulesEngineBehaviorString
		beh.SetName("run_function")
		beh.SetTarget("")

		req.SetBehaviors([]sdk.RulesEngineBehaviorEntry{
			{
				RulesEngineBehaviorString: &beh,
			},
		})

		var criteria sdk.RulesEngineCriteria
		criteria.SetConditional("if")
		criteria.SetOperator("starts_with")
		criteria.SetVariable("${uri}")
		criteria.SetInputValue(utils.Concat(".edge/", v.InputValue))

		cri[i][i] = criteria

		req.SetCriteria(cri)
		requestList = append(requestList, apiapp.RequestsRulesEngine{
			Request: req.CreateRulesEngineRequest,
			Phase:   "response",
		})
	}

	return requestList
}

func prepareRequestCachePolicyRulesEngine(cacheID int64, template, mode string) apiapp.RequestsRulesEngine {
	req := apiapp.CreateRulesEngineRequest{}
	req.SetName("cache policy")

	var beh sdk.RulesEngineBehaviorString
	beh.SetName("set_cache_policy")
	beh.SetTarget(fmt.Sprintf("%d", cacheID))

	req.SetBehaviors([]sdk.RulesEngineBehaviorEntry{
		{
			RulesEngineBehaviorString: &beh,
		},
	})

	cri := make([][]sdk.RulesEngineCriteria, 1)
	for i := 0; i < 1; i++ {
		cri[i] = make([]sdk.RulesEngineCriteria, 1)
	}

	cri[0][0].SetConditional("if")
	cri[0][0].SetVariable("${uri}")
	cri[0][0].SetOperator("starts_with")

	if template == "Next" && strings.ToLower(mode) == "compute" {
		cri[0][0].SetInputValue("/_next/static")
	} else {
		cri[0][0].SetInputValue("/")
	}
	req.SetCriteria(cri)

	return apiapp.RequestsRulesEngine{
		Request: req.CreateRulesEngineRequest,
		Phase:   "request",
	}
}

func prepareRequestEnableGZipRulesEngine() apiapp.RequestsRulesEngine {
	req := apiapp.CreateRulesEngineRequest{}
	req.SetName("enable gzip")

	var beh sdk.RulesEngineBehaviorString
	beh.SetName("enable_gzip")
	beh.SetTarget("")

	req.SetBehaviors([]sdk.RulesEngineBehaviorEntry{
		{
			RulesEngineBehaviorString: &beh,
		},
	})

	cri := make([][]sdk.RulesEngineCriteria, 1)
	for i := 0; i < 1; i++ {
		cri[i] = make([]sdk.RulesEngineCriteria, 1)
	}

	cri[0][0].SetConditional("if")
	cri[0][0].SetVariable("${request_uri}")
	cri[0][0].SetOperator("exists")
	cri[0][0].SetInputValue("")
	req.SetCriteria(cri)

	return apiapp.RequestsRulesEngine{
		Request: req.CreateRulesEngineRequest,
		Phase:   "response",
	}
}
