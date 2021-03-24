package testing

import (
	"fmt"
	"github.com/gorilla/mux"
	types "github.com/queueup-dev/qup-types"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"testing"
	"time"
)

type Logger interface {
	Log(string)
}

type StdLogger int
type RouteHandler func(w http.ResponseWriter, r *http.Request)

type routeAssert struct {
}

func (l StdLogger) Log(s string) {
	log.Println(fmt.Sprintf("Logger %v : "+s, l))
}

type routesFunctions struct {
	expects  func(w http.ResponseWriter, r *http.Request)
	response func(w http.ResponseWriter, r *http.Request)
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
			builder.log.Log("failed assertion")
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
	h.inputValue = "response"

	return h.assertion
}

type DummyAPI struct {
	routes        map[string]*routesFunctions
	router        *mux.Router
	logger        Logger
	waitGroup     *sync.WaitGroup
	assertBuilder HttpAssertBuilder
}

func (api *DummyAPI) Assert(t *testing.T) *HttpAssertBuilder {
	return &HttpAssertBuilder{
		httpAssertions: []*HttpAssertion{},
		t:              t,
		log:            api.logger,
	}
}

func (api *DummyAPI) Expects(t *testing.T, uri string, method string, result string) {
	api.waitGroup.Add(1)

	api.addRoute(uri, method, func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)

		if err != nil {
			api.logger.Log(err.Error())
			t.Fail()
		}

		if string(b) != result {
			api.logger.Log(string(b) + " \n\n does not match the expected string : \n\n " + result)
			t.Fail()
		}

		api.waitGroup.Done()
	})
}

func (api *DummyAPI) addRoute(uri string, method string, handler RouteHandler) *mux.Route {
	return api.router.HandleFunc(uri, handler).Methods(method)
}

func (api *DummyAPI) RegisterReflectionRoute(uri string) {
	api.waitGroup.Add(1)

	api.routes[uri] = api.routes[uri].append(routesFunctions{
		response: func(w http.ResponseWriter, r *http.Request) {
			b, err := ioutil.ReadAll(r.Body)

			if err != nil {
				api.writeErrorResponse(err, w)
				return
			}
			w.WriteHeader(200)
			w.Write(b)

			api.waitGroup.Done()
		},
	})
}

func (api *DummyAPI) writeErrorResponse(err error, w http.ResponseWriter) {
	w.WriteHeader(500)
	w.Write([]byte(err.Error()))

	api.logger.Log(err.Error())
}

func (api *DummyAPI) RegisterRoute(uri string, response types.PayloadWriter, statusCode int) {

	api.waitGroup.Add(1)

	api.routes[uri] = api.routes[uri].append(routesFunctions{
		response: func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			log.Print(vars)

			response, err := response.ToString()

			if err != nil {
				api.writeErrorResponse(err, w)
				return
			}

			w.WriteHeader(statusCode)

			w.Write([]byte(*response))

			api.waitGroup.Done()
		},
	})

	api.logger.Log(strconv.Itoa(len(api.routes)))
}

func (f *routesFunctions) append(functions routesFunctions) *routesFunctions {
	if f == nil {
		return &functions
	}

	if f.response != nil && functions.response == nil {
		functions.response = f.response
	}

	if f.expects != nil && functions.expects == nil {
		functions.expects = f.expects
	}

	return &functions
}

func (f routesFunctions) compose() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f.expects(w, r)
		f.response(w, r)
	}
}

func (api *DummyAPI) Listen(address string) {
	r := mux.NewRouter()

	for uri, handler := range api.routes {
		api.logger.Log("called" + uri)
		r.HandleFunc(uri, handler.compose())
	}

	http.Handle("/", r)
	srv := &http.Server{
		Handler: r,
		Addr:    address,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	srv.ListenAndServe()
}

func NewDummyApi(l Logger, wg *sync.WaitGroup) DummyAPI {
	return DummyAPI{
		routes:        map[string]*routesFunctions{},
		logger:        l,
		waitGroup:     wg,
		assertBuilder: HttpAssertBuilder{},
	}
}
