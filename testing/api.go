package testing

import (
	"github.com/gorilla/mux"
	qupHttp "github.com/queueup-dev/qup-io/http"
	types "github.com/queueup-dev/qup-types"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
)

const (
	inputTypeRequestBody = "REQUEST_BODY"
)

type Logger interface {
	Print(...interface{})
}

type StdLogger int

func (l StdLogger) Print(v ...interface{}) {
	log.Print(v...)
}

type HttpMockBuilder struct {
	mocks []*HttpMock
}

func (m *HttpMockBuilder) When(uri string, method string) *HttpMock {
	newMock := &HttpMock{
		routeUri:    uri,
		routeMethod: method,
		response:    nil,
	}

	m.mocks = append(m.mocks, newMock)

	return newMock
}

type HttpMock struct {
	routeUri    string
	routeMethod string
	response    *HttpMockResponse
}

func (h *HttpMock) RespondWith(body types.PayloadWriter, headers qupHttp.Headers, statusCode int) {
	response := &HttpMockResponse{
		headers:    headers,
		body:       body,
		statusCode: statusCode,
	}

	h.response = response
}

type HttpMockResponse struct {
	headers    qupHttp.Headers
	body       types.PayloadWriter
	statusCode int
}

type HttpAssertBuilder struct {
	httpAssertions []*HttpAssertion
	t              *testing.T
	log            Logger
}

func (builder *HttpAssertBuilder) That(uri string, method string) *HttpAssertion {
	assertion := HttpAssertion{
		routeUri:    uri,
		routeMethod: method,
		assertion:   &AssertInstance{},
	}

	builder.httpAssertions = append(builder.httpAssertions, &assertion)

	assertInstance := builder.httpAssertions[len(builder.httpAssertions)-1]
	return assertInstance
}

func (builder *HttpAssertBuilder) execute(input interface{}) bool {
	for _, httpAssertion := range builder.httpAssertions {
		if !httpAssertion.assertion.Execute(input) {
			builder.log.Print("failed assertion")
			builder.t.Fail()
		}
	}

	return true
}

type HttpAssertion struct {
	routeUri    string
	routeMethod string
	inputValue  string
	assertion   *AssertInstance
}

func (h *HttpAssertion) RequestBody() *AssertInstance {
	h.inputValue = inputTypeRequestBody

	return h.assertion
}

type MockAPI struct {
	router        *mux.Router
	logger        Logger
	waitGroup     *sync.WaitGroup
	assertBuilder HttpAssertBuilder
	mockBuilder   HttpMockBuilder
}

func (api *MockAPI) Assert() *HttpAssertBuilder {
	return &api.assertBuilder
}

func (api *MockAPI) Mock() *HttpMockBuilder {
	return &api.mockBuilder
}

func (api *MockAPI) composeAssertion(assertion *HttpAssertion) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if assertion.routeMethod != r.Method {
			// not for this method, just return
			return
		}

		switch assertion.inputValue {
		case inputTypeRequestBody:
			read, _ := ioutil.ReadAll(r.Body)
			if !assertion.assertion.Execute(string(read)) {
				api.assertBuilder.t.Fail()
			}
		}
	}
}

func (api *MockAPI) composeMock(mock *HttpMock) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if mock.routeMethod != r.Method {
			// not for this method, just return
			return
		}

		for key, val := range mock.response.headers {
			w.Header().Add(key, val)
		}

		w.WriteHeader(mock.response.statusCode)

		response, _ := mock.response.body.Marshal()
		w.Write(response.([]byte))
	}
}

func (api *MockAPI) compose(callbacks []func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	api.waitGroup.Add(1)
	return func(w http.ResponseWriter, r *http.Request) {
		for _, callback := range callbacks {
			callback(w, r)
		}

		api.waitGroup.Done()
	}
}

func (api *MockAPI) Listen(address string) error {
	r := mux.NewRouter()

	routes := map[string][]func(http.ResponseWriter, *http.Request){}
	for _, httpAssertion := range api.assertBuilder.httpAssertions {
		routes[httpAssertion.routeUri] = append(routes[httpAssertion.routeUri], api.composeAssertion(httpAssertion))
	}
	for _, httpMock := range api.mockBuilder.mocks {
		routes[httpMock.routeUri] = append(routes[httpMock.routeUri], api.composeMock(httpMock))
	}

	for route, callbacks := range routes {
		r.HandleFunc(route, api.compose(callbacks))
	}

	http.Handle("/", r)
	srv := &http.Server{
		Handler: r,
		Addr:    address,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return srv.ListenAndServe()
}

func NewMockApi(t *testing.T, l Logger, wg *sync.WaitGroup) MockAPI {
	return MockAPI{
		logger:    l,
		waitGroup: wg,
		assertBuilder: HttpAssertBuilder{
			log:            l,
			t:              t,
			httpAssertions: []*HttpAssertion{},
		},
		mockBuilder: HttpMockBuilder{
			mocks: []*HttpMock{},
		},
	}
}
