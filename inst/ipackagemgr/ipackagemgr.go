package ipackagemgr

//IPackageMgr is the interface for implementing functions for different package managers
type IPackageMgr interface {
	IsInstalled(packageName string) error
	GetInstalledVersion(packageName string, parseVersion bool) (string, error)
}
