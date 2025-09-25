// Filename: cgi_query.c

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// URL decode function (simple)
void url_decode(char *src, char *dest) {
    char *pstr = src, *pbuf = dest;
    while (*pstr) {
        if (*pstr == '%') {
            if (*(pstr + 1) && *(pstr + 2)) {
                char hex[3] = { *(pstr + 1), *(pstr + 2), 0 };
                *pbuf++ = (char) strtol(hex, NULL, 16);
                pstr += 3;
            } else {
                pstr++;
            }
        } else if (*pstr == '+') {
            *pbuf++ = ' ';
            pstr++;
        } else {
            *pbuf++ = *pstr++;
        }
    }
    *pbuf = '\0';
}

int main(void) {
    // Get the QUERY_STRING environment variable (from URL)
    char *query = getenv("QUERY_STRING");

    // Output HTTP header
    printf("Content-Type: text/html\n\n");

    printf("<html><head><title>CGI Query Params</title></head><body>\n");
    printf("<h1>Received Query Parameters</h1>\n");

    if (query == NULL || strlen(query) == 0) {
        printf("<p>No query parameters received.</p>\n");
    } else {
        printf("<ul>\n");

        // Make a writable copy of query
        char query_copy[1024];
        strncpy(query_copy, query, sizeof(query_copy)-1);
        query_copy[sizeof(query_copy)-1] = '\0';

        char *param = strtok(query_copy, "&");
        while (param != NULL) {
            char key[512], value[512];
            char *eq = strchr(param, '=');
            if (eq) {
                *eq = '\0';
                url_decode(param, key);
                url_decode(eq + 1, value);
                printf("<li><b>%s</b> = %s</li>\n", key, value);
            } else {
                url_decode(param, key);
                printf("<li><b>%s</b> (no value)</li>\n", key);
            }
            param = strtok(NULL, "&");
        }
        printf("</ul>\n");
    }

    printf("</body></html>\n");
    return 0;
}
