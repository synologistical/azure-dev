// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package project

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/machinelearning/armmachinelearning/v3"
	"github.com/azure/azure-dev/cli/azd/pkg/ai"
	"github.com/azure/azure-dev/cli/azd/pkg/async"
	"github.com/azure/azure-dev/cli/azd/pkg/environment"
	"github.com/azure/azure-dev/cli/azd/pkg/tools"
)

// aiEndpointTarget is a ServiceTarget implementation for deploying to Azure ML online endpoints
type aiEndpointTarget struct {
	env        *environment.Environment
	envManager environment.Manager
	aiHelper   AiHelper
}

// NewAiEndpointTarget creates a new aiEndpointTarget instance
func NewAiEndpointTarget(
	env *environment.Environment,
	envManager environment.Manager,
	aiHelper AiHelper,
) ServiceTarget {
	return &aiEndpointTarget{
		env:        env,
		envManager: envManager,
		aiHelper:   aiHelper,
	}
}

// AiEndpointDeploymentResult is a struct to hold the result of an online endpoint deployment
type AiEndpointDeploymentResult struct {
	Environment *armmachinelearning.EnvironmentVersion
	Model       *armmachinelearning.ModelVersion
	Flow        *ai.Flow
	Deployment  *armmachinelearning.OnlineDeployment
}

// Initialize initializes the aiEndpointTarget
func (m *aiEndpointTarget) Initialize(ctx context.Context, serviceConfig *ServiceConfig) error {
	return nil
}

// RequiredExternalTools returns the required external tools for the machineLearningEndpointTarget
func (m *aiEndpointTarget) RequiredExternalTools(ctx context.Context, serviceConfig *ServiceConfig) []tools.ExternalTool {
	return m.aiHelper.RequiredExternalTools(ctx)
}

// Package packages the service for deployment to an Azure ML online endpoint
// This method is a no-op since the actual packaging is handled by the underlying AI/ML Python SDKs
func (m *aiEndpointTarget) Package(
	ctx context.Context,
	serviceConfig *ServiceConfig,
	frameworkPackageOutput *ServicePackageResult,
	progress *async.Progress[ServiceProgress],
) (*ServicePackageResult, error) {
	return &ServicePackageResult{}, nil
}

// Deploy deploys the service to an Azure ML online endpoint
func (m *aiEndpointTarget) Deploy(
	ctx context.Context,
	serviceConfig *ServiceConfig,
	servicePackage *ServicePackageResult,
	targetResource *environment.TargetResource,
	progress *async.Progress[ServiceProgress],
) (*ServiceDeployResult, error) {
	endpointConfig, err := ai.ParseConfig[ai.EndpointDeploymentConfig](serviceConfig.Config)
	if err != nil {
		return nil, err
	}

	workspaceScope, err := m.getWorkspaceScope(serviceConfig, targetResource)
	if err != nil {
		return nil, err
	}

	deployResult := &AiEndpointDeploymentResult{}

	// Initialize the AI project that will be used for the python bridge
	progress.SetProgress(NewServiceProgress("Initializing AI project"))
	if err := m.aiHelper.Initialize(ctx); err != nil {
		return nil, fmt.Errorf("failed initializing AI project: %w", err)
	}

	// Ensure the workspace is valid
	if err := m.aiHelper.ValidateWorkspace(ctx, workspaceScope); err != nil {
		return nil, fmt.Errorf("workspace '%s' was not found within subscription '%s' and resource group '%s': %w",
			workspaceScope.Workspace(),
			workspaceScope.SubscriptionId(),
			workspaceScope.ResourceGroup(),
			err,
		)
	}

	// Deploy flow
	if endpointConfig.Flow != nil {
		progress.SetProgress(NewServiceProgress("Deploying AI Prompt Flow"))
		flow, err := m.aiHelper.CreateFlow(ctx, workspaceScope, serviceConfig, endpointConfig.Flow)
		if err != nil {
			return nil, err
		}

		deployResult.Flow = flow
	}

	// Deploy environment
	if endpointConfig.Environment != nil {
		progress.SetProgress(NewServiceProgress("Configuring AI environment"))
		envVersion, err := m.aiHelper.CreateEnvironmentVersion(
			ctx,
			workspaceScope,
			serviceConfig,
			endpointConfig.Environment,
		)
		if err != nil {
			return nil, err
		}

		deployResult.Environment = envVersion
	}

	// Deploy model
	if endpointConfig.Model != nil {
		progress.SetProgress(NewServiceProgress("Configuring AI model"))
		modelVersion, err := m.aiHelper.CreateModelVersion(ctx, workspaceScope, serviceConfig, endpointConfig.Model)
		if err != nil {
			return nil, err
		}

		deployResult.Model = modelVersion
	}

	// Deploy to endpoint
	if endpointConfig.Deployment != nil {
		progress.SetProgress(NewServiceProgress("Deploying to AI Online Endpoint"))
		endpointName := filepath.Base(targetResource.ResourceName())
		onlineDeployment, err := m.aiHelper.DeployToEndpoint(
			ctx,
			workspaceScope,
			serviceConfig,
			endpointName,
			endpointConfig,
		)
		if err != nil {
			return nil, err
		}

		if onlineDeployment == nil {
			return nil, fmt.Errorf("unexpected response from deployToEndpoint: deployment is nil")
		}
		if onlineDeployment.Name == nil {
			return nil, fmt.Errorf("unexpected response from deployToEndpoint: deployment name is nil")
		}

		deploymentName := *onlineDeployment.Name
		progress.SetProgress(NewServiceProgress("Updating traffic"))
		_, err = m.aiHelper.UpdateTraffic(ctx, workspaceScope, endpointName, deploymentName)
		if err != nil {
			return nil, fmt.Errorf("failed updating traffic: %w", err)
		}

		progress.SetProgress(NewServiceProgress("Removing old deployments"))
		if err := m.aiHelper.DeleteDeployments(
			ctx, workspaceScope, endpointName, []string{deploymentName}); err != nil {
			return nil, fmt.Errorf("failed deleting previous deployments: %w", err)
		}

		deployResult.Deployment = onlineDeployment
	}

	endpoints, err := m.Endpoints(ctx, serviceConfig, targetResource)
	if err != nil {
		return nil, err
	}

	if err := m.envManager.Save(ctx, m.env); err != nil {
		return nil, fmt.Errorf("failed saving environment: %w", err)
	}

	return &ServiceDeployResult{
		Details:   deployResult,
		Package:   servicePackage,
		Endpoints: endpoints,
	}, nil
}

// Endpoints returns the endpoints for the service
func (m *aiEndpointTarget) Endpoints(
	ctx context.Context,
	serviceConfig *ServiceConfig,
	targetResource *environment.TargetResource,
) ([]string, error) {
	if err := m.aiHelper.Initialize(ctx); err != nil {
		return nil, fmt.Errorf("failed initializing AI project: %w", err)
	}

	tenantId, has := m.env.LookupEnv(environment.TenantIdEnvVarName)
	if !has {
		return nil, fmt.Errorf(
			"tenant ID not found. Ensure %s has been set in the environment.",
			environment.TenantIdEnvVarName,
		)
	}

	workspaceScope, err := m.getWorkspaceScope(serviceConfig, targetResource)
	if err != nil {
		return nil, err
	}

	workspaceLink := ai.AiStudioWorkspaceLink(
		tenantId,
		workspaceScope.SubscriptionId(),
		workspaceScope.ResourceGroup(),
		workspaceScope.Workspace(),
	)

	endpoints := []string{
		fmt.Sprintf("Workspace: %s", workspaceLink),
	}

	endpointName := filepath.Base(targetResource.ResourceName())
	onlineEndpoint, err := m.aiHelper.GetEndpoint(ctx, workspaceScope, endpointName)
	if err != nil {
		return nil, err
	}

	var deploymentName string
	for key, value := range onlineEndpoint.Properties.Traffic {
		if *value == 100 {
			deploymentName = key
			break
		}
	}

	if deploymentName != "" {
		deploymentLink := ai.AiStudioDeploymentLink(
			tenantId,
			workspaceScope.SubscriptionId(),
			workspaceScope.ResourceGroup(),
			workspaceScope.Workspace(),
			endpointName,
			deploymentName,
		)

		endpoints = append(endpoints, fmt.Sprintf("Deployment: %s", deploymentLink))
	}

	endpoints = append(
		endpoints,
		fmt.Sprintf("Scoring: %s", *onlineEndpoint.Properties.ScoringURI),
		fmt.Sprintf("Swagger: %s", *onlineEndpoint.Properties.SwaggerURI),
	)

	return endpoints, nil
}

// getWorkspaceScope returns the scope for the workspace
func (m *aiEndpointTarget) getWorkspaceScope(
	serviceConfig *ServiceConfig,
	targetResource *environment.TargetResource,
) (*ai.Scope, error) {
	endpointConfig, err := ai.ParseConfig[ai.EndpointDeploymentConfig](serviceConfig.Config)
	if err != nil {
		return nil, err
	}

	workspaceName, err := endpointConfig.Workspace.Envsubst(m.env.Getenv)
	if err != nil {
		return nil, err
	}

	// Workspace name can come from the following:
	// 1. The workspace field in the endpoint service config
	// 2. The AZUREAI_PROJECT_NAME environment variable
	if workspaceName == "" {
		workspaceName = m.env.Getenv(AiProjectNameEnvVarName)
	}

	if workspaceName == "" {
		return nil, fmt.Errorf("workspace name is required")
	}

	return ai.NewScope(
		m.env.GetSubscriptionId(),
		targetResource.ResourceGroupName(),
		workspaceName,
	), nil
}
