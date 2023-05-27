// Description of things common to the domain layer
package domain

// Provide access to the business logic of the service
type Service struct {
}

// Create a new Service instance
func New() *Service {
	return &Service{}
}
