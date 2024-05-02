package operations

import (
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/noahc3/hatch-cli/src/log"
	"github.com/noahc3/hatch-cli/src/types"
)

func generateFail(name string, rule string, msg string) bool {
	log.Error("  Variable '%s' validation rules failed: (%s) %s", name, rule, msg)
	return true
}

func validateVar(variable types.Variable, value string) bool {
	rules := strings.Split(variable.Rules, "|")
	failed := false

	for _, rule := range rules {
		if rule == "required" {
			if value == "" {
				failed = generateFail(variable.EnvVariable, rule, "Value is not defined and the Egg does not provide a default")
			}
		} else if rule == "string" {
			// PASS
		} else if strings.HasPrefix(rule, "max:") {
			numstr := strings.Split(rule, ":")[1]
			count, err := strconv.Atoi(numstr)
			if err != nil {
				failed = generateFail(variable.EnvVariable, rule, "BAD RULE (invalid max value)")
				continue
			}

			if len(value) > count {
				failed = generateFail(variable.EnvVariable, rule, "Value is too long")
			}
		} else if strings.HasPrefix(rule, "regex:") {
			pattern := strings.Split(rule, ":")[1]
			pattern = pattern[1 : len(pattern)-1]
			re := regexp.MustCompile(pattern)

			if !re.MatchString(value) {
				failed = generateFail(variable.EnvVariable, rule, "Value does not match regex pattern")
			}
		}
	}

	return failed
}

func validateVariables(egg *types.Egg, variables map[string]string) bool {
	failed := false

	for _, variable := range egg.Variables {
		failed = failed || validateVar(variable, variables[variable.EnvVariable])
	}

	return !failed
}

func GetEnvVariables(egg *types.Egg) (map[string]string, error) {
	variables := make(map[string]string)

	log.Info("Loaded environment variables:\n")

	for _, env := range egg.Variables {
		value, res := os.LookupEnv(env.EnvVariable)

		if !res {
			log.Info("  %s: '%s' (default)\n", env.EnvVariable, env.DefaultValue)
			variables[env.EnvVariable] = env.DefaultValue
			continue
		}

		log.Info("  %s: %s\n", env.EnvVariable, value)
		variables[env.EnvVariable] = value
	}

	log.Nl()

	log.Info("Checking variable validation rules\n")

	pass := validateVariables(egg, variables)

	if !pass {
		return nil, log.Error("Environment variables failed validation\n")
	}

	log.Info("  Validation passed\n")

	return variables, nil
}

func CheckVars(egg *types.Egg) error {
	_, err := GetEnvVariables(egg)

	if err != nil {
		return err
	}

	return nil
}
