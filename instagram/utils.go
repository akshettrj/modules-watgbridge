package instagram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"golang.org/x/exp/slices"
)

var client *http.Client

func GetMediaType(body []byte) int {
	var mt struct {
		Items []*struct {
			MediaType int `json:"media_type"`
		} `json:"items"`
	}

	err := json.Unmarshal(body, &mt)
	if err != nil || mt.Items == nil || len(mt.Items) == 0 {
		return MediaTypeInvalid
	}

	return mt.Items[0].MediaType
}

func IsInstagramLink(link string) bool {
	parsedURL, err := url.Parse(link)
	if err != nil {
		return false
	}

	hostname := strings.ToLower(parsedURL.Hostname())
	return slices.Contains(InstagramHostnames, hostname)
}

func IsSupportedLink(link string) bool {
	parsedURL, err := url.Parse(link)
	if err != nil {
		return false
	}

	hostname := strings.ToLower(parsedURL.Hostname())

	var (
		isInstagramLink = slices.Contains(InstagramHostnames, hostname)
		isPostLink      = strings.HasPrefix(parsedURL.Path, "/p/")
		isReelLink      = strings.HasPrefix(parsedURL.Path, "/reel/")
		isTVLink        = strings.HasPrefix(parsedURL.Path, "/tv/")
		isStoriesLink   = InstagramStoriesRegexp.MatchString(parsedURL.Path)
	)

	return isInstagramLink && (isPostLink || isReelLink || isTVLink || isStoriesLink)
}

func IsStoriesLink(link string) bool {
	parsedURL, err := url.Parse(link)
	if err != nil {
		return false
	}

	match := InstagramStoriesRegexp.MatchString(parsedURL.Path)
	return match
}

func AddQueries(req *http.Request) {
	q := req.URL.Query()
	q.Set("__a", "1")
	q.Set("__d", "1")
	req.URL.RawQuery = q.Encode()
}

func AddHeaders(req *http.Request) {
	headers := instaConfig.Headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}
}

func AddCookies(req *http.Request) {
	cookies := instaConfig.Cookies
	for k, v := range cookies {
		req.AddCookie(&http.Cookie{Name: k, Value: v})
	}
}

func SaveCookies(res *http.Response) {
	for _, cookie := range res.Cookies() {
		if _, present := instaConfig.Cookies[cookie.Name]; present {
			instaConfig.Cookies[cookie.Name] = cookie.Value
		}
	}
	instaConfig.SaveConfig()
}

func DownloadFile(req *http.Request) ([]byte, error) {
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error making request: %s", err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Received status '%s'", res.Status)
	}

	return io.ReadAll(res.Body)
}

func init() {
	jar, _ := cookiejar.New(&cookiejar.Options{})
	client = &http.Client{Jar: jar}
}
