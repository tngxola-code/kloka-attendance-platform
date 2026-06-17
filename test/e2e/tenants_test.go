package e2e

import (
    "bytes"
    "encoding/json"
    "net/http"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

const baseURL = "http://localhost:8080/api/v1"

func TestTenantCreateAndGet(t *testing.T) {
    // This test requires the server to be running with a real DB and platform key.
    // For CI, you'd start a test database. For now, we assume manual run.
    platformKey := "test-platform-key"

    // Create tenant
    reqBody := map[string]interface{}{
        "name":          "E2E Tenant",
        "billing_email": "e2e@example.com",
        "country_code":  "ZA",
    }
    body, _ := json.Marshal(reqBody)
    req, _ := http.NewRequest("POST", baseURL+"/tenants", bytes.NewReader(body))
    req.Header.Set("X-Platform-Key", platformKey)
    req.Header.Set("Content-Type", "application/json")

    resp, err := http.DefaultClient.Do(req)
    require.NoError(t, err)
    defer resp.Body.Close()
    assert.Equal(t, http.StatusCreated, resp.StatusCode)

    var created struct {
        TenantID  string `json:"tenant_id"`
        Name      string `json:"name"`
        TenantKey string `json:"tenant_key"`
    }
    err = json.NewDecoder(resp.Body).Decode(&created)
    require.NoError(t, err)
    assert.NotEmpty(t, created.TenantID)
    assert.NotEmpty(t, created.TenantKey)

    // Get tenant by ID
    getResp, err := http.Get(baseURL + "/tenants/" + created.TenantID)
    require.NoError(t, err)
    defer getResp.Body.Close()
    assert.Equal(t, http.StatusOK, getResp.StatusCode)

    var getResult map[string]interface{}
    err = json.NewDecoder(getResp.Body).Decode(&getResult)
    require.NoError(t, err)
    assert.Equal(t, created.TenantID, getResult["id"])
}

func TestTenantList(t *testing.T) {
    // Requires server running and at least one tenant
    resp, err := http.Get(baseURL + "/tenants?limit=5")
    require.NoError(t, err)
    defer resp.Body.Close()
    assert.Equal(t, http.StatusOK, resp.StatusCode)

    var list struct {
        Data   []interface{} `json:"data"`
        Limit  int           `json:"limit"`
        Offset int           `json:"offset"`
        HasMore bool         `json:"has_more"`
    }
    err = json.NewDecoder(resp.Body).Decode(&list)
    require.NoError(t, err)
    assert.GreaterOrEqual(t, list.Limit, 0)
}
