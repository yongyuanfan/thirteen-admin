package textx

import "testing"

func TestStripDocumentArtifactReferencesRemovesKnownURLAndFileName(t *testing.T) {
	url := "http://192.168.2.161:9000/clash-admin/presentations/demo-20260609143551.pptx"
	content := "已为你生成完成，下载地址：" + url + "，文件名：demo-20260609143551.pptx"
	files := []DocumentArtifactRef{{
		StorageURL:  url,
		FileName:    "demo-20260609143551.pptx",
		DisplayName: "demo-20260609143551.pptx",
	}}

	got := StripDocumentArtifactReferences(content, files)
	want := "已为你生成完成"

	if got != want {
		t.Fatalf("unexpected content:\nwant: %q\ngot:  %q", want, got)
	}
}

func TestStripDocumentArtifactReferencesRemovesMarkdownLink(t *testing.T) {
	url := "https://example.com/excels/generated.xlsx"
	content := "表格已生成，[点击下载](" + url + ")"
	files := []DocumentArtifactRef{{
		StorageURL: url,
		FileName:   "generated.xlsx",
	}}

	got := StripDocumentArtifactReferences(content, files)
	want := "表格已生成"

	if got != want {
		t.Fatalf("unexpected content:\nwant: %q\ngot:  %q", want, got)
	}
}

func TestStripDocumentArtifactReferencesRemovesArtifactURLWithoutKnownFiles(t *testing.T) {
	content := "下载链接：https://example.com/documents/report-20260609.docx"
	got := StripDocumentArtifactReferences(content, nil)
	if got != "" {
		t.Fatalf("unexpected content: %q", got)
	}
}

func TestStripDocumentArtifactReferencesKeepsNonArtifactURL(t *testing.T) {
	content := "官网链接：https://example.com/files/demo.pptx"
	got := StripDocumentArtifactReferences(content, nil)
	if got != content {
		t.Fatalf("expected non-artifact link to be preserved, got %q", got)
	}
}

func TestStripDocumentArtifactReferencesRemovesOrphanLabelLine(t *testing.T) {
	content := "文档已生成。\n下载地址：\n文件名：\n请查看附件。"
	got := StripDocumentArtifactReferences(content, nil)
	want := "文档已生成。\n请查看附件。"

	if got != want {
		t.Fatalf("unexpected content:\nwant: %q\ngot:  %q", want, got)
	}
}

func TestStripDocumentArtifactReferencesKeepsNormalText(t *testing.T) {
	content := "这是一份关于消防安全的 Word 文档说明，包含标题与正文结构。"
	got := StripDocumentArtifactReferences(content, nil)
	if got != content {
		t.Fatalf("expected normal text to be preserved, got %q", got)
	}
}
