package routes

import (
	"net/http"
)

// ErrorHandler is used to common way of providing an error handler to all the routes
type ErrorHandler func(status int, err error, w http.ResponseWriter, r *http.Request)

// Route defines the required methods needed to be eligible to act as an endpoint for a request
type Route interface {
	HandleRequest(w http.ResponseWriter, r *http.Request)
}

// Setup is used to establish all the routes providing an error handler that will be used in the event of a problem
func Setup(errorHandler ErrorHandler) {
	routes := make(map[string]Route)

	routes["/verifyIPAddressInCountries"] = &VerifyIPAddressInCountriesRoute{
		HandleError: errorHandler,
	}

	for pattern, route := range routes {
		http.HandleFunc(pattern, route.HandleRequest)
	}
}