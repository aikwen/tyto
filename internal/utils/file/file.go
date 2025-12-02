package file

import (
	"os"
	"path/filepath"
	"strings"
	"fmt"
	"slices"
	"encoding/json"

	"tyto/internal/utils/hash"
	"tyto/internal/utils/markdown"
	"tyto/internal/model"
)

const MetaFileName = "meta.json"

type MdFileInfo struct {
	NameWithoutExt string    // 文件名不带后缀
	Abspath string // 文件绝对路径
	ID   string    // 文件内容 xxhash 值
}

type DirInfo struct {
	Name string
	Title string
	Category string
	Files []*MdFileInfo
}

type MetaInfo struct {
	Title string    `json:"title"`
	Category string `json:"category"`
	Order []string  `json:"order"`
}

func ParseMDFileInfo(abspath string) (*MdFileInfo, error) {
	contentHash, err := hash.HashFile(abspath)
	if err != nil {
		return nil, err
	}

	fileName := filepath.Base(abspath)
	nameWithoutExt := TrimSuffix(fileName, filepath.Ext(fileName))
	return &MdFileInfo{
		NameWithoutExt: nameWithoutExt,
		Abspath: abspath,
		ID: fmt.Sprintf("%x", contentHash),
	}, nil
}

func ParseMetaFromJson(abspath string) (*MetaInfo, error) {
	data, err := os.ReadFile(abspath)
	if err != nil {
		return nil, err
	}

	var meta MetaInfo
	err = json.Unmarshal(data, &meta)
	if err != nil {
		return nil, err
	}

	return &meta, nil
}

// TrimSuffix 去掉文件后缀
func TrimSuffix(name, ext string) string {
    // 简单粗暴且高效：直接切掉后缀长度
	return name[:len(name)-len(ext)]
}

func GetDirInfo(path string) (*DirInfo, error) {
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	// 转成绝对路径
	absDirPath, err := filepath.Abs(path)
	if err != nil {
		absDirPath = path
	}

	f := map[string]*MdFileInfo{}
	var meta *MetaInfo
	// 遍历每一个文件
	for _, entry := range dirEntries {
		if !entry.IsDir() {
			// 获取文件名和文件绝对路径
			fullPath := filepath.Join(absDirPath, entry.Name())
			ext := filepath.Ext(fullPath)
			ext = strings.ToLower(ext)
			switch ext {
			 	case ".md":
					i, err := ParseMDFileInfo(fullPath)
					if err != nil {
						continue
					}
					f[i.NameWithoutExt] = i
				case ".json":
					if strings.ToLower(entry.Name()) != MetaFileName{
						continue
					}

					meta, err = ParseMetaFromJson(fullPath)
					if err != nil {
						continue
					}
				default:
					// 忽略其他文件
			}
		}
	}

	category := ""
	name := filepath.Base(absDirPath)
	title := name
	files := make([]*MdFileInfo, 0, len(f))
	if meta != nil {
		// 使用 meta 信息
		if meta.Category != ""{
			category = meta.Category
		}
		if meta.Title != ""{
			title = meta.Title
		}

		// 遍历order
		for _, name := range meta.Order {
			if info, ok := f[name]; ok {
				files = append(files, info)
				delete(f, name)
			}
		}
	}

	// 添加剩余的文件
	tempfiles := make([]*MdFileInfo, 0, len(f))
	for _, info := range f {
		tempfiles = append(tempfiles, info)
	}
	// 按名称排序
	slices.SortFunc(tempfiles, func(a, b *MdFileInfo) int {
		return strings.Compare(a.NameWithoutExt, b.NameWithoutExt)
	})
	files = append(files, tempfiles...)

	return &DirInfo{
		Name: name,
		Title: title,
		Category: category,
		Files: files,
	}, nil
}


func ParseDir(path string) ([]*DirInfo, error) {
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	// 转成绝对路径
	absDirPath, err := filepath.Abs(path)
	if err != nil {
		absDirPath = path
	}

	var dirs []*DirInfo
	for _, entry := range dirEntries {
		if entry.IsDir() {
			name := entry.Name()
			if strings.HasPrefix(name, ".") {
				continue
			}
			fullPath := filepath.Join(absDirPath, entry.Name())
			dirInfo, err := GetDirInfo(fullPath)
			if err != nil {
				continue
			}

			if len(dirInfo.Files) == 0 {
				continue
			}
			
			dirs = append(dirs, dirInfo)
		}
	}
	return dirs, nil
}

func DirInfoToStoreData(dirs []*DirInfo, cache map[string]string, mdConverter *markdown.MarkdownConverter) (*model.StoreData, error) {
	categoryItemMap := map[string][]model.CategoryTreeParent{}
	contents := map[string]string{}
	categoryItems := map[string]string{}
	// 遍历每个文件夹
	for _, dir := range dirs {
		// 遍历文件夹下的每个文件并将每个文件夹转化为 categoryTreeParent
		categoryTreeParent := model.CategoryTreeParent{
			Title: dir.Title,
			Files: make([]model.CategoryTreeItem, 0),
		}

		for _, file := range dir.Files {
			categoryTreeParent.Files = append(categoryTreeParent.Files, model.CategoryTreeItem{
				ID:   file.ID,
				File: file.NameWithoutExt,
			})
			// cache 中寻找值
			content := ""
			ok := false
			if cache != nil {
				content, ok = cache[file.ID]
			}
			if  ok {
				contents[file.ID] = content
			}else{
					data, err := os.ReadFile(file.Abspath)
					if err != nil {
						contents[file.ID] = ""
						fmt.Println("Read file error:", file.Abspath,err)
						continue
					}
					html, err := mdConverter.MdToHtml(data)
					if err != nil {
						contents[file.ID] = ""
						fmt.Println("Convert md to html error:", file.Abspath, err)
						continue
					}
					contents[file.ID] = html
				}

		}

		// 获取当前文件夹的分类
		categoryName := dir.Category
		if categoryName == "" {
			categoryName = "未分类" // 或者 "Uncategorized"
		}
		id := fmt.Sprintf("%x", hash.Hash64([]byte(categoryName)))
		categoryItems[id] = categoryName
		// 将 categoryTreeParent 添加到对应的分类中
		categoryItemMap[id] = append(categoryItemMap[id], categoryTreeParent)
	}
	// 构建 CategoryItem 列表
	categories := []model.CategoryItem{}
	for id, name := range categoryItems {
		categories = append(categories, model.CategoryItem{
			ID:   id,
			Name: name,
		})
	}
	// 对分类列表进行排序
	slices.SortFunc(categories, func(a, b model.CategoryItem) int {
		return strings.Compare(a.Name, b.Name)
	})
	// 对每个分类下的文件夹进行排序
	for id := range categoryItemMap {
		slices.SortFunc(categoryItemMap[id], func(a, b model.CategoryTreeParent) int {
			return strings.Compare(a.Title, b.Title)
		})
	}

	return &model.StoreData{
		Categories: categories,
		CategoryTree: categoryItemMap,
		Contents: contents,
	}, nil
}