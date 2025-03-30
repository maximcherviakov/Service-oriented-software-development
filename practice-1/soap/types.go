package soap

import "encoding/xml"

const (
	USD Currency = "USD"
	UAH Currency = "UAH"
)

// Currency represents a currency type
type Currency string

// SOAPEnvelope represents the SOAP envelope
type SOAPEnvelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Body    SOAPBody
}

// SOAPBody represents the SOAP body
type SOAPBody struct {
	XMLName  xml.Name                 `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
	Fault    *SOAPFault               `xml:",omitempty"`
	Request  *ConvertCurrencyRequest  `xml:",omitempty"`
	Response *ConvertCurrencyResponse `xml:",omitempty"`
}

// SOAPFault represents a SOAP fault
type SOAPFault struct {
	XMLName     xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Fault"`
	FaultCode   string   `xml:"faultcode"`
	FaultString string   `xml:"faultstring"`
	Detail      string   `xml:"detail,omitempty"`
}

// ConvertCurrencyRequest represents a currency conversion request
type ConvertCurrencyRequest struct {
	XMLName      xml.Name `xml:"ConvertCurrencyRequest"`
	Amount       float64  `xml:"amount"`
	FromCurrency Currency `xml:"fromCurrency"`
	ToCurrency   Currency `xml:"toCurrency"`
}

// ConvertCurrencyResponse represents a currency conversion response
type ConvertCurrencyResponse struct {
	XMLName         xml.Name `xml:"ConvertCurrencyResponse"`
	ConvertedAmount float64  `xml:"convertedAmount"`
	FromCurrency    Currency `xml:"fromCurrency"`
	ToCurrency      Currency `xml:"toCurrency"`
	Rate            float64  `xml:"rate"`
}
