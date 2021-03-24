package testing

import (
	"fmt"
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
	dummyAPI := NewDummyApi(StdLogger(1), &wg)

	assertInstance := dummyAPI.Assert(t)
	assertInstance.That("test_uri", "GET").RequestBody().Eq("test123")

	fmt.Println(*assertInstance.httpAssertions[0])

}
