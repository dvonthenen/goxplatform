package common

//IInit is the interface for implementing system init types
type IInit interface {
	Start(serviceName string) error
	StartEx(serviceName string, successRegex string) error
	Restart(serviceName string) error
	RestartEx(serviceName string, successStopRegex string, successStartRegex string) error
	Status(serviceName string) (bool, error)
	StatusEx(serviceName string, successRegex string) (bool, error)
	Stop(serviceName string) error
	StopEx(serviceName string, successRegex string) error

	Enable(serviceName string) error
	Disable(serviceName string) error

	AddDependentService(serviceName string, depName string) error
	RemoveDependentService(serviceName string, depName string) error
}
