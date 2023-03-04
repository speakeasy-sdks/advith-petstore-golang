package operations

import (
	"github.com/speakeasy-sdks/advith-petstore-golang/pkg/models/shared"
	"net/http"
)

type ShowPetByIDPathParams struct {
	PetID string `pathParam:"style=simple,explode=false,name=petId"`
}

type ShowPetByIDRequest struct {
	PathParams ShowPetByIDPathParams
}

type ShowPetByIDResponse struct {
	ContentType string
	Error       *shared.Error
	Pet         *shared.Pet
	StatusCode  int
	RawResponse *http.Response
}
