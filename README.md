# Backend
Repository that holds best practices for structuring a backend micro-service.
"Best practices" are subjective to my personal experience and opinion.
At this point in time, the repository holds best practices for logging, testing, error handling, and using `Twirp` (i.e. grpc, protobufs).

## Dependencies
- [Twirp](https://twitchtv.github.io/twirp/docs/intro.html)
- [Cobra](https://github.com/spf13/cobra)

## Table of Contents
* [Requests](#requests)
* [Logs](#logs)
* [Testing](#testing)


## [Requests](#requests)
Using protobuf clients and [Twirp](https://twitchtv.github.io/twirp/docs/intro.html) will shield the request header from the method that will handle the endpoint being called, and will force us to use `context` instead.
This is not a huge issue, and allows for us to add an important piece of information to the context, an unique trace id.
With each request that comes through, we should be able to trace all the actions that the specific request entices.
To do this, let's take a look at some code inside of our [server](https://github.com/ykamo001/backend/blob/master/cmd/server.go#L33):

```go
router := mux.NewRouter().StrictSlash(false)
...
err := http.ListenAndServe(":8080", request.WithRequestHeaders(router))
```

Our router is wrapped around a function that will add a unique id, which we can refer to as a `trace id`, to the context so that each request can be identified.
This function serves as a [hook](https://github.com/ykamo001/backend/blob/master/request/hooks.go#L12) on the router and will inject this before the request gets to the endpoint handler function
```go
id, err := uuid.NewRandom()
if err == nil {
    ctx = context.WithValue(ctx, "id", id.ID())
}
```

Let's take a look at a function that handles a request on our server, in particular the `FillIn` [endpoint](https://github.com/ykamo001/backend/blob/master/internal/paint/provider.go#L24).
```go
func (p provider) FillIn(ctx context.Context, request *paintservice.FillInRequest) (*paintservice.FillInResponse, error) {
	entry := p.logger.WithFields(logrus.Fields{
		"id": ctx.Value("id"),
	})
    ...
}
```

We can see that each request that gets kicked off will now have an unique trace id we'll be able to use in our logging and debugging purposes when an error is encountered.
This is incredibly helpful since `Twirp` utilizes `grpc` and each request is invoked in its individual go-routine.

## [Logs](#logs)
Logs should help discern what the error was, where it occurred, what state the program was in such that it can be reproducible, and which object it effected (i.e. user, account, index in database, etc.).
A package that can be utilized to achieve all the aforementioned conditions is [logrus](https://github.com/sirupsen/logrus).
More importantly, using instances of a `logger`, and sharing and using those `loggers` instances within the entire program.

The motivation behind using instances of a logger instead of just using the package level `logrus` to log errors, is that each logger can can be customized for the needs of the specific portion of the program.
By adding hooks, for example, we can have our logs that are at the `error` or `warning` level automatically be pushed into [datadog](https://www.datadoghq.com/), where we already have settings such that `N` number of errors on `X` endpoint triggers [pagerduty](https://www.pagerduty.com/) for the on-call engineer to investigate.
This flow of reporting can be crucial when dealing with time-sensitive events. 
Moreover, if the method of logging is followed, the engineer tackling the error will easily be able to get to the root cause of the error. 

Before jumping into how to structure logs, make sure to have reviewed [requests](#requests) on utilizing `context`.

Let's take a look at the endpoint `FillIn`, which serves as a parody to `MS Paint Fill In` functionality.
`picture` in this example is a new instance of the `picture` struct.
```go
func (p provider) FillIn(ctx context.Context, request *paintservice.FillInRequest) (*paintservice.FillInResponse, error) {
	
    ...
	picture := NewPicture(matrix, p.logger)
	err := picture.FillIn(ctx, request.Value, request.X, request.Y)
	if err != nil {
		entry.Error("FillIn")
		return &paintservice.FillInResponse{}, twirp.NewError(twirp.Internal, "internal error")
	}
    ...
```
We can see that when the initial endpoint is hit, a call to `FillIn` is made on the `picture` struct.
```go
func (p picture) FillIn(ctx context.Context, value string, x, y int64) error {
	entry := p.logger.WithFields(logrus.Fields{
		"id":    ctx.Value("id"),
		"x":     x,
		"y":     y,
		"value": value,
	})

	if x < 0 || x >= int64(len(p.values)) {
		entry.WithField("maxRows", len(p.values)).Error("x is out of range")
		return errors.New("x is out of range")
	}
    ...
```

When this particular error is encountered, we should log it, and we will messages of the following kind
```go
{"file":".../backend/internal/paint/paint.go:32","func":"github.com/ykamo001/backend/internal/paint.picture.FillIn","id":704462431,"level":"error","maxRows":10,"msg":"x is out of range","time":"2019-11-03T21:40:47-08:00","value":"y","x":-2,"y":2}
{"file":".../backend/internal/paint/provider.go:41","func":"github.com/ykamo001/backend/internal/paint.provider.FillIn","id":704462431,"level":"error","msg":"FillIn","time":"2019-11-03T21:40:47-08:00"}
```

We can tell exactly where this error happened, what the state of the server and function was, what time it happened, and which request it can be tied back to. 
## [Testing](#testing)
Testing should be a crucial part of development and no code should be shipped unless there are tests written to cover their use.
We should also be able to let any developer test any particular test case without having them run all or none, much like table testing.
We can still leverage the advantages of table-testing, but restructure the way the tests are invoked. 

Let's take a look at the integration tests for the `MS Paint` [server handler](https://github.com/ykamo001/backend/blob/master/internal/paint/provider_integration_test.go).
Each one of these tests can be invoked separately along with all of them at once.

For running all tests on the provider:
`go test ./internal/paint -v --tags=integration -run Provider`

For running a specific test case on the provider:
`go test ./internal/paint -v --tags=integration -run Provider/success`
`go test ./internal/paint -v --tags=integration -run Provider/invalid_x_input`

For running all test:
`go test ./... -tags integration`