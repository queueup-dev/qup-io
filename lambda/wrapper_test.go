package lambda

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"testing"
)

type TestInput struct {
	Foo string
}

type TestOutput struct {
	Bar string
}

var (
	ctx = context.Background()
)

type TestHandler struct{}

func (t TestHandler) Invoke(ctx context.Context, payload TestInput) (*TestOutput, error) {
	return &TestOutput{Bar: payload.Foo}, nil
}

func (t TestHandler) ErrorInvoke(ctx context.Context, payload TestInput) (*TestOutput, error) {
	return nil, errors.New("foo error")
}

func TestWrapHandler(t *testing.T) {
	handler := wrapHandler(TestHandler{}.Invoke)

	testCargo := "{ \"Foo\": \"Hello World!\"}"
	result, err := handler.Invoke(ctx, []byte(testCargo))

	if err != nil {
		log.Print(err)
		t.Fail()
	}

	var unmarshalledResponse TestOutput
	err = json.Unmarshal(result, &unmarshalledResponse)

	if err != nil {
		log.Print(err)
		t.Fail()
	}

	if unmarshalledResponse.Bar != "Hello World!" {
		log.Print("values of the test handler do not match")
		t.Fail()
	}
}

func TestWrapHandlerWithLogicError(t *testing.T) {
	handler := wrapHandler(TestHandler{}.ErrorInvoke)
	testCargo := "{ \"Foo\": \"Hello World!\"}"

	result, err := handler.Invoke(ctx, []byte(testCargo))

	if err == nil {
		log.Print("error expected but none received")
		t.Fail()
	}

	if result != nil {
		log.Print("nil cargo expected")
		t.Fail()
	}

	if err == nil || err.Error() != "foo error" {
		log.Print("unexpected error received")
		t.Fail()
	}
}

func TestWrapHandlerFailInvalidJson(t *testing.T) {
	handler := wrapHandler(TestHandler{}.Invoke)
	brokenCargo := "{ \"Foo\" \"Hello World!\"}"

	result, err := handler.Invoke(ctx, []byte(brokenCargo))

	if err == nil {
		log.Print("broken json body should not succeed")
		t.Fail()
	}

	if err != nil && err.Error() != "invalid character '\"' after object key" {
		log.Print("invalid error returned")
		t.Fail()
	}

	if result != nil {
		log.Print("broken json body should not provide result")
		t.Fail()
	}
}

func TestWrapEpsagon(t *testing.T) {
	testCargo := "{ \"Foo\": \"Hello World!\"}"
	epsagonWrapped, err := wrapEpsagon(TestHandler{}.Invoke)(ctx, []byte(testCargo))

	if err != nil {
		log.Print(err)
		t.Fail()
	}

	var unmarshalledResponse TestOutput
	err = json.Unmarshal(epsagonWrapped, &unmarshalledResponse)

	if err != nil {
		log.Print(err)
		t.Fail()
	}

	if unmarshalledResponse.Bar != "Hello World!" {
		log.Print("values of the test handler do not match")
		t.Fail()
	}
}

func init() {
	_ = os.Setenv("EPSAGON_APP_ID", "xxx")
	_ = os.Setenv("EPSAGON_TOKEN", "xxx")
}
