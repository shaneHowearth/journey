// Package structure -
package structure

import (
	"sync"
)

// Blog - settings that are used for template execution
type Blog struct {
	sync.RWMutex
	URL             []byte
	Title           []byte
	Description     []byte
	Logo            []byte
	Cover           []byte
	AssetPath       []byte
	PostCount       int64
	PostsPerPage    int64
	ActiveTheme     string
	NavigationItems []Navigation
}
