package environments

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
)

func TestViewEnvironment(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	orgID := testConstants["orgId"]
	environmentName := testConstants["environmentName"]
	environmentID := "5d8f6f1a-9c2b-4d7e-8a41-3f0e2b7c9d10"

	tests := []utils.MesheryCommandTest{
		{
			Name:             "given too many arguments when running environment view then return error",
			Args:             []string{"view", environmentID, "extra-arg", "--orgId", orgID},
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			URL:              fmt.Sprintf("/%s", environmentApiPath),
			Fixture:          "view.environment.api.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			ExpectedError:    utils.ErrInvalidArgument(errors.New("please provide at most one environment name or ID\n\nUsage: mesheryctl environment view [environment-name|environment-id] --orgId [orgId]\nRun 'mesheryctl environment view --help' to see detailed help message")),
		},
		{
			Name:             "given no orgId when running environment view then return error",
			Args:             []string{"view", environmentID},
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			URL:              fmt.Sprintf("/%s", environmentApiPath),
			Fixture:          "view.environment.api.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			ExpectedError:    utils.ErrInvalidArgument(errors.New("[ orgId ] isn't specified\n\nUsage: mesheryctl environment view [environment-name|environment-id] --orgId [orgId]\nRun 'mesheryctl environment view --help' to see detailed help message")),
		},
		{
			Name:             "given invalid orgId when running environment view then return error",
			Args:             []string{"view", environmentID, "--orgId", "not-a-uuid"},
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			URL:              fmt.Sprintf("/%s", environmentApiPath),
			Fixture:          "view.environment.api.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			ExpectedError:    utils.ErrInvalidUUID(fmt.Errorf("invalid orgId: %s", "not-a-uuid")),
		},
		{
			Name:             "given environment ID when running environment view then display environment",
			Args:             []string{"view", environmentID, "--orgId", orgID},
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			URL:              fmt.Sprintf("/%s/%s?orgId=%s", environmentApiPath, environmentID, orgID),
			Fixture:          "view.environment.api.response.golden",
			ExpectedResponse: "view.environment.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "given environment ID from another organization when running environment view then return not found error",
			Args:             []string{"view", environmentID, "--orgId", orgID},
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			URL:              fmt.Sprintf("/%s/%s?orgId=%s", environmentApiPath, environmentID, orgID),
			Fixture:          "view.environment.other.org.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			ExpectedError:    utils.ErrNotFound(fmt.Errorf("No environment found with ID %s in organization: %s", environmentID, orgID)),
		},
		{
			Name:             "given environment name with a single match when running environment view then display environment",
			Args:             []string{"view", environmentName, "--orgId", orgID},
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			URL:              fmt.Sprintf("/%s?orgId=%s&page=0&pagesize=10&search=%s", environmentApiPath, orgID, environmentName),
			Fixture:          "view.environment.search.single.response.golden",
			ExpectedResponse: "view.environment.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "given no environment argument and a single environment when running environment view then display environment",
			Args:             []string{"view", "--orgId", orgID},
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			URL:              fmt.Sprintf("/%s?orgId=%s&page=0&pagesize=10", environmentApiPath, orgID),
			Fixture:          "view.environment.search.single.response.golden",
			ExpectedResponse: "view.environment.output.golden",
			ExpectError:      false,
		},
		{
			// The server reports the total as camelCase totalCount (v1beta3). It must be
			// read as 11, not 0, so the command pages instead of reporting nothing found.
			Name:             "given more environments than one page when running environment view without a terminal then report every match",
			Args:             []string{"view", "--orgId", orgID},
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			URL:              fmt.Sprintf("/%s?orgId=%s&page=0&pagesize=10", environmentApiPath, orgID),
			Fixture:          "view.environment.first.page.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			ExpectedError:    display.ErrAmbiguousSelection(11),
		},
		{
			Name:             "given unknown environment name when running environment view then return not found error",
			Args:             []string{"view", "does-not-exist", "--orgId", orgID},
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			URL:              fmt.Sprintf("/%s?orgId=%s&page=0&pagesize=10&search=%s", environmentApiPath, orgID, "does-not-exist"),
			Fixture:          "view.environment.search.empty.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			ExpectedError:    utils.ErrNotFound(errors.New("No environment(s) found with name: does-not-exist")),
		},
	}

	utils.InvokeMesheryctlTestCommand(t, update, EnvironmentCmd, tests, currDir, "environment")
}
