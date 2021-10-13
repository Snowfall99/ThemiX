#include <stdio.h>
#include <stdlib.h>
#include <curl/curl.h>
#include <sys/time.h>

int main(int argc, char** argv) {
    int requests = atoi(argv[1]);
    CURL *curl0, *curl1, *curl2, *curl3, *curl4;
    CURLcode res;
    time_t start;
    time_t end;

    FILE *fp = NULL;
    fp = fopen("coordinator.txt", "w+");
    if (!fp) {
        printf("open file error");
        return 0;
    }
    
    curl_global_init(CURL_GLOBAL_ALL);

    curl0 = curl_easy_init();
    curl1 = curl_easy_init();
    curl2 = curl_easy_init();
    curl3 = curl_easy_init();
    curl4 = curl_easy_init();
    
    if (curl0 && curl1 && curl2 && curl3 && curl4) {
        curl_easy_setopt(curl0, CURLOPT_URL, "http://localhost:11200/client");
        curl_easy_setopt(curl0, CURLOPT_POSTFIELDS, "a");
        curl_easy_setopt(curl1, CURLOPT_URL, "http://localhost:11210/client");
        curl_easy_setopt(curl1, CURLOPT_POSTFIELDS, "a");
        curl_easy_setopt(curl2, CURLOPT_URL, "http://localhost:11220/client");
        curl_easy_setopt(curl2, CURLOPT_POSTFIELDS, "a");
        curl_easy_setopt(curl3, CURLOPT_URL, "http://localhost:11230/client");
        curl_easy_setopt(curl3, CURLOPT_POSTFIELDS, "a");
        curl_easy_setopt(curl4, CURLOPT_URL, "http://localhost:11240/client");
        curl_easy_setopt(curl4, CURLOPT_POSTFIELDS, "a");

        time(&start);
        for (int i = 0; i < requests; i++) {
            struct timeval tv;
            char start_str[100], end_str[100];

            gettimeofday(&tv, NULL);
            sprintf(start_str, "%d start: %ld", 5*i, (tv.tv_sec*1000 + tv.tv_usec/1000));
            fprintf(fp, "%s\n", start_str);
            res = curl_easy_perform(curl0);
            gettimeofday(&tv, NULL);
            sprintf(end_str, "%d end: %ld", 5*i, (tv.tv_sec*1000 + tv.tv_usec/1000));
            fprintf(fp, "%s\n", end_str);
            if (res != CURLE_OK) {
                fprintf(stderr, "curl_easy_perform failed: %s\n", curl_easy_strerror(res));
            }

            gettimeofday(&tv, NULL);
            sprintf(start_str, "%d start: %ld", 5*i+1, (tv.tv_sec*1000 + tv.tv_usec/1000));
            fprintf(fp, "%s\n", start_str);
            res = curl_easy_perform(curl1);
            gettimeofday(&tv, NULL);
            sprintf(end_str, "%d end: %ld", 5*i+1, (tv.tv_sec*1000 + tv.tv_usec/1000));
            fprintf(fp, "%s\n", end_str);
            if (res != CURLE_OK) {
                fprintf(stderr, "curl_easy_perform failed: %s\n", curl_easy_strerror(res));
            }

            gettimeofday(&tv, NULL);
            sprintf(start_str, "%d start: %ld", 5*i+2, (tv.tv_sec*1000 + tv.tv_usec/1000));
            fprintf(fp, "%s\n", start_str);
            res = curl_easy_perform(curl2);
            gettimeofday(&tv, NULL);
            sprintf(end_str, "%d end: %ld", 5*i+2, (tv.tv_sec*1000 + tv.tv_usec/1000));
            fprintf(fp, "%s\n", end_str);
            if (res != CURLE_OK) {
                fprintf(stderr, "curl_easy_perform failed: %s\n", curl_easy_strerror(res));
            }

            gettimeofday(&tv, NULL);
            sprintf(start_str, "%d start: %ld", 5*i+3, (tv.tv_sec*1000 + tv.tv_usec/1000));
            fprintf(fp, "%s\n", start_str);
            res = curl_easy_perform(curl3);
            gettimeofday(&tv, NULL);
            sprintf(end_str, "%d end: %ld", 5*i+3, (tv.tv_sec*1000 + tv.tv_usec/1000));
            fprintf(fp, "%s\n", end_str);
            if (res != CURLE_OK) {
                fprintf(stderr, "curl_easy_perform failed: %s\n", curl_easy_strerror(res));
            }

            gettimeofday(&tv, NULL);
            sprintf(start_str, "%d start: %ld", 5*i+4, (tv.tv_sec*1000 + tv.tv_usec/1000));
            fprintf(fp, "%s\n", start_str);
            res = curl_easy_perform(curl4);
            gettimeofday(&tv, NULL);
            sprintf(end_str, "%d end: %ld", 5*i+4, (tv.tv_sec*1000 + tv.tv_usec/1000));
            fprintf(fp, "%s\n", end_str);
            if (res != CURLE_OK) {
                fprintf(stderr, "curl_easy_perform failed: %s\n", curl_easy_strerror(res));
            }
        }
        time(&end);
        curl_easy_cleanup(curl0);
        curl_easy_cleanup(curl1);
        curl_easy_cleanup(curl2);
        curl_easy_cleanup(curl3);
        curl_easy_cleanup(curl4);
    }
    curl_global_cleanup();
    printf("duration is %ld seconds\n", (end - start));
    fclose(fp);
    return 0;
}