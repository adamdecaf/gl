/*
 * Accounts API
 *
 * Moov Accounts is an HTTP service which represents both a general ledger and chart of accounts for customers. The service is designed to abstract over various core systems and provide a uniform API for developers.
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

// Phone struct for Phone
type Phone struct {
	// phone number
	Number string `json:"number,omitempty"`
	// phone number has been validated to connect with customer
	Valid bool   `json:"valid,omitempty"`
	Type  string `json:"type,omitempty"`
}
