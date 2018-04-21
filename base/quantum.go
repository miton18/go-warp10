package base

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// GetQuantumLink return an URL on Quantum instance with configured Backend and script
func (c *Client) GetQuantumLink(warpScript string) string {
	return fmt.Sprintf("%s/#/warpscript/%s/%s", c.QuantumInstance, getBackend(c), getWarpScript(warpScript))
}

func getBackend(c *Client) string {
	body, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	return base64.RawStdEncoding.EncodeToString(body)
}

func getWarpScript(ws string) string {
	return base64.RawStdEncoding.EncodeToString([]byte(ws))
}
