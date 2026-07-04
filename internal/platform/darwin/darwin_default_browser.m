#import <Foundation/Foundation.h>
#import <ApplicationServices/ApplicationServices.h>

NSString* safeString(const char* input) {
    NSString *result = nil;
    if (input != nil) {
        result = [NSString stringWithUTF8String:input];
    }
    return result;
}

void SetDefaultBrowser(char *url_scheme, char *handler) {
    LSSetDefaultHandlerForURLScheme(
        (__bridge CFStringRef) safeString(url_scheme),
        (__bridge CFStringRef) safeString(handler)
    );
}