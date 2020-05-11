// Package winfsd interposes on the CreateFileW Windows API call that underlies
// syscall.CreateFile, adding the FILE_SHARE_DELETE flag to CreateFileW's
// dwShareMode parameter.  It does nothing on other OSes.
//
// To use the package, link it into your program:
//      import _ "github.com/nwt/winfsd"
package winfsd
