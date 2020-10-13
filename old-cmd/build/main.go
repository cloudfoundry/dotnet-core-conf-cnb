package main

import (
	"fmt"
	"os"

	"github.com/buildpack/libbuildpack/buildpackplan"
	"github.com/cloudfoundry/libcfbuildpack/build"
	"github.com/paketo-buildpacks/dotnet-core-conf/conf"
)

func main() {
	context, err := build.DefaultBuild()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to create a default build context: %s", err)
		os.Exit(101)
	}

	code, err := runBuild(context)
	if err != nil {
		context.Logger.Info(err.Error())
	}

	os.Exit(code)

}

func runBuild(context build.Build) (int, error) {
	context.Logger.Title(context.Buildpack)

	contributor, willContribute, err := conf.NewContributor(context)
	if err != nil {
		return 102, err
	}

	if willContribute {
		if err := contributor.Contribute(); err != nil {
			return 103, err
		}
	}

	return context.Success(buildpackplan.Plan{})
}
