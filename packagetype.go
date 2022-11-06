package msi

import "github.com/google/uuid"

type PackageType int

const (
	// An installer package, which installs a new application.
	PackageTypeInstaller PackageType = iota
	// A patch package, which provides an update to an application.
	PackageTypePatch
	// A transform, which is a collection of changes applied to an installation.
	PackageTypeTransform
)

func PackageTypeFromCLSID(clsid uuid.UUID) PackageType {

	switch clsid {
	case uuid.MustParse(INSTALLER_PACKAGE_CLSID):
		return PackageTypeInstaller
	case uuid.MustParse(PATCH_PACKAGE_CLSID):
		return PackageTypePatch
	case uuid.MustParse(TRANSFORM_PACKAGE_CLSID):
		return PackageTypeTransform
	default:
		return -1
	}
}

func (p PackageType) CLSID() uuid.UUID {
	switch p {
	case PackageTypeInstaller:
		return uuid.MustParse(INSTALLER_PACKAGE_CLSID)
	case PackageTypePatch:
		return uuid.MustParse(PATCH_PACKAGE_CLSID)
	case PackageTypeTransform:
		return uuid.MustParse(TRANSFORM_PACKAGE_CLSID)
	default:
		return uuid.Nil
	}
}

func (p PackageType) String() string {
	switch p {
	case PackageTypeInstaller:
		return "Installation Database"
	case PackageTypePatch:
		return "Patch"
	case PackageTypeTransform:
		return "Transform"
	default:
		return "Unknown"
	}
}
