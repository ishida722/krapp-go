package frontmatter

import (
	"bytes"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// HasFrontMatter checks if markdown text has yaml frontmatter at the top
func HasFrontMatter(md string) bool {
	md = strings.TrimLeft(md, "\ufeff\n\r ") // BOMや空白除去
	return strings.HasPrefix(md, "---\n")
}

// ExtractFrontMatter returns (yaml, body, found)
func ExtractFrontMatter(md string) (string, string, bool) {
	md = strings.TrimLeft(md, "\ufeff\n\r ")
	if !strings.HasPrefix(md, "---\n") {
		return "", md, false
	}
	end := strings.Index(md[4:], "---")
	if end == -1 {
		return "", md, false
	}
	end += 4
	// ---\n ...yaml... \n---
	yamlBlock := md[4:end]
	body := md[end+3:]
	body = strings.TrimLeft(body, "\r\n")
	return yamlBlock, body, true
}

// ParseFrontMatter parses yaml frontmatter into map
func ParseFrontMatter(md string) (map[string]interface{}, bool) {
	yamlBlock, _, ok := ExtractFrontMatter(md)
	if !ok {
		return nil, false
	}
	m := make(map[string]interface{})
	if err := yaml.Unmarshal([]byte(yamlBlock), &m); err != nil {
		return nil, false
	}
	return m, true
}

// ReplaceFrontMatter returns markdown with new yaml frontmatter
func ReplaceFrontMatter(md string, data map[string]interface{}) (string, error) {
	var buf bytes.Buffer
	buf.WriteString("---\n")
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(data); err != nil {
		return "", err
	}
	buf.WriteString("---\n")
	_, body, found := ExtractFrontMatter(md)
	if !found {
		body = md
	}
	buf.WriteString(strings.TrimLeft(body, "\r\n"))
	return buf.String(), nil
}

// RemoveFrontMatter removes yaml frontmatter if present
func RemoveFrontMatter(md string) string {
	_, body, found := ExtractFrontMatter(md)
	if found {
		return body
	}
	return md
}

// AddTag adds a tag to the frontmatter (creates frontmatter/tags if needed)
func AddTag(md string, tag string) (string, error) {
	fm, _ := ParseFrontMatter(md)
	if fm == nil {
		fm = make(map[string]interface{})
	}
	// tags: []string or string
	tags, ok := fm["tags"]
	if !ok {
		fm["tags"] = []string{tag}
	} else {
		switch v := tags.(type) {
		case []interface{}:
			found := false
			for _, t := range v {
				if ts, ok := t.(string); ok && ts == tag {
					found = true
					break
				}
			}
			if !found {
				// append and convert to []string
				strs := make([]string, 0, len(v)+1)
				for _, t := range v {
					if ts, ok := t.(string); ok {
						strs = append(strs, ts)
					}
				}
				strs = append(strs, tag)
				fm["tags"] = strs
			}
		case []string:
			found := false
			for _, ts := range v {
				if ts == tag {
					found = true
					break
				}
			}
			if !found {
				fm["tags"] = append(v, tag)
			}
		case string:
			if v != tag {
				fm["tags"] = []string{v, tag}
			}
		}
	}
	return ReplaceFrontMatter(md, fm)
}

// AddCreated adds or updates the Created field in the frontmatter to today (YYYY-MM-DD)
func AddCreated(md string) (string, error) {
	fm, _ := ParseFrontMatter(md)
	if fm == nil {
		fm = make(map[string]interface{})
	}
	fm["Created"] = time.Now().Format("2006-01-02")
	return ReplaceFrontMatter(md, fm)
}
