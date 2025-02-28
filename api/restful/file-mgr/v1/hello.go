// Copyright @2025-2028 <SieYuan> . All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package v1

import "github.com/gogf/gf/v2/util/gmeta"

type IsAutoUpdateReq struct {
	gmeta.Meta `path:"/version/isAutoupdate" tags:"版本管理" method:"post" summary:"配置是否自动更新"`
	IsAuto     bool `p:"is_auto" dc:"ture:自动更新"`
}
type IsAutoUpdateRes struct {
	gmeta.Meta `mime:"application/json"`
}
type GetConfigReq struct {
	gmeta.Meta `path:"/version/configures" tags:"版本管理" method:"post" summary:"查询是否自动更新"`
	Ver        int `p:"ver" v:"required|in:1,2" dc:"配置的版本号"`
}
type GetConfigRes struct {
	gmeta.Meta `mime:"application/json"`
	IsAuto     bool   `json:"is_auto"     dc:"ture:自动更新"`
	ExeVersion string `json:"exe_version" dc:"本exe的版本号"`
	VerCommit  string `json:"ver_commit"  dc:"版本描述"`
}
