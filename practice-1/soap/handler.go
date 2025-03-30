package soap

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Current exchange rates (for demo purposes)
const (
	UAHtoUSDRate = 0.025 // 1 UAH = 0.025 USD
	USDtoUAHRate = 40.0  // 1 USD = 40 UAH
)

// HandleCurrencyConversion processes a SOAP currency conversion request
func HandleCurrencyConversion(c *gin.Context) {
	// Check content type
	if c.GetHeader("Content-Type") != "text/xml" {
		c.XML(http.StatusBadRequest, SOAPEnvelope{
			Body: SOAPBody{
				Fault: &SOAPFault{
					FaultCode:   "Client",
					FaultString: "Invalid Content-Type. Expected text/xml",
				},
			},
		})
		return
	}

	// Read the request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.XML(http.StatusBadRequest, SOAPEnvelope{
			Body: SOAPBody{
				Fault: &SOAPFault{
					FaultCode:   "Client",
					FaultString: "Failed to read request body",
					Detail:      err.Error(),
				},
			},
		})
		return
	}

	// Parse the SOAP envelope
	var envelope SOAPEnvelope
	if err := xml.Unmarshal(body, &envelope); err != nil {
		c.XML(http.StatusBadRequest, SOAPEnvelope{
			Body: SOAPBody{
				Fault: &SOAPFault{
					FaultCode:   "Client",
					FaultString: "Failed to parse SOAP envelope",
					Detail:      err.Error(),
				},
			},
		})
		return
	}

	// Check if we have a valid request
	if envelope.Body.Request == nil {
		c.XML(http.StatusBadRequest, SOAPEnvelope{
			Body: SOAPBody{
				Fault: &SOAPFault{
					FaultCode:   "Client",
					FaultString: "Missing currency conversion request",
				},
			},
		})
		return
	}

	convRequest := envelope.Body.Request

	// Validate currencies
	if !isValidCurrencyPair(convRequest.FromCurrency, convRequest.ToCurrency) {
		c.XML(http.StatusBadRequest, SOAPEnvelope{
			Body: SOAPBody{
				Fault: &SOAPFault{
					FaultCode:   "Client",
					FaultString: "Invalid currency pair",
					Detail:      fmt.Sprintf("Conversion from %s to %s is not supported", convRequest.FromCurrency, convRequest.ToCurrency),
				},
			},
		})
		return
	}

	// Perform the conversion
	rate, convertedAmount := convertCurrency(convRequest.Amount, convRequest.FromCurrency, convRequest.ToCurrency)

	// Create and send the response
	c.Header("Content-Type", "text/xml")
	c.XML(http.StatusOK, SOAPEnvelope{
		Body: SOAPBody{
			Response: &ConvertCurrencyResponse{
				ConvertedAmount: convertedAmount,
				FromCurrency:    convRequest.FromCurrency,
				ToCurrency:      convRequest.ToCurrency,
				Rate:            rate,
			},
		},
	})
}

// isValidCurrencyPair checks if the currency pair is supported
func isValidCurrencyPair(from, to Currency) bool {
	return (from == UAH && to == USD) || (from == USD && to == UAH)
}

// convertCurrency performs the currency conversion
func convertCurrency(amount float64, from, to Currency) (rate, convertedAmount float64) {
	switch {
	case from == UAH && to == USD:
		rate = UAHtoUSDRate
		convertedAmount = amount * UAHtoUSDRate
	case from == USD && to == UAH:
		rate = USDtoUAHRate
		convertedAmount = amount * USDtoUAHRate
	}
	return rate, convertedAmount
}
