package textx

import (
	"testing"

	"github.com/cloudwego/eino/schema"
)

func TestExtractAssistantDisplayContentPreservesMarkdownStreamingChunks(t *testing.T) {
	msg := &schema.Message{Content: "`unzip` 是类Unix系统（Linux、macOS等）中用于解压`.zip`格式压缩包的常用命令，以下是它的常见用法和参数说明：\n\n### 基础语法\n```bash\nunzip [选项] 压缩包名称.zip\n```\n\n1. **直接解压到当前目录**\n"}

	got := ExtractAssistantDisplayContent(msg)

	if got != msg.Content {
		t.Fatalf("expected content to be preserved, got %q", got)
	}
}

func TestExtractAssistantDisplayContentPreservesChunkWhitespace(t *testing.T) {
	msg := &schema.Message{
		AssistantGenMultiContent: []schema.MessageOutputPart{
			{Type: schema.ChatMessagePartTypeText, Text: "### 基础语法\n"},
			{Type: schema.ChatMessagePartTypeText, Text: "```bash\nunzip test.zip\n```\n"},
			{Type: schema.ChatMessagePartTypeText, Text: "\n1. **直接解压到当前目录**\n"},
		},
	}

	got := ExtractAssistantDisplayContent(msg)
	want := "### 基础语法\n```bash\nunzip test.zip\n```\n\n1. **直接解压到当前目录**\n"

	if got != want {
		t.Fatalf("unexpected content:\nwant: %q\ngot:  %q", want, got)
	}
}

func TestExtractAssistantDisplayContentSeparatesMarkdownImages(t *testing.T) {
	url := "https://example.com/demo.png"
	msg := &schema.Message{
		AssistantGenMultiContent: []schema.MessageOutputPart{
			{Type: schema.ChatMessagePartTypeText, Text: "图片如下："},
			{Type: schema.ChatMessagePartTypeImageURL, Image: &schema.MessageOutputImage{MessagePartCommon: schema.MessagePartCommon{URL: &url}}},
			{Type: schema.ChatMessagePartTypeText, Text: "\n请查看说明。"},
		},
	}

	got := ExtractAssistantDisplayContent(msg)
	want := "图片如下：\n\n![](https://example.com/demo.png)\n请查看说明。"

	if got != want {
		t.Fatalf("unexpected image markdown:\nwant: %q\ngot:  %q", want, got)
	}
}
