package io

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Tenant struct {
	Globals struct {
		Deployment struct {
			Type string `yaml:"type"`
		} `yaml:"deployment"`
	} `yaml:"globals"`
	Realms struct {
		Nonprod struct {
			Deployment struct {
				Type string `yaml:"type"`
			} `yaml:"deployment"`
		} `yaml:"nonprod"`
		Prod struct {
			Deployment struct {
				Type string `yaml:"type"`
			} `yaml:"deployment"`
		} `yaml:"prod"`
	} `yaml:"realms"`
}

func GetYamlFiles(directory string) ([]string, error) {
	var files []string
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".yaml") {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func ParseDeploymentTypeAndMigration(fileContent string) (string, string) {
	var tenant Tenant

	err := yaml.Unmarshal([]byte(fileContent), &tenant)
	if err != nil {
		fmt.Printf("Error parsing YAML content: %v\n", err)
		return "Error", "Error"
	}

	deploymentTypes := make(map[string]bool)

	// Collect deployment types
	globalType := strings.Trim(tenant.Globals.Deployment.Type, "\" ")
	if globalType != "" {
		deploymentTypes[globalType] = true
	}

	nonprodType := strings.Trim(tenant.Realms.Nonprod.Deployment.Type, "\" ")
	if nonprodType != "" {
		deploymentTypes[nonprodType] = true
	}

	prodType := strings.Trim(tenant.Realms.Prod.Deployment.Type, "\" ")
	if prodType != "" {
		deploymentTypes[prodType] = true
	}

	// Determine migration status
	_, hasArgocd := deploymentTypes["argocd"]
	_, hasGha := deploymentTypes["gha"]

	var migrationStatus string
	if hasArgocd && hasGha {
		migrationStatus = "in_progress"
	} else if hasArgocd {
		migrationStatus = "complete"
	} else if hasGha {
		migrationStatus = "not_started"
	} else {
		migrationStatus = "unknown"
	}

	/*
	  This is being used to track live and complete migrations to ArgoCD
	  On average these migrations have taken less than a week to complete
	  and the detailed progress is tracked by migrationStatus
	*/
	deploymentType := ""
	if hasArgocd {
		deploymentType = "argocd"
	} else if hasGha {
		deploymentType = "gha"
	} else {
		deploymentType = "none"
	}

	return deploymentType, migrationStatus
}
