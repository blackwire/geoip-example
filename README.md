# geoip-example
A small go server that uses a geoIP database file to figure out if the country the provided IP address is from is in the list of allowed countries provided. Deals in JSON.

# example-request
{
	"ipAddress": "70.240.237.116",
	"allowedCountries": ["United Kingdom", "United States"]
}

# example-response-200
{
  "ipAddress": "70.240.237.116",
  "allowed": true
}

# example-response-error
## 501
{
  "httpStatusCode": 501,
  "errorMessage": "That request type is not available for this endpoint. Please select an available request type and try again."
}
## 400
{
  "httpStatusCode": 400,
  "errorMessage": "Failed to parse request. Please ensure request JSON is valid and try your request again."
}
