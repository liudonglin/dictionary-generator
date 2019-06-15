package core

// Pager 分页基类
type Pager struct {
	// Index 从零开始
	Index   int    `json:"index"`
	Size    int    `json:"size"`
	OrderBy string `json:"order_by"`
}
