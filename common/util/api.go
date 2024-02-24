package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func CallAPI(method, urlStr string, queryParams map[string]string) ([]byte, error) {
	if len(queryParams) > 0 {
		query := url.Values{}
		for key, value := range queryParams {
			query.Add(key, value)
		}
		urlStr += "?" + query.Encode()
	}

	req, err := http.NewRequest(strings.ToUpper(method), urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	return body, nil
}

func RateLimiter(rateLimiter *limiter.Limiter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ipClient := ctx.ClientIP()
		limiterCtx, err := rateLimiter.Get(ctx, ipClient)
		if err != nil {
			return
		}
		if limiterCtx.Reached {
			ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
