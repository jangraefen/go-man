// +build windows

package manager

import (
	"fmt"
	"os"
	"os/exec"

	"golang.org/x/sys/windows"

	"github.com/NoizeMe/go-man/internal/fileutil"
)

func link(sourceDirectory, targetDirectory string) error {
	if fileutil.PathExists(targetDirectory) {
		return fmt.Errorf("%s: file or directory already exists", sourceDirectory)
	}

	if admin, adminErr := hasAdmin(); admin && adminErr == nil {
		return os.Symlink(sourceDirectory, targetDirectory)
	}

	return exec.Command("cmd", "/c", "mklink", "/J", targetDirectory, sourceDirectory).Run()
}

func unlink(directory string) error {
	if !fileutil.PathExists(directory) {
		return fmt.Errorf("%s: no such file or directory", directory)
	}

	return os.RemoveAll(directory)
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

	defer func() {
		_ = windows.FreeSid(sid)
	}()

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
