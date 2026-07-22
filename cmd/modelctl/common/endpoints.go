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
	openAiListenerKey = "openai"
)

type Listeners map[string]Listener

type Listener struct {
	Url        string `json:"url" yaml:"url"`
	ApiSpec    string `json:"api-spec,omitempty" yaml:"api-spec,omitempty"`
	UnixSocket string `json:"unix-socket,omitempty" yaml:"unix-socket,omitempty"`
}

func ServerListeners(ctx *Context) (Listeners, error) {
	activeEngineName, err := ctx.Cache.GetActiveEngine()
	if err != nil {
		return nil, fmt.Errorf("%s: %v", LookingUpActiveEngine, err)
	}
	activeEngineManifest, err := engines.LoadManifest(ctx.EnginesDir, activeEngineName)
	if err != nil {
		return nil, fmt.Errorf("loading active engine manifest: %v", err)
	}

	// If the engine does not list a runtime, return no listeners
	if activeEngineManifest.Runtime == "" {
		return nil, nil
	}

	runtimeManifest, err := runtimes.LoadManifest(ctx.RuntimesDir, activeEngineManifest.Runtime)
	if err != nil {
		return nil, fmt.Errorf("loading runtime manifest: %v", err)
	}

	listeners := make(Listeners)

	for serverName, serverSettings := range runtimeManifest.Servers {
		var listener *Listener
		var err error

		switch serverSettings.Protocol {
		case "http", "https":
			listener, err = serverHttpListener(ctx, serverSettings)
			if err != nil {
				return nil, fmt.Errorf("getting server HTTP listener: %v", err)
			}
		case "ws":
			listener, err = serverWsListener(ctx, serverName, serverSettings)
			if err != nil {
				return nil, fmt.Errorf("getting server WebSocket listener: %v", err)
			}
		case "ws+unix":
			listener, err = serverWsOverUnixListener(ctx, serverName, serverSettings)
			if err != nil {
				return nil, fmt.Errorf("getting server WebSocket Unix listener: %v", err)
			}
		default:
			return nil, fmt.Errorf("unsupported protocol %q for server %q in runtime %q",
				serverSettings.Protocol, serverName, activeEngineManifest.Runtime)
		}

		listeners[serverName] = *listener
	}

	// If builtin webui is enabled, also list it as well
	if WebUiEnabled() {
		webUiUrl, err := UiServerHttpUrl(ctx)
		if err != nil {
			return nil, fmt.Errorf("getting web UI url: %v", err)
		}
		listeners["webui"] = Listener{Url: webUiUrl}
	}

	return listeners, nil
}

func serverHttpListener(ctx *Context, server runtimes.Server) (*Listener, error) {
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

	httpHost, err := listenerHost(ctx, confHost)
	if err != nil {
		return nil, err
	}
	listenerUrl := url.URL{
		Scheme: server.Protocol,
		Host:   net.JoinHostPort(httpHost, fmt.Sprint(httpPort)),
		Path:   basePath,
	}

	return &Listener{Url: listenerUrl.String()}, nil
}

func serverWsListener(ctx *Context, serverName string, server runtimes.Server) (*Listener, error) {
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

	wsHost, err := listenerHost(ctx, confWsHost)
	if err != nil {
		return nil, err
	}

	listenerUrl := url.URL{
		Scheme: "ws",
		Host:   net.JoinHostPort(wsHost, fmt.Sprint(wsPort)),
		Path:   basePath,
	}

	return &Listener{Url: listenerUrl.String()}, nil
}

func serverWsOverUnixListener(ctx *Context, serverName string, server runtimes.Server) (*Listener, error) {
	// For ws+unix, use the same host/port configuration as regular ws
	// but also populate the UnixSocket field
	listener, err := serverWsListener(ctx, serverName, server)
	if err != nil {
		return nil, err
	}

	socketMap, err := ctx.Config.Get("ws.unix-socket")
	if err == nil {
		listener.UnixSocket = fmt.Sprint(socketMap["ws.unix-socket"])
	}

	return listener, nil
}

func OpenAiBaseUrl(ctx *Context) (string, error) {
	listeners, err := ServerListeners(ctx)
	if err != nil {
		return "", fmt.Errorf("getting server listeners: %v", err)
	}
	listener, found := listeners[openAiListenerKey]
	if !found {
		return "", fmt.Errorf("%q not found in server listeners", openAiListenerKey)
	}
	return listener.Url, nil
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

	httpHost, err := listenerHost(ctx, confWebuiHost)
	if err != nil {
		return "", err
	}

	listenerUrl := url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(httpHost, fmt.Sprint(httpPort)),
		Path:   defaultBasePath,
	}

	return listenerUrl.String(), nil
}

func listenerHost(ctx *Context, hostConfigKey string) (string, error) {
	hostMap, err := ctx.Config.Get(hostConfigKey)
	if err != nil {
		return "", fmt.Errorf("getting config %q: %v", hostConfigKey, err)
	}
	host := fmt.Sprint(hostMap[hostConfigKey])

	host = strings.TrimSpace(host)

	return host, nil
}
