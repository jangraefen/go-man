// +build windows

package selectors

import (
	"github.com/otiai10/copy"
	"golang.org/x/sys/windows"
	"os"
)

func symlink(sourceDirectory, targetDirectory string) error {
	if admin, adminErr := hasAdmin(); admin && adminErr == nil {
		return os.Symlink(sourceDirectory, targetDirectory)
	}

	if _, statErr := os.Stat(targetDirectory); statErr != nil && os.IsNotExist(statErr) {
		if removeErr := os.RemoveAll(targetDirectory); removeErr != nil {
			return statErr
		}
	}

	return copy.Copy(sourceDirectory, targetDirectory)
}

// For more details on this function, check https://coolaj86.com/articles/golang-and-windows-and-admins-oh-my/
func hasAdmin() (bool, error) {
	var sid *windows.SID

	// Although this looks scary, it is directly copied from the
	// official windows documentation. The Go API for this is a
	// direct wrap around the official C++ API.
	// See https://docs.microsoft.com/en-us/windows/desktop/api/securitybaseapi/nf-securitybaseapi-checktokenmembership
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid)
	if err != nil {
		return false, err
	}
	defer windows.FreeSid(sid)

	// This appears to cast a null pointer so I'm not sure why this
	// works, but this guy says it does and it Works for Me™:
	// https://github.com/golang/go/issues/28804#issuecomment-438838144
	token := windows.Token(0)

	member, err := token.IsMember(sid)
	if err != nil {
		return false, err
	}

	return token.IsElevated() || member, nil
}
