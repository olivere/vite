package vite

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type metadataKeyType string

var metadataKey = metadataKeyType("metadata")

// MetadataFromContext returns the metadata from the context.
// Use [MetadataToContext] to set the metadata in the context.
func MetadataFromContext(ctx context.Context) *Metadata {
	if md, ok := ctx.Value(metadataKey).(*Metadata); ok {
		return md
	}
	return nil
}

// MetadataToContext sets the metadata in the context.
// It is the inverse of [MetadataFromContext].
func MetadataToContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, metadataKey, &md)
}

type TitleData struct {
	Template string
	Default  string
	Absolute string
}

type Author struct {
	Name string
	URL  string
}

type FormatDetection struct {
	Email     bool
	Address   bool
	Telephone bool
}

type OpenGraph struct {
	Title         string
	Description   string
	URL           string
	SiteName      string
	Images        []OpenGraphImage
	Locale        string
	Type          string
	PublishedTime time.Time
	Authors       []string
}

type OpenGraphImage struct {
	URL    string
	Width  int
	Height int
	Alt    string
}

type Twitter struct {
	Card        string // e.g. "summary_large_image"
	Title       string
	Description string
	SiteID      string
	Creator     string
	CreatorID   string
	Images      []string
	App         *TwitterApp
}

type TwitterApp struct {
	Name string
	ID   *TwitterAppID
	URL  *TwitterAppURL
}

type TwitterAppID struct {
	IPhone     string
	IPad       string
	GooglePlay string
}

type TwitterAppURL struct {
	IPhone string
	IPad   string
}

type Robots struct {
	Index     bool
	Follow    bool
	NoCache   bool
	GoogleBot *GoogleBot
}

type GoogleBot struct {
	Index           bool
	Follow          bool
	NoImageIndex    bool
	MaxVideoPreview int    // e.g. -1
	MaxImagePreview string // e.g. "large"
	MaxSnippet      int    // e.g. -1
}

type Icons struct {
	Icon     []Icon
	Shortcut []string
	Apple    []AppleIcon
	Other    []OtherIcon
}

type Icon struct {
	URL   string
	Media string
	Type  string
}

type AppleIcon struct {
	URL   string
	Sizes []string
	Type  string
}

type OtherIcon struct {
	Rel string
	URL string
}

type Viewport struct {
	ThemeColor   []ThemeColor
	Width        string
	InitialScale float64
	MaximumScale float64
	UserScalable *bool
	ColorScheme  string
}

type ThemeColor struct {
	Name  string
	Color string
	Media string
}

type Metadata struct {
	Title       string
	TitleFunc   func() TitleData
	Description string

	Generator       string
	ApplicationName string
	Referrer        string
	Keywords        []string
	Authors         []Author
	Creator         string
	Publisher       string
	FormatDetection *FormatDetection

	Canonical string
	Languages map[string]string // "en-US": "/en-US"

	OpenGraph *OpenGraph
	Twitter   *Twitter
	Robots    *Robots
	Icons     *Icons

	Viewport *Viewport

	Manifest string

	// Verification map[string]string
	// AppleWebApp
	// Alternates
	// AppLinks
	// Archives
	// Assets
	// Bookmarks
	// Category

	Other map[string]string
}

// String output for the metadata.
func (m Metadata) String() string {
	var sb strings.Builder

	// Title
	title := m.Title
	if m.TitleFunc != nil {
		titleData := m.TitleFunc()
		if titleData.Absolute != "" {
			title = titleData.Absolute
		} else if titleData.Template != "" {
			title = fmt.Sprintf(titleData.Template, m.Title)
		} else if titleData.Default != "" {
			title = titleData.Default
		} else {
			title = m.Title
		}
	}
	sb.WriteString("<title>")
	sb.WriteString(title)
	sb.WriteString("</title>")
	sb.WriteString("\n")

	// Description
	if m.Description != "" {
		sb.WriteString(`<meta name="description" content="`)
		sb.WriteString(m.Description)
		sb.WriteString(`" />`)
		sb.WriteString("\n")
	}

	// Viewport
	if m.Viewport != nil {
		// Width
		if m.Viewport.Width != "" {
			sb.WriteString(`<meta name="viewport" content="width=`)
			sb.WriteString(m.Viewport.Width)
			if m.Viewport.InitialScale > 0 {
				sb.WriteString(`,initial-scale=`)
				sb.WriteString(fmt.Sprint(m.Viewport.InitialScale))
			}
			if m.Viewport.MaximumScale > 0 {
				sb.WriteString(`,maximum-scale=`)
				sb.WriteString(fmt.Sprint(m.Viewport.MaximumScale))
			}
			if m.Viewport.UserScalable != nil {
				if *m.Viewport.UserScalable {
					sb.WriteString(`,user-scalable=yes`)
				} else {
					sb.WriteString(`,user-scalable=no`)
				}
			}
			if m.Viewport.ColorScheme != "" {
				sb.WriteString(`,color-scheme=`)
				sb.WriteString(m.Viewport.ColorScheme)
			}
			sb.WriteString(`" />`)
			sb.WriteString("\n")
		}
		// ThemeColor
		for _, themeColor := range m.Viewport.ThemeColor {
			sb.WriteString(`<meta name="theme-color" content="`)
			sb.WriteString(themeColor.Color)
			if themeColor.Media != "" {
				sb.WriteString(`" media="`)
				sb.WriteString(themeColor.Media)
			}
			sb.WriteString(`" />`)
			sb.WriteString("\n")
		}
		// ColorScheme
		if m.Viewport.ColorScheme != "" {
			sb.WriteString(`<meta name="color-scheme" content="`)
			sb.WriteString(m.Viewport.ColorScheme)
			sb.WriteString(`" />`)
			sb.WriteString("\n")
		}
	}

	// Generator
	if m.Generator != "" {
		sb.WriteString(`<meta name="generator" content="`)
		sb.WriteString(m.Generator)
		sb.WriteString(`" />`)
		sb.WriteString("\n")
	}

	// ApplicationName
	if m.ApplicationName != "" {
		sb.WriteString(`<meta name="application-name" content="`)
		sb.WriteString(m.ApplicationName)
		sb.WriteString(`" />`)
		sb.WriteString("\n")
	}

	// Referrer
	if m.Referrer != "" {
		sb.WriteString(`<meta name="referrer" content="`)
		sb.WriteString(m.Referrer)
		sb.WriteString(`" />`)
		sb.WriteString("\n")
	}

	// Keywords
	if len(m.Keywords) > 0 {
		sb.WriteString(`<meta name="keywords" content="`)
		sb.WriteString(strings.Join(m.Keywords, ","))
		sb.WriteString(`" />`)
		sb.WriteString("\n")
	}

	// Authors
	for _, author := range m.Authors {
		if author.Name != "" {
			sb.WriteString(`<meta name="author" content="`)
			sb.WriteString(author.Name)
			sb.WriteString(`" />`)
			sb.WriteString("\n")
		}
		if author.URL != "" {
			sb.WriteString(`<link rel="author" href="`)
			sb.WriteString(author.URL)
			sb.WriteString(`" />`)
			sb.WriteString("\n")
		}
	}

	// Creator
	if m.Creator != "" {
		sb.WriteString(`<meta name="creator" content="`)
		sb.WriteString(m.Creator)
		sb.WriteString(`" />`)
		sb.WriteString("\n")
	}

	// Publisher
	if m.Publisher != "" {
		sb.WriteString(`<meta name="publisher" content="`)
		sb.WriteString(m.Publisher)
		sb.WriteString(`" />`)
		sb.WriteString("\n")
	}

	// FormatDetection
	if m.FormatDetection != nil {
		sb.WriteString(`<meta name="format-detection" content="`)
		if m.FormatDetection.Email {
			sb.WriteString("email=no")
		} else {
			sb.WriteString("email=yes")
		}
		sb.WriteString(",")
		if m.FormatDetection.Address {
			sb.WriteString("address=no")
		} else {
			sb.WriteString("address=yes")
		}
		sb.WriteString(",")
		if m.FormatDetection.Telephone {
			sb.WriteString("telephone=no")
		} else {
			sb.WriteString("telephone=yes")
		}
		sb.WriteString(`" />`)
		sb.WriteString("\n")
	}

	// Canonical
	if m.Canonical != "" {
		sb.WriteString(`<link rel="canonical" href="`)
		sb.WriteString(m.Canonical)
		sb.WriteString(`" />`)
		sb.WriteString("\n")
	}

	// Languages
	for lang, href := range m.Languages {
		sb.WriteString(`<link rel="alternate" hreflang="`)
		sb.WriteString(lang)
		sb.WriteString(`" href="`)
		sb.WriteString(href)
		sb.WriteString(`" />`)
		sb.WriteString("\n")
	}

	// OpenGraph
	if m.OpenGraph != nil {
		if m.OpenGraph.Title != "" {
			sb.WriteString(`<meta property="og:title" content="`)
			sb.WriteString(m.OpenGraph.Title)
			sb.WriteString(`" />`)
			sb.WriteString("\n")
		}
		if m.OpenGraph.Description != "" {
			sb.WriteString(`<meta property="og:description" content="`)
			sb.WriteString(m.OpenGraph.Description)
			sb.WriteString(`" />`)
			sb.WriteString("\n")
		}
		if m.OpenGraph.URL != "" {
			sb.WriteString(`<meta property="og:url" content="`)
			sb.WriteString(m.OpenGraph.URL)
			sb.WriteString(`" />`)
			sb.WriteString("\n")
		}
		if m.OpenGraph.SiteName != "" {
			sb.WriteString(`<meta property="og:site_name" content="`)
			sb.WriteString(m.OpenGraph.SiteName)
			sb.WriteString(`" />`)
			sb.WriteString("\n")
		}
		for _, image := range m.OpenGraph.Images {
			sb.WriteString(`<meta property="og:image" content="`)
			sb.WriteString(image.URL)
			sb.WriteString(`" />`)
			sb.WriteString("\n")
			if image.Width > 0 {
				sb.WriteString(`<meta property="og:image:width" content="`)
				sb.WriteString(fmt.Sprint(image.Width))
				sb.WriteString(`" />`)
				sb.WriteString("\n")
			}
			if image.Height > 0 {
				sb.WriteString(`<meta property="og:image:height" content="`)
				sb.WriteString(fmt.Sprint(image.Height))
				sb.WriteString(`" />`)
				sb.WriteString("\n")
			}
			if image.Alt != "" {
				sb.WriteString(`<meta property="og:image:alt" content="`)
				sb.WriteString(image.Alt)
				sb.WriteString(`" />`)
				sb.WriteString("\n")
			}
		}
		if m.OpenGraph.Locale != "" {
			sb.WriteString(`<meta property="og:locale" content="`)
			sb.WriteString(m.OpenGraph.Locale)
			sb.WriteString(`" />`)
			sb.WriteString("\n")
		}
		if m.OpenGraph.Type != "" {
			sb.WriteString(`<meta property="og:type" content="`)
			sb.WriteString(m.OpenGraph.Type)
			sb.WriteString(`" />`)
			sb.WriteString("\n")
		}
		if !m.OpenGraph.PublishedTime.IsZero() {
			sb.WriteString(`<meta property="article:published_time" content="`)
			sb.WriteString(m.OpenGraph.PublishedTime.Format(time.RFC3339))
			sb.WriteString(`" />`)
			sb.WriteString("\n")
		}
		for _, author := range m.OpenGraph.Authors {
			sb.WriteString(`<meta property="article:author" content="`)
			sb.WriteString(author)
			sb.WriteString(`" />`)
			sb.WriteString("\n")
		}
	}

	// Twitter
	if m.Twitter != nil {
		if m.Twitter.Card != "" {
			sb.WriteString(`<meta name="twitter:card" content="`)
			sb.WriteString(m.Twitter.Card)
			sb.WriteString(`" />`)
			sb.WriteString("\n")
		}
		if m.Twitter.Title != "" {
			sb.WriteString(`<meta name="twitter:title" content="`)
			sb.WriteString(m.Twitter.Title)
			sb.WriteString(`" />`)
			sb.WriteString("\n")
		}
		if m.Twitter.Description != "" {
			sb.WriteString(`<meta name="twitter:description" content="`)
			sb.WriteString(m.Twitter.Description)
			sb.WriteString(`" />`)
			sb.WriteString("\n")
		}
		if m.Twitter.SiteID != "" {
			sb.WriteString(`<meta name="twitter:site:id" content="`)
			sb.WriteString(m.Twitter.SiteID)
			sb.WriteString(`" />`)
			sb.WriteString("\n")
		}
		if m.Twitter.Creator != "" {
			sb.WriteString(`<meta name="twitter:creator" content="`)
			sb.WriteString(m.Twitter.Creator)
			sb.WriteString(`" />`)
			sb.WriteString("\n")
		}
		if m.Twitter.CreatorID != "" {
			sb.WriteString(`<meta name="twitter:creator:id" content="`)
			sb.WriteString(m.Twitter.CreatorID)
			sb.WriteString(`" />`)
			sb.WriteString("\n")
		}
		for _, image := range m.Twitter.Images {
			sb.WriteString(`<meta name="twitter:image" content="`)
			sb.WriteString(image)
			sb.WriteString(`" />`)
			sb.WriteString("\n")
		}
		if m.Twitter.App != nil {
			if m.Twitter.App.Name != "" {
				sb.WriteString(`<meta name="twitter:app:name" content="`)
				sb.WriteString(m.Twitter.App.Name)
				sb.WriteString(`" />`)
				sb.WriteString("\n")
			}
			if m.Twitter.App.ID != nil {
				if m.Twitter.App.ID.IPhone != "" {
					sb.WriteString(`<meta name="twitter:app:id:iphone" content="`)
					sb.WriteString(m.Twitter.App.ID.IPhone)
					sb.WriteString(`" />`)
					sb.WriteString("\n")
				}
				if m.Twitter.App.ID.IPad != "" {
					sb.WriteString(`<meta name="twitter:app:id:ipad" content="`)
					sb.WriteString(m.Twitter.App.ID.IPad)
					sb.WriteString(`" />`)
					sb.WriteString("\n")
				}
				if m.Twitter.App.ID.GooglePlay != "" {
					sb.WriteString(`<meta name="twitter:app:id:googleplay" content="`)
					sb.WriteString(m.Twitter.App.ID.GooglePlay)
					sb.WriteString(`" />`)
					sb.WriteString("\n")
				}
			}
			if m.Twitter.App.URL != nil {
				if m.Twitter.App.URL.IPhone != "" {
					sb.WriteString(`<meta name="twitter:app:url:iphone" content="`)
					sb.WriteString(m.Twitter.App.URL.IPhone)
					sb.WriteString(`" />`)
					sb.WriteString("\n")
				}
				if m.Twitter.App.URL.IPad != "" {
					sb.WriteString(`<meta name="twitter:app:url:ipad" content="`)
					sb.WriteString(m.Twitter.App.URL.IPad)
					sb.WriteString(`" />`)
					sb.WriteString("\n")
				}
			}
		}
	}

	// Robots
	if m.Robots != nil {
		sb.WriteString(`<meta name="robots" content="`)
		if m.Robots.Index {
			sb.WriteString(`index`)
		} else {
			sb.WriteString(`noindex`)
		}
		if m.Robots.Follow {
			sb.WriteString(`,follow`)
		} else {
			sb.WriteString(`,nofollow`)
		}
		if m.Robots.NoCache {
			sb.WriteString(`,nocache`)
		} else {
			sb.WriteString(`,cache`)
		}
		sb.WriteString(`" />`)
		sb.WriteString("\n")

		if m.Robots.GoogleBot != nil {
			sb.WriteString(`<meta name="googlebot" content="`)
			if m.Robots.GoogleBot.Index {
				sb.WriteString(`index`)
			} else {
				sb.WriteString(`noindex`)
			}
			if m.Robots.GoogleBot.Follow {
				sb.WriteString(`,follow`)
			} else {
				sb.WriteString(`,nofollow`)
			}
			if m.Robots.GoogleBot.NoImageIndex {
				sb.WriteString(`,noimageindex`)
			} else {
				sb.WriteString(`,imageindex`)
			}
			if m.Robots.GoogleBot.MaxVideoPreview >= 0 {
				sb.WriteString(`,max-video-preview:`)
				sb.WriteString(fmt.Sprint(m.Robots.GoogleBot.MaxVideoPreview))
			}
			if m.Robots.GoogleBot.MaxImagePreview != "" {
				sb.WriteString(`,max-image-preview:`)
				sb.WriteString(m.Robots.GoogleBot.MaxImagePreview)
			}
			if m.Robots.GoogleBot.MaxSnippet >= 0 {
				sb.WriteString(`,max-snippet:`)
				sb.WriteString(fmt.Sprint(m.Robots.GoogleBot.MaxSnippet))
			}
			sb.WriteString(`" />`)
			sb.WriteString("\n")
		}
	}

	// Icons
	if m.Icons != nil {
		for _, icon := range m.Icons.Icon {
			sb.WriteString(`<link rel="icon" href="`)
			sb.WriteString(icon.URL)
			if icon.Type != "" {
				sb.WriteString(`" type="`)
				sb.WriteString(icon.Type)
			}
			if icon.Media != "" {
				sb.WriteString(`" media="`)
				sb.WriteString(icon.Media)
			}
			sb.WriteString(`" />`)
			sb.WriteString("\n")
		}
		for _, shortcut := range m.Icons.Shortcut {
			sb.WriteString(`<link rel="shortcut icon" href="`)
			sb.WriteString(shortcut)
			sb.WriteString(`" />`)
			sb.WriteString("\n")
		}
		for _, apple := range m.Icons.Apple {
			sb.WriteString(`<link rel="apple-touch-icon" href="`)
			sb.WriteString(apple.URL)
			if len(apple.Sizes) > 0 {
				sb.WriteString(`" sizes="`)
				sb.WriteString(strings.Join(apple.Sizes, " "))
			}
			if apple.Type != "" {
				sb.WriteString(`" type="`)
				sb.WriteString(apple.Type)
			}
			sb.WriteString(`" />`)
			sb.WriteString("\n")
		}
		for _, other := range m.Icons.Other {
			sb.WriteString(`<link rel="`)
			sb.WriteString(other.Rel)
			sb.WriteString(`" href="`)
			sb.WriteString(other.URL)
			sb.WriteString(`" />`)
			sb.WriteString("\n")
		}
	}

	// Manifest
	if m.Manifest != "" {
		sb.WriteString(`<link rel="manifest" href="`)
		sb.WriteString(m.Manifest)
		sb.WriteString(`" />`)
		sb.WriteString("\n")
	}

	// Other
	for name, content := range m.Other {
		sb.WriteString(`<meta name="`)
		sb.WriteString(name)
		sb.WriteString(`" content="`)
		sb.WriteString(content)
		sb.WriteString(`" />`)
		sb.WriteString("\n")
	}

	return sb.String()
}
