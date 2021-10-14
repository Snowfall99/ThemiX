#include <stdio.h>
#include <stdlib.h>
#include <curl/curl.h>
#include <sys/time.h>

void curl_init(CURL *curl, char *addr, char* payload) {
    curl_easy_setopt(curl, CURLOPT_URL, addr);
    curl_easy_setopt(curl, CURLOPT_POSTFIELDS, payload);
}

int main(int argc, char** argv) {
    int requests = atoi(argv[1]);
    char* url0 = argv[2];
    char* url1 = argv[3];
    char* url2 = argv[4];
    char* url3 = argv[5];
    char* url4 = argv[6];
    printf("requets: %d\n", requests);
    CURLcode res;
    time_t start;
    time_t end;

    FILE *fp = NULL;
    fp = fopen("coordinator.txt", "w+");
    if (!fp) {
        printf("open file error\n");
        return 0;
    }
    
    time(&start);
    for (int i = 0; i < requests; i++) {
        curl_global_init(CURL_GLOBAL_ALL);
        CURLM* curlm = curl_multi_init();
        CURL *curl0, *curl1, *curl2, *curl3, *curl4;

        curl0 = curl_easy_init();
        if (curl0) {
            curl_init(curl0, url0, "a");
            curl_multi_add_handle(curlm, curl0);
        }
        curl1 = curl_easy_init();
        if (curl1) {
            curl_init(curl1, url1, "a");
            curl_multi_add_handle(curlm, curl1);
        }
        curl2 = curl_easy_init();
        if (curl2) {
            curl_init(curl2, url2, "a");
            curl_multi_add_handle(curlm, curl2);
        }
        curl3 = curl_easy_init();
        if (curl3) {
            curl_init(curl3, url3, "a");
            curl_multi_add_handle(curlm, curl3);
        }
        curl4 = curl_easy_init();
        if (curl4) {
            curl_init(curl4, url4, "a");
            curl_multi_add_handle(curlm, curl4);
        }

        struct timeval tv;
        char start_str[100], end_str[100];
        gettimeofday(&tv, NULL);
        sprintf(start_str, "%d start: %ld", i, (tv.tv_sec*1000 + tv.tv_usec/1000));
        fprintf(fp, "%s\n", start_str);
        int running_handlers = 0;
        do {
            curl_multi_wait(curlm, NULL, 0, 2000, NULL);
            curl_multi_perform(curlm, &running_handlers);
        } while (running_handlers > 0);
        gettimeofday(&tv, NULL);
        sprintf(end_str, "%d end: %ld", i, (tv.tv_sec*1000 + tv.tv_usec/1000));
        fprintf(fp, "%s\n", end_str);

        curl_easy_cleanup(curl0);
        curl_easy_cleanup(curl1);
        curl_easy_cleanup(curl2);
        curl_easy_cleanup(curl3);
        curl_easy_cleanup(curl4);

        curl_multi_cleanup(curlm);
    }
    time(&end);

    printf("duration is %ld seconds\n", (end - start));

    // if (curl0 && curl1 && curl2 && curl3 && curl4) {
    //     time(&start);
    //     for (int i = 0; i < requests; i++) {
    //         struct timeval tv;
    //         char start_str[100], end_str[100];

    //         gettimeofday(&tv, NULL);
    //         sprintf(start_str, "%03d start: %ld", 5*i, (tv.tv_sec*1000 + tv.tv_usec/1000));
    //         fprintf(fp, "%s\n", start_str);
    //         res = curl_easy_perform(curl0);
    //         gettimeofday(&tv, NULL);
    //         sprintf(end_str, "%03d end: %ld", 5*i, (tv.tv_sec*1000 + tv.tv_usec/1000));
    //         fprintf(fp, "%s\n", end_str);
    //         if (res != CURLE_OK) {
    //             fprintf(stderr, "curl_easy_perform failed: %s\n", curl_easy_strerror(res));
    //         }

    //         gettimeofday(&tv, NULL);
    //         sprintf(start_str, "%03d start: %ld", 5*i+1, (tv.tv_sec*1000 + tv.tv_usec/1000));
    //         fprintf(fp, "%s\n", start_str);
    //         res = curl_easy_perform(curl1);
    //         gettimeofday(&tv, NULL);
    //         sprintf(end_str, "%03d end: %ld", 5*i+1, (tv.tv_sec*1000 + tv.tv_usec/1000));
    //         fprintf(fp, "%s\n", end_str);
    //         if (res != CURLE_OK) {
    //             fprintf(stderr, "curl_easy_perform failed: %s\n", curl_easy_strerror(res));
    //         }

    //         gettimeofday(&tv, NULL);
    //         sprintf(start_str, "%03d start: %ld", 5*i+2, (tv.tv_sec*1000 + tv.tv_usec/1000));
    //         fprintf(fp, "%s\n", start_str);
    //         res = curl_easy_perform(curl2);
    //         gettimeofday(&tv, NULL);
    //         sprintf(end_str, "%03d end: %ld", 5*i+2, (tv.tv_sec*1000 + tv.tv_usec/1000));
    //         fprintf(fp, "%s\n", end_str);
    //         if (res != CURLE_OK) {
    //             fprintf(stderr, "curl_easy_perform failed: %s\n", curl_easy_strerror(res));
    //         }

    //         gettimeofday(&tv, NULL);
    //         sprintf(start_str, "%03d start: %ld", 5*i+3, (tv.tv_sec*1000 + tv.tv_usec/1000));
    //         fprintf(fp, "%s\n", start_str);
    //         res = curl_easy_perform(curl3);
    //         gettimeofday(&tv, NULL);
    //         sprintf(end_str, "%03d end: %ld", 5*i+3, (tv.tv_sec*1000 + tv.tv_usec/1000));
    //         fprintf(fp, "%s\n", end_str);
    //         if (res != CURLE_OK) {
    //             fprintf(stderr, "curl_easy_perform failed: %s\n", curl_easy_strerror(res));
    //         }

    //         gettimeofday(&tv, NULL);
    //         sprintf(start_str, "%03d start: %ld", 5*i+4, (tv.tv_sec*1000 + tv.tv_usec/1000));
    //         fprintf(fp, "%s\n", start_str);
    //         res = curl_easy_perform(curl4);
    //         gettimeofday(&tv, NULL);
    //         sprintf(end_str, "%03d end: %ld", 5*i+4, (tv.tv_sec*1000 + tv.tv_usec/1000));
    //         fprintf(fp, "%s\n", end_str);
    //         if (res != CURLE_OK) {
    //             fprintf(stderr, "curl_easy_perform failed: %s\n", curl_easy_strerror(res));
    //         }
    //     }
    //     time(&end);
    //     curl_easy_cleanup(curl0);
    //     curl_easy_cleanup(curl1);
    //     curl_easy_cleanup(curl2);
    //     curl_easy_cleanup(curl3);
    //     curl_easy_cleanup(curl4);
    // }
    // curl_global_cleanup();
    // printf("duration is %ld seconds\n", (end - start));
    fclose(fp);
    return 0;
}