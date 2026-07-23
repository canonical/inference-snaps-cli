package common

import (
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
	"github.com/canonical/inference-snaps-cli/v2/pkg/runtimes"
)

const (
	openAiEntrypointKey = "openai"
)

type Entrypoints map[string]Entrypoint

type Entrypoint struct {
	Url        string `json:"url" yaml:"url"`
	ApiSpec    string `json:"api-spec,omitempty" yaml:"api-spec,omitempty"`
	UnixSocket string `json:"unix-socket,omitempty" yaml:"unix-socket,omitempty"`
}

func ServerEntrypoints(ctx *Context) (Entrypoints, error) {
	activeEngineName, err := ctx.Cache.GetActiveEngine()
	if err != nil {
		return nil, fmt.Errorf("%s: %v", LookingUpActiveEngine, err)
	}
	activeEngineManifest, err := engines.LoadManifest(ctx.EnginesDir, activeEngineName)
	if err != nil {
		return nil, fmt.Errorf("loading active engine manifest: %v", err)
	}

	// If the engine does not list a runtime, return no entrypoints
	if activeEngineManifest.Runtime == "" {
		return nil, nil
	}

	runtimeManifest, err := runtimes.LoadManifest(ctx.RuntimesDir, activeEngineManifest.Runtime)
	if err != nil {
		return nil, fmt.Errorf("loading runtime manifest: %v", err)
	}

	entrypoints := make(Entrypoints)

	for serverName, serverSettings := range runtimeManifest.Servers {
		var entrypoint *Entrypoint
		var err error

		switch serverSettings.Protocol {
		case "http", "https":
			entrypoint, err = serverHttpEntrypoint(ctx, serverSettings)
			if err != nil {
				return nil, fmt.Errorf("getting server HTTP entrypoint: %v", err)
			}
		case "ws":
			entrypoint, err = serverWsEntrypoint(ctx, serverName, serverSettings)
			if err != nil {
				return nil, fmt.Errorf("getting server WebSocket entrypoint: %v", err)
			}
		case "ws+unix":
			entrypoint, err = serverWsOverUnixEntrypoint(ctx, serverName, serverSettings)
			if err != nil {
				return nil, fmt.Errorf("getting server WebSocket Unix entrypoint: %v", err)
			}
		default:
			return nil, fmt.Errorf("unsupported protocol %q for server %q in runtime %q",
				serverSettings.Protocol, serverName, activeEngineManifest.Runtime)
		}

		entrypoints[serverName] = *entrypoint
	}

	// If builtin webui is enabled, also list it as well
	if WebUiEnabled() {
		webUiUrl, err := UiServerHttpUrl(ctx)
		if err != nil {
			return nil, fmt.Errorf("getting web UI url: %v", err)
		}
		entrypoints["webui"] = Entrypoint{Url: webUiUrl}
	}

	return entrypoints, nil
}

func serverHttpEntrypoint(ctx *Context, server runtimes.Server) (*Entrypoint, error) {
	const (
		confHttpPort    = "http.port"
		defaultBasePath = "/"
		confHost        = "http.host"
	)

	httpPortMap, err := ctx.Config.Get(confHttpPort)
	if err != nil {
		return nil, fmt.Errorf("getting config %q: %v", confHttpPort, err)
	}
	httpPort := httpPortMap[confHttpPort]

	basePath := server.BasePath
	if basePath == "" {
		basePath = defaultBasePath
	}

	httpHost, err := entrypointHost(ctx, confHost)
	if err != nil {
		return nil, err
	}
	entrypointUrl := url.URL{
		Scheme: server.Protocol,
		Host:   net.JoinHostPort(httpHost, fmt.Sprint(httpPort)),
		Path:   basePath,
	}

	return &Entrypoint{Url: entrypointUrl.String()}, nil
}

func serverWsEntrypoint(ctx *Context, serverName string, server runtimes.Server) (*Entrypoint, error) {
	const (
		confWsPort      = "ws.port"
		defaultBasePath = "/"
		confWsHost      = "ws.host"
	)

	wsPortMap, err := ctx.Config.Get(confWsPort)
	if err != nil {
		return nil, fmt.Errorf("getting config %q: %v", confWsPort, err)
	}
	wsPort := wsPortMap[confWsPort]

	basePath := server.BasePath
	if basePath == "" {
		basePath = defaultBasePath
	}

	wsHost, err := entrypointHost(ctx, confWsHost)
	if err != nil {
		return nil, err
	}

	entrypointUrl := url.URL{
		Scheme: "ws",
		Host:   net.JoinHostPort(wsHost, fmt.Sprint(wsPort)),
		Path:   basePath,
	}

	return &Entrypoint{Url: entrypointUrl.String()}, nil
}

func serverWsOverUnixEntrypoint(ctx *Context, serverName string, server runtimes.Server) (*Entrypoint, error) {
	// For ws+unix, use the same host/port configuration as regular ws
	// but also populate the UnixSocket field
	entrypoint, err := serverWsEntrypoint(ctx, serverName, server)
	if err != nil {
		return nil, err
	}

	socketMap, err := ctx.Config.Get("ws.unix-socket")
	if err == nil {
		entrypoint.UnixSocket = fmt.Sprint(socketMap["ws.unix-socket"])
	}

	return entrypoint, nil
}

func OpenAiBaseUrl(ctx *Context) (string, error) {
	entrypoints, err := ServerEntrypoints(ctx)
	if err != nil {
		return "", fmt.Errorf("getting server entrypoints: %v", err)
	}
	entrypoint, found := entrypoints[openAiEntrypointKey]
	if !found {
		return "", fmt.Errorf("%q not found in server entrypoints", openAiEntrypointKey)
	}
	return entrypoint.Url, nil
}

func UiServerHttpUrl(ctx *Context) (string, error) {
	const (
		confWebuiHttpPort = "webui.http.port"
		defaultBasePath   = "/"
		confWebuiHost     = "webui.http.host"
	)

	httpPortMap, err := ctx.Config.Get(confWebuiHttpPort)
	if err != nil {
		return "", fmt.Errorf("getting config %q: %v", confWebuiHttpPort, err)
	}
	httpPort := httpPortMap[confWebuiHttpPort]

	httpHost, err := entrypointHost(ctx, confWebuiHost)
	if err != nil {
		return "", err
	}

	entrypointUrl := url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(httpHost, fmt.Sprint(httpPort)),
		Path:   defaultBasePath,
	}

	return entrypointUrl.String(), nil
}

func entrypointHost(ctx *Context, hostConfigKey string) (string, error) {
	hostMap, err := ctx.Config.Get(hostConfigKey)
	if err != nil {
		return "", fmt.Errorf("getting config %q: %v", hostConfigKey, err)
	}
	host := fmt.Sprint(hostMap[hostConfigKey])

	host = strings.TrimSpace(host)

	return host, nil
}
