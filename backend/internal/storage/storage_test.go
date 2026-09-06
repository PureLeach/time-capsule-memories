package storage

import "testing"

func TestAttachmentName(t *testing.T) {
	cases := []struct {
		name        string
		key         string
		contentType string
		want        string
	}{
		{name: "image type becomes extension", key: "folder/abc", contentType: "image/jpeg", want: "abc.jpeg"},
		{name: "charset parameter is dropped", key: "folder/abc", contentType: "image/png; charset=binary", want: "abc.png"},
		{name: "non image type is left alone", key: "folder/abc", contentType: "application/pdf", want: "abc"},
		{name: "missing type is left alone", key: "folder/abc", want: "abc"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := attachmentName(tc.key, tc.contentType); got != tc.want {
				t.Fatalf("attachmentName(%q, %q) = %q, want %q", tc.key, tc.contentType, got, tc.want)
			}
		})
	}
}
