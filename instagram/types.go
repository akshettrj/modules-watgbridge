package instagram

import (
	"fmt"
	"path"
)

const (
	MediaTypeInvalid = -1

	MediaTypeImage    = 1
	MediaTypeVideo    = 2
	MediaTypeCarousel = 8
)

type InstagramUserProfile struct {
	Graphql struct {
		User struct {
			Username        string `json:"username"`
			FullName        string `json:"full_name"`
			Biography       string `json:"biography"`
			ID              string `json:"id"`
			ProfilePicURLHD string `json:"profile_pic_url_hd"`
			ProfilePicURL   string `json:"profile_pic_url"`
			EdgeFollowedBy  struct {
				Count int64 `json:"count"`
			} `json:"edge_followed_by"`
			EdgeFollow struct {
				Count int64 `json:"count"`
			} `json:"edge_follow"`
			IsProfessionalAccount bool `json:"is_professional_account"`
			HasBlockedViewer      bool `json:"has_blocked_viewer"`
			FollowsViewer         bool `json:"follows_viewer"`
			FollowedByViewer      bool `json:"followed_by_viewer"`
			BlockedByViewer       bool `json:"blocked_by_viewer"`
			IsVerified            bool `json:"is_verified"`
			IsPrivate             bool `json:"is_private"`
			IsBusinessAccount     bool `json:"is_business_account"`
		} `json:"user"`
	} `json:"graphql"`
}

type InstagramUser struct {
	Username      string `json:"username,omitempty"`
	FullName      string `json:"full_name,omitempty"`
	ProfilePicURL string `json:"profile_pic_url,omitempty"`

	IsPrivate  bool `json:"is_private,omitempty"`
	IsVerified bool `json:"is_verified,omitempty"`

	FriendshipStatus struct {
		Following       bool `json:"following,omitempty"`
		OutgoingRequest bool `json:"outgoing_request,omitempty"`
		IsBestie        bool `json:"is_bestie,omitempty"`
		IsRestricted    bool `json:"is_restricted,omitempty"`
		IsFeedFavorite  bool `json:"is_feed_favorite,omitempty"`
	} `json:"friendship_status,omitempty"`
}

type InstagramCaption struct {
	Text string        `json:"text,omitempty"`
	User InstagramUser `json:"user,omitempty"`
}

type InstagramImage struct {
	Items []struct {
		Code          string           `json:"code,omitempty"`
		ID            string           `json:"id,omitempty"`
		Caption       InstagramCaption `json:"caption,omitempty"`
		User          InstagramUser    `json:"user,omitempty"`
		TopLikers     []string         `json:"top_likers,omitempty"`
		VideoVersions []VideoVersion   `json:"video_versions,omitempty"`
		ImageVersions struct {
			Candidates []ImageVersion `json:"candidates,omitempty"`
		} `json:"image_versions2,omitempty"`
		MediaType                  int   `json:"media_type,omitempty"`
		CommentCount               int64 `json:"comment_count,omitempty"`
		LikeCount                  int64 `json:"like_count,omitempty"`
		Height                     int32 `json:"original_height,omitempty"`
		Width                      int32 `json:"original_width,omitempty"`
		HaveLiked                  bool  `json:"has_liked,omitempty"`
		IsPhotoOfMe                bool  `json:"photo_of_you,omitempty"`
		IsLikeAndViewCountDisabled bool  `json:"like_and_view_counts_disabled,omitempty"`
		IsCaptionEdited            bool  `json:"caption_is_edited,omitempty"`
	} `json:"items,omitempty"`
}

type VideoItem struct {
	ImageVersions struct {
		AdditionalCandidates struct {
			IGTVFirstFrame ImageVersion `json:"igtv_first_frame,omitempty"`
			FirstFrame     ImageVersion `json:"first_frame,omitempty"`
		} `json:"additional_candidates,omitempty"`
		Candidates []ImageVersion `json:"candidates,omitempty"`
	} `json:"image_versions2,omitempty"`
	Code                       string           `json:"code,omitempty"`
	ID                         string           `json:"id,omitempty"`
	Caption                    InstagramCaption `json:"caption,omitempty"`
	User                       InstagramUser    `json:"user,omitempty"`
	VideoVersions              []VideoVersion   `json:"video_versions,omitempty"`
	TopLikers                  []string         `json:"top_likers,omitempty"`
	MediaType                  int              `json:"media_type,omitempty"`
	PK                         int64            `json:"pk"`
	CommentCount               int64            `json:"comment_count,omitempty"`
	LikeCount                  int64            `json:"like_count,omitempty"`
	ViewCount                  int64            `json:"view_count,omitempty"`
	PlayCount                  int64            `json:"play_count,omitempty"`
	VideoDuration              float64          `json:"video_duration,omitempty"`
	Height                     int32            `json:"original_height,omitempty"`
	Width                      int32            `json:"original_width,omitempty"`
	HaveLiked                  bool             `json:"has_liked,omitempty"`
	IsPhotoOfMe                bool             `json:"photo_of_you,omitempty"`
	IsLikeAndViewCountDisabled bool             `json:"like_and_view_counts_disabled,omitempty"`
	IsCaptionEdited            bool             `json:"caption_is_edited,omitempty"`
}

type InstagramReel struct {
	Items []VideoItem `json:"items,omitempty"`
}

type InstagramCarousel struct {
	Items []struct {
		Code                       string           `json:"code,omitempty"`
		ID                         string           `json:"id,omitempty"`
		Caption                    InstagramCaption `json:"caption,omitempty"`
		User                       InstagramUser    `json:"user,omitempty"`
		FacepileTopLikers          []InstagramUser  `json:"facepile_top_likers,omitempty"`
		CarouselMedia              []CarouselMedia  `json:"carousel_media,omitempty"`
		TopLikers                  []string         `json:"top_likers,omitempty"`
		MediaType                  int              `json:"media_type,omitempty"`
		CommentCount               int64            `json:"comment_count,omitempty"`
		LikeCount                  int64            `json:"like_count,omitempty"`
		CarouselMediaCount         int32            `json:"carousel_media_count,omitempty"`
		IsPhotoOfMe                bool             `json:"photo_of_you,omitempty"`
		HaveLiked                  bool             `json:"has_liked,omitempty"`
		IsLikeAndViewCountDisabled bool             `json:"like_and_view_counts_disabled,omitempty"`
		IsCaptionEdited            bool             `json:"caption_is_edited,omitempty"`
	} `json:"items,omitempty"`
}

type InstagramStoryPublic struct {
	User struct {
		ID            string `json:"id"`
		Username      string `json:"username"`
		ProfilePicURL string `json:"profile_pic_url"`
	} `json:"user"`
}

type InstagramStory struct {
	User               InstagramUser `json:"user"`
	MediaIDs           []int64       `json:"media_ids"`
	Items              []VideoItem   `json:"items"`
	ID                 int64         `json:"id"`
	CanReply           bool          `json:"can_reply"`
	CanGIFQuickReply   bool          `json:"can_gif_quick_reply"`
	CanReshare         bool          `json:"can_reshare"`
	CanReactWithAvatar bool          `json:"can_react_with_avatar"`
}

type ImageVersion struct {
	URL    string `json:"url,omitempty"`
	Width  int32  `json:"width,omitempty"`
	Height int32  `json:"height,omitempty"`
}

type VideoVersion struct {
	ID     string `json:"id,omitempty"`
	URL    string `json:"url,omitempty"`
	Type   int32  `json:"type,omitempty"`
	Width  int32  `json:"width,omitempty"`
	Height int32  `json:"height,omitempty"`
}

type CarouselMedia struct {
	ID            string `json:"id,omitempty"`
	ImageVersions struct {
		Candidates []ImageVersion `json:"candidates,omitempty"`
	} `json:"image_versions2,omitempty"`
	VideoVersions []VideoVersion `json:"video_versions,omitempty"`
	MediaType     int            `json:"media_type,omitempty"`
	VideoDuration float64        `json:"video_duration,omitempty"`
	Height        int32          `json:"original_height"`
	Width         int32          `json:"original_weight"`
}

func (ii InstagramImage) Caption() string {
	if len(ii.Items) == 0 {
		return ""
	}
	item := ii.Items[0]

	caption := item.Caption.Text

	if item.IsLikeAndViewCountDisabled {
		caption += fmt.Sprintln("\n\n*Counts disabled by OP!*")
	} else {
		caption += fmt.Sprintf("\n\n*👍 : %v*\n", item.LikeCount)
		caption += fmt.Sprintf("*💬 : %v*\n", item.CommentCount)
	}

	username := item.User.Username
	fullname := item.User.FullName
	if fullname != "" {
		caption += fmt.Sprintf("*👤 : %s [@%s]*", fullname, username)
	} else {
		caption += fmt.Sprintf("*👤 : @%s*", username)
	}

	return caption
}

func (ir InstagramReel) Caption() string {
	if len(ir.Items) == 0 {
		return ""
	}
	item := ir.Items[0]

	caption := item.Caption.Text

	if item.IsLikeAndViewCountDisabled {
		caption += fmt.Sprintln("\n\n*Counts disabled by OP!*")
	} else {
		caption += fmt.Sprintf("\n\n*👍 : %v*\n", item.LikeCount)
		caption += fmt.Sprintf("*💬 : %v*\n", item.CommentCount)
		caption += fmt.Sprintf("*👀 : %v*\n", item.ViewCount)
	}

	username := item.User.Username
	fullname := item.User.FullName
	if fullname != "" {
		caption += fmt.Sprintf("*👤 : %s [@%s]*", fullname, username)
	} else {
		caption += fmt.Sprintf("*👤 : @%s*", username)
	}

	return caption
}

func (ic InstagramCarousel) Caption() string {
	if len(ic.Items) == 0 {
		return ""
	}
	item := ic.Items[0]

	caption := item.Caption.Text

	if item.IsLikeAndViewCountDisabled {
		caption += fmt.Sprintln("\n\n*Counts disabled by OP!*")
	} else {
		caption += fmt.Sprintf("\n\n*👍 : %v*\n", item.LikeCount)
		caption += fmt.Sprintf("*💬 : %v*\n", item.CommentCount)
	}

	username := item.User.Username
	fullname := item.User.FullName
	if fullname != "" {
		caption += fmt.Sprintf("*👤 : %s [@%s]*", fullname, username)
	} else {
		caption += fmt.Sprintf("*👤 : @%s*", username)
	}

	return caption
}

func (is InstagramStory) Caption() string {
	if len(is.Items) == 0 {
		return ""
	}

	username := is.User.Username
	fullname := is.User.FullName
	if fullname != "" {
		return fmt.Sprintf("*👤 : %s [@%s]*", fullname, username)
	} else {
		return fmt.Sprintf("*👤 : @%s*", username)
	}
}

func (ii InstagramImage) DownloadPath() string {
	if len(ii.Items) == 0 {
		return ""
	}
	return path.Join("downloads", "instagram", ii.Items[0].Code)
}

func (ir InstagramReel) DownloadPath() string {
	if len(ir.Items) == 0 {
		return ""
	}
	return path.Join("downloads", "instagram", ir.Items[0].Code)
}

func (ic InstagramCarousel) DownloadPath() string {
	if len(ic.Items) == 0 {
		return ""
	}
	return path.Join("downloads", "instagram", ic.Items[0].Code)
}

func (is InstagramStory) DownloadPath() string {
	if len(is.Items) == 0 {
		return ""
	}
	return path.Join("downloads", "instagram", is.User.Username)
}

func (ii InstagramImage) DownloadLink() string {
	if len(ii.Items) == 0 {
		return ""
	}
	item := ii.Items[0]
	width, height := item.Width, item.Height
	for _, candidate := range item.ImageVersions.Candidates {
		if candidate.Height == height && candidate.Width == width {
			return candidate.URL
		}
	}
	var (
		currentMax int32  = 0
		link       string = ""
	)
	for _, candidate := range item.ImageVersions.Candidates {
		resolution := candidate.Height * candidate.Width
		if resolution > currentMax {
			currentMax = resolution
			link = candidate.URL
		}
	}
	return link
}

func (ir InstagramReel) DownloadLink() string {
	if len(ir.Items) == 0 {
		return ""
	}
	item := ir.Items[0]
	width, height := item.Width, item.Height
	for _, candidate := range item.VideoVersions {
		if candidate.Height == height && candidate.Width == width {
			return candidate.URL
		}
	}
	var (
		currentMax int32  = 0
		link       string = ""
	)
	for _, candidate := range item.VideoVersions {
		resolution := candidate.Height * candidate.Width
		if resolution > currentMax {
			currentMax = resolution
			link = candidate.URL
		}
	}
	return link
}

func (is InstagramStory) DownloadLink(mediaID int64) string {
	if len(is.Items) == 0 {
		return ""
	}
	for _, item := range is.Items {
		if item.PK != mediaID {
			continue
		}
		width, height := item.Width, item.Height
		for _, candidate := range item.VideoVersions {
			if candidate.Height == height && candidate.Width == width {
				return candidate.URL
			}
		}
		var (
			currentMax int32  = 0
			link       string = ""
		)
		for _, candidate := range item.VideoVersions {
			resolution := candidate.Height * candidate.Width
			if resolution > currentMax {
				currentMax = resolution
				link = candidate.URL
			}
		}
		return link
	}
	return ""
}

func (cm CarouselMedia) DownloadLink() string {
	if cm.MediaType == MediaTypeImage {

		width, height := cm.Width, cm.Height
		for _, candidate := range cm.ImageVersions.Candidates {
			if candidate.Height == height && candidate.Width == width {
				return candidate.URL
			}
		}
		var (
			currentMax int32  = 0
			link       string = ""
		)
		for _, candidate := range cm.ImageVersions.Candidates {
			resolution := candidate.Height * candidate.Width
			if resolution > currentMax {
				currentMax = resolution
				link = candidate.URL
			}
		}
		return link

	} else if cm.MediaType == MediaTypeVideo {

		width, height := cm.Width, cm.Height
		for _, candidate := range cm.VideoVersions {
			if candidate.Height == height && candidate.Width == width {
				return candidate.URL
			}
		}
		var (
			currentMax int32  = 0
			link       string = ""
		)
		for _, candidate := range cm.VideoVersions {
			resolution := candidate.Height * candidate.Width
			if resolution > currentMax {
				currentMax = resolution
				link = candidate.URL
			}
		}
		return link
	}

	return ""
}

func (iup InstagramUserProfile) Followers() int64 {
	return iup.Graphql.User.EdgeFollowedBy.Count
}

func (iup InstagramUserProfile) Following() int64 {
	return iup.Graphql.User.EdgeFollow.Count
}

func (iup InstagramUserProfile) ProfilePicURL() string {
	return iup.Graphql.User.ProfilePicURL
}

func (iup InstagramUserProfile) ProfilePicURLHD() string {
	return iup.Graphql.User.ProfilePicURLHD
}

func (iup InstagramUserProfile) Caption() string {
	var (
		fullName = iup.Graphql.User.FullName
		username = iup.Graphql.User.Username
		bio      = iup.Graphql.User.Biography
	)

	if fullName == "" {
		fullName = "-"
	}

	if bio == "" {
		bio = "-"
	}

	caption := fmt.Sprintf("*Name* : %s\n", fullName)
	caption += fmt.Sprintf("*Username* : @%s\n", username)
	caption += fmt.Sprintf("*Bio* : %s\n", bio)
	caption += fmt.Sprintf("*Followers* : %v\n", iup.Followers())
	caption += fmt.Sprintf("*Following* : %v\n\n", iup.Following())

	if iup.Graphql.User.BlockedByViewer {
		caption += "• *Blocked by me* ✅\n"
	}
	if iup.Graphql.User.HasBlockedViewer {
		caption += "• *Blocks me* 💀\n"
	}
	if iup.Graphql.User.IsPrivate {
		caption += "• *Private Acc.* ✅\n"
	}
	if iup.Graphql.User.IsVerified {
		caption += "• *Verfied* ✅\n"
	}
	if iup.Graphql.User.IsProfessionalAccount {
		caption += "• *Professional Acc.* ✅\n"
	}
	if iup.Graphql.User.IsBusinessAccount {
		caption += "• *Bussiness Acc.* ✅\n"
	}
	if iup.Graphql.User.FollowsViewer {
		caption += "• *Follows Me* 😍\n"
	}
	if iup.Graphql.User.FollowedByViewer {
		caption += "• *Followed by me* 🤔\n"
	}

	return caption
}
