package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"sync/atomic"

	"gorm.io/gorm"
)

// JSONMap 是可以映射到 JSON 文本列的自定义 Go 类型。
type JSONMap map[string]any

// Scan 从数据库值解码 JSONMap。
func (m *JSONMap) Scan(value any) error {
	if value == nil {
		*m = nil
		return nil
	}
	var raw []byte
	switch typed := value.(type) {
	case []byte:
		raw = typed
	case string:
		raw = []byte(typed)
	default:
		return errors.New("unsupported JSON value type")
	}
	return json.Unmarshal(raw, m)
}

// Value 将 JSONMap 编码为数据库驱动可以写入的 JSON 字节。
func (m JSONMap) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

// MetadataRecord 演示带有自定义 JSON 字段的模型。
type MetadataRecord struct {
	ID       uint `gorm:"primaryKey"`
	Name     string
	Metadata JSONMap `gorm:"type:json"`
}

// QueryMetricsPlugin 统计经过 GORM 查询回调的次数。
type QueryMetricsPlugin struct {
	Queries atomic.Int64
}

// Name 返回插件名称。
func (p *QueryMetricsPlugin) Name() string {
	return "query_metrics"
}

// Initialize 注册查询完成后的统计回调。
func (p *QueryMetricsPlugin) Initialize(db *gorm.DB) error {
	return db.Callback().Query().After("gorm:query").Register("query_metrics:after_query", func(_ *gorm.DB) {
		p.Queries.Add(1)
	})
}

// InstallQueryMetrics 安装查询统计插件。
func InstallQueryMetrics(db *gorm.DB, plugin *QueryMetricsPlugin) error {
	return db.Use(plugin)
}
