package httpx

import (
    "encoding/json"
    "net/http"
)

type Problem struct {
    Type     string `json:"type"`
    Title    string `json:"title"`
    Status   int    `json:"status"`
    Detail   string `json:"detail,omitempty"`
    Instance string `json:"instance,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

func WriteProblem(w http.ResponseWriter, r *http.Request, status int, title, detail string) {
    prob := Problem{
        Type:     "about:blank",
        Title:    title,
        Status:   status,
        Detail:   detail,
        Instance: r.URL.Path,
    }
    w.Header().Set("Content-Type", "application/problem+json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(prob)
}
