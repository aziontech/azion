package delete

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/config/delete"
	app "github.com/aziontech/azion-cli/pkg/api/applications"
	cachesetting "github.com/aziontech/azion-cli/pkg/api/cache_setting"
	connector "github.com/aziontech/azion-cli/pkg/api/connector"
	firewall "github.com/aziontech/azion-cli/pkg/api/firewall"
	firewallrules "github.com/aziontech/azion-cli/pkg/api/firewall_rules"
	function "github.com/aziontech/azion-cli/pkg/api/function"
	functioninstance "github.com/aziontech/azion-cli/pkg/api/function_instance"
	rulesengine "github.com/aziontech/azion-cli/pkg/api/rules_engine"
	"github.com/aziontech/azion-cli/pkg/api/storage"
	workload "github.com/aziontech/azion-cli/pkg/api/workloads"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// isNotFoundError checks if the error is a 404 Not Found error
func isNotFoundError(err error) bool {
	return errors.Is(err, utils.ErrorNotFound404)
}

type DeleteCmd struct {
	Io         *iostreams.IOStreams
	f          *cmdutil.Factory
	GetAzion   func(confPath string) (*contracts.AzionApplicationOptions, error)
	AskInput   func(string) (string, error)
	WriteFile  func(filename string, data []byte, perm fs.FileMode) error
	GetWorkDir func() (string, error)
}

type Fields struct {
	ConfigDir string
	Force     bool
}

func NewDeleteCmd(f *cmdutil.Factory) *DeleteCmd {
	return &DeleteCmd{
		Io:         f.IOStreams,
		f:          f,
		GetAzion:   utils.GetAzionJsonContent,
		AskInput:   utils.AskInput,
		WriteFile:  os.WriteFile,
		GetWorkDir: os.Getwd,
	}
}

func NewCobraCmd(delete *DeleteCmd) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion config delete
		$ azion config delete --force
		$ azion config delete --config-dir ./my-project
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return delete.Run(fields)
		},
	}

	cmd.Flags().StringVar(&fields.ConfigDir, "config-dir", ".", msg.FlagConfigDir)
	cmd.Flags().BoolVar(&fields.Force, "force", false, msg.FlagForce)
	cmd.Flags().BoolP("help", "h", false, msg.FlagHelp)

	return cmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDeleteCmd(f))
}

func (del *DeleteCmd) Run(fields *Fields) error {
	ctx := context.Background()
	logger.Debug("Running config delete command")

	azionJson, err := del.GetAzion(fields.ConfigDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return msg.ErrorMissingAzionJson
		}
		return err
	}

	if azionJson.Application.ID == 0 && len(azionJson.Firewalls) == 0 && len(azionJson.Function) == 0 && azionJson.Workloads.Id == 0 {
		logger.FInfo(del.Io.Out, "No resources found in azion.json to delete\n")
		return nil
	}

	if !fields.Force {
		answer, err := del.AskInput(msg.ConfirmDeletion)
		if err != nil {
			return err
		}
		if strings.ToLower(answer) != "y" && strings.ToLower(answer) != "yes" {
			logger.FInfo(del.Io.Out, msg.DeletionAborted)
			return nil
		}
	}

	logger.FInfo(del.Io.Out, msg.DeletingResources)

	clientApp := app.NewClient(del.f.HttpClient, del.f.Config.GetString("api_v4_url"), del.f.Config.GetString("token"))
	clientFunc := function.NewClient(del.f.HttpClient, del.f.Config.GetString("api_v4_url"), del.f.Config.GetString("token"))
	clientWorkload := workload.NewClient(del.f.HttpClient, del.f.Config.GetString("api_v4_url"), del.f.Config.GetString("token"))
	clientFirewall := firewall.NewClient(del.f.HttpClient, del.f.Config.GetString("api_v4_url"), del.f.Config.GetString("token"))
	clientRulesEngine := rulesengine.NewClient(del.f.HttpClient, del.f.Config.GetString("api_v4_url"), del.f.Config.GetString("token"))
	clientCacheSetting := cachesetting.NewClientV4(del.f.HttpClient, del.f.Config.GetString("api_v4_url"), del.f.Config.GetString("token"))
	clientFuncInstance := functioninstance.NewClient(del.f.HttpClient, del.f.Config.GetString("api_v4_url"), del.f.Config.GetString("token"))
	clientFwRules := firewallrules.NewClient(del.f.HttpClient, del.f.Config.GetString("api_v4_url"), del.f.Config.GetString("token"))
	clientConnector := connector.NewClient(del.f.HttpClient, del.f.Config.GetString("api_v4_url"), del.f.Config.GetString("token"))
	clientStorage := storage.NewClient(del.f.HttpClient, del.f.Config.GetString("storage_url"), del.f.Config.GetString("token"))

	var errs []string
	successCount := 0
	failCount := 0

	// Delete ALL Application Rules Engine rules from server first (not just tracked ones)
	// This is critical because rules may reference function instances or cache settings
	if azionJson.Application.ID != 0 {
		// List and delete ALL request phase rules
		listOpts := &contracts.ListOptions{PageSize: 100}
		for {
			requestRules, err := clientApp.ListRulesEngineRequest(ctx, listOpts, azionJson.Application.ID)
			if err != nil {
				// If application doesn't exist (404), there are no rules to delete - skip silently
				if isNotFoundError(err) {
					logger.Debug("Application not found, skipping request phase rules deletion", zap.Int64("applicationID", azionJson.Application.ID))
					break
				}
				errs = append(errs, fmt.Sprintf("Failed to list request phase rules for application %d: %v", azionJson.Application.ID, err))
				break
			}
			if requestRules == nil || len(requestRules.Results) == 0 {
				break
			}
			for _, rule := range requestRules.Results {
				ruleName := rule.GetName()
				ruleId := rule.GetId()
				logger.FInfo(del.Io.Out, fmt.Sprintf(msg.DeletingRulesEngineApp, ruleName, ruleId))
				err := clientRulesEngine.DeleteRequest(ctx, azionJson.Application.ID, ruleId)
				if err != nil {
					errs = append(errs, fmt.Sprintf("Failed to delete Rules Engine rule '%s' (ID: %d): %v", ruleName, ruleId, err))
					logger.FInfo(del.Io.Out, fmt.Sprintf("Failed to delete Rules Engine rule '%s': %v\n", ruleName, err))
					failCount++
				} else {
					successCount++
				}
			}
			if requestRules.Results == nil || len(requestRules.Results) < int(listOpts.PageSize) {
				break
			}
			listOpts.Page++
		}

		// List and delete ALL response phase rules
		listOpts = &contracts.ListOptions{PageSize: 100}
		for {
			responseRules, err := clientApp.ListRulesEngineResponse(ctx, listOpts, azionJson.Application.ID)
			if err != nil {
				// If application doesn't exist (404), there are no rules to delete - skip silently
				if isNotFoundError(err) {
					logger.Debug("Application not found, skipping response phase rules deletion", zap.Int64("applicationID", azionJson.Application.ID))
					break
				}
				errs = append(errs, fmt.Sprintf("Failed to list response phase rules for application %d: %v", azionJson.Application.ID, err))
				break
			}
			if responseRules == nil || len(responseRules.Results) == 0 {
				break
			}
			for _, rule := range responseRules.Results {
				ruleName := rule.GetName()
				ruleId := rule.GetId()
				logger.FInfo(del.Io.Out, fmt.Sprintf(msg.DeletingRulesEngineApp, ruleName, ruleId))
				err := clientRulesEngine.DeleteResponse(ctx, azionJson.Application.ID, ruleId)
				if err != nil {
					errs = append(errs, fmt.Sprintf("Failed to delete Rules Engine rule '%s' (ID: %d): %v", ruleName, ruleId, err))
					logger.FInfo(del.Io.Out, fmt.Sprintf("Failed to delete Rules Engine rule '%s': %v\n", ruleName, err))
					failCount++
				} else {
					successCount++
				}
			}
			if responseRules.Results == nil || len(responseRules.Results) < int(listOpts.PageSize) {
				break
			}
			listOpts.Page++
		}
	}

	// Delete ALL Firewall Rules from server (not just tracked ones)
	if len(azionJson.Firewalls) > 0 {
		for _, fw := range azionJson.Firewalls {
			listOpts := &contracts.ListOptions{PageSize: 100}
			for {
				fwRules, err := clientFwRules.List(ctx, listOpts, fw.Id)
				if err != nil {
					// If firewall doesn't exist (404), there are no rules to delete - skip silently
					if isNotFoundError(err) {
						logger.Debug("Firewall not found, skipping firewall rules deletion", zap.Int64("firewallID", fw.Id))
						break
					}
					errs = append(errs, fmt.Sprintf("Failed to list rules for firewall %d: %v", fw.Id, err))
					break
				}
				if fwRules == nil || len(fwRules.Results) == 0 {
					break
				}
				for _, rule := range fwRules.Results {
					ruleName := rule.GetName()
					ruleId := rule.GetId()
					logger.FInfo(del.Io.Out, fmt.Sprintf(msg.DeletingRulesEngineFw, ruleName, ruleId, fw.Id))
					err := clientFwRules.Delete(ctx, fw.Id, ruleId)
					if err != nil {
						errs = append(errs, fmt.Sprintf("Failed to delete Firewall rule '%s' (ID: %d) from firewall %d: %v", ruleName, ruleId, fw.Id, err))
						logger.FInfo(del.Io.Out, fmt.Sprintf("Failed to delete Firewall rule '%s': %v\n", ruleName, err))
						failCount++
					} else {
						successCount++
					}
				}
				if fwRules.Results == nil || len(fwRules.Results) < int(listOpts.PageSize) {
					break
				}
				listOpts.Page++
			}
		}
	}

	// Now delete ALL Function Instances from server (not just tracked ones)
	// This must happen AFTER rules are deleted
	if azionJson.Application.ID != 0 {
		listOpts := &contracts.ListOptions{PageSize: 100}
		for {
			funcInstances, err := clientApp.EdgeFuncInstancesList(ctx, listOpts, azionJson.Application.ID)
			if err != nil {
				// If application doesn't exist (404), there are no function instances to delete - skip silently
				if isNotFoundError(err) {
					logger.Debug("Application not found, skipping function instances deletion", zap.Int64("applicationID", azionJson.Application.ID))
					break
				}
				errs = append(errs, fmt.Sprintf("Failed to list function instances for application %d: %v", azionJson.Application.ID, err))
				break
			}
			if funcInstances == nil || len(funcInstances.Results) == 0 {
				break
			}
			for _, fn := range funcInstances.Results {
				fnName := fn.GetName()
				fnId := fn.GetId()
				logger.FInfo(del.Io.Out, fmt.Sprintf(msg.DeletingFuncInstanceApp, fnName, fnId))
				err := clientFuncInstance.Delete(ctx, azionJson.Application.ID, fnId)
				if err != nil {
					errs = append(errs, fmt.Sprintf("Failed to delete Function Instance '%s' (ID: %d): %v", fnName, fnId, err))
					logger.FInfo(del.Io.Out, fmt.Sprintf("Failed to delete Function Instance '%s': %v\n", fnName, err))
					failCount++
				} else {
					successCount++
				}
			}
			if funcInstances.Results == nil || len(funcInstances.Results) < int(listOpts.PageSize) {
				break
			}
			listOpts.Page++
		}
	}

	// Now delete ALL Cache Settings from server (not just tracked ones)
	// This must happen AFTER rules are deleted
	if azionJson.Application.ID != 0 {
		listOpts := &contracts.ListOptions{PageSize: 100}
		for {
			cacheSettings, err := clientCacheSetting.List(ctx, listOpts, azionJson.Application.ID)
			if err != nil {
				// If application doesn't exist (404), there are no cache settings to delete - skip silently
				if isNotFoundError(err) {
					logger.Debug("Application not found, skipping cache settings deletion", zap.Int64("applicationID", azionJson.Application.ID))
					break
				}
				errs = append(errs, fmt.Sprintf("Failed to list cache settings for application %d: %v", azionJson.Application.ID, err))
				break
			}
			if cacheSettings == nil {
				break
			}
			results := cacheSettings.GetResults()
			if len(results) == 0 {
				break
			}
			for _, cs := range results {
				csName := cs.GetName()
				csId := cs.GetId()
				logger.FInfo(del.Io.Out, fmt.Sprintf(msg.DeletingCacheSetting, csName, csId))
				_, err := clientCacheSetting.Delete(ctx, azionJson.Application.ID, csId)
				if err != nil {
					errs = append(errs, fmt.Sprintf("Failed to delete Cache Setting '%s' (ID: %d): %v", csName, csId, err))
					logger.FInfo(del.Io.Out, fmt.Sprintf("Failed to delete Cache Setting '%s': %v\n", csName, err))
					failCount++
				} else {
					successCount++
				}
			}
			if len(results) < int(listOpts.PageSize) {
				break
			}
			listOpts.Page++
		}
	}

	// Finally delete the Application itself
	if azionJson.Application.ID != 0 {
		logger.FInfo(del.Io.Out, fmt.Sprintf(msg.DeletingApplication, azionJson.Application.Name, azionJson.Application.ID))
		err := clientApp.Delete(ctx, azionJson.Application.ID)
		if err != nil {
			errs = append(errs, fmt.Sprintf("Failed to delete Application '%s' (ID: %d): %v", azionJson.Application.Name, azionJson.Application.ID, err))
			logger.FInfo(del.Io.Out, fmt.Sprintf("Failed to delete Application '%s': %v\n", azionJson.Application.Name, err))
			failCount++
		} else {
			successCount++
		}
	}

	if len(azionJson.Firewalls) > 0 {
		for _, fw := range azionJson.Firewalls {
			logger.FInfo(del.Io.Out, fmt.Sprintf(msg.DeletingFirewall, fw.Name, fw.Id))
			err := clientFirewall.Delete(ctx, fw.Id)
			if err != nil {
				errs = append(errs, fmt.Sprintf("Failed to delete Firewall '%s' (ID: %d): %v", fw.Name, fw.Id, err))
				logger.FInfo(del.Io.Out, fmt.Sprintf("Failed to delete Firewall '%s': %v\n", fw.Name, err))
				failCount++
			} else {
				successCount++
			}
		}
	}

	if len(azionJson.Function) > 0 {
		for _, fn := range azionJson.Function {
			if fn.ID != 0 {
				logger.FInfo(del.Io.Out, fmt.Sprintf(msg.DeletingFunction, fn.Name, fn.ID))
				err := clientFunc.Delete(ctx, fn.ID)
				if err != nil {
					errs = append(errs, fmt.Sprintf("Failed to delete Function '%s' (ID: %d): %v", fn.Name, fn.ID, err))
					logger.FInfo(del.Io.Out, fmt.Sprintf("Failed to delete Function '%s': %v\n", fn.Name, err))
					failCount++
				} else {
					successCount++
				}
			}
		}
	}

	if azionJson.Workloads.Id != 0 {
		logger.FInfo(del.Io.Out, fmt.Sprintf(msg.DeletingWorkload, azionJson.Workloads.Name, azionJson.Workloads.Id))
		err := clientWorkload.Delete(ctx, azionJson.Workloads.Id)
		if err != nil {
			errs = append(errs, fmt.Sprintf("Failed to delete Workload '%s' (ID: %d): %v", azionJson.Workloads.Name, azionJson.Workloads.Id, err))
			logger.FInfo(del.Io.Out, fmt.Sprintf("Failed to delete Workload '%s': %v\n", azionJson.Workloads.Name, err))
			failCount++
		} else {
			successCount++
		}
	}

	if azionJson.Bucket != "" {
		logger.FInfo(del.Io.Out, fmt.Sprintf(msg.DeletingBucket, azionJson.Bucket))
		err := clientStorage.DeleteBucket(ctx, azionJson.Bucket)
		if err != nil {
			errs = append(errs, fmt.Sprintf("Failed to delete Storage Bucket '%s': %v", azionJson.Bucket, err))
			logger.FInfo(del.Io.Out, fmt.Sprintf("Failed to delete Storage Bucket '%s': %v\n", azionJson.Bucket, err))
			failCount++
		} else {
			successCount++
		}
	}

	if len(azionJson.Connectors) > 0 {
		for _, conn := range azionJson.Connectors {
			if conn.Id != 0 {
				logger.FInfo(del.Io.Out, fmt.Sprintf(msg.DeletingConnector, conn.Name, conn.Id))
				err := clientConnector.Delete(ctx, conn.Id)
				if err != nil {
					errs = append(errs, fmt.Sprintf("Failed to delete Connector '%s' (ID: %d): %v", conn.Name, conn.Id, err))
					logger.FInfo(del.Io.Out, fmt.Sprintf("Failed to delete Connector '%s': %v\n", conn.Name, err))
					failCount++
				} else {
					successCount++
				}
			}
		}
	}

	logger.FInfo(del.Io.Out, msg.ResettingConfig)
	err = del.resetAzionJson(fields.ConfigDir)
	if err != nil {
		errs = append(errs, fmt.Sprintf("Failed to reset azion.json: %v", err))
		logger.FInfo(del.Io.Out, fmt.Sprintf("Failed to reset azion.json: %v\n", err))
		failCount++
	} else {
		logger.FInfo(del.Io.Out, msg.ConfigResetSuccess)
	}

	logger.FInfo(del.Io.Out, msg.DeletionSummary)
	logger.FInfo(del.Io.Out, fmt.Sprintf(msg.ResourcesDeleted, successCount))
	logger.FInfo(del.Io.Out, fmt.Sprintf(msg.ResourcesFailed, failCount))

	if len(errs) > 0 {
		logger.FInfo(del.Io.Out, msg.ErrorsDuringDeletion)
		for _, e := range errs {
			logger.FInfo(del.Io.Out, fmt.Sprintf("  - %s\n", e))
		}
		return fmt.Errorf(msg.ErrorPartialDeletion.Error(), len(errs))
	}

	deleteOut := output.GeneralOutput{
		Msg:   msg.DeleteSuccess,
		Out:   del.Io.Out,
		Flags: del.f.Flags,
	}
	return output.Print(&deleteOut)
}

func deleteRulesEngineByPhase(ctx context.Context, client *rulesengine.Client, applicationID, ruleID int64, phase string) error {
	logger.Debug("Deleting Rules Engine rule", zap.Int64("applicationID", applicationID), zap.Int64("ruleID", ruleID), zap.String("phase", phase))

	switch phase {
	case "request":
		return client.DeleteRequest(ctx, applicationID, ruleID)
	case "response":
		return client.DeleteResponse(ctx, applicationID, ruleID)
	default:
		return client.DeleteRequest(ctx, applicationID, ruleID)
	}
}

func (del *DeleteCmd) resetAzionJson(configDir string) error {
	wd, err := del.GetWorkDir()
	if err != nil {
		return err
	}

	var configPath string
	if path.IsAbs(configDir) {
		configPath = path.Join(configDir, "azion.json")
	} else {
		configPath = path.Join(wd, configDir, "azion.json")
	}

	azionJson := &contracts.AzionApplicationOptions{}
	data, err := json.MarshalIndent(azionJson, "", "  ")
	if err != nil {
		logger.Debug("Error marshaling azion.json", zap.Error(err))
		return err
	}

	if err := del.WriteFile(configPath, data, 0644); err != nil {
		logger.Debug("Error creating config file", zap.Error(err))
		return err
	}

	return nil
}
