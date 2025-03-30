# Go SOAP Service

A Go web service that provides SOAP endpoints for currency conversion, built with the Gin framework.

## Features

- SOAP endpoint for currency conversion (UAH ↔ USD)
- REST health check endpoint
- Swagger documentation
- Environment-based configuration

## Project Structure

```
practice-1/
├── main.go           # Main application entry point
├── go.mod           # Go module definition
├── .env             # Environment variables (optional)
├── .gitignore       # Git ignore rules
├── soap/
│   ├── types.go     # SOAP request/response types
│   └── handler.go   # SOAP request handlers
└── README.md        # Project documentation
```

## Prerequisites

- Go 1.21 or higher
- Git

## Setup

1. Clone the repository:

```bash
git clone <repository-url>
cd practice-1
```

2. Install dependencies:

```bash
go mod download
```

3. Create a `.env` file in the root directory (optional):

```bash
PORT=8080
GIN_MODE=debug
```

## Running the Application

Development mode:

```bash
go run main.go
```

Production mode:

```bash
GIN_MODE=release go run main.go
```

## Available Endpoints

### REST Endpoints

- `GET /api/v1/health` - Health check endpoint
- `GET /swagger/*` - Swagger documentation

### SOAP Endpoints

Currency Conversion Service

- Endpoint: `POST /soap/convert-currency`
- Content-Type: `text/xml`

Example Request:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">
   <Body>
      <ConvertCurrencyRequest xmlns="http://practice-1/soap">
         <amount>1000</amount>
         <fromCurrency>UAH</fromCurrency>
         <toCurrency>USD</toCurrency>
      </ConvertCurrencyRequest>
   </Body>
</Envelope>
```

Example Response:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">
   <Body>
      <ConvertCurrencyResponse xmlns="http://practice-1/soap">
         <convertedAmount>25</convertedAmount>
         <fromCurrency>UAH</fromCurrency>
         <toCurrency>USD</toCurrency>
         <rate>0.025</rate>
      </ConvertCurrencyResponse>
   </Body>
</Envelope>
```

## Exchange Rates

The service uses fixed exchange rates for demonstration:

- 1 UAH = 0.025 USD
- 1 USD = 40 UAH

## Error Handling

The service returns SOAP faults in the following cases:

- Invalid Content-Type (not text/xml)
- Malformed SOAP request
- Invalid currency pair
- Missing or invalid request parameters

## License

MIT License
