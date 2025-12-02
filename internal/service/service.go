package service

import (
	"log"
	"tyto/internal/model"
	"tyto/internal/store"
	"tyto/internal/utils/file"
	"tyto/internal/utils/git"
	"tyto/internal/utils/markdown"

	"sync"
	"fmt"
)

type ServiceInterface interface {
	GetCategoriesService() []model.CategoryItem
	GetCategoryTreeService(string) []model.CategoryTreeParent
	GetContentService(string) string
	GetAllContentService() map[string]string
	UpdateStore(*model.StoreData)
	SyncData() error
}

type Service struct {
	store store.StoreInterface
	gitRepoURL string
	repositoryDir string
	converter *markdown.MarkdownConverter
	mu sync.RWMutex
}

func NewService(s store.StoreInterface,
	gitRepourl string, repositoryDir string) *Service {
	return &Service{
		store: s,
		gitRepoURL: gitRepourl,
		repositoryDir: repositoryDir,
		converter: markdown.NewMarkdownConverter(),
	}
}

func (s *Service) GetCategoriesService() []model.CategoryItem {
	return s.store.GetCategories()
}

func (s *Service) GetCategoryTreeService(categoryID string) []model.CategoryTreeParent {
	return s.store.GetCategoryTree(categoryID)
}

func (s *Service) GetContentService(contentID string) string {
	return s.store.GetContent(contentID)
}

func (s *Service) GetAllContentService() map[string]string {
	return s.store.GetAllContent()
}

func (s *Service) UpdateStore(newData *model.StoreData) {
	s.store.Update(newData)
}

func (s *Service) SyncData() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// git 拉取
	output, err := git.Sync(s.gitRepoURL, s.repositoryDir)
	if err != nil {
		return fmt.Errorf("git sync error: %w", err)
	}
	log.Printf("[Service] Git output: %s\n", output)
	// 解析文件
	dirs, err := file.ParseDir(s.repositoryDir)
	if err != nil {
		return fmt.Errorf("parse dir error: %w", err)
	}
	// 转换为存储数据
	storeData, err := file.DirInfoToStoreData(dirs,
				s.store.GetAllContent(),
				s.converter)
	if err != nil {
		return fmt.Errorf("dir info to store data error: %w", err)
	}
	// 更新存储
	s.UpdateStore(storeData)
	log.Println("[Service] Sync completed successfully.")
	return nil
}