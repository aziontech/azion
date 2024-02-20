package deploy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	msg "github.com/aziontech/azion-cli/messages/deploy"
	apiEdgeApplications "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"
	"go.uber.org/zap"

	thoth "github.com/aziontech/go-thoth"
)

type Manifest struct {
	Routes []Routes `json:"routes"`
	Fs     []any    `json:"fs"`
}

type Routes struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Priority int    `json:"priority"`
	Type     string `json:"type"`
}

var manifestFilePath = "/.edge/manifest.json"

func readManifest(cmd *DeployCmd) (*Manifest, error) {
	logger.Debug("read manifest")
	pathWorkingDir, err := cmd.GetWorkDir()
	if err != nil {
		return nil, err
	}

	b, err := cmd.FileReader(utils.Concat(pathWorkingDir, manifestFilePath))
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

// Interpreted TODO: better interpreted, removed flows edge application, domain
func (manifest *Manifest) Interpreted(f *cmdutil.Factory, cmd *DeployCmd, conf *contracts.AzionApplicationOptions, clients *Clients) error {
	logger.Debug("Execute manifest")
	ctx := context.Background()

	err := cmd.doApplication(clients.EdgeApplication, ctx, conf)
	if err != nil {
		return err
	}

	domainName, err := cmd.doDomain(clients.Domain, ctx, conf)
	if err != nil {
		return err
	}

	err = cmd.doBucket(clients.Bucket, ctx, conf)
	if err != nil {
		return err
	}

	// skip upload when type = javascript, typescript (storage folder does not exist in these cases)
	if conf.Template != "javascript" && conf.Template != "typescript" {
		err = cmd.uploadFiles(f, conf)
		if err != nil {
			return err
		}
	}

	// cacheID created for "compute" in the SSR, will be used to create the function and configure the caching policy.
	var cacheID int64 = 0

	for _, route := range manifest.Routes {
		if route.From == "/_next/data/" {
			continue
		}

		if route.Type == "compute" {
			conf.Function.File = ".edge/worker.js"
			err := cmd.doFunction(clients, ctx, conf)
			if err != nil {
				return err
			}	

			var reqCache apiEdgeApplications.CreateCacheSettingsRequest
			reqCache.SetName("function policy")
			reqCache.SetBrowserCacheSettings("honor")
			reqCache.SetCdnCacheSettings("honor")
			reqCache.SetCdnCacheSettingsMaximumTtl(0)
			reqCache.SetCacheByQueryString("all")
			reqCache.SetCacheByCookies("all")

			// create cache to function next
			cache, err := clients.EdgeApplication.CreateCacheEdgeApplication(ctx, &reqCache, conf.Application.ID)
			if err != nil {
				logger.Debug("Error while creating cache settings", zap.Error(err))
				return err
			}
			cacheID = cache.GetId()
			logger.FInfo(cmd.F.IOStreams.Out, msg.CacheSettingsSuccessful)	
		}

		err = cmd.doOrigin(clients.EdgeApplication, clients.Origin, ctx, conf)
		if err != nil {
			logger.Debug("Error while creating origin", zap.Error(err))
			return err
		}
		logger.FInfo(cmd.F.IOStreams.Out, msg.OriginsSuccessful)

		// TODO: Update default rule engine is being run multiple times
		ruleDefaultID, err := clients.EdgeApplication.GetRulesDefault(ctx, conf.Application.ID, "request")
		if err != nil {
			logger.Debug("Error while getting default rules engine", zap.Error(err))
			return err
		}

		if strings.ToLower(conf.Template) == "javascript" || strings.ToLower(conf.Template) == "typescript" {
			reqRules := apiEdgeApplications.UpdateRulesEngineRequest{}
			reqRules.IdApplication = conf.Application.ID

			_, err := clients.EdgeApplication.UpdateRulesEnginePublish(ctx, &reqRules, conf.Function.InstanceID)
			if err != nil {
				return err
			}
			continue
		} else {
			behaviors := make([]sdk.RulesEngineBehaviorEntry, 0)

			var behString sdk.RulesEngineBehaviorString
			behString.SetName("set_origin")

			if route.Type == "compute" {
				behString.SetTarget(strconv.Itoa(int(conf.Origin.SingleOriginID)))
			} else {
				behString.SetTarget(strconv.Itoa(int(conf.Origin.StorageOriginID)))
			}

			behaviors = append(behaviors, sdk.RulesEngineBehaviorEntry{
				RulesEngineBehaviorString: &behString,
			})

			reqUpdateRulesEngine := apiEdgeApplications.UpdateRulesEngineRequest{
				IdApplication: conf.Application.ID,
				Phase:         "request",
				Id:            ruleDefaultID,
			}

			reqUpdateRulesEngine.SetBehaviors(behaviors)

			_, err = clients.EdgeApplication.UpdateRulesEngine(ctx, &reqUpdateRulesEngine)
			if err != nil {
				logger.Debug("Error while updating default rules engine", zap.Error(err))
				return err
			}
		}

		// check if the rules engines have not been created.
		if !conf.RulesEngine.Created {
			// create rules engines to compute else delivery
			if strings.ToLower(conf.Mode) == "compute" {
				requestRules, err := requestRulesEngineManifest(conf, route, cacheID)
				if err != nil {
					return err
				}

				_, err = clients.EdgeApplication.CreateRulesEngine(ctx, conf.Application.ID, "request", &requestRules)
				if err != nil {
					return err
				}

			} else {

				reqDeliver := apiEdgeApplications.CreateRulesEngineRequest{}
				reqDeliver.SetName("rule_rewrite_deliver")

				behaviors := make([]sdk.RulesEngineBehaviorEntry, 0)

				var behRewriteRequest sdk.RulesEngineBehaviorString

				behRewriteRequest.SetName("rewrite_request")

				var target = fmt.Sprintf("${uri}%sindex.html", "")
				if strings.ToLower(conf.Template) == "html" && len(Path) > 0 {
					Path = strings.ReplaceAll(Path, "/", "")
					Path = strings.ReplaceAll(Path, ".", "")
					target = fmt.Sprintf("${uri}/%s/index.html", Path)
				}
				behRewriteRequest.SetTarget(target)

				behaviors = append(behaviors, sdk.RulesEngineBehaviorEntry{
					RulesEngineBehaviorString: &behRewriteRequest,
				})

				reqDeliver.SetBehaviors(behaviors)

				criteria := make([][]sdk.RulesEngineCriteria, 1)
				for i := 0; i < 1; i++ {
					criteria[i] = make([]sdk.RulesEngineCriteria, 1)
				}

				criteria[0][0].SetConditional("if")
				criteria[0][0].SetVariable("${uri}")
				criteria[0][0].SetOperator("matches")
				criteria[0][0].SetInputValue(".*/$")
				reqDeliver.SetCriteria(criteria)

				_, err = clients.EdgeApplication.CreateRulesEngine(ctx, conf.Application.ID, "request", &reqDeliver)
				if err != nil {
					return err
				}

				reqDeliverRoot := apiEdgeApplications.CreateRulesEngineRequest{}
				reqDeliverRoot.SetName("rule_rewrite_deliver_root")

				behaviorsRoot := make([]sdk.RulesEngineBehaviorEntry, 0)

				var behRewriteRequestRoot sdk.RulesEngineBehaviorString

				behRewriteRequestRoot.SetName("rewrite_request")
				behRewriteRequestRoot.SetTarget("${uri}/index.html")

				behaviorsRoot = append(behaviorsRoot, sdk.RulesEngineBehaviorEntry{
					RulesEngineBehaviorString: &behRewriteRequestRoot,
				})

				reqDeliverRoot.SetBehaviors(behaviorsRoot)

				criteriaRoot := make([][]sdk.RulesEngineCriteria, 1)
				for i := 0; i < 1; i++ {
					criteriaRoot[i] = make([]sdk.RulesEngineCriteria, 1)
				}

				criteriaRoot[0][0].SetConditional("if")
				criteriaRoot[0][0].SetVariable("${uri}")
				criteriaRoot[0][0].SetOperator("matches")
				regexPattern := `^(?!.*\/$)(?![\s\S]*\.[a-zA-Z0-9]+$).*`
				criteriaRoot[0][0].SetInputValue(regexPattern)
				reqDeliverRoot.SetCriteria(criteriaRoot)

				_, err = clients.EdgeApplication.CreateRulesEngine(ctx, conf.Application.ID, "request", &reqDeliverRoot)
				if err != nil {
					return err
				}
			}
		}
	}

	conf.RulesEngine.Created = true

	err = cmd.WriteAzionJsonContent(conf)
	if err != nil {
		logger.Debug("Error while writing azion.json file", zap.Error(err))
		return err
	}

	logger.FInfo(cmd.F.IOStreams.Out, msg.DeploySuccessful)
	logger.FInfo(cmd.F.IOStreams.Out, fmt.Sprintf(msg.DeployOutputDomainSuccess, utils.Concat("https://", domainName)))
	logger.FInfo(cmd.F.IOStreams.Out, msg.DeployPropagation)
	return nil
}

func requestRulesEngineManifest(conf *contracts.AzionApplicationOptions, routes Routes, cacheID int64) (apiEdgeApplications.CreateRulesEngineRequest, error) {
	logger.Debug("Create Rules Engine set origin")

	req := apiEdgeApplications.CreateRulesEngineRequest{}
	req.SetName(fmt.Sprintf("rules_manifest_%s", thoth.GenerateName()))

	behaviors := make([]sdk.RulesEngineBehaviorEntry, 0)

	if routes.Type == "compute" {
		var behFunction sdk.RulesEngineBehaviorString
		behFunction.SetName("run_function")
		behFunction.SetTarget(fmt.Sprintf("%d", conf.Function.InstanceID))
		behaviors = append(behaviors, sdk.RulesEngineBehaviorEntry{
			RulesEngineBehaviorString: &behFunction,
		})

		if strings.ToLower(conf.Template) == "next" {
			var behForwardCookies sdk.RulesEngineBehaviorString
			behForwardCookies.SetName("forward_cookies")
			behaviors = append(behaviors, sdk.RulesEngineBehaviorEntry{
				RulesEngineBehaviorString: &behForwardCookies,
			})
		}

		var behCache sdk.RulesEngineBehaviorString
		behCache.SetName("set_cache_policy")
		behCache.SetTarget(fmt.Sprintf("%d", cacheID))
		behaviors = append(behaviors, sdk.RulesEngineBehaviorEntry{
			RulesEngineBehaviorString: &behCache,
		})
	} else {
		var behOrigin sdk.RulesEngineBehaviorString
		behOrigin.SetName("set_origin")
		behOrigin.SetTarget(fmt.Sprintf("%d", conf.Origin.StorageOriginID))

		behaviors = append(behaviors, sdk.RulesEngineBehaviorEntry{
			RulesEngineBehaviorString: &behOrigin,
		})

		var behDeliver sdk.RulesEngineBehaviorString
		behDeliver.SetName("deliver")
		behaviors = append(behaviors, sdk.RulesEngineBehaviorEntry{
			RulesEngineBehaviorString: &behDeliver,
		})
	}

	req.SetBehaviors(behaviors)

	criteria := make([][]sdk.RulesEngineCriteria, 1)
	for i := 0; i < 1; i++ {
		criteria[i] = make([]sdk.RulesEngineCriteria, 1)
	}

	operator, err := checkFieldFrom(routes.From)
	if err != nil {
		return req, err
	}

	criteria[0][0].SetConditional("if")
	criteria[0][0].SetVariable("${uri}")
	criteria[0][0].SetOperator(operator)
	criteria[0][0].SetInputValue(routes.From)
	req.SetCriteria(criteria)

	return req, nil
}

func checkFieldFrom(from string) (string, error) {
	if strings.HasPrefix(from, "/") {
		startWith := "starts_with"
		return startWith, nil
	} else if strings.HasPrefix(from, "\\.") && strings.HasSuffix(from, "$") {
		doesNotMatch := "matches"
		return doesNotMatch, nil
	}

	return "", errors.New("the value of 'from' could not be recognized")
}
