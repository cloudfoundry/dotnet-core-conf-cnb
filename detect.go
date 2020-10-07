package dotnetcoreexecute

import (
	"fmt"
	"path/filepath"

	"github.com/paketo-buildpacks/packit"
)

//go:generate faux --interface Parser --output fakes/parser.go
type Parser interface {
	ParseProjectPath(path string) (projectPath string, err error)
}

func Detect(buildpackYMLParser Parser) packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		projRoot := context.WorkingDir

		// Checking if there is a buildpack.yml file that contains a project path to use as the working dir
		bpYMLProjPath, err := buildpackYMLParser.ParseProjectPath(filepath.Join(context.WorkingDir, "buildpack.yml"))
		if err != nil {
			return packit.DetectResult{}, fmt.Errorf("failed to parse buildpack.yml: %w", err)
		}

		if bpYMLProjPath != "" {
			projRoot = filepath.Join(projRoot, bpYMLProjPath)
		}

		runtimeConfigFiles, err := filepath.Glob(filepath.Join(projRoot, "*.runtimeconfig.json"))
		if err != nil {
			return packit.DetectResult{}, fmt.Errorf("failed checking pattern *.runtimeconfig.json: %w", err)
		}
		if len(runtimeConfigFiles) > 1 {
			return packit.DetectResult{}, packit.Fail.WithMessage("multiple *.runtimeconfig.json files present")
		}

		projFiles, err := filepath.Glob(filepath.Join(projRoot, "*.*sproj"))
		if err != nil {
			return packit.DetectResult{}, fmt.Errorf("failed checking pattern *.*sproj: %w", err)
		}

		if len(runtimeConfigFiles)+len(projFiles) == 0 {
			return packit.DetectResult{}, packit.Fail
		}

		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{
					{
						Name: "dotnet-core-execute",
					},
				},
				Requires: []packit.BuildPlanRequirement{
					{
						Name: "dotnet-core-execute",
						Metadata: map[string]interface{}{
							"build": true,
						},
					},
					{
						Name: "icu",
						Metadata: map[string]interface{}{
							"launch": true,
						},
					},
				},
			},
		}, nil
	}
}
