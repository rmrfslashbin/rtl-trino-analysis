package useragent

import (
	"strings"

	"github.com/mssola/user_agent"
)

/*
(*user_agent.UserAgent)(0xc0001eaf00)({
 ua: (string) (len=139) "Mozilla/5.0 (iPhone; CPU iPhone OS 14_4_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Mobile/15E148 Safari/604.1",
 mozilla: (string) (len=3) "5.0",
 platform: (string) (len=6) "iPhone",
 os: (string) (len=34) "CPU iPhone OS 14_4_2 like Mac OS X",
 localization: (string) "",
 browser: (user_agent.Browser) {
  Engine: (string) (len=11) "AppleWebKit",
  EngineVersion: (string) (len=8) "605.1.15",
  Name: (string) (len=6) "Safari",
  Version: (string) (len=6) "14.0.3"
 },
 bot: (bool) false,
 mobile: (bool) true,
 undecided: (bool) false
})
*/

// Record is the data returned from the useragent package.
type Record struct {
	Raw                  string `json:"raw"`
	BrowserEngine        string `json:"browser_engine"`
	BrowserEngineVersion string `json:"browser_engine_version"`
	BrowserName          string `json:"browser_name"`
	BrowserVersion       string `json:"browser_version"`
	Mozilla              string `json:"mozilla"`
	Platform             string `json:"platform"`
	OS                   string `json:"os"`
	Localization         string `json:"localization"`
	Bot                  bool   `json:"bot"`
	Mobile               bool   `json:"mobile"`
}

// ParseFile parses the given file.
func Parse(line string) (*Record, error) {
	client := user_agent.New(line)

	engineName, engineVersion := client.Engine()
	browserName, browserVersion := client.Browser()

	return &Record{
		Raw:                  strings.TrimSpace(line),
		Mozilla:              strings.TrimSpace(client.Mozilla()),
		Platform:             strings.TrimSpace(client.Platform()),
		OS:                   strings.TrimSpace(client.OS()),
		Localization:         strings.TrimSpace(client.Localization()),
		Bot:                  client.Bot(),
		Mobile:               client.Mobile(),
		BrowserEngine:        strings.TrimSpace(engineName),
		BrowserEngineVersion: strings.TrimSpace(engineVersion),
		BrowserName:          strings.TrimSpace(browserName),
		BrowserVersion:       strings.TrimSpace(browserVersion),
	}, nil
}
