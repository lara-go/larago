package http

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/lara-go/larago/logger"
	"github.com/olekukonko/tablewriter"

	"github.com/urfave/cli"
)

// CommandRoutes to apply DB changes.
type CommandRoutes struct {
	Router *Router
	Logger *logger.Logger
}

// GetCommand for the cli to register.
func (c *CommandRoutes) GetCommand() cli.Command {
	return cli.Command{
		Name:     "http:routes",
		Usage:    "List available HTTP routes",
		Category: "HTTP server",
	}
}

// Handle command.
func (c *CommandRoutes) Handle(args cli.Args) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Method", "Path", "Name", "Middleware"})
	table.SetColWidth(200) // Set max col width to 200. There may be lots of middleware or long path.

	routes := c.Router.GetRoutes()
	routesLen := len(routes)
	globalMiddleware := c.Router.GetMiddleware()

	for i := 0; i < routesLen; i++ {
		route := routes[i]
		table.Append(c.makeRow(route, globalMiddleware))
	}

	table.Render()

	return nil
}

func (c *CommandRoutes) makeRow(route *Route, globalMiddleware []Middleware) []string {
	allMiddleware := append(globalMiddleware, route.Middlewares...)

	middlewareLen := len(allMiddleware)
	middleware := make([]string, middlewareLen)

	for i := 0; i < middlewareLen; i++ {
		middleware[i] = fmt.Sprintf("%s", reflect.TypeOf(allMiddleware[i]).Elem())
	}

	return []string{route.Method, route.Path, route.Name, strings.Join(middleware, ", ")}
}
