package pagination

// Pagination 페이지네이션 정보
type Pagination struct {
	Limit  int
	Offset int
	Page   int
}

// NewPagination 페이지네이션 생성
func NewPagination(page, limit int) *Pagination {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20 // 기본값
	}
	if limit > 100 {
		limit = 100 // 최대값
	}
	
	return &Pagination{
		Limit:  limit,
		Offset: (page - 1) * limit,
		Page:   page,
	}
}

// TotalPages 전체 페이지 수 계산
func TotalPages(total, limit int) int {
	if limit <= 0 {
		return 0
	}
	pages := total / limit
	if total%limit > 0 {
		pages++
	}
	return pages
}

