package routes

import (
	"anonymousoverflow/config"
	"anonymousoverflow/src/types"
	"anonymousoverflow/src/utils"
	"fmt"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetHome(c *gin.Context) {
	theme := utils.GetThemeFromEnv()
	c.HTML(200, "home.html", gin.H{
		"version": config.Version,
		"theme":   theme,
	})
}

type urlConversionRequest struct {
	URL string `form:"url" binding:"required"`
}

var coreRegex = regexp.MustCompile(`(?:https?://)?(?:www\.)?([^/]+)(/(?:questions|q|a)/.+)`)

// Will return `nil` if `rawUrl` is invalid.
func translateUrl(rawUrl string) string {
	coreMatches := coreRegex.FindStringSubmatch(rawUrl)
	if coreMatches == nil {
		return ""
	}

	domain := coreMatches[1] // Extract the domain
	rest := coreMatches[2]   // Extract the rest of the URL path

	exchange := ""
	// Check if the domain matches "stackoverflow.com" first, as a default
	if domain == "stackoverflow.com" {
		// No exchange parameter needed for stackoverflow.com
	} else {
		for _, exchangeDomain := range types.ExchangeDomains {
			if sub, found := strings.CutSuffix(domain, "."+exchangeDomain+".com"); found {
				if sub == "" {
					return ""
				} else if strings.Contains(sub, ".") {
					// Anything containing dots is interpreted as a full domain, so we use the correct full domain.
					exchange = domain
				} else {
					exchange = sub
				}
			} else {
				exchange = domain
			}
		}
	}

	// Ensure we properly format the return string to avoid double slashes
	if exchange == "" {
		return rest
	} else {
		return fmt.Sprintf("/exchange/%s%s", exchange, rest)
	}
}

func PostHome(c *gin.Context) {
	body := urlConversionRequest{}

	if err := c.ShouldBind(&body); err != nil {
		c.HTML(400, "home.html", gin.H{
			"errorMessage": "Invalid request body",
		})
		return
	}

	translated := translateUrl(body.URL)

	if translated == "" {
		theme := utils.GetThemeFromEnv()
		c.HTML(400, "home.html", gin.H{
			"errorMessage": "Invalid stack overflow/exchange URL",
			"theme":        theme,
		})
		return
	}

	c.Redirect(302, translated)
}
