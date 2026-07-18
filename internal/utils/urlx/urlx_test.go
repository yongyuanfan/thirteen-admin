package urlx

import "testing"

func TestCanonicalStorageURLKeyEncodedAndDecodedMatch(t *testing.T) {
	encoded := "https://vision-1444271626.cos.ap-chengdu.myqcloud.com/presentations/%E9%87%91%E6%B2%99%E5%B0%8F%E5%AD%A6%E5%90%91%E9%98%B3%E8%80%8C%E7%94%9F%E4%B8%BB%E9%A2%98%E6%B4%BB%E5%8A%A8%E6%BC%94%E7%A4%BA%E7%A8%BF_1782371779.pptx"
	decoded := "https://vision-1444271626.cos.ap-chengdu.myqcloud.com/presentations/金沙小学向阳而生主题活动演示稿_1782371779.pptx"
	if CanonicalStorageURLKey(encoded) != CanonicalStorageURLKey(decoded) {
		t.Fatalf("encoded and decoded URLs should share canonical key")
	}
}

func TestCanonicalStorageURLKeyStripsQuery(t *testing.T) {
	unsigned := "https://cos.example.com/documents/a.docx"
	signed := "https://cos.example.com/documents/a.docx?q-sign-algorithm=sha1&q-ak=abc"
	if CanonicalStorageURLKey(unsigned) != CanonicalStorageURLKey(signed) {
		t.Fatalf("signed and unsigned URLs should share canonical key")
	}
}
