package operations

import (
	"github.com/speakeasy-sdks/advith-petstore-golang/pkg/models/shared"
	"net/http"
)

type CreatePetsResponse struct {
	ContentType string
	Error       *shared.Error
	StatusCode  int
	RawResponse *http.Response
}
