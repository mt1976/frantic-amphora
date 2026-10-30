package database

import (
	"fmt"

	"github.com/mt1976/frantic-core/logHandler"
)

// PageInfo describes the position of a page within a paged result set, e.g. for
// rendering "Page 1 of 10" in a GUI.
type PageInfo struct {
	Page         int // Current page number (1-based)
	PageSize     int // Maximum number of records per page
	TotalRecords int // Total number of records matching the query
	TotalPages   int // Total number of pages (always at least 1)
}

// HasPrevious reports whether there is a page before the current page.
func (p PageInfo) HasPrevious() bool {
	return p.Page > 1
}

// HasNext reports whether there is a page after the current page.
func (p PageInfo) HasNext() bool {
	return p.Page < p.TotalPages
}

// String returns a display friendly representation, e.g. "Page 1 of 10".
func (p PageInfo) String() string {
	return fmt.Sprintf("Page %d of %d", p.Page, p.TotalPages)
}

// TotalPages returns the number of pages required to hold totalRecords records at
// pageSize records per page. An empty result set is reported as a single page so that
// "Page 1 of 1" can still be displayed.
func TotalPages(totalRecords, pageSize int) int {
	if pageSize < 1 || totalRecords < 1 {
		return 1
	}
	return (totalRecords + pageSize - 1) / pageSize
}

// NewPageInfo builds a PageInfo for the given page, page size and total record count.
func NewPageInfo(page, pageSize, totalRecords int) PageInfo {
	return PageInfo{
		Page:         page,
		PageSize:     pageSize,
		TotalRecords: totalRecords,
		TotalPages:   TotalPages(totalRecords, pageSize),
	}
}

// validatePaging normalises the requested page and validates the page size.
// Pages below 1 are treated as page 1. It returns the normalised page and the
// number of records to skip.
func validatePaging(operation string, page, pageSize int) (int, int, error) {
	if pageSize < 1 {
		logHandler.Error.Printf("[%v] Invalid page size %d - must be at least 1", operation, pageSize)
		return 0, 0, fmt.Errorf("invalid page size %d for %v: must be at least 1", pageSize, operation)
	}
	if page < 1 {
		logHandler.Warning.Printf("[%v] Invalid page %d - defaulting to page 1", operation, page)
		page = 1
	}
	return page, (page - 1) * pageSize, nil
}
