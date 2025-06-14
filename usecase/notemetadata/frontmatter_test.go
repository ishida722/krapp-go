package notemetadata

import (
	"testing"
)

func TestHasFrontMatter(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		expect bool
	}{
		{"has frontmatter", "---\ntitle: test\n---\nbody", true},
		{"no frontmatter", "body only", false},
		{"frontmatter not at top", "body\n---\ntitle: test\n---", false},
	}
	for _, c := range cases {
		if got := HasFrontMatter(c.input); got != c.expect {
			t.Errorf("%s: want %v got %v", c.name, c.expect, got)
		}
	}
}

func TestExtractFrontMatter(t *testing.T) {
	md := "---\ntitle: test\ntag: foo\n---\nbody text"
	yaml, body, found := ExtractFrontMatter(md)
	if !found {
		t.Fatal("frontmatter not found")
	}
	if yaml != "title: test\ntag: foo\n" {
		t.Errorf("yaml wrong: %q", yaml)
	}
	if body != "body text" {
		t.Errorf("body wrong: %q", body)
	}
}

func TestParseFrontMatter(t *testing.T) {
	md := "---\ntitle: test\ntag: foo\n---\nbody"
	m, ok := ParseFrontMatter(md)
	if !ok {
		t.Fatal("parse failed")
	}
	if m["title"] != "test" || m["tag"] != "foo" {
		t.Errorf("unexpected map: %#v", m)
	}
}

func TestReplaceFrontMatter(t *testing.T) {
	md := "body text"
	m := map[string]interface{}{"title": "new", "tag": "bar"}
	out, err := ReplaceFrontMatter(md, m)
	if err != nil {
		t.Fatal(err)
	}
	if !HasFrontMatter(out) {
		t.Error("frontmatter not added")
	}
	m2, ok := ParseFrontMatter(out)
	if !ok || m2["title"] != "new" || m2["tag"] != "bar" {
		t.Errorf("unexpected: %#v", m2)
	}
}

func TestRemoveFrontMatter(t *testing.T) {
	md := "---\ntitle: test\n---\nbody"
	out := RemoveFrontMatter(md)
	if out != "body" {
		t.Errorf("got %q", out)
	}
	md2 := "body only"
	if RemoveFrontMatter(md2) != md2 {
		t.Error("should not change if no frontmatter")
	}
}

func TestAddTag(t *testing.T) {
	cases := []struct {
		name       string
		input      string
		tag        string
		expectTags []string
	}{
		{
			"no frontmatter",
			"本文だけ",
			"foo",
			[]string{"foo"},
		},
		{
			"frontmatter no tags",
			"---\ntitle: test\n---\nbody",
			"bar",
			[]string{"bar"},
		},
		{
			"frontmatter with tags (string)",
			"---\ntags: hoge\n---\nbody",
			"fuga",
			[]string{"hoge", "fuga"},
		},
		{
			"frontmatter with tags ([]string)",
			"---\ntags:\n  - a\n  - b\n---\nbody",
			"c",
			[]string{"a", "b", "c"},
		},
		{
			"frontmatter with tags ([]interface{})",
			"---\ntags:\n  - x\n  - y\n---\nbody",
			"z",
			[]string{"x", "y", "z"},
		},
		{
			"already has tag",
			"---\ntags:\n  - foo\n  - bar\n---\nbody",
			"foo",
			[]string{"foo", "bar"},
		},
	}
	for _, c := range cases {
		out, err := AddTag(c.input, c.tag)
		if err != nil {
			t.Errorf("%s: unexpected error: %v", c.name, err)
			continue
		}
		m, ok := ParseFrontMatter(out)
		if !ok {
			t.Errorf("%s: no frontmatter", c.name)
			continue
		}
		tags, ok := m["tags"]
		if !ok {
			t.Errorf("%s: no tags field", c.name)
			continue
		}
		var got []string
		switch v := tags.(type) {
		case []interface{}:
			for _, t := range v {
				if ts, ok := t.(string); ok {
					got = append(got, ts)
				}
			}
		case []string:
			got = v
		case string:
			got = []string{v}
		}
		if len(got) != len(c.expectTags) {
			t.Errorf("%s: tags length mismatch: got %v want %v", c.name, got, c.expectTags)
			continue
		}
		for i := range got {
			if got[i] != c.expectTags[i] {
				t.Errorf("%s: tags[%d] = %q, want %q", c.name, i, got[i], c.expectTags[i])
			}
		}
	}
}

func TestAddCreated(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"no frontmatter", "本文だけ"},
		{"frontmatter exists", "---\ntitle: test\n---\nbody"},
	}
	for _, c := range cases {
		out, err := AddCreated(c.input)
		if err != nil {
			t.Errorf("%s: unexpected error: %v", c.name, err)
			continue
		}
		m, ok := ParseFrontMatter(out)
		if !ok {
			t.Errorf("%s: no frontmatter", c.name)
			continue
		}
		created, ok := m["Created"]
		if !ok {
			t.Errorf("%s: no Created field", c.name)
			continue
		}
		// 日付の形式だけチェック
		if s, ok := created.(string); !ok || len(s) != 10 {
			t.Errorf("%s: Created field format invalid: %v", c.name, created)
		}
	}
}
