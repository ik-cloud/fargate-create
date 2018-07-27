package cmd

import (
	"errors"
	"fmt"
	"strings"
)

func parseInputVars(tf string) (string, string, string, string, error) {
	app := ""
	environment := ""
	profile := ""
	region := ""

	//look for variables
	lines := strings.Split(tf, "\n")
	inTags := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		//ignore whitespace and comments
		if trimmed == "" || strings.HasPrefix(trimmed, "//") {
			continue
		}
		if trimmed == "tags = {" {
			inTags = true
			continue
		}
		if inTags {
			if trimmed == "}" {
				inTags = false
			}
			continue
		}
		//key = "value"
		parts := strings.Split(trimmed, "=")
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			value = strings.Replace(value, `"`, "", -1)
			if key == "app" {
				app = value
			}
			if key == "environment" {
				environment = value
			}
			if key == "aws_profile" {
				profile = value
			}
			if key == "region" {
				region = value
			}
		}
	}

	//did we find it?
	if app == "" {
		return "", "", "", "", errors.New(`missing variable: "app"`)
	}
	if environment == "" {
		return "", "", "", "", errors.New(`missing variable: "environment"`)
	}
	if profile == "" {
		return "", "", "", "", errors.New(`missing variable: "profile"`)
	}

	return app, environment, profile, region, nil
}

func updateTerraformBackend(tf string, profile string, app string, env string) string {
	//update terraform.backend (which doesn't support dynamic variables)
	// profile = ""
	// bucket  = ""
	// key     = "dev.terraform.tfstate"
	tmp := strings.Split(tf, "\n")
	newTf := ""
	for _, line := range tmp {
		updatedLine := line
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, `profile = ""`) {
			updatedLine = fmt.Sprintf(`    profile = "%s"`, profile)
		}
		if strings.HasPrefix(trimmed, "bucket") {
			updatedLine = fmt.Sprintf(`    bucket  = "tf-state-%s"`, app)
		}
		if strings.HasPrefix(trimmed, "key") {
			updatedLine = fmt.Sprintf(`    key     = "%s.terraform.tfstate"`, env)
		}
		newTf += updatedLine + "\n"
	}
	return newTf
}
