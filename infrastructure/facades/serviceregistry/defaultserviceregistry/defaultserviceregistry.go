package defaultserviceregistry

type DefaultServiceRegistry struct {
}

func NewDefaultServiceRegistry() *DefaultServiceRegistry {
	r := DefaultServiceRegistry{}

	return &r
}

func (r *DefaultServiceRegistry) Register() error {
	return nil

}

func (r *DefaultServiceRegistry) Deregister() error {
	return nil

}
