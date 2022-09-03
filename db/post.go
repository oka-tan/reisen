//Package db provides entities and methods
//for database access
package db

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/uptrace/bun"
)

//Post is a post in the db
type Post struct {
	bun.BaseModel `bun:"table:post,alias:post"`

	Board                 string     `bun:"board,pk"`
	PostNumber            int64      `bun:"post_number,pk"`
	ThreadNumber          int64      `bun:"thread_number"`
	Op                    bool       `bun:"op"`
	Deleted               bool       `bun:"deleted"`
	Hidden                bool       `bun:"hidden"`
	TimePosted            time.Time  `bun:"time_posted"`
	LastModified          time.Time  `bun:"last_modified"`
	CreatedAt             time.Time  `bun:"created_at"`
	Name                  *string    `bun:"name"`
	Tripcode              *string    `bun:"tripcode"`
	Capcode               *string    `bun:"capcode"`
	PosterID              *string    `bun:"poster_id"`
	Country               *string    `bun:"country"`
	Flag                  *string    `bun:"flag"`
	Email                 *string    `bun:"email"`
	Subject               *string    `bun:"subject"`
	Comment               *string    `bun:"comment"`
	HasMedia              bool       `bun:"has_media"`
	MediaDeleted          *bool      `bun:"media_deleted"`
	TimeMediaDeleted      *time.Time `bun:"time_media_deleted"`
	MediaTimestamp        *int64     `bun:"media_timestamp"`
	Media4chanHash        *[]byte    `bun:"media_4chan_hash"`
	MediaInternalHash     *[]byte    `bun:"media_internal_hash"`
	ThumbnailInternalHash *[]byte    `bun:"thumbnail_internal_hash"`
	MediaExtension        *string    `bun:"media_extension"`
	MediaFileName         *string    `bun:"media_file_name"`
	MediaSize             *int       `bun:"media_size"`
	MediaHeight           *int16     `bun:"media_height"`
	MediaWidth            *int16     `bun:"media_width"`
	ThumbnailHeight       *int16     `bun:"thumbnail_height"`
	ThumbnailWidth        *int16     `bun:"thumbnail_width"`
	Spoiler               *bool      `bun:"spoiler"`
	CustomSpoiler         *int16     `bun:"custom_spoiler"`
	Sticky                *bool      `bun:"sticky"`
	Closed                *bool      `bun:"closed"`
	Posters               *int16     `bun:"posters"`
	Replies               *int16     `bun:"replies"`
	Since4Pass            *int16     `bun:"since4pass"`
	OekakiInternalHash    *[]byte    `bun:"oekaki_internal_hash"`
}

//We need all of these goofy methods
//because of how mustache works

//FormatName returns the user name
//when a name is available on the db
//and 'Anonymous' otherwise
func (p *Post) FormatName() string {
	if p.Name != nil {
		return *p.Name
	}
	return "Anonymous"
}

//FormatTime formats the post time
//as a string
func (p *Post) FormatTime() string {
	return p.TimePosted.UTC().Format("Mon 2 Jan 2006 15:04:05")
}

//SubjectIsNil returns a boolean indicating
//whether the Subject field is nil or not
func (p *Post) SubjectIsNil() bool {
	return p.Subject == nil
}

//DerefSubject derefs the subject field
//to return a string
func (p *Post) DerefSubject() string {
	return *(p.Subject)
}

//CommentIsNil returns a boolean indicating
//whether or not the comment field is nil
func (p *Post) CommentIsNil() bool {
	return p.Comment == nil
}

//DerefComment derefs the post comment field
//to return a string
func (p *Post) DerefComment() string {
	return *(p.Comment)
}

//DerefThumbnailWidth derefs the thumbnail width
//field
func (p *Post) DerefThumbnailWidth() int16 {
	return *(p.ThumbnailWidth)
}

//DerefThumbnailHeight derefs the thumbnail height
//field
func (p *Post) DerefThumbnailHeight() int16 {
	return *(p.ThumbnailHeight)
}

//DerefMediaWidth derefs the media width
//field
func (p *Post) DerefMediaWidth() int16 {
	return *(p.MediaWidth)
}

//DerefMediaHeight derefs the media height
//field
func (p *Post) DerefMediaHeight() int16 {
	return *(p.MediaHeight)
}

//DerefMediaFileName derefs the media file name
func (p *Post) DerefMediaFileName() string {
	if p.MediaFileName == nil {
		return ""
	}

	return *(p.MediaFileName)
}

//DerefMediaFileNameShort derefs the media file
//name field and truncates it for diplay purposes
func (p *Post) DerefMediaFileNameShort() string {
	if p.MediaFileName == nil {
		return ""
	}

	s := *(p.MediaFileName)

	if len(s) < 20 {
		return s
	}

	runes := []rune(s)
	runesLen := len(runes)

	if runesLen < 17 {
		return s
	}

	return fmt.Sprintf("%s...", string(runes[:17]))
}

//MediaTimestampIsNil returns a boolean
//indicating whether or not the media timestamp
//field is nil
func (p *Post) MediaTimestampIsNil() bool {
	return p.MediaTimestamp == nil
}

//DerefMediaTimestamp derefs the media timestamp
//field
func (p *Post) DerefMediaTimestamp() int64 {
	return *(p.MediaTimestamp)
}

//DerefMediaExtension derefs the media extension
//field
func (p *Post) DerefMediaExtension() string {
	return *(p.MediaExtension)
}

//DerefThumbnailInternalHash derefs the thumbnail
//internal hash field
func (p *Post) DerefThumbnailInternalHash() string {
	return base64.URLEncoding.EncodeToString(*(p.ThumbnailInternalHash))
}

//DerefMediaInternalHash derefs the media internal
//hash field
func (p *Post) DerefMediaInternalHash() string {
	return base64.URLEncoding.EncodeToString(*(p.MediaInternalHash))
}

//DerefMedia4chanHash derefs the media 4chan hash
//field
func (p *Post) DerefMedia4chanHash() string {
	return base64.URLEncoding.EncodeToString(*(p.Media4chanHash))
}

//MediaAvailable returns a boolean indicating whether
//a post's media has been downloaded and is available
//for download by users
func (p *Post) MediaAvailable() bool {
	return p.MediaInternalHash != nil
}

//ThumbnailAvailable returns a boolean indicating whether
//a post's thumbnail has been downloaded and is available
//for visualization by users
func (p *Post) ThumbnailAvailable() bool {
	return p.ThumbnailInternalHash != nil
}

//OekakiAvailable returns a boolean indicating whether
//a post has a tegaki replay, it has been downloaded and
//it is available for visualization by users
func (p *Post) OekakiAvailable() bool {
	return p.OekakiInternalHash != nil
}

//DerefOekakiInternalHash derefs the oekaki internal
//hash field
func (p *Post) DerefOekakiInternalHash() string {
	return base64.URLEncoding.EncodeToString(*(p.OekakiInternalHash))
}

//HasCountry returns a boolean indicating
//whether or not a post specifies a country
func (p *Post) HasCountry() bool {
	return p.Country != nil
}

//DerefCountry derefs the country field
//and lowercases it
func (p *Post) DerefCountry() string {
	return strings.ToLower(*(p.Country))
}

//HasFlag returns a boolean indicating
//whether or not a post specifies a flag
func (p *Post) HasFlag() bool {
	return p.Flag != nil
}

//DerefFlag derefs the flag field
//and lowercases it
func (p *Post) DerefFlag() string {
	return strings.ToLower(*(p.Flag))
}
