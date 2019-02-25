package e2e

import (
	"fmt"
	"path/filepath"
	"regexp"
	"testing"

	"gotest.tools/assert"
	"gotest.tools/fs"
	"gotest.tools/icmd"
)

const (
	config = `{
		"cliPluginsExtraDirs": ["%s"]
}`
	help = `Usage:	docker app COMMAND

Build and deploy Docker Application Packages.

Commands:
  bundle      Create a CNAB invocation image and bundle.json for the application.
  completion  Generates completion scripts for the specified shell (bash or zsh)
  init        Start building a Docker application
  inspect     Shows metadata, parameters and a summary of the compose file for a given application
  install     Install an application
  merge       Merge a multi-file application into a single file
  push        Push the application to a registry
  render      Render the Compose file for the application
  split       Split a single-file application into multiple files
  status      Get the installation status. If the installation is a docker application, the status shows the stack services.
  uninstall   Uninstall an application
  upgrade     Upgrade an installed application
  validate    Checks the rendered application is syntactically correct
  version     Print version information

Run 'docker app COMMAND --help' for more information on a command.`
)

func TestInvokePluginFromCLI(t *testing.T) {
	dir := fs.NewDir(t, t.Name())
	defer dir.Remove()
	relPath, err := filepath.Rel(dir.Path(), dockerApp)
	assert.NilError(t, err)
	fs.Apply(t, dir,
		fs.WithFile("config.json", fmt.Sprintf(config, dir.Path())),
		fs.WithSymlink(filepath.Base(dockerApp), relPath))

	// docker --help should list app as a top command
	icmd.RunCommand(dockerCli, "--config", dir.Path(), "--help").Assert(t, icmd.Expected{
		Out: "app*        Docker Application Packages (Docker Inc.,",
	})

	// docker app --help prints docker-app help
	icmd.RunCommand(dockerCli, "--config", dir.Path(), "app", "--help").Assert(t, icmd.Expected{
		Out: help,
	})

	// docker info should print app version and short description
	re := regexp.MustCompile(`app: \(.*, Docker Inc\.\) Docker Application Packages`)
	output := icmd.RunCommand(dockerCli, "--config", dir.Path(), "info").Assert(t, icmd.Success).Combined()
	assert.Assert(t, re.MatchString(output))
}
