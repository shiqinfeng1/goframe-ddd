package migration

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func InitDB() {
	ctx := gctx.New()
	db := g.DB()
	if err := createTable(ctx, db); err != nil {
		g.Log().Fatal(ctx, err)
	}
	// l := glog.New()
	// l.SetConfig(glog.Config{})
	// db.SetLogger(l)
}

func createTable(ctx context.Context, db gdb.DB) error {
	// 定义建表 SQL
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS "users" (
        "user_id" TEXT NOT NULL,
        "username" TEXT NOT NULL,
        "email" TEXT NOT NULL,
        "mobile_phone" TEXT NOT NULL,
        "password_hash" TEXT NOT NULL,
        "is_locked" INTEGER NOT NULL DEFAULT 0,
        "locked_until" TEXT,
        "failed_attempts" INTEGER NOT NULL DEFAULT 0,
		PRIMARY KEY ("user_id"),
        UNIQUE ("username"),
        UNIQUE ("email"),
        UNIQUE ("mobile_phone")
    )STRICT;
	CREATE TABLE IF NOT EXISTS "user_tokens" (
		"user_id" TEXT NOT NULL,
		"refresh_id" TEXT NOT NULL,
		PRIMARY KEY ("user_id")
	)STRICT;
	CREATE INDEX IF NOT EXISTS "idx_token_refresh_id" ON "user_tokens" ("refresh_id");
    `
	// 执行建表语句
	_, err := db.Exec(ctx, createTableSQL)
	g.Log().Debug(ctx, "migrate table 'users','user_tokens' success")
	return err
}
