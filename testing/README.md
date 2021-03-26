## Usage

### Configuring a Mock API
This will allow for doing a full integration test using a configurable mock API.
```go
package test

import (
    qupTest "github.com/queueup-dev/qup-io/testing"
    "github.com/queueup-dev/qup-io/writer"
)

var (
    wg      sync.WaitGroup
    logger     = qupTest.StdLogger(1)
    httpClient = qupHttp.NewDefaultClient()
)

func TestFunction(t *testing.T) {
    mockAPI := qupTest.NewMockApi(t, logger, &wg)
    
    mockAPI.Mock().When("/products", "GET").RespondWith(
        writer.NewJsonWriter(struct{Hello string}{ Hello: "World"}),
        qupHttp.Headers{},
        200
    )
    
    go mockAPI.Listen("localhost:8000")
    
    // ... do your calls and asserts
}

func TestFunctionError(t *testing.T) {
    mockAPI := qupTest.NewMockApi(t, logger, &wg)
    
    mockAPI.Mock().When("/products", "GET").RespondWith(
    writer.NewJsonWriter(nil),
    qupHttp.Headers{},
    503
    )
    
    go mockAPI.Listen("localhost:8000")

    // ... do your calls and asserts
}
```

### Setting up Assertions
This will allow you to test your client calls on the integration (API) level.
```go
package test

import (
	qupTest "github.com/queueup-dev/qup-io/testing"
	"github.com/queueup-dev/qup-io/writer"
)

func TestFunctionError(t *testing.T) {
    mockAPI := qupTest.NewMockApi(t, logger, &wg)
    
    mockAPI.Assert().That("/products", "POST").RequestBody().Eq(
    	"{\"Hello\": \"World\"}",
    )
    
    go mockAPI.Listen("localhost:8000")
    
    // ... do your calls and the assert will fire on call
}
```