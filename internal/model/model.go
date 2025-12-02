package model

// CategoryItem 代表一个分类项
type CategoryItem struct {
	ID   string   `json:"id"`
	Name string   `json:"name"`
}

// CategoryTreeItem 代表分类树中的一个文件项
type CategoryTreeItem struct {
	ID   string   `json:"id"`
	File string   `json:"file"`
}

// CategoryTreeParent 代表分类树中的一个父项
type CategoryTreeParent struct {
	Title   string             `json:"title"`
	Files []CategoryTreeItem   `json:"files"`
}

type StoreData struct {
	Contents map[string]string
	Categories []CategoryItem
	CategoryTree map[string][]CategoryTreeParent
}

