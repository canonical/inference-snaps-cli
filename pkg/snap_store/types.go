package snap_store

type SnapInfo struct {
	ChannelMap   []interface{} `json:"channel-map"`
	DefaultTrack interface{}   `json:"default-track"`
	Name         string        `json:"name"`
	Snap         struct {
		License string `json:"license"`
		Name    string `json:"name"`
		Prices  struct {
		} `json:"prices"`
		Publisher struct {
			DisplayName string `json:"display-name"`
			Id          string `json:"id"`
			Username    string `json:"username"`
			Validation  string `json:"validation"`
		} `json:"publisher"`
		SnapId   string `json:"snap-id"`
		StoreUrl string `json:"store-url"`
		Summary  string `json:"summary"`
		Title    string `json:"title"`
	} `json:"snap"`
	SnapId string `json:"snap-id"`
}

type SnapRefreshRequest struct {
	Context []SnapRefreshContext `json:"context"`
	Actions []SnapRefreshActions `json:"actions"`
	Fields  []string             `json:"fields"`
}

type SnapRefreshContext struct {
	SnapId          string `json:"snap-id"`
	InstanceKey     string `json:"instance-key"`
	Revision        int    `json:"revision"`
	TrackingChannel string `json:"tracking-channel"`
}

type SnapRefreshActions struct {
	Action      string `json:"action"`
	InstanceKey string `json:"instance-key"`
	SnapId      string `json:"snap-id"`
	Revision    int    `json:"revision"`
}

type SnapRefreshResponse struct {
	ErrorList []interface{} `json:"error-list"`
	Results   []struct {
		InstanceKey string      `json:"instance-key"`
		Name        string      `json:"name"`
		ReleasedAt  interface{} `json:"released-at"`
		Result      string      `json:"result"`
		Snap        struct {
			Resources []SnapResources `json:"resources"`
		} `json:"snap"`
		SnapId string `json:"snap-id"`
	} `json:"results"`
}

type SnapResources struct {
	Architectures []string `json:"architectures"`
	CreatedAt     string   `json:"created-at"`
	Description   string   `json:"description"`
	Download      struct {
		Sha3384 string `json:"sha3-384"`
		Size    int64  `json:"size"`
		Url     string `json:"url"`
	} `json:"download"`
	Name     string `json:"name"`
	Revision int    `json:"revision"`
	Type     string `json:"type"`
	Version  string `json:"version"`
}
