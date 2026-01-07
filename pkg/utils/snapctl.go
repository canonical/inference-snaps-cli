package utils

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// Add the service struct
type service struct {
	Enabled string
	Active  string
	Notes   string
}

func GetServices() (map[string]service, error) {
	socketPath := "/run/snapd-snap.socket"
	url := "http://localhost/v2/snapctl"

	snapContext := os.Getenv("SNAP_CONTEXT")
	if snapContext == "" {
		return nil, fmt.Errorf("snap context not set")
	}

	payload := map[string]interface{}{
		"context-id": snapContext,
		"args":       []string{"services"},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, err error) {
				return (&net.Dialer{}).DialContext(ctx, "unix", socketPath)
			},
		},
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Language", "en-US")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("snapctl request failed: %s", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result struct {
		Result struct {
			Stdout string `json:"stdout"`
		} `json:"result"`
	}
	if err = json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse snapctl response: %s", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(result.Result.Stdout))
	// throw away the header
	scanner.Scan()

	services := make(map[string]service)
	re := regexp.MustCompile("[[:space:]]+")
	for scanner.Scan() {
		line := scanner.Text()
		cells := re.Split(line, 4)
		if len(cells) != 4 {
			continue // skip lines not matching expected format
		}
		serviceName := cells[0]
		services[serviceName] = service{
			Enabled: cells[1],
			Active:  cells[2],
			Notes:   cells[3],
		}
	}
	return services, nil
}
