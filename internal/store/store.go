package store

import (
	"tyto/internal/model"
	"sync/atomic"
	"maps"
)

type StoreInterface interface {
	GetCategories() []model.CategoryItem
	GetCategoryTree(string) []model.CategoryTreeParent
	GetContent(string) string
	GetAllContent() map[string]string
	Update(*model.StoreData)
}

type Store struct {
	ptr atomic.Pointer[model.StoreData]
}

func NewStore() StoreInterface {
	s := &Store{}
	s.ptr.Store(&model.StoreData{
		Contents:     make(map[string]string),
		Categories:   make([]model.CategoryItem, 0),
		CategoryTree: make(map[string][]model.CategoryTreeParent),
	})
	return s
}

func (s *Store) GetCategories() []model.CategoryItem {
	ptr := s.ptr.Load()
	return ptr.Categories
}

func (s *Store) GetCategoryTree(categoryID string) []model.CategoryTreeParent {
	ptr := s.ptr.Load()
	categoryTree, ok := ptr.CategoryTree[categoryID]
	if !ok {
		return []model.CategoryTreeParent{}
	}
	return categoryTree
}

func (s *Store) GetContent(contentID string) string {
	ptr := s.ptr.Load()
	content, ok := ptr.Contents[contentID]
	if !ok {
		return ""
	}
	return content
}

func (s *Store) GetAllContent() map[string]string {
	ptr := s.ptr.Load()
	return maps.Clone(ptr.Contents)
}

func (s *Store) Update(newData *model.StoreData) {
	s.ptr.Store(newData)
}