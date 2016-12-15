package archive

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// SnapshotURL is the URL that the snapshot request is sent to.
const SnapshotURL = "https://pragma.archivelab.org"

// JSONError is an error returned by a JSON endpoint.
type JSONError struct {
	Error string `json:"error"`
}

// SnapshotRequest is the format the snapshot API accepts.
type SnapshotRequest struct {
	URL string `json:"url"`
}

// SnapshotResponse is the response from the snapshot API.
type SnapshotResponse struct {
	AnnotationID string `json:"annotation_id"`
	Domain       string `json:"domain"`
	ID           int    `json:"id"`
	Path         string `json:"path"`
	Protocol     string `json:"protocol"`
	WaybackID    string `json:"wayback_id"`
}

// Snapshot requests the internet archive to record the URL.
func Snapshot(url string) (SnapshotResponse, error) {
	req := SnapshotRequest{
		URL: url,
	}
	var resp SnapshotResponse
	err := postJSON(SnapshotURL, req, &resp)
	return resp, err
}

func postJSON(url string, req, resp interface{}) error {
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return err
	}
	post, err := http.NewRequest("POST", url, bytes.NewBuffer(reqJSON))
	post.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	postResp, err := client.Do(post)
	if err != nil {
		return err
	}
	defer postResp.Body.Close()

	if postResp.StatusCode != http.StatusOK {
		return errors.Errorf("expected status OK: got %d", postResp.StatusCode)
	}

	body, _ := ioutil.ReadAll(postResp.Body)

	var jsonErr JSONError
	if err := json.Unmarshal(body, &jsonErr); err != nil {
		return err
	}
	if len(jsonErr.Error) > 0 {
		return errors.New(jsonErr.Error)
	}

	return json.Unmarshal(body, resp)
}
