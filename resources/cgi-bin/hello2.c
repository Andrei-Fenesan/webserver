// Filename: hello.c

#include <stdio.h>

int main(void) {
    // Every CGI program must output the Content-Type header first
    printf("Content-Type: text/html\n\n");

    // Now output the HTML content
    printf("<html>\n");
    printf("<head><title>Hello from CGI</title></head>\n");
    printf("<body>\n");
    printf("<h1>Hello from a C CGI script!</h1>\n");
    printf("</body>\n");
    printf("</html>\n");

    return 0;
}
