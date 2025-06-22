package search

import (
	"strings"

	"matu7/pkg/models"
)

// FuzzySearchOfflineTools 模糊搜索离线工具
func FuzzySearchOfflineTools(tools []models.OfflineTool, query string) []models.OfflineTool {
	if query == "" {
		return tools
	}

	queries := strings.Split(query, ",")
	var results []models.OfflineTool

	for _, tool := range tools {
		match := true
		for _, q := range queries {
			q = strings.TrimSpace(q)
			if q == "" {
				continue
			}

			// 在名称、标签和描述中搜索
			nameMatch := strings.Contains(strings.ToLower(tool.Name), strings.ToLower(q))
			descMatch := strings.Contains(strings.ToLower(tool.Description), strings.ToLower(q))

			// 检查标签匹配
			tagMatch := false
			for _, tag := range tool.Tags {
				if strings.Contains(strings.ToLower(tag), strings.ToLower(q)) {
					tagMatch = true
					break
				}
			}

			// 如果都不匹配，则排除此工具
			if !nameMatch && !descMatch && !tagMatch {
				match = false
				break
			}
		}

		if match {
			results = append(results, tool)
		}
	}

	return results
}

// FuzzySearchWebTools 模糊搜索网页工具
func FuzzySearchWebTools(tools []models.WebTool, query string) []models.WebTool {
	if query == "" {
		return tools
	}

	queries := strings.Split(query, ",")
	var results []models.WebTool

	for _, tool := range tools {
		match := true
		for _, q := range queries {
			q = strings.TrimSpace(q)
			if q == "" {
				continue
			}

			// 在名称、标签和描述中搜索
			nameMatch := strings.Contains(strings.ToLower(tool.Name), strings.ToLower(q))
			descMatch := strings.Contains(strings.ToLower(tool.Description), strings.ToLower(q))

			// 检查标签匹配
			tagMatch := false
			for _, tag := range tool.Tags {
				if strings.Contains(strings.ToLower(tag), strings.ToLower(q)) {
					tagMatch = true
					break
				}
			}

			// 如果都不匹配，则排除此工具
			if !nameMatch && !descMatch && !tagMatch {
				match = false
				break
			}
		}

		if match {
			results = append(results, tool)
		}
	}

	return results
}

// FuzzySearchNotes 模糊搜索笔记
func FuzzySearchNotes(notes []models.Note, query string) []models.Note {
	if query == "" {
		return notes
	}

	queries := strings.Split(query, ",")
	var results []models.Note

	for _, note := range notes {
		match := true
		for _, q := range queries {
			q = strings.TrimSpace(q)
			if q == "" {
				continue
			}

			// 在标题、标签和笔记内容中搜索
			titleMatch := strings.Contains(strings.ToLower(note.Title), strings.ToLower(q))
			noteMatch := strings.Contains(strings.ToLower(note.Note), strings.ToLower(q))
			sourceMatch := strings.Contains(strings.ToLower(note.Source), strings.ToLower(q))

			// 检查标签匹配
			tagMatch := false
			for _, tag := range note.Tags {
				if strings.Contains(strings.ToLower(tag), strings.ToLower(q)) {
					tagMatch = true
					break
				}
			}

			// 如果都不匹配，则排除此笔记
			if !titleMatch && !noteMatch && !tagMatch && !sourceMatch {
				match = false
				break
			}
		}

		if match {
			results = append(results, note)
		}
	}

	return results
}
