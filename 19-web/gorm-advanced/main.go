package main

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Author 是文章作者模型，一个作者可以拥有多篇文章。
type Author struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Posts []Post
}

// Post 是文章模型，通过 post_tags 与标签建立多对多关系。
type Post struct {
	ID        uint `gorm:"primaryKey"`
	Title     string
	AuthorID  uint
	Tags      []*Tag         `gorm:"many2many:post_tags;"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Tag 是文章标签模型，一个标签可以关联多篇文章。
type Tag struct {
	ID    uint    `gorm:"primaryKey"`
	Name  string  `gorm:"uniqueIndex;not null"`
	Posts []*Post `gorm:"many2many:post_tags;"`
}

// NewDB 创建内存 SQLite 数据库并完成本章模型迁移。
func NewDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	if err := db.AutoMigrate(&Author{}, &Post{}, &Tag{}); err != nil {
		return nil, fmt.Errorf("migrate association models: %w", err)
	}
	return db, nil
}

// SeedGraph 在一个事务中创建作者、文章、标签及其关联关系。
func SeedGraph(db *gorm.DB) (Author, error) {
	var author Author
	err := db.Transaction(func(tx *gorm.DB) error {
		author = Author{Name: "Ada"}
		if err := tx.Create(&author).Error; err != nil {
			return err
		}

		posts := []Post{
			{Title: "Go 事务", AuthorID: author.ID},
			{Title: "GORM 查询", AuthorID: author.ID},
		}
		if err := tx.Create(&posts).Error; err != nil {
			return err
		}

		tags := []Tag{{Name: "go"}, {Name: "database"}}
		if err := tx.Create(&tags).Error; err != nil {
			return err
		}
		if err := tx.Model(&posts[0]).Association("Tags").Append(&tags[0]); err != nil {
			return err
		}
		return tx.Model(&posts[1]).Association("Tags").Append(&tags[0], &tags[1])
	})
	return author, err
}

// LoadAuthorGraph 使用嵌套预加载读取作者、文章和文章标签。
func LoadAuthorGraph(db *gorm.DB, authorID uint) (Author, error) {
	var author Author
	err := db.Preload("Posts.Tags").First(&author, authorID).Error
	return author, err
}

// ReplacePostTags 用新的标签集合替换文章的多对多关联。
func ReplacePostTags(db *gorm.DB, postID uint, tags ...*Tag) error {
	var post Post
	if err := db.First(&post, postID).Error; err != nil {
		return err
	}
	return db.Model(&post).Association("Tags").Replace(tags)
}

func main() {
	db, err := NewDB()
	if err != nil {
		log.Fatal(err)
	}
	author, err := SeedGraph(db)
	if err != nil {
		log.Fatal(err)
	}
	loaded, err := LoadAuthorGraph(db, author.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("author=%s posts=%d tags-in-first-post=%d\n", loaded.Name, len(loaded.Posts), len(loaded.Posts[0].Tags))
}
