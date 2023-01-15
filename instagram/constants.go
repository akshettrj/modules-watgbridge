package instagram

import "regexp"

var InstagramHostnames = []string{
	"www.instagram.com",
	"instagram.com",
}

var (
	InstagramUsernameRegexp = regexp.MustCompile(`@([a-zA-Z0-9._]+)`)
	InstagramStoriesRegexp  = regexp.MustCompile(`^/stories/[a-zA-Z0-9._]+/\d+`)
)
