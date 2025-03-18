package comments

import "errors"

var (
	ErrCommentNotFound = errors.New("comment not found")
	ErrPostNotFound    = errors.New("post not found")
	ErrUnauthorized    = errors.New("unauthorized to modify this comment")
	ErrEmptyContent    = errors.New("comment content cannot be empty")
)
