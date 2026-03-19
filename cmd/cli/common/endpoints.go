package common

import (
	"fmt"
	"net/url"
)

const (
	OpenAiEndpointKey = "openai"
	protocolKey       = "http"
	basePathKey       = "base-path"
)

func ServerEndpoints(ctx *Context) (map[string]string, error) {
	componentConfigs, err := LoadEngineEnvironment(ctx)
	if err != nil {
		return nil, fmt.Errorf("error loading engine environment: %v", err)
	}
	return serverEndpoints(ctx, componentConfigs)
}

func serverEndpoints(ctx *Context, componentConfigs []ComponentConfig) (map[string]string, error) {
	endpoints := make(map[string]string)
	for _, componentConfig := range componentConfigs {
		for serverName, serverConfig := range componentConfig.Servers {
			switch serverConfig[protocolKey] {
			case "http", "https":
				httpUrl, err := serverHttpUrl(ctx, serverConfig)
				if err != nil {
					return nil, fmt.Errorf("error getting server http endpoints: %v", err)
				}
				endpoints[serverName] = httpUrl
			default:
				return nil, fmt.Errorf("unsupported protocol %q for server %q", serverConfig["protocol"], serverName)
			}
		}
	}

	return endpoints, nil
}

func serverHttpUrl(ctx *Context, serverConfig map[string]string) (string, error) {
	const (
		confHttpPort    = "http.port"
		defaultBasePath = "/"
	)

	httpPortMap, err := ctx.Config.Get(confHttpPort)
	if err != nil {
		return "", fmt.Errorf("error getting %q: %v", confHttpPort, err)
	}
	httpPort := httpPortMap[confHttpPort]

	basePath, found := serverConfig[basePathKey]
	if !found {
		basePath = defaultBasePath
	}

	endpointUrl := url.URL{
		Scheme: serverConfig[protocolKey],
		Host:   fmt.Sprintf("localhost:%v", httpPort),
		Path:   basePath,
	}

	return endpointUrl.String(), nil
}
