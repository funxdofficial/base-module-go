package funbuild

// Config holds CLI options for service generation.
type Config struct {
	PkgName     string
	Output      string
	ServiceType string
}

const (
	ServiceTypeREST     = "rest"
	ServiceTypeConsumer = "consumer"
)

func normalizeServiceType(serviceType string) (string, error) {
	if serviceType == "" {
		return ServiceTypeREST, nil
	}
	switch serviceType {
	case ServiceTypeREST:
		return ServiceTypeREST, nil
	case ServiceTypeConsumer, "cons":
		return ServiceTypeConsumer, nil
	default:
		return "", errInvalidServiceType(serviceType)
	}
}

type invalidServiceTypeError string

func (e invalidServiceTypeError) Error() string {
	return "invalid service type " + string(e) + " (use rest or consumer)"
}

func errInvalidServiceType(serviceType string) error {
	return invalidServiceTypeError(serviceType)
}
