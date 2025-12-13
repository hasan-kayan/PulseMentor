package httpx

import (
	"encoding/json"
	"io"
	"net/http"
)

func BindJSON(r *http.Request, v interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, v); err != nil {
		return err
	}

	return nil
}

