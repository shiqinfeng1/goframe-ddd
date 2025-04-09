package migration

import (
	"context"

	"github.com/shiqinfeng1/goframe-ddd/internal/adapters/ent"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	_ "github.com/mattn/go-sqlite3"
)

func NewEntClient(ctx context.Context) *ent.Client {
	if !gfile.IsDir("./data") {
		gfile.Mkdir("./data")
	}
	// 连接文件模式的 SQLite 数据库
	client, err := ent.Open("sqlite3", "file:./data/gomg.db?_fk=1")
	if err != nil {
		g.Log().Fatalf(ctx, "failed opening connection to sqlite: %v", err)
	}
	// 自动迁移数据库，创建表结构
	if err := client.Schema.Create(ctx); err != nil {
		g.Log().Fatalf(ctx, "Failed to create schema: %v", err)
	}
	return client
}
