package testing

import (
	"fmt"
	"github.com/queueup-dev/qup-io/writer"
	"sync"
	"testing"
)

func TestAssert_Eq(t *testing.T) {
	a := Assert{Input: int64(1)}

	fmt.Print(a.Eq(int64(1)))
}

func TestAssert_Same(t *testing.T) {
	a := Assert{Input: float64(1)}

	fmt.Print(a.Same(int64(1)))
}

func TestAssertInstance(t *testing.T) {
	a := AssertInstance{}

	a.Eq("test")

	fmt.Print(a.Execute("test"))
}

func TestDummyAPI_Assert(t *testing.T) {
	var wg sync.WaitGroup
	dummyAPI := NewDummyApi(t, StdLogger(1), &wg)

	dummyAPI.Assert().That("test_uri", "GET").RequestBody().Eq("test123")

	dummyAPI.Mock().When("test_uri", "GET").RespondWith(
		writer.NewJsonWriter(struct{ Hello string }{Hello: "world"}), nil, 200,
	)

}
