package search

import (
	"sort"
	"strings"

	"matu7/pkg/models"
)

// TagCount 标签及对应的工具数量
type TagCount struct {
	Tag   string
	Count int
}

// GetAllOfflineToolTags 获取所有离线工具的标签
func GetAllOfflineToolTags(tools []models.OfflineTool) []string {
	tagMap := make(map[string]bool)

	for _, tool := range tools {
		for _, tag := range tool.Tags {
			tagMap[tag] = true
		}
	}

	var tags []string
	for tag := range tagMap {
		tags = append(tags, tag)
	}

	// 对标签进行排序
	sort.Strings(tags)

	return tags
}

// GetAllOfflineToolTagsWithCount 获取所有离线工具的标签及其数量
func GetAllOfflineToolTagsWithCount(tools []models.OfflineTool) []TagCount {
	tagMap := make(map[string]int)

	for _, tool := range tools {
		for _, tag := range tool.Tags {
			tagMap[tag]++
		}
	}

	var tagCounts []TagCount
	for tag, count := range tagMap {
		tagCounts = append(tagCounts, TagCount{
			Tag:   tag,
			Count: count,
		})
	}

	// 对标签进行排序
	sort.Slice(tagCounts, func(i, j int) bool {
		return tagCounts[i].Tag < tagCounts[j].Tag
	})

	return tagCounts
}

// GetAllWebToolTags 获取所有网页工具的标签
func GetAllWebToolTags(tools []models.WebTool) []string {
	tagMap := make(map[string]bool)

	for _, tool := range tools {
		for _, tag := range tool.Tags {
			tagMap[tag] = true
		}
	}

	var tags []string
	for tag := range tagMap {
		tags = append(tags, tag)
	}

	// 对标签进行排序
	sort.Strings(tags)

	return tags
}

// GetAllWebToolTagsWithCount 获取所有网页工具的标签及其数量
func GetAllWebToolTagsWithCount(tools []models.WebTool) []TagCount {
	tagMap := make(map[string]int)

	for _, tool := range tools {
		for _, tag := range tool.Tags {
			tagMap[tag]++
		}
	}

	var tagCounts []TagCount
	for tag, count := range tagMap {
		tagCounts = append(tagCounts, TagCount{
			Tag:   tag,
			Count: count,
		})
	}

	// 对标签进行排序
	sort.Slice(tagCounts, func(i, j int) bool {
		return tagCounts[i].Tag < tagCounts[j].Tag
	})

	return tagCounts
}

// SearchOfflineToolsByTag 根据标签搜索离线工具
func SearchOfflineToolsByTag(tools []models.OfflineTool, tag string) []models.OfflineTool {
	if tag == "" {
		return tools
	}

	var results []models.OfflineTool
	for _, tool := range tools {
		for _, t := range tool.Tags {
			if strings.Contains(strings.ToLower(t), strings.ToLower(tag)) {
				results = append(results, tool)
				break
			}
		}
	}

	return results
}

// SearchWebToolsByTag 根据标签搜索网页工具
func SearchWebToolsByTag(tools []models.WebTool, tag string) []models.WebTool {
	if tag == "" {
		return tools
	}

	var results []models.WebTool
	for _, tool := range tools {
		for _, t := range tool.Tags {
			if strings.Contains(strings.ToLower(t), strings.ToLower(tag)) {
				results = append(results, tool)
				break
			}
		}
	}

	return results
}

// GetAllNotesTags 获取所有笔记的标签
func GetAllNotesTags(notes []models.Note) []string {
	tagMap := make(map[string]bool)

	for _, note := range notes {
		for _, tag := range note.Tags {
			tagMap[tag] = true
		}
	}

	var tags []string
	for tag := range tagMap {
		tags = append(tags, tag)
	}

	// 对标签进行排序
	sort.Strings(tags)

	return tags
}

// GetAllNotesTagsWithCount 获取所有笔记的标签及其数量
func GetAllNotesTagsWithCount(notes []models.Note) []TagCount {
	tagMap := make(map[string]int)

	for _, note := range notes {
		for _, tag := range note.Tags {
			tagMap[tag]++
		}
	}

	var tagCounts []TagCount
	for tag, count := range tagMap {
		tagCounts = append(tagCounts, TagCount{
			Tag:   tag,
			Count: count,
		})
	}

	// 对标签进行排序
	sort.Slice(tagCounts, func(i, j int) bool {
		return tagCounts[i].Tag < tagCounts[j].Tag
	})

	return tagCounts
}

// SearchNotesByTag 根据标签搜索笔记
func SearchNotesByTag(notes []models.Note, tag string) []models.Note {
	if tag == "" {
		return notes
	}

	var results []models.Note
	for _, note := range notes {
		for _, t := range note.Tags {
			if strings.Contains(strings.ToLower(t), strings.ToLower(tag)) {
				results = append(results, note)
				break
			}
		}
	}

	return results
}
