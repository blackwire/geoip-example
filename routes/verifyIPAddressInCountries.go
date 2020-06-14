package routes

import (
	"encoding/json"
	"net"
	"net/http"
	"errors"
	geoip "github.com/oschwald/geoip2-golang"
	"os"
	"fmt"
)

const geoIPLanguage = "en"
const geoIPDatabaseFilePath = "data/geoipCountries.mmdb"

type verifyIPAddressInCountriesRequest struct {
	IPAddress string `json:"ipAddress"`
	AllowedCountries []string `json:"allowedCountries"`
}

type verifyIPAddressInCountriesResponse struct {
	IPAddress string `json:"ipAddress"`
	Allowed bool `json:"allowed"`
}

// VerifyIPAddressInCountriesRoute is a route that takes an IP Address and list of Countries and returns a JSON boolean value as to whether or not that IP is in that list of countries
type VerifyIPAddressInCountriesRoute struct {
	Route
	HandleError ErrorHandler
	dataFilePath string
}

// HandleRequest is the request and response handler for this route
func (rt *VerifyIPAddressInCountriesRoute) HandleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		err := errors.New("endpoint for VerifyIPAddressInCountriesRoute only accepts GET requests")
		rt.HandleError(http.StatusNotImplemented, err, w, r)
		return
	}

	ipAddress, allowedCountries, err := rt.parseRequest(r)
	if err != nil {
		rt.HandleError(http.StatusBadRequest, err, w, r)
		return
	}

	country, err := rt.getCountryUsing(ipAddress)
	if err != nil {
		rt.HandleError(http.StatusInternalServerError, err, w, r)
		return
	}

	response := &verifyIPAddressInCountriesResponse{
		IPAddress: ipAddress.String(),
		Allowed: false,
	}

	for _, allowedCountry := range allowedCountries {
		if country == allowedCountry {
			response.Allowed = true
		}
	}

	bytes, err := json.Marshal(response)
	if err != nil {
		rt.HandleError(http.StatusInternalServerError, err, w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func (rt *VerifyIPAddressInCountriesRoute) parseRequest(r *http.Request) (net.IP, []string, error) {
	decoder := json.NewDecoder(r.Body)
    var request verifyIPAddressInCountriesRequest
    err := decoder.Decode(&request)
    if err != nil {
		return nil, nil, err
	}

	ip := net.ParseIP(request.IPAddress)
	if ip == nil {
		return nil, nil, errors.New("ip address provided is not valid. update ipAddress field and try again")
	}

	return ip, request.AllowedCountries, nil
}

func (rt *VerifyIPAddressInCountriesRoute) getCountryUsing(ipAddress net.IP) (string, error) {
	err := rt.ensureGeoIPDataFilePath()
	if err != nil {
		return "", err
	}

	db, err := geoip.Open(rt.dataFilePath)
	if err != nil {
		return "", err
	}
	defer db.Close()

	record, err := db.Country(ipAddress)
	if err != nil {
		return "", err
	}

	return record.Country.Names[geoIPLanguage], nil
}

func (rt *VerifyIPAddressInCountriesRoute) ensureGeoIPDataFilePath() error {
	if rt.dataFilePath == "" {
		directory, err := os.Getwd()
		if err != nil {
			return err
		}

		rt.dataFilePath = fmt.Sprintf("%s/routes/%s", directory, geoIPDatabaseFilePath)
	}

	return nil
}