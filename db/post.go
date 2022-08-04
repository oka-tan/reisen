package db

import (
	"encoding/base64"
	"fmt"
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

func (p *Post) FormatName() string {
	if p.Name != nil {
		return *p.Name
	} else {
		return "Anonymous"
	}
}

func (p *Post) FormatTime() string {
	return p.TimePosted.Format("Mon 2 Jan 2006 15:04:05")
}

func (p *Post) SubjectIsNil() bool {
	return p.Subject == nil
}

func (p *Post) DerefSubject() string {
	return *(p.Subject)
}

func (p *Post) CommentIsNil() bool {
	return p.Comment == nil
}

func (p *Post) DerefComment() string {
	return *(p.Comment)
}

func (p *Post) DerefThumbnailWidth() int16 {
	return *(p.ThumbnailWidth)
}

func (p *Post) DerefThumbnailHeight() int16 {
	return *(p.ThumbnailHeight)
}

func (p *Post) DerefMediaWidth() int16 {
	return *(p.MediaWidth)
}

func (p *Post) DerefMediaHeight() int16 {
	return *(p.MediaHeight)
}

func (p *Post) DerefMediaFileName() string {
	if p.MediaFileName == nil {
		return ""
	}

	return *(p.MediaFileName)
}

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

func (p *Post) MediaTimestampIsNil() bool {
	return p.MediaTimestamp == nil
}

func (p *Post) DerefMediaTimestamp() int64 {
	return *(p.MediaTimestamp)
}

func (p *Post) DerefMediaExtension() string {
	return *(p.MediaExtension)
}

func (p *Post) DerefThumbnailInternalHash() string {
	return base64.URLEncoding.EncodeToString(*(p.ThumbnailInternalHash))
}

func (p *Post) DerefMediaInternalHash() string {
	return base64.URLEncoding.EncodeToString(*(p.MediaInternalHash))
}

func (p *Post) MediaAvailable() bool {
	return p.MediaInternalHash != nil
}

func (p *Post) ThumbnailAvailable() bool {
	return p.ThumbnailInternalHash != nil
}

func (p *Post) OekakiAvailable() bool {
	return p.OekakiInternalHash != nil
}

func (p *Post) DerefOekakiInternalHash() string {
	return base64.URLEncoding.EncodeToString(*(p.OekakiInternalHash))
}
