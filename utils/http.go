package utils

import (
	"bytes"
	"io"

	"github.com/mikespook/possum/log"
)

// BodyTrace logs the content of an HTTP request body and restores it for further processing.
func BodyTrace(body io.ReadCloser) io.ReadCloser {
	// Read and log the original request body
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		log.Trace().Err(err).Send()
		return nil
	} // Log the raw request body
	log.Trace().Str("raw_body", string(bodyBytes)).Send()
	// Restore the request body for further processing
	return io.NopCloser(bytes.NewBuffer(bodyBytes))
}
