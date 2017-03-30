package http_test

import (
	"fmt"
	"io/ioutil"
	"log"
	net_http "net/http"
	"testing"

	ozzo "github.com/go-ozzo/ozzo-validation"
	"github.com/lara-go/larago"
	"github.com/lara-go/larago/container"
	"github.com/lara-go/larago/database"
	"github.com/lara-go/larago/http"
	"github.com/lara-go/larago/http/errors"
	"github.com/lara-go/larago/http/responses"
	"github.com/lara-go/larago/logger"
	"github.com/lara-go/larago/support/testsuite"
	"github.com/lara-go/larago/validation"
)

func factory() *http.Router {
	container := container.New()
	logger := &logger.Logger{
		DateTimeFormat: larago.DateTimeFormat,
		DebugMode:      true,
		Logger:         log.New(ioutil.Discard, "", 0),
		// Logger: log.New(os.Stdout, "", 0),
	}

	router := http.NewRouter()
	router.Logger = logger
	router.Container = container
	router.ErrorsHandler = &http.ErrorsHandler{
		Logger: logger,
		ValidationErrorsConverter: &validation.OzzoErrorsConverter{},
	}

	return router
}

func TestCommonRequestRouter(t *testing.T) {
	router := factory()

	// Simple request.
	router.GET("/").Action(func(request *http.Request) string {
		return fmt.Sprintf("%s: Index", request.Method())
	})

	// Test POST.
	router.POST("/").Action(func(request *http.Request) int {
		return 12345
	})

	// Test PATCH.
	router.PATCH("/").Action(func(request *http.Request) float64 {
		return 1.2
	})

	// Test PUT.
	router.PUT("/").Action(func(request *http.Request) bool {
		return true
	})

	// Test DELETE.
	router.DELETE("/").Action(func(request *http.Request) bool {
		return false
	})

	// Simple request with custom middleware.
	router.GET("/foo").Action(func(request *http.Request) responses.Response {
		return responses.NewHTML(200, "").
			WithHeader("X-FOO", "bar").
			WithCookies(&net_http.Cookie{Name: "test", Value: "foobar"})
	})

	// Simple request with custom middleware.
	router.GET("/json").Action(func(request *http.Request) map[string]string {
		return map[string]string{"foo": "bar"}
	})

	e := testsuite.NewHTTPExpect(router.Bootstrap().GetHTTPRouter(), t)

	e.GET("/").
		Expect().Status(200).
		ContentType("text/plain", "utf-8").
		Body().Equal("GET: Index")

	e.POST("/").Expect().Status(200).Body().Equal("12345")
	e.PATCH("/").Expect().Status(200).Body().Equal("1.200000e+00")
	e.PUT("/").Expect().Status(200).Body().Equal("true")
	e.DELETE("/").Expect().Status(200).Body().Equal("false")

	resp := e.GET("/foo").
		Expect().Status(200).
		ContentType("text/html", "utf-8")
	resp.Cookies().Contains("test")
	resp.Headers().ValueEqual("X-Foo", []string{"bar"})

	e.GET("/json").
		Expect().Status(200).
		ContentType("application/json", "utf-8").
		JSON().Object().ValueEqual("foo", "bar")
}

func TestRedirect(t *testing.T) {
	router := factory()

	// Redirect.
	router.GET("/redirect").Action(func() responses.Response {
		return responses.NewRedirect(302).Route("foobar", map[string]string{"bar": "baz"})
	})

	// Aliased route with param
	router.GET("/foo/:bar").As("foobar").Action(func(bar string) responses.Response {
		return responses.NewText(200, "Foo Bar: %s", bar)
	})

	e := testsuite.NewHTTPExpect(router.Bootstrap().GetHTTPRouter(), t)

	e.GET("/redirect").Expect().Status(200).Body().Equal("Foo Bar: baz")
}

func TestErrors(t *testing.T) {
	router := factory()

	router.GET("/http-error").Action(func() *errors.HTTPError {
		return errors.NotFoundHTTPError().
			WithMeta(map[string]string{"foo": "Bar"}).
			WithContext("Custom context. Can be anything.")
	})

	router.GET("/custom-error").Action(func() error {
		return fmt.Errorf("Custom error")
	})

	router.GET("/panic").Action(func() {
		panic("Panic!")
	})

	router.GET("/model-not-found").Action(func() error {
		panic(&database.ModelNotFoundError{})
	})

	e := testsuite.NewHTTPExpect(router.Bootstrap().GetHTTPRouter(), t)

	e.GET("/http-error").Expect().Status(404)
	e.POST("/http-error").Expect().Status(405)

	e.GET("/custom-error").Expect().Status(500).Body().Equal("The server encountered an internal error or misconfiguration and was unable to complete your request.")
	e.GET("/custom-error").WithHeader("Accept", "application/json").
		Expect().Status(500).JSON().Object().ContainsKey("error")

	e.GET("/panic").Expect().Status(500)
	e.GET("/model-not-found").Expect().Status(404)
	e.GET("/not-found").Expect().Status(404)
}

type PostNewsJSON struct {
	Text string `json:"text"`
}

// ValidateJSON data.
func (r *PostNewsJSON) ValidateJSON() error {
	return r.Validate()
}

// Validate me.
func (r *PostNewsJSON) Validate() error {
	return ozzo.ValidateStruct(r,
		ozzo.Field(&r.Text, ozzo.Required),
	)
}

type PostNewsForm struct {
	Text string `schema:"text"`
}

// ValidateForm data.
func (r *PostNewsForm) ValidateForm() error {
	return r.Validate()
}

// Validate me.
func (r *PostNewsForm) Validate() error {
	return ozzo.ValidateStruct(r,
		ozzo.Field(&r.Text, ozzo.Required),
	)
}

type PostNewsQuery struct {
	Text string `schema:"text"`
}

// ValidateQuery data.
func (r *PostNewsQuery) ValidateQuery() error {
	return r.Validate()
}

// Validate me.
func (r *PostNewsQuery) Validate() error {
	return ozzo.ValidateStruct(r,
		ozzo.Field(&r.Text, ozzo.Required),
	)
}

type Form struct {
	Text string `schema:"text"`
}

func TestFormRequests(t *testing.T) {
	router := factory()
	router.SetArgsInjectors(&http.RequestsInjector{})

	// BarBaz accepts param, requesting it from unmarshalling to struct.
	router.GET("/bar/:text").Action(func(request *http.Request) responses.Response {
		var form Form
		request.ReadParams(&form)

		return responses.NewText(200, "Bar Baz: %s", form.Text)
	})

	// Form accepts POST form request.
	router.POST("/form").Action(func(request *http.Request) responses.Response {
		var form Form
		request.ReadForm(&form)

		return responses.NewHTML(200, "Form: <br>%+v", form)
	})

	// BarBaz accepts param, requesting it from unmarshalling to struct.
	router.GET("/query").Action(func(request *PostNewsQuery) responses.Response {
		return responses.NewHTML(200, "Query: %s", request.Text)
	}).Validate(&PostNewsQuery{})

	// JSON returns JSON formatted application/json response.
	router.POST("/json").Action(func(request *PostNewsJSON) responses.Response {
		return responses.NewJSON(200, map[string]string{"foo": request.Text})
	}).Validate(&PostNewsJSON{})

	e := testsuite.NewHTTPExpect(router.Bootstrap().GetHTTPRouter(), t)

	e.GET("/bar/text1").
		Expect().Status(200).
		Body().Contains("text1")

	e.POST("/form").WithFormField("text", "Text input").
		Expect().Status(200).
		Body().Contains("Text input")

	e.GET("/query").WithQuery("text", "Text").
		Expect().Status(200).
		Body().Contains("Query: Text")

	e.POST("/json").WithHeader("Accept", "application/json").WithJSON(map[string]interface{}{"text": "Text value"}).
		Expect().Status(200).
		JSON().Object().ValueEqual("foo", "Text value")

	e.POST("/json").WithHeader("Accept", "application/json").WithJSON(map[string]interface{}{"text": ""}).
		Expect().Status(422)
}

type ZeroMiddleware struct{}

func (m *ZeroMiddleware) Handle(request *http.Request, next http.Handler) responses.Response {
	response := next(request)

	return responses.NewText(response.Status(), fmt.Sprintf("%s Zero", response.Body()))
}

type FirstMiddleware struct{}

func (m *FirstMiddleware) Handle(request *http.Request, next http.Handler) responses.Response {
	response := next(request)

	return responses.NewText(response.Status(), fmt.Sprintf("%s First", response.Body()))
}

type SecondMiddleware struct{}

func (m *SecondMiddleware) Handle(request *http.Request, next http.Handler) responses.Response {
	response := next(request)

	return responses.NewText(response.Status(), fmt.Sprintf("%s Second", response.Body()))
}

func TestGroupsAndMiddleware(t *testing.T) {
	router := factory()

	router.Middleware(&ZeroMiddleware{})

	router.Group("/group1", func() {
		router.GET("/").Action(func() string {
			return "Group1 index"
		})

		router.GET("/path1").Action(func() string {
			return "Group1 path"
		}).Middleware(&SecondMiddleware{})

		router.Group("/group2", func() {
			router.GET("/path2").Action(func() string {
				return "Group2 path"
			})
		})
	}, &FirstMiddleware{})

	e := testsuite.NewHTTPExpect(router.Bootstrap().GetHTTPRouter(), t)
	e.GET("/group1").Expect().Status(200).Body().Equal("Group1 index First Zero")
	e.GET("/group1/path1").Expect().Status(200).Body().Equal("Group1 path Second First Zero")
	e.GET("/group1/group2/path2").Expect().Status(200).Body().Equal("Group2 path First Zero")
}

func TestBindings(t *testing.T) {
	router := factory()

	router.GET("/foo/:bar").Action(func(bar string) string {
		return bar
	})

	router.Bind("bar", func(param string) (interface{}, error) {
		if param == "bar" {
			return "baz", nil
		}

		return "bar", nil
	})

	e := testsuite.NewHTTPExpect(router.Bootstrap().GetHTTPRouter(), t)
	e.GET("/foo/bar").Expect().Status(200).Body().Equal("baz")
	e.GET("/foo/baz").Expect().Status(200).Body().Equal("bar")
}
