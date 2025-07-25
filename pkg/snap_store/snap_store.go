package snap_store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func GetComponentsOfCurrentSnap() ([]SnapResources, error) {
	snapName := os.Getenv("SNAP_NAME")
	if snapName == "" {
		return nil, fmt.Errorf("error: SNAP_NAME must be set - likely not inside a snap")
	}
	snapInfo, err := Info(os.Getenv("SNAP_NAME"))
	if err != nil {
		return nil, fmt.Errorf("error getting snap info: %v", err)
	}
	snapRevision, err := strconv.ParseInt(os.Getenv("SNAP_REVISION"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing snap revision: %v", err)
	}
	components, err := Components(snapInfo.SnapId, int(snapRevision), os.Getenv("SNAP_ARCH"))
	if err != nil {
		return nil, fmt.Errorf("error getting components: %v", err)
	}
	return components, nil
}

func Info(snapName string) (SnapInfo, error) {
	info := SnapInfo{}
	// curl -H 'Snap-Device-Series: 24' http://api.snapcraft.io/v2/snaps/info/$SNAP_NAME
	url := "https://api.snapcraft.io/v2/snaps/info/" + snapName

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return info, fmt.Errorf("error creating new http request: %v", err)
	}

	req.Header.Add("Snap-Device-Series", "16")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return info, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Check for a successful HTTP status code
	if resp.StatusCode != http.StatusOK {
		return info, fmt.Errorf("received non-OK HTTP status: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return info, fmt.Errorf("error decoding JSON: %v", err)
	}

	return info, nil
}

// SNAP_REVISION SNAP_ARCH
func Components(snapId string, revision int, snapArch string) ([]SnapResources, error) {
	request := SnapRefreshRequest{
		Context: []SnapRefreshContext{
			{
				SnapId:          snapId,
				InstanceKey:     snapId,
				Revision:        revision,
				TrackingChannel: "",
			}},
		Actions: []SnapRefreshActions{
			{
				Action:      "refresh",
				InstanceKey: snapId,
				SnapId:      snapId,
				Revision:    revision,
			},
		},
		Fields: []string{"resources"},
	}
	requestJson, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request json: %v", err)
	}

	url := "https://api.snapcraft.io/v2/snaps/refresh"

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestJson))
	if err != nil {
		return nil, fmt.Errorf("error creating new http request: %v", err)
	}

	req.Header.Add("Snap-Device-Series", "16")
	req.Header.Add("Snap-Device-Architecture", snapArch)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Check for a successful HTTP status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK HTTP status: %d", resp.StatusCode)
	}

	var response SnapRefreshResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	if len(response.Results) == 0 {
		return nil, fmt.Errorf("no results found for snap id %s", snapId)
	}

	return response.Results[0].Snap.Resources, nil
}

func ComponentSize(components []SnapResources, componentName string) (int64, error) {
	for _, component := range components {
		if component.Name == componentName {
			return component.Download.Size, nil
		}
	}
	return 0, fmt.Errorf("component %s not found", componentName)
}
