package notemetadata

import (
	"fmt"
	"time"

	"github.com/ishida722/krapp-go/models"
)

// add tags は仕様が複雑なので整理してからまた考える
// func AddTag(md string, tags []string) (string, error) {
// 	note, err := models.ParseNote(md)
// 	if err != nil {
// 		return md, err
// 	}
// 	// tags: []string or string
// 	// frontmatter がない場合
// 	if note.FrontMatter == nil {
// 		note.FrontMatter = models.FrontMatter{
// 			"tags": tags,
// 		}
// 		return note.ToString()
// 	}
// 	// tagsを取得
// 	old_tags, ok := note.FrontMatter["tags"]
// 	// tagsがない場合
// 	if !ok {
// 		note.FrontMatter["tags"] = tags
// 		return note.ToString()
// 	}
// 	// tagsがある場合, tagsの型によって処理を分ける
// 	switch v := old_tags.(type) {
// 	// tagsがstring配列の場合
// 	case []string:
// 		found := false
// 		for _, t := range v {
// 			if ts, ok := t.(string); ok && ts == tag {
// 				found = true
// 				break
// 			}
// 		}
// 		if !found {
// 			// append and convert to []string
// 			strs := make([]string, 0, len(v)+1)
// 			for _, t := range v {
// 				if ts, ok := t.(string); ok {
// 					strs = append(strs, ts)
// 				}
// 			}
// 			strs = append(strs, tag)
// 			fm["tags"] = strs
// 		}
// 	case string:
// 		if v != tag {
// 			note.FrontMatter["tags"] = []string{v, tag}
// 		}
// 	}
// 	}
// 	return ReplaceFrontMatter(md, fm)
// }

// AddCreated adds or updates the Created field in the frontmatter to today (YYYY-MM-DD)
func AddCreated(md string) (string, error) {
	note, err := models.ParseNote(md)
	if err != nil {
		return md, fmt.Errorf("failed to parse note: %w", err)
	}
	if note.FrontMatter == nil {
		note.FrontMatter = models.FrontMatter{}
	}
	note.FrontMatter["Created"] = time.Now().Format("2006-01-02")
	return note.ToString()
}

func AddLabel(md string, label string) (string, error) {
	note, err := models.ParseNote(md)
	if err != nil {
		return md, err
	}
	// frontmatter がない場合
	if note.FrontMatter == nil {
		note.FrontMatter = models.FrontMatter{}
	}
	// labelを更新
	note.FrontMatter["label"] = label
	return note.ToString()
}
