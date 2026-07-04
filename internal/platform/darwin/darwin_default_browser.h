//
//  darwin_default_browser.h
//  test
//
//  Created by icaliskanoglu on 11/09/23.
//

#ifndef darwin_default_browser_h
#define darwin_default_browser_h

#import <Foundation/Foundation.h>
#import <ApplicationServices/ApplicationServices.h>

void SetDefaultBrowser(const char *url_scheme, const char *handler);

#endif