package textx

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/cloudwego/eino/schema"
)

// DocumentArtifactRef 表示文档类 MCP 产物引用，用于从助手正文中剥离链接与文件名。
type DocumentArtifactRef struct {
	StorageURL  string
	FileName    string
	DisplayName string
}

var (
	documentArtifactURLPattern               = regexp.MustCompile(`https?://[^\s)\]>"]+\.(?:pptx|docx|xlsx)(?:\?[^\s)\]>"]*)?`)
	documentArtifactMarkdownLinkPattern      = regexp.MustCompile(`\[[^\]]*\]\((https?://[^\s)]+)\)`)
	documentArtifactEmptyMarkdownLinkPattern = regexp.MustCompile(`\[[^\]]*\]\(\s*\)`)
	documentArtifactOrphanLabelPattern       = regexp.MustCompile(`(?m)^[^\S\n]*(?:下载地址|下载链接|文件地址|文件名)[：:][^\S\n]*$`)
	documentArtifactInlineLabelPattern       = regexp.MustCompile(`(?:下载地址|下载链接|文件地址|文件名)[：:]\s*`)
)

// StripDocumentArtifactReferences 从助手正文中移除文档产物下载链接与文件名。
func StripDocumentArtifactReferences(content string, files []DocumentArtifactRef) string {
	content = strings.TrimSpace(content)
	if content == "" {
		return ""
	}

	knownURLs := make(map[string]struct{}, len(files)+4)
	for _, file := range files {
		if url := strings.TrimSpace(file.StorageURL); url != "" {
			knownURLs[url] = struct{}{}
		}
	}

	content = documentArtifactMarkdownLinkPattern.ReplaceAllStringFunc(content, func(match string) string {
		submatches := documentArtifactMarkdownLinkPattern.FindStringSubmatch(match)
		if len(submatches) < 2 {
			return match
		}
		url := strings.TrimSpace(submatches[1])
		if _, ok := knownURLs[url]; ok || looksLikeDocumentArtifactURL(url) {
			return ""
		}
		return match
	})

	for url := range knownURLs {
		content = strings.ReplaceAll(content, url, "")
	}
	for _, file := range files {
		for _, name := range []string{file.FileName, file.DisplayName} {
			name = strings.TrimSpace(name)
			if name == "" {
				continue
			}
			content = strings.ReplaceAll(content, name, "")
		}
	}

	content = documentArtifactURLPattern.ReplaceAllStringFunc(content, func(rawURL string) string {
		trimmed := strings.TrimSpace(strings.TrimRight(rawURL, ".,;"))
		if looksLikeDocumentArtifactURL(trimmed) {
			return ""
		}
		if _, ok := knownURLs[trimmed]; ok {
			return ""
		}
		return rawURL
	})

	content = documentArtifactEmptyMarkdownLinkPattern.ReplaceAllString(content, "")
	content = documentArtifactInlineLabelPattern.ReplaceAllString(content, "")
	content = documentArtifactOrphanLabelPattern.ReplaceAllString(content, "")
	content = collapseBlankLines(content)
	return strings.TrimSpace(content)
}

func looksLikeDocumentArtifactURL(rawURL string) bool {
	value := strings.ToLower(strings.TrimSpace(rawURL))
	if value == "" {
		return false
	}
	return strings.Contains(value, "/presentations/") ||
		strings.Contains(value, "/documents/") ||
		strings.Contains(value, "/excels/")
}

func collapseBlankLines(content string) string {
	lines := strings.Split(content, "\n")
	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		filtered = append(filtered, strings.TrimRight(line, " \t，,"))
	}
	return strings.Join(filtered, "\n")
}

// ExtractAssistantDisplayContent 提取模型返回的可展示内容，支持文本和图片多模态输出。
func ExtractAssistantDisplayContent(message *schema.Message) string {
	if message == nil {
		return ""
	}

	parts := make([]string, 0, 4)
	if message.Content != "" {
		parts = append(parts, message.Content)
	}

	for _, item := range message.AssistantGenMultiContent {
		switch item.Type {
		case schema.ChatMessagePartTypeText:
			if item.Text != "" {
				parts = append(parts, item.Text)
			}
		case schema.ChatMessagePartTypeImageURL:
			if item.Image == nil {
				continue
			}
			if markdown := buildImageMarkdown(item.Image); markdown != "" {
				parts = append(parts, markdown)
			}
		}
	}

	return joinAssistantDisplayParts(parts)
}

func joinAssistantDisplayParts(parts []string) string {
	if len(parts) == 0 {
		return ""
	}

	var builder strings.Builder
	for _, part := range parts {
		if part == "" {
			continue
		}
		if builder.Len() > 0 && shouldInsertMarkdownSeparator(builder.String(), part) {
			builder.WriteString("\n\n")
		}
		builder.WriteString(part)
	}
	return builder.String()
}

func shouldInsertMarkdownSeparator(current string, next string) bool {
	lastRune, _ := utf8.DecodeLastRuneInString(current)
	if lastRune == utf8.RuneError && len(current) == 0 {
		return false
	}
	if lastRune == '\n' {
		if strings.HasPrefix(next, "\n") {
			return false
		}
		return false
	}
	if strings.HasPrefix(next, "\n") {
		return false
	}
	if looksLikeMarkdownImage(current) || looksLikeMarkdownImage(next) {
		return true
	}
	return false
}

func looksLikeMarkdownImage(value string) bool {
	trimmed := strings.TrimSpace(value)
	return strings.HasPrefix(trimmed, "![](") && strings.HasSuffix(trimmed, ")")
}

func buildImageMarkdown(image *schema.MessageOutputImage) string {
	if image == nil {
		return ""
	}

	if image.URL != nil {
		url := strings.TrimSpace(*image.URL)
		if url != "" {
			return fmt.Sprintf("![](%s)", url)
		}
	}

	if image.Base64Data != nil {
		data := strings.TrimSpace(*image.Base64Data)
		if data == "" {
			return ""
		}
		mimeType := strings.TrimSpace(image.MIMEType)
		if mimeType == "" {
			mimeType = "image/png"
		}
		return fmt.Sprintf("![](data:%s;base64,%s)", mimeType, data)
	}

	return ""
}

// TrimSurroundingQuotes 移除首尾成对引号。
func TrimSurroundingQuotes(content string) string {
	title := strings.TrimSpace(content)
	if title == "" {
		return ""
	}

	quotePairs := [][2]string{
		{"\"", "\""},
		{"'", "'"},
		{"`", "`"},
		{"“", "”"},
		{"‘", "’"},
		{"「", "」"},
		{"『", "』"},
		{"《", "》"},
	}

	for {
		trimmed := strings.TrimSpace(title)
		updated := false
		for _, pair := range quotePairs {
			if strings.HasPrefix(trimmed, pair[0]) && strings.HasSuffix(trimmed, pair[1]) {
				trimmed = strings.TrimSpace(strings.TrimPrefix(strings.TrimSuffix(trimmed, pair[1]), pair[0]))
				updated = true
				break
			}
		}
		title = trimmed
		if !updated {
			break
		}
	}

	return title
}

// NormalizeTitle 清洗会话标题文本。
func NormalizeTitle(raw string) string {
	title := strings.TrimSpace(raw)
	if title == "" {
		return ""
	}

	title = strings.ReplaceAll(title, "\r\n", "\n")
	title = strings.ReplaceAll(title, "\r", "\n")
	title = strings.ReplaceAll(title, "\n", " ")
	title = strings.Join(strings.Fields(title), " ")
	title = TrimSurroundingQuotes(title)

	return strings.TrimSpace(TruncateRunes(title, 60, "..."))
}

// TruncateRunes 按 rune 截断文本并追加后缀。
func TruncateRunes(value string, maxLen int, suffix string) string {
	if maxLen <= 0 {
		return ""
	}
	runes := []rune(value)
	if len(runes) <= maxLen {
		return value
	}
	return string(runes[:maxLen]) + suffix
}

// BuildTitleFromContent 根据文本构造标题，空值时返回兜底标题。
func BuildTitleFromContent(content string, maxLen int, fallback string) string {
	title := strings.TrimSpace(content)
	if title == "" {
		return fallback
	}
	return TruncateRunes(title, maxLen, "...")
}
