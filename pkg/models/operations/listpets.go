package operations

import (
	"github.com/speakeasy-sdks/advith-petstore-golang/pkg/models/shared"
	"net/http"
)

type ListPetsQueryParams struct {
	Limit *int `queryParam:"style=form,explode=true,name=limit"`
}

type ListPetsRequest struct {
	QueryParams ListPetsQueryParams
}

type ListPetsResponse struct {
	ContentType string
	Error       *shared.Error
	Pets        []shared.Pet
	StatusCode  int
	RawResponse *http.Response
}
