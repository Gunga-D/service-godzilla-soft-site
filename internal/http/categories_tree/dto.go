package categories_tree

type CategoryDTO struct {
	ID       int64         `json:"id"`
	Name     string        `json:"name"`
	Children []CategoryDTO `json:"children"`
}
