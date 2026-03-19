package common

import (
	"fmt"
	"maps"
	"net/url"
	"os"
	"strings"
)

const OpenAiEndpointKey = "openai"

func ServerEndpoints(ctx *Context) (endpoints map[string]string, err error) {
	const httpEndpointsEnv = "HTTP_ENDPOINTS"
	endpoints = make(map[string]string)

	err = LoadEngineEnvironment(ctx)
	if err != nil {
		return nil, fmt.Errorf("error loading engine environment: %v", err)
	}

	if deprecatedBasePath, found := os.LookupEnv("OPENAI_BASE_PATH"); found {
		fmt.Printf("Warning: OPENAI_BASE_PATH environment variable is deprecated, use HTTP_OPENAI_BASE_PATH instead\n")

		// Migrate the values to new environment variables for backward compatibility
		err := os.Setenv("HTTP_OPENAI_BASE_PATH", deprecatedBasePath)
		if err != nil {
			return nil, fmt.Errorf("error setting HTTP_OPENAI_BASE_PATH: %v", err)
		}
		err = os.Setenv("HTTP_ENDPOINTS", "openai")
		if err != nil {
			return nil, fmt.Errorf("error setting HTTP_OPENAI_BASE_PATH: %v", err)
		}
	}

	httpEndpointNames, found := os.LookupEnv(httpEndpointsEnv)
	if found {
		httpEndpoints, err := serverHTTPEndpoints(ctx, httpEndpointNames)
		if err != nil {
			return nil, fmt.Errorf("error getting server http endpoints: %v", err)
		}
		maps.Copy(endpoints, httpEndpoints)
	}

	return endpoints, nil

}

func serverHTTPEndpoints(ctx *Context, httpEndpointNames string) (map[string]string, error) {
	const (
		confHttpPort            = "http.port"
		httpBasePathEnvTemplate = "HTTP_%s_BASE_PATH"
		defaultBasePath         = "/"
	)

	httpPortMap, err := ctx.Config.Get(confHttpPort)
	if err != nil {
		return nil, fmt.Errorf("error getting %q: %v", confHttpPort, err)
	}
	httpPort := httpPortMap[confHttpPort]

	endpoints := make(map[string]string)
	for endpointName := range strings.SplitSeq(httpEndpointNames, ",") {
		endpointName = strings.TrimSpace(endpointName)
		basePathEnv := fmt.Sprintf(httpBasePathEnvTemplate, strings.ToUpper(strings.TrimSpace(endpointName)))

		basePath, found := os.LookupEnv(basePathEnv)
		if !found {
			basePath = defaultBasePath
		}

		endpointUrl := url.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("localhost:%v", httpPort),
			Path:   basePath,
		}
		endpoints[endpointName] = endpointUrl.String()
	}

	return endpoints, nil
}
