package shared

import "fmt"

var Messages = struct {
	EnvironmentErrorNotFound        string
	EnvironmentErrorNeedOne         func(string) string
	ConnectSuccess                  func(string) string
	AuthorizeFailedMissingParameter func(string) string
	AuthorizeFailedBadParameter     string
	AuthorizeFailedWrongParameter   func(string) string
	AuthorizeFailedNoUser           string
	AuthorizeFailedWrongPassword    string
	AuthorizeFailedWrongToken       string
	AuthorizeFailedUserExist        string
	ArticleFailedBadContent         string
}{
	EnvironmentErrorNotFound: fmt.Sprintf(
		"[环境变量错误] 找不到 .env 文件。",
	),

	EnvironmentErrorNeedOne: func(payload string) string {
		return fmt.Sprintf(
			"[环境变量错误] 请在 .env 文件中设置名为 `%s` 的环境变量。",
			payload,
		)
	},

	ConnectSuccess: func(payload string) string {
		return fmt.Sprintf(
			"[连接成功] 成功连接至%s！",
			payload,
		)
	},

	AuthorizeFailedMissingParameter: func(payload string) string {
		return fmt.Sprintf(
			"[认证失败] 缺少%s参数。",
			payload,
		)
	},

	AuthorizeFailedBadParameter: "[认证失败] 验证信息格式不正确",

	AuthorizeFailedWrongParameter: func(payload string) string {
		return fmt.Sprintf(
			"[认证失败] 参数%s错误。",
			payload,
		)
	},

	AuthorizeFailedNoUser:        "[认证失败] 找不到用户。",
	AuthorizeFailedWrongPassword: "[认证失败] 密码不一致。",
	AuthorizeFailedWrongToken:    "[认证失败] 错误的令牌。",
	AuthorizeFailedUserExist:     "[认证失败] 用户已存在。",
	ArticleFailedBadContent:      "[创建帖子失败] 文章内容不正确",
}
