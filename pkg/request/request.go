package request

import "resty.dev/v3"

const (
	UserAgentDesktopChrome  = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36"
	UserAgentMicroMessenger = "Mozilla/5.0 (iPhone; CPU iPhone OS 15_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.31(0x18001f37) NetType/WIFI Language/zh_CN"
	UserAgentTelegramBot    = "TelegramBot (like TwitterBot)"
)

type Client struct {
	*resty.Client
}

func New() *Client {
	return &Client{
		Client: resty.New(),
	}
}

func WithDefaultHeaders() *Client {
	client := resty.New().
		SetHeader("priority", "u=0, i").
		SetHeader("sec-ch-ua", `"Google Chrome";v="111", "Not(A:Brand";v="8", "Chromium";v="111"`).
		SetHeader("sec-ch-ua-mobile", "?0").
		SetHeader("sec-ch-ua-platform", "Windows").
		SetHeader("sec-fetch-dest", "document").
		SetHeader("sec-fetch-mode", "navigate").
		SetHeader("sec-fetch-site", "same-origin").
		SetHeader("sec-fetch-user", "?1").
		SetHeader("upgrade-insecure-requests", "1").
		SetHeader("user-agent", UserAgentDesktopChrome)
	return &Client{
		Client: client,
	}
}
