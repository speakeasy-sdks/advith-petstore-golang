# github.com/speakeasy-sdks/advith-petstore-golang

<!-- Start SDK Installation -->
## SDK Installation

```bash
go get github.com/speakeasy-sdks/advith-petstore-golang
```
<!-- End SDK Installation -->

## SDK Example Usage
<!-- Start SDK Example Usage -->
```go
package main

import (
    "context"
    "log"
    "github.com/speakeasy-sdks/advith-petstore-golang"
    "github.com/speakeasy-sdks/advith-petstore-golang/pkg/models/shared"
    "github.com/speakeasy-sdks/advith-petstore-golang/pkg/models/operations"
)

func main() {
    s := sdk.New()

    ctx := context.Background()
    res, err := s.Pets.CreatePets(ctx)
    if err != nil {
        log.Fatal(err)
    }

    if res.StatusCode == http.StatusOK {
        // handle response
    }
}
```
<!-- End SDK Example Usage -->

<!-- Start SDK Available Operations -->
## SDK Available Operations


### Pets

* `CreatePets` - Create a pet
* `ListPets` - List all pets
* `ShowPetByID` - Info for a specific pet
<!-- End SDK Available Operations -->

### SDK Generated by [Speakeasy](https://docs.speakeasyapi.dev/docs/using-speakeasy/client-sdks)