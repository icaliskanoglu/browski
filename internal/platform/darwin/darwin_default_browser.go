//go:build darwin
// +build darwin

package darwin

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework Cocoa -framework WebKit
#import <Foundation/Foundation.h>
#import "darwin_default_browser.h"

#include <stdlib.h>
*/
import "C"

func SetDefaultBrowser(bundleId string) error {
	c := NewCalloc()
	defer c.Free()
	bundleIdC := c.String(bundleId)
	http := c.String("http")
	https := c.String("https")
	C.SetDefaultBrowser(http, bundleIdC)
	C.SetDefaultBrowser(https, bundleIdC)
	return nil
}
