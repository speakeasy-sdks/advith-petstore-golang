# Pets

### Available Operations

* [CreatePets](#createpets) - Create a pet
* [ListPets](#listpets) - List all pets
* [ShowPetByID](#showpetbyid) - Info for a specific pet

## CreatePets

Create a pet

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/speakeasy-sdks/advith-petstore-golang"
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

## ListPets

List all pets

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/speakeasy-sdks/advith-petstore-golang"
	"github.com/speakeasy-sdks/advith-petstore-golang/pkg/models/operations"
)

func main() {
    s := sdk.New()

    ctx := context.Background()
    res, err := s.Pets.ListPets(ctx, operations.ListPetsRequest{
        Limit: sdk.Int(548814),
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.Pets != nil {
        // handle response
    }
}
```

## ShowPetByID

Info for a specific pet

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/speakeasy-sdks/advith-petstore-golang"
	"github.com/speakeasy-sdks/advith-petstore-golang/pkg/models/operations"
)

func main() {
    s := sdk.New()

    ctx := context.Background()
    res, err := s.Pets.ShowPetByID(ctx, operations.ShowPetByIDRequest{
        PetID: "provident",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.Pet != nil {
        // handle response
    }
}
```
