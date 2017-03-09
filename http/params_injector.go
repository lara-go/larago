package http

// RouteParamsInjector injects raw route params.
type RouteParamsInjector struct{}

// Inject custom params to for the action.
func (i *RouteParamsInjector) Inject(params []interface{}, request *Request) ([]interface{}, error) {
	for _, param := range request.Route.Params {
		params = append(params, param.Value)
	}

	return params, nil
}
