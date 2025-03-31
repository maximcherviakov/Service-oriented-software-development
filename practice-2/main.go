package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"practice-2/currency"
)

// CurrencyService implements the SOAP service
type CurrencyService struct{}

// ConvertCurrency implements the currency conversion functionality
func (s *CurrencyService) ConvertCurrency(request *currency.ConvertCurrencyRequest) (*currency.ConvertCurrencyResponse, error) {
	// Simple rates for demonstration purposes
	rates := map[string]map[string]float64{
		"USD": {"EUR": 0.93, "GBP": 0.79, "JPY": 152.0, "UAH": 41.5},
		"EUR": {"USD": 1.07, "GBP": 0.85, "JPY": 163.0, "UAH": 44.6},
		"GBP": {"USD": 1.26, "EUR": 1.18, "JPY": 192.0, "UAH": 52.5},
		"JPY": {"USD": 0.0066, "EUR": 0.0061, "GBP": 0.0052, "UAH": 0.27},
		"UAH": {"USD": 0.024, "EUR": 0.022, "GBP": 0.019, "JPY": 3.7},
	}

	from := request.FromCurrency
	to := request.ToCurrency
	amount := request.Amount

	rate, ok := rates[from][to]
	if !ok {
		return nil, fmt.Errorf("conversion rate not found for %s to %s", from, to)
	}

	convertedAmount := amount * rate

	return &currency.ConvertCurrencyResponse{
		ConvertedAmount: convertedAmount,
		FromCurrency:    from,
		ToCurrency:      to,
		Rate:            rate,
	}, nil
}

// ConvertCurrencyContext implements the context-aware version of the conversion functionality
func (s *CurrencyService) ConvertCurrencyContext(ctx context.Context, request *currency.ConvertCurrencyRequest) (*currency.ConvertCurrencyResponse, error) {
	return s.ConvertCurrency(request)
}

// SOAPEnvelope is the root element for SOAP requests/responses
type SOAPEnvelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Body    SOAPBody
}

// SOAPBody contains the actual SOAP message
type SOAPBody struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
	Content interface{}
	Fault   *SOAPFault `xml:",omitempty"`
}

// SOAPFault represents a SOAP error
type SOAPFault struct {
	XMLName     xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Fault"`
	FaultCode   string   `xml:"faultcode"`
	FaultString string   `xml:"faultstring"`
	Detail      string   `xml:"detail,omitempty"`
}

// SOAPHandler handles SOAP requests for currency conversion
func (s *CurrencyService) SOAPHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Basic request validation
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "text/xml; charset=utf-8")

		// 2. Read request body
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			sendSOAPFault(w, "Failed to read request body", err.Error())
			return
		}

		// 3. Parse request with flexible structure
		var requestData struct {
			XMLName xml.Name `xml:"Envelope"`
			Body    struct {
				XMLName xml.Name `xml:"Body"`
				Request struct {
					XMLName      xml.Name `xml:"ConvertCurrencyRequest"`
					Amount       float64  `xml:"amount"`
					FromCurrency string   `xml:"fromCurrency"`
					ToCurrency   string   `xml:"toCurrency"`
				}
			}
		}

		if err := xml.Unmarshal(body, &requestData); err != nil {
			sendSOAPFault(w, "Failed to parse request", err.Error())
			return
		}

		// 4. Create properly typed request
		request := &currency.ConvertCurrencyRequest{
			Amount:       requestData.Body.Request.Amount,
			FromCurrency: requestData.Body.Request.FromCurrency,
			ToCurrency:   requestData.Body.Request.ToCurrency,
		}

		// 5. Process the request
		response, err := s.ConvertCurrency(request)
		if err != nil {
			sendSOAPFault(w, "Failed to process request", err.Error())
			return
		}

		// 6. Send response
		sendSOAPResponse(w, response)
	}
}

// sendSOAPResponse sends a successful SOAP response
func sendSOAPResponse(w http.ResponseWriter, content interface{}) {
	responseEnvelope := SOAPEnvelope{
		Body: SOAPBody{
			Content: content,
		},
	}

	output, err := xml.MarshalIndent(responseEnvelope, "", "  ")
	if err != nil {
		sendSOAPFault(w, "Failed to encode response", err.Error())
		return
	}

	w.Write([]byte(xml.Header + string(output)))
}

// sendSOAPFault sends a SOAP fault response
func sendSOAPFault(w http.ResponseWriter, faultString, detail string) {
	fault := SOAPFault{
		FaultCode:   "Server",
		FaultString: faultString,
		Detail:      detail,
	}

	envelope := SOAPEnvelope{
		Body: SOAPBody{
			Fault: &fault,
		},
	}

	output, err := xml.MarshalIndent(envelope, "", "  ")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(xml.Header + string(output)))
}

// WSDLFileServer serves static WSDL files
func WSDLFileServer(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the file name from the URL
		fileName := filepath.Base(r.URL.Path)
		filePath := filepath.Join(dir, fileName)

		// Check if the file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}

		// Read the file
		data, err := os.ReadFile(filePath)
		if err != nil {
			http.Error(w, "Error reading WSDL file", http.StatusInternalServerError)
			return
		}

		// Set content type and write the file
		w.Header().Set("Content-Type", "text/xml; charset=utf-8")
		w.Write(data)
	}
}

func main() {
	// Create and register the currency service
	currencyService := &CurrencyService{}

	// Register the SOAP handler for the currency service
	http.HandleFunc("/soap/convert-currency", currencyService.SOAPHandler())

	// Serve WSDL files
	http.HandleFunc("/wsdl/", WSDLFileServer("wsdl"))

	// Start the HTTP server
	fmt.Println("Starting SOAP server at http://localhost:8080")
	fmt.Println("WSDL: http://localhost:8080/wsdl/currency.wsdl")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
