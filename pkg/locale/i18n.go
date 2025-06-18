package locale

import (
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/net/ghttp"
)

// en	英文
// en-US	美式英文
// en-GB	英式英文
// zh	中文
// zh-CN	简体中文（中国）
// zh-HK	中文（中国香港）
// zh-TW	繁体中文（中国台湾）
// es	西班牙语
// fr	法语
// ja	日语
// ko	韩语
// de	德语
// ru	俄语
var Lang = map[string]*gi18n.Manager{
	"en": func() *gi18n.Manager {
		i18n := gi18n.New()
		i18n.SetPath("config/i18n")
		i18n.SetLanguage("en")
		return i18n
	}(),
	"zh": func() *gi18n.Manager {
		i18n := gi18n.New()
		i18n.SetPath("config/i18n")
		i18n.SetLanguage("zh-CN")
		return i18n
	}(),
}

// I18n 国际化中间件，支持多源获取语言参数
func Locale(r *ghttp.Request) {
	var (
		lang        string
		defaultLang = "en" // 默认语言
	)

	// 1. 从查询参数获取（优先级最高）
	lang = r.GetQuery("lang").String()
	if lang != "" {
		r.SetCtxVar("lang", lang)
		r.Middleware.Next()
		return
	}

	// 2. 从请求头获取（标准头或自定义头）
	lang = r.Header.Get("X-Language")
	if lang == "" {
		lang = r.Header.Get("Accept-Language")
	}
	// 处理标准 Accept-Language 格式（如 "zh-CN,zh;q=0.9"）
	if len(lang) >= 2 {
		lang = lang[:2] // 提取主语言代码（如 "zh-CN" → "zh"）
		r.SetCtxVar("lang", lang)
		r.Middleware.Next()
		return
	}

	// 3. 从Cookie获取
	lang = r.Cookie.Get("lang").String()
	if lang != "" {
		r.SetCtxVar("lang", lang)
		r.Middleware.Next()
		return
	}

	// 4. 从请求体获取（仅针对POST/PUT请求）
	var req struct {
		Lang string `json:"lang" form:"lang"`
	}
	if err := r.Parse(&req); err == nil && req.Lang != "" {
		r.SetCtxVar("lang", req.Lang)
		r.Middleware.Next()
		return
	}

	// 5. 使用默认语言
	r.SetCtxVar("lang", defaultLang)
	r.Middleware.Next()
}
