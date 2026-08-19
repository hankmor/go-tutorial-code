package main

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Author 是文章作者模型，一个作者可以拥有多篇文章。
type Author struct {
	ID     uint `gorm:"primaryKey"`
	Name   string
	Posts  []Post
	Orders []Order
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

// Order 是作者订单模型，用于演示聚合、连接和原生 SQL 查询。
type Order struct {
	ID       uint    `gorm:"primaryKey"`
	AuthorID uint    `gorm:"index"`
	Amount   float64 `gorm:"not null"`
	Status   string  `gorm:"size:20;index"`
}

// AuthorOrderSummary 是作者订单聚合查询的结果模型。
type AuthorOrderSummary struct {
	AuthorID   uint
	OrderCount int64
	Total      float64
}

// NewDB 创建内存 SQLite 数据库并完成本章模型迁移。
func NewDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	if err := db.AutoMigrate(&Author{}, &Post{}, &Tag{}, &Order{}); err != nil {
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
		if err := tx.Model(&posts[1]).Association("Tags").Append(&tags[0], &tags[1]); err != nil {
			return err
		}
		orders := []Order{
			{AuthorID: author.ID, Amount: 30, Status: "paid"},
			{AuthorID: author.ID, Amount: 80, Status: "paid"},
			{AuthorID: author.ID, Amount: 120, Status: "pending"},
		}
		return tx.Create(&orders).Error
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

// FindOrdersAboveAverage 查询金额高于全部订单平均值的订单。
func FindOrdersAboveAverage(db *gorm.DB) ([]Order, error) {
	average := db.Model(&Order{}).Select("AVG(amount)")
	var orders []Order
	err := db.Where("amount > (?)", average).Order("amount ASC").Find(&orders).Error
	return orders, err
}

// SummarizeAuthorOrders 按作者统计订单数量和订单总额。
func SummarizeAuthorOrders(db *gorm.DB, minimumTotal float64) ([]AuthorOrderSummary, error) {
	var summaries []AuthorOrderSummary
	err := db.Model(&Order{}).
		Select("author_id, COUNT(*) AS order_count, SUM(amount) AS total").
		Group("author_id").
		Having("SUM(amount) >= ?", minimumTotal).
		Order("total DESC").
		Scan(&summaries).Error
	return summaries, err
}

// FindAuthorsWithOrders 查询至少有一笔订单的作者，避免重复作者行。
func FindAuthorsWithOrders(db *gorm.DB) ([]Author, error) {
	var authors []Author
	err := db.Model(&Author{}).
		Joins("JOIN orders ON orders.author_id = authors.id").
		Group("authors.id").
		Order("authors.id ASC").
		Find(&authors).Error
	return authors, err
}

// FindAuthorsWithoutOrders 使用左连接查询没有订单的作者。
func FindAuthorsWithoutOrders(db *gorm.DB) ([]Author, error) {
	var authors []Author
	err := db.Model(&Author{}).
		Joins("LEFT JOIN orders ON orders.author_id = authors.id").
		Where("orders.id IS NULL").
		Order("authors.id ASC").
		Find(&authors).Error
	return authors, err
}

// RawAuthorOrderTotals 使用参数化原生 SQL 查询作者订单汇总。
func RawAuthorOrderTotals(db *gorm.DB, minimumTotal float64) ([]AuthorOrderSummary, error) {
	var summaries []AuthorOrderSummary
	err := db.Raw(`
		SELECT author_id, COUNT(*) AS order_count, SUM(amount) AS total
		FROM orders
		GROUP BY author_id
		HAVING SUM(amount) >= ?
		ORDER BY total DESC`, minimumTotal).Scan(&summaries).Error
	return summaries, err
}

// UpdateOrderStatus 使用 Exec 更新订单状态，并区分没有匹配记录的情况。
func UpdateOrderStatus(db *gorm.DB, orderID uint, status string) error {
	result := db.Exec("UPDATE orders SET status = ? WHERE id = ?", status, orderID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
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
	summaries, err := SummarizeAuthorOrders(db, 100)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("author=%s posts=%d tags-in-first-post=%d summaries=%d\n", loaded.Name, len(loaded.Posts), len(loaded.Posts[0].Tags), len(summaries))
}
