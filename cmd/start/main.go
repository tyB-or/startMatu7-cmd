package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/c-bata/go-prompt"

	"matu7/internal/config"
	"matu7/internal/launcher"
	"matu7/internal/search"
	"matu7/pkg/models"
)

// 全局表格宽度常量
const (
	// 表格列宽
	NameColWidth   = 22
	TagsColWidth   = 26
	DescColWidth   = 32
	TitleColWidth  = 32
	SourceColWidth = 22

	// 表格总宽度
	TableTotalWidth = 90
)

var cfg *config.Config

// TableColumn 定义表格列的属性
type TableColumn struct {
	Title string
	Width int
	Color string
}

// TableRow 定义表格行的数据
type TableRow struct {
	Columns []string
}

// Table 定义表格结构
type Table struct {
	Title         string
	TitleColor    string
	BorderColor   string
	HeaderColor   string
	CellColor     string
	Columns       []TableColumn
	Rows          []TableRow
	CategoryTitle string
}

func main() {
	// 处理命令行参数
	if len(os.Args) > 1 {
		// 检查是否是添加配置路径命令
		if os.Args[1] == "--add-path" && len(os.Args) > 2 {
			configPath := os.Args[2]
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				fmt.Printf("错误: 配置路径不存在: %s\n", configPath)
				os.Exit(1)
			}

			if err := config.SaveConfigPath(configPath); err != nil {
				fmt.Printf("保存配置路径失败: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("配置路径已保存: %s\n", configPath)
			os.Exit(0)
		}

		// 获取配置路径
		configPath, err := config.GetConfigPath()
		if err != nil {
			fmt.Println("错误: 未设置配置路径，请使用 --add-path 命令添加配置路径")
			fmt.Println("示例: start --add-path /path/to/config")
			os.Exit(1)
		}

		// 加载配置
		cfg, err = config.LoadConfig(configPath)
		if err != nil {
			fmt.Printf("加载配置失败: %v\n", err)
			os.Exit(1)
		}

		// 处理其他命令
		handleCommandLine(os.Args[1:])
		return
	}

	// 获取配置路径
	configPath, err := config.GetConfigPath()
	if err != nil {
		fmt.Println("错误: 未设置配置路径，请使用 --add-path 命令添加配置路径")
		fmt.Println("示例: start --add-path /path/to/config")
		os.Exit(1)
	}

	// 加载配置
	cfg, err = config.LoadConfig(configPath)
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 进入交互模式
	fmt.Println("欢迎使用 Matu7 工具启动器")
	fmt.Println("输入 'help' 获取帮助")

	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix("matu7> "),
		prompt.OptionTitle("Matu7 工具启动器"),
	)
	p.Run()
}

func handleCommandLine(args []string) {
	if len(args) == 0 {
		return
	}

	switch args[0] {
	case "-t":
		if len(args) < 2 {
			handleOfflineTool("")
			return
		}
		handleOfflineTool(args[1])
	case "-tm":
		if len(args) < 2 {
			displayAllOfflineToolTags()
			return
		}
		handleOfflineToolByTag(args[1])
	case "-w":
		if len(args) < 2 {
			handleWebTool("")
			return
		}
		handleWebTool(args[1])
	case "-wm":
		if len(args) < 2 {
			displayAllWebToolTags()
			return
		}
		handleWebToolByTag(args[1])
	case "-n":
		if len(args) < 2 {
			displayNotes()
			return
		}
		handleNoteSearch(args[1])
	case "-nm":
		if len(args) < 2 {
			displayAllNoteTags()
			return
		}
		handleNoteByTag(args[1])
	case "help":
		displayHelp()
	default:
		fmt.Println("未知命令，输入 'help' 获取帮助")
	}
}

func executor(input string) {
	input = strings.TrimSpace(input)
	if input == "" {
		return
	}

	args := strings.Fields(input)
	handleCommandLine(args)
}

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "-t", Description: "显示所有或搜索启动离线工具"},
		{Text: "-tm", Description: "根据标签搜索离线工具"},
		{Text: "-w", Description: "显示所有或搜索打开网页工具"},
		{Text: "-wm", Description: "根据标签搜索网页工具"},
		{Text: "-n", Description: "显示所有或搜索网页笔记"},
		{Text: "-nm", Description: "根据标签搜索网页笔记"},
		{Text: "help", Description: "显示帮助信息"},
	}

	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func handleOfflineTool(query string) {
	if query == "" {
		// 不加参数时显示所有离线工具
		fmt.Println("\n显示所有离线工具:")
		listOfflineTools()
		return
	}

	// 处理查询条件
	results := search.FuzzySearchOfflineTools(cfg.OfflineTools.Tools, query)

	if len(results) == 0 {
		fmt.Println("未找到匹配的工具")
		// 提供一些建议的标签
		fmt.Println("您可以尝试以下热门标签:")
		displayTopTags(cfg.OfflineTools.Tools, 5)
		return
	} else if len(results) == 1 {
		// 只有一个结果，直接启动
		tool := results[0]
		fmt.Printf("正在启动: %s\n", tool.Name)
		if err := launcher.LaunchOfflineTool(tool); err != nil {
			fmt.Printf("启动失败: %v\n", err)
		}
	} else {
		// 检查是否有名称完全匹配的工具
		var exactMatch *models.OfflineTool
		for i, tool := range results {
			if strings.EqualFold(tool.Name, query) {
				exactMatch = &results[i]
				break
			}
		}

		// 如果有名称完全匹配的工具，询问用户是否直接启动
		if exactMatch != nil {
			fmt.Printf("找到完全匹配的工具: %s\n", exactMatch.Name)
			fmt.Print("是否直接启动? (y/n): ")
			var answer string
			fmt.Scanln(&answer)
			if strings.ToLower(answer) == "y" || strings.ToLower(answer) == "yes" {
				fmt.Printf("正在启动: %s\n", exactMatch.Name)
				if err := launcher.LaunchOfflineTool(*exactMatch); err != nil {
					fmt.Printf("启动失败: %v\n", err)
				}
				return
			}
			fmt.Println() // 添加空行，提高可读性
		}

		// 设置颜色
		const borderColor = "\033[1;36m"
		const headerColor = "\033[1;36m"
		const nameColor = "\033[1;37m"
		const tagsColor = "\033[0;33m"
		const descColor = "\033[0;37m"
		const indexColor = "\033[1;33m"
		const categoryColor = "\033[1;33m"

		// 打印标题
		titleBorder := strings.Repeat("─", TableTotalWidth)
		title := "离线工具搜索结果"
		titleLen := runeWidth(title)
		titlePadding := (TableTotalWidth - titleLen) / 2
		leftPadding := strings.Repeat(" ", titlePadding)
		rightPadding := strings.Repeat(" ", TableTotalWidth-titleLen-titlePadding)

		fmt.Printf("%s┌%s┐\033[0m\n", borderColor, titleBorder)
		fmt.Printf("%s│%s%s%s│\033[0m\n", borderColor, leftPadding, title, rightPadding)
		fmt.Printf("%s└%s┘\033[0m\n", borderColor, titleBorder)

		// 按分类对工具进行分组
		categoryMap := make(map[string][]models.OfflineTool)
		var categories []string

		for _, tool := range results {
			if _, exists := categoryMap[tool.Category]; !exists {
				categories = append(categories, tool.Category)
			}
			categoryMap[tool.Category] = append(categoryMap[tool.Category], tool)
		}

		// 对分类进行排序
		sort.Strings(categories)

		// 创建序号到工具的映射
		indexMap := make(map[int]models.OfflineTool)
		currentIndex := 1

		// 按分类输出工具
		for _, category := range categories {
			tools := categoryMap[category]

			// 对每个分类中的工具按名称排序
			sort.Slice(tools, func(i, j int) bool {
				return tools[i].Name < tools[j].Name
			})

			// 创建分类表格
			categoryTable := Table{
				CategoryTitle: category,
				BorderColor:   categoryColor,
				HeaderColor:   headerColor,
				CellColor:     nameColor,
				Columns: []TableColumn{
					{Title: "序号", Width: 8, Color: indexColor},
					{Title: "名称", Width: NameColWidth, Color: nameColor},
					{Title: "标签", Width: TagsColWidth, Color: tagsColor},
					{Title: "描述", Width: DescColWidth - 10, Color: descColor},
				},
			}

			// 添加数据行
			for _, tool := range tools {
				tags := strings.Join(tool.Tags, ", ")
				tags = truncateString(cleanString(tags), TagsColWidth)

				desc := cleanString(tool.Description)
				desc = truncateString(desc, DescColWidth-10)

				name := cleanString(tool.Name)
				name = truncateString(name, NameColWidth)

				index := fmt.Sprintf("[%d]", currentIndex)
				indexMap[currentIndex] = tool
				currentIndex++

				categoryTable.Rows = append(categoryTable.Rows, TableRow{
					Columns: []string{index, name, tags, desc},
				})
			}

			// 打印分类表格
			printTable(categoryTable)
		}

		// 输出工具总数
		fmt.Printf("\n%s总计: %d 个工具, %d 个分类%s\n", borderColor, len(results), len(categories), "\033[0m")

		// 增加交互性的选择
		fmt.Print("\n请选择要启动的工具 (输入序号或 'q' 退出): ")
		var input string
		fmt.Scanln(&input)

		input = strings.ToLower(input)
		if input == "q" || input == "quit" {
			return
		}

		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("无效的输入")
			return
		}

		if choice > 0 && choice < currentIndex {
			tool := indexMap[choice]
			fmt.Printf("正在启动: %s\n", tool.Name)
			if err := launcher.LaunchOfflineTool(tool); err != nil {
				fmt.Printf("启动失败: %v\n", err)
			}
		} else {
			fmt.Println("无效的选择")
		}
	}
}

// 显示最常用的标签
func displayTopTags(tools []models.OfflineTool, count int) {
	tagCounts := make(map[string]int)

	// 统计标签出现次数
	for _, tool := range tools {
		for _, tag := range tool.Tags {
			tagCounts[tag]++
		}
	}

	// 转换为切片并排序
	type tagCount struct {
		Tag   string
		Count int
	}

	var tagList []tagCount
	for tag, count := range tagCounts {
		tagList = append(tagList, tagCount{Tag: tag, Count: count})
	}

	// 按出现次数排序
	sort.Slice(tagList, func(i, j int) bool {
		return tagList[i].Count > tagList[j].Count
	})

	// 显示前N个标签
	fmt.Print("  ")
	for i := 0; i < count && i < len(tagList); i++ {
		fmt.Printf("\033[0;33m%s\033[0m ", tagList[i].Tag)
	}
	fmt.Println("\n使用 -tm <标签> 命令可按标签搜索工具")
}

func handleOfflineToolByTag(tag string) {
	if tag == "" {
		displayAllOfflineToolTags()
		return
	}

	// 从标签中模糊搜索，而不是完全匹配
	results := []models.OfflineTool{}
	for _, tool := range cfg.OfflineTools.Tools {
		for _, t := range tool.Tags {
			if strings.Contains(strings.ToLower(t), strings.ToLower(tag)) {
				results = append(results, tool)
				break
			}
		}
	}

	if len(results) == 0 {
		fmt.Println("未找到匹配标签的工具")
		return
	}

	// 设置颜色
	const borderColor = "\033[1;36m"
	const headerColor = "\033[1;36m"
	const nameColor = "\033[1;37m"
	const tagsColor = "\033[0;33m"
	const descColor = "\033[0;37m"
	const indexColor = "\033[1;33m"
	const categoryColor = "\033[1;33m"

	// 打印标题
	titleBorder := strings.Repeat("─", TableTotalWidth)
	title := fmt.Sprintf("标签 '%s' 下的离线工具", tag)
	titleLen := runeWidth(title)
	titlePadding := (TableTotalWidth - titleLen) / 2
	leftPadding := strings.Repeat(" ", titlePadding)
	rightPadding := strings.Repeat(" ", TableTotalWidth-titleLen-titlePadding)

	fmt.Printf("%s┌%s┐\033[0m\n", borderColor, titleBorder)
	fmt.Printf("%s│%s%s%s│\033[0m\n", borderColor, leftPadding, title, rightPadding)
	fmt.Printf("%s└%s┘\033[0m\n", borderColor, titleBorder)

	// 按分类对工具进行分组
	categoryMap := make(map[string][]models.OfflineTool)
	var categories []string

	for _, tool := range results {
		if _, exists := categoryMap[tool.Category]; !exists {
			categories = append(categories, tool.Category)
		}
		categoryMap[tool.Category] = append(categoryMap[tool.Category], tool)
	}

	// 对分类进行排序
	sort.Strings(categories)

	// 创建序号到工具的映射
	indexMap := make(map[int]models.OfflineTool)
	currentIndex := 1

	// 按分类输出工具
	for _, category := range categories {
		tools := categoryMap[category]

		// 对每个分类中的工具按名称排序
		sort.Slice(tools, func(i, j int) bool {
			return tools[i].Name < tools[j].Name
		})

		// 创建分类表格
		categoryTable := Table{
			CategoryTitle: category,
			BorderColor:   categoryColor,
			HeaderColor:   headerColor,
			CellColor:     nameColor,
			Columns: []TableColumn{
				{Title: "序号", Width: 8, Color: indexColor},
				{Title: "名称", Width: NameColWidth, Color: nameColor},
				{Title: "描述", Width: DescColWidth + TagsColWidth - 8, Color: descColor},
			},
		}

		// 添加数据行
		for _, tool := range tools {
			desc := cleanString(tool.Description)
			desc = truncateString(desc, DescColWidth+TagsColWidth-8)

			name := cleanString(tool.Name)
			name = truncateString(name, NameColWidth)

			index := fmt.Sprintf("[%d]", currentIndex)
			indexMap[currentIndex] = tool
			currentIndex++

			categoryTable.Rows = append(categoryTable.Rows, TableRow{
				Columns: []string{index, name, desc},
			})
		}

		// 打印分类表格
		printTable(categoryTable)
	}

	// 输出工具总数
	fmt.Printf("\n%s总计: %d 个工具, %d 个分类%s\n", borderColor, len(results), len(categories), "\033[0m")

	for {
		fmt.Print("\n请选择要启动的工具 (输入序号, 输入q退出): ")
		var input string
		fmt.Scanln(&input)

		input = strings.ToLower(input)
		if input == "q" || input == "quit" || input == "exit" {
			fmt.Println("已退出")
			return
		}

		// 尝试将输入转换为数字
		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("无效输入，请输入有效的数字或q退出")
			continue
		}

		if choice > 0 && choice < currentIndex {
			tool := indexMap[choice]
			fmt.Printf("正在启动: %s\n", tool.Name)
			if err := launcher.LaunchOfflineTool(tool); err != nil {
				fmt.Printf("启动失败: %v\n", err)
			}

			// 工具运行结束后询问用户是否还需要启动其他工具
			fmt.Print("\n是否继续选择其他工具? (y/n): ")
			var continueChoice string
			fmt.Scanln(&continueChoice)

			if strings.ToLower(continueChoice) != "y" && strings.ToLower(continueChoice) != "yes" {
				return
			}
		} else {
			fmt.Println("无效的选择")
		}
	}
}

func displayAllOfflineToolTags() {
	tagCounts := search.GetAllOfflineToolTagsWithCount(cfg.OfflineTools.Tools)

	// 设置颜色
	const borderColor = "\033[1;36m"
	const headerColor = "\033[1;36m"
	const nameColor = "\033[1;37m"
	const countColor = "\033[0;33m"

	// 打印标题
	titleBorder := strings.Repeat("─", TableTotalWidth)
	title := "离线工具标签列表"
	titleLen := runeWidth(title)
	titlePadding := (TableTotalWidth - titleLen) / 2
	leftPadding := strings.Repeat(" ", titlePadding)
	rightPadding := strings.Repeat(" ", TableTotalWidth-titleLen-titlePadding)

	fmt.Printf("%s┌%s┐\033[0m\n", borderColor, titleBorder)
	fmt.Printf("%s│%s%s%s│\033[0m\n", borderColor, leftPadding, title, rightPadding)
	fmt.Printf("%s└%s┘\033[0m\n", borderColor, titleBorder)

	// 横排布局，每行显示4个标签
	const tagsPerRow = 4
	const tagWidth = TableTotalWidth/tagsPerRow - 2 // 减去分隔符宽度

	// 创建表格
	fmt.Printf("%s┌", borderColor)
	for i := 0; i < tagsPerRow; i++ {
		fmt.Print(strings.Repeat("─", tagWidth))
		if i < tagsPerRow-1 {
			fmt.Print("┬")
		}
	}
	fmt.Print("┐\033[0m\n")

	// 打印行
	for i := 0; i < len(tagCounts); i += tagsPerRow {
		fmt.Printf("%s│", borderColor)
		for j := 0; j < tagsPerRow; j++ {
			if i+j < len(tagCounts) {
				tag := tagCounts[i+j].Tag
				count := tagCounts[i+j].Count
				displayText := fmt.Sprintf("%s (%d)", tag, count)
				paddedText := padString(displayText, tagWidth-2)
				fmt.Printf(" %s%s %s│", nameColor, paddedText, borderColor)
			} else {
				fmt.Printf(" %s%s %s│", nameColor, padString("", tagWidth-2), borderColor)
			}
		}
		fmt.Print("\033[0m\n")

		// 打印行间分隔符，除了最后一行
		if i+tagsPerRow < len(tagCounts) {
			fmt.Printf("%s├", borderColor)
			for j := 0; j < tagsPerRow; j++ {
				fmt.Print(strings.Repeat("─", tagWidth))
				if j < tagsPerRow-1 {
					fmt.Print("┼")
				}
			}
			fmt.Print("┤\033[0m\n")
		}
	}

	// 打印底部边框
	fmt.Printf("%s└", borderColor)
	for i := 0; i < tagsPerRow; i++ {
		fmt.Print(strings.Repeat("─", tagWidth))
		if i < tagsPerRow-1 {
			fmt.Print("┴")
		}
	}
	fmt.Print("┘\033[0m\n")

	// 输出标签总数
	fmt.Printf("\n%s总计: %d 个标签%s\n", borderColor, len(tagCounts), "\033[0m")
}

func handleWebTool(query string) {
	if query == "" {
		// 不加参数时显示所有网页工具
		fmt.Println("\n显示所有网页工具:")
		listWebTools()
		return
	}

	results := search.FuzzySearchWebTools(cfg.WebTools.Tools, query)

	if len(results) == 0 {
		fmt.Println("未找到匹配的网页工具")
		// 提供一些建议的标签
		fmt.Println("您可以尝试以下热门标签:")
		displayTopWebTags(cfg.WebTools.Tools, 5)
		return
	} else if len(results) == 1 {
		// 只有一个结果，直接打开
		tool := results[0]
		fmt.Printf("正在打开: %s\n", tool.Name)
		if err := launcher.LaunchWebTool(tool); err != nil {
			fmt.Printf("打开失败: %v\n", err)
		}
	} else {
		// 检查是否有名称完全匹配的工具
		var exactMatch *models.WebTool
		for i, tool := range results {
			if strings.EqualFold(tool.Name, query) {
				exactMatch = &results[i]
				break
			}
		}

		// 如果有名称完全匹配的工具，询问用户是否直接启动
		if exactMatch != nil {
			fmt.Printf("找到完全匹配的工具: %s\n", exactMatch.Name)
			fmt.Print("是否直接打开? (y/n): ")
			var answer string
			fmt.Scanln(&answer)
			if strings.ToLower(answer) == "y" || strings.ToLower(answer) == "yes" {
				fmt.Printf("正在打开: %s\n", exactMatch.Name)
				if err := launcher.LaunchWebTool(*exactMatch); err != nil {
					fmt.Printf("打开失败: %v\n", err)
				}
				return
			}
			fmt.Println() // 添加空行，提高可读性
		}

		// 设置颜色
		const borderColor = "\033[1;35m"
		const headerColor = "\033[1;35m"
		const nameColor = "\033[1;37m"
		const tagsColor = "\033[0;33m"
		const descColor = "\033[0;37m"
		const indexColor = "\033[1;33m"
		const categoryColor = "\033[1;33m"

		// 打印标题
		titleBorder := strings.Repeat("─", TableTotalWidth)
		title := "网页工具搜索结果"
		titleLen := runeWidth(title)
		titlePadding := (TableTotalWidth - titleLen) / 2
		leftPadding := strings.Repeat(" ", titlePadding)
		rightPadding := strings.Repeat(" ", TableTotalWidth-titleLen-titlePadding)

		fmt.Printf("%s┌%s┐\033[0m\n", borderColor, titleBorder)
		fmt.Printf("%s│%s%s%s│\033[0m\n", borderColor, leftPadding, title, rightPadding)
		fmt.Printf("%s└%s┘\033[0m\n", borderColor, titleBorder)

		// 按分类对工具进行分组
		categoryMap := make(map[string][]models.WebTool)
		var categories []string

		for _, tool := range results {
			if _, exists := categoryMap[tool.Category]; !exists {
				categories = append(categories, tool.Category)
			}
			categoryMap[tool.Category] = append(categoryMap[tool.Category], tool)
		}

		// 对分类进行排序
		sort.Strings(categories)

		// 创建序号到工具的映射
		indexMap := make(map[int]models.WebTool)
		currentIndex := 1

		// 按分类输出工具
		for _, category := range categories {
			tools := categoryMap[category]

			// 对每个分类中的工具按名称排序
			sort.Slice(tools, func(i, j int) bool {
				return tools[i].Name < tools[j].Name
			})

			// 创建分类表格
			categoryTable := Table{
				CategoryTitle: category,
				BorderColor:   categoryColor,
				HeaderColor:   headerColor,
				CellColor:     nameColor,
				Columns: []TableColumn{
					{Title: "序号", Width: 8, Color: indexColor},
					{Title: "名称", Width: NameColWidth, Color: nameColor},
					{Title: "标签", Width: TagsColWidth, Color: tagsColor},
					{Title: "描述", Width: DescColWidth - 10, Color: descColor},
				},
			}

			// 添加数据行
			for _, tool := range tools {
				tags := strings.Join(tool.Tags, ", ")
				tags = truncateString(cleanString(tags), TagsColWidth)

				desc := cleanString(tool.Description)
				desc = truncateString(desc, DescColWidth-10)

				name := cleanString(tool.Name)
				name = truncateString(name, NameColWidth)

				index := fmt.Sprintf("[%d]", currentIndex)
				indexMap[currentIndex] = tool
				currentIndex++

				categoryTable.Rows = append(categoryTable.Rows, TableRow{
					Columns: []string{index, name, tags, desc},
				})
			}

			// 打印分类表格
			printTable(categoryTable)
		}

		// 输出工具总数
		fmt.Printf("\n%s总计: %d 个网页工具, %d 个分类%s\n", borderColor, len(results), len(categories), "\033[0m")

		// 增加交互性的选择
		fmt.Print("\n请选择要打开的工具 (输入序号或 'q' 退出): ")
		var input string
		fmt.Scanln(&input)

		input = strings.ToLower(input)
		if input == "q" || input == "quit" {
			return
		}

		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("无效的输入")
			return
		}

		if choice > 0 && choice < currentIndex {
			tool := indexMap[choice]
			fmt.Printf("正在打开: %s\n", tool.Name)
			if err := launcher.LaunchWebTool(tool); err != nil {
				fmt.Printf("打开失败: %v\n", err)
			}
		} else {
			fmt.Println("无效的选择")
		}
	}
}

// 显示最常用的网页工具标签
func displayTopWebTags(tools []models.WebTool, count int) {
	tagCounts := make(map[string]int)

	// 统计标签出现次数
	for _, tool := range tools {
		for _, tag := range tool.Tags {
			tagCounts[tag]++
		}
	}

	// 转换为切片并排序
	type tagCount struct {
		Tag   string
		Count int
	}

	var tagList []tagCount
	for tag, count := range tagCounts {
		tagList = append(tagList, tagCount{Tag: tag, Count: count})
	}

	// 按出现次数排序
	sort.Slice(tagList, func(i, j int) bool {
		return tagList[i].Count > tagList[j].Count
	})

	// 显示前N个标签
	fmt.Print("  ")
	for i := 0; i < count && i < len(tagList); i++ {
		fmt.Printf("\033[0;33m%s\033[0m ", tagList[i].Tag)
	}
	fmt.Println("\n使用 -wm <标签> 命令可按标签搜索网页工具")
}

func handleWebToolByTag(tag string) {
	if tag == "" {
		displayAllWebToolTags()
		return
	}

	// 从标签中模糊搜索，而不是完全匹配
	results := []models.WebTool{}
	for _, tool := range cfg.WebTools.Tools {
		for _, t := range tool.Tags {
			if strings.Contains(strings.ToLower(t), strings.ToLower(tag)) {
				results = append(results, tool)
				break
			}
		}
	}

	if len(results) == 0 {
		fmt.Println("未找到匹配标签的工具")
		return
	}

	// 设置颜色
	const borderColor = "\033[1;36m"
	const headerColor = "\033[1;36m"
	const nameColor = "\033[1;37m"
	const tagsColor = "\033[0;33m"
	const descColor = "\033[0;37m"
	const indexColor = "\033[1;33m"
	const categoryColor = "\033[1;33m"

	// 打印标题
	titleBorder := strings.Repeat("─", TableTotalWidth)
	title := fmt.Sprintf("标签 '%s' 下的网页工具", tag)
	titleLen := runeWidth(title)
	titlePadding := (TableTotalWidth - titleLen) / 2
	leftPadding := strings.Repeat(" ", titlePadding)
	rightPadding := strings.Repeat(" ", TableTotalWidth-titleLen-titlePadding)

	fmt.Printf("%s┌%s┐\033[0m\n", borderColor, titleBorder)
	fmt.Printf("%s│%s%s%s│\033[0m\n", borderColor, leftPadding, title, rightPadding)
	fmt.Printf("%s└%s┘\033[0m\n", borderColor, titleBorder)

	// 按分类对工具进行分组
	categoryMap := make(map[string][]models.WebTool)
	var categories []string

	for _, tool := range results {
		if _, exists := categoryMap[tool.Category]; !exists {
			categories = append(categories, tool.Category)
		}
		categoryMap[tool.Category] = append(categoryMap[tool.Category], tool)
	}

	// 对分类进行排序
	sort.Strings(categories)

	// 创建序号到工具的映射
	indexMap := make(map[int]models.WebTool)
	currentIndex := 1

	// 按分类输出工具
	for _, category := range categories {
		tools := categoryMap[category]

		// 对每个分类中的工具按名称排序
		sort.Slice(tools, func(i, j int) bool {
			return tools[i].Name < tools[j].Name
		})

		// 创建分类表格
		categoryTable := Table{
			CategoryTitle: category,
			BorderColor:   categoryColor,
			HeaderColor:   headerColor,
			CellColor:     nameColor,
			Columns: []TableColumn{
				{Title: "序号", Width: 8, Color: indexColor},
				{Title: "名称", Width: NameColWidth, Color: nameColor},
				{Title: "描述", Width: DescColWidth + TagsColWidth - 8, Color: descColor},
			},
		}

		// 添加数据行
		for _, tool := range tools {
			desc := cleanString(tool.Description)
			desc = truncateString(desc, DescColWidth+TagsColWidth-8)

			name := cleanString(tool.Name)
			name = truncateString(name, NameColWidth)

			index := fmt.Sprintf("[%d]", currentIndex)
			indexMap[currentIndex] = tool
			currentIndex++

			categoryTable.Rows = append(categoryTable.Rows, TableRow{
				Columns: []string{index, name, desc},
			})
		}

		// 打印分类表格
		printTable(categoryTable)
	}

	// 输出工具总数
	fmt.Printf("\n%s总计: %d 个工具, %d 个分类%s\n", borderColor, len(results), len(categories), "\033[0m")

	for {
		fmt.Print("\n请选择要打开的工具 (输入序号, 输入q退出): ")
		var input string
		fmt.Scanln(&input)

		input = strings.ToLower(input)
		if input == "q" || input == "quit" || input == "exit" {
			fmt.Println("已退出")
			return
		}

		// 尝试将输入转换为数字
		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("无效输入，请输入有效的数字或q退出")
			continue
		}

		if choice > 0 && choice < currentIndex {
			tool := indexMap[choice]
			fmt.Printf("正在打开: %s\n", tool.Name)
			if err := launcher.LaunchWebTool(tool); err != nil {
				fmt.Printf("打开失败: %v\n", err)
			}

			// 工具运行结束后询问用户是否还需要启动其他工具
			fmt.Print("\n是否继续选择其他工具? (y/n): ")
			var continueChoice string
			fmt.Scanln(&continueChoice)

			if strings.ToLower(continueChoice) != "y" && strings.ToLower(continueChoice) != "yes" {
				return
			}
		} else {
			fmt.Println("无效的选择")
		}
	}
}

func displayAllWebToolTags() {
	tagCounts := search.GetAllWebToolTagsWithCount(cfg.WebTools.Tools)

	// 设置颜色
	const borderColor = "\033[1;36m"
	const headerColor = "\033[1;36m"
	const nameColor = "\033[1;37m"
	const countColor = "\033[0;33m"

	// 打印标题
	titleBorder := strings.Repeat("─", TableTotalWidth)
	title := "网页工具标签列表"
	titleLen := runeWidth(title)
	titlePadding := (TableTotalWidth - titleLen) / 2
	leftPadding := strings.Repeat(" ", titlePadding)
	rightPadding := strings.Repeat(" ", TableTotalWidth-titleLen-titlePadding)

	fmt.Printf("%s┌%s┐\033[0m\n", borderColor, titleBorder)
	fmt.Printf("%s│%s%s%s│\033[0m\n", borderColor, leftPadding, title, rightPadding)
	fmt.Printf("%s└%s┘\033[0m\n", borderColor, titleBorder)

	// 横排布局，每行显示4个标签
	const tagsPerRow = 4
	const tagWidth = TableTotalWidth/tagsPerRow - 2 // 减去分隔符宽度

	// 创建表格
	fmt.Printf("%s┌", borderColor)
	for i := 0; i < tagsPerRow; i++ {
		fmt.Print(strings.Repeat("─", tagWidth))
		if i < tagsPerRow-1 {
			fmt.Print("┬")
		}
	}
	fmt.Print("┐\033[0m\n")

	// 打印行
	for i := 0; i < len(tagCounts); i += tagsPerRow {
		fmt.Printf("%s│", borderColor)
		for j := 0; j < tagsPerRow; j++ {
			if i+j < len(tagCounts) {
				tag := tagCounts[i+j].Tag
				count := tagCounts[i+j].Count
				displayText := fmt.Sprintf("%s (%d)", tag, count)
				paddedText := padString(displayText, tagWidth-2)
				fmt.Printf(" %s%s %s│", nameColor, paddedText, borderColor)
			} else {
				fmt.Printf(" %s%s %s│", nameColor, padString("", tagWidth-2), borderColor)
			}
		}
		fmt.Print("\033[0m\n")

		// 打印行间分隔符，除了最后一行
		if i+tagsPerRow < len(tagCounts) {
			fmt.Printf("%s├", borderColor)
			for j := 0; j < tagsPerRow; j++ {
				fmt.Print(strings.Repeat("─", tagWidth))
				if j < tagsPerRow-1 {
					fmt.Print("┼")
				}
			}
			fmt.Print("┤\033[0m\n")
		}
	}

	// 打印底部边框
	fmt.Printf("%s└", borderColor)
	for i := 0; i < tagsPerRow; i++ {
		fmt.Print(strings.Repeat("─", tagWidth))
		if i < tagsPerRow-1 {
			fmt.Print("┴")
		}
	}
	fmt.Print("┘\033[0m\n")

	// 输出标签总数
	fmt.Printf("\n%s总计: %d 个标签%s\n", borderColor, len(tagCounts), "\033[0m")
}

func displayNotes() {
	if len(cfg.WebNotes.Notes) == 0 {
		fmt.Println("没有可用的笔记")
		return
	}

	// 设置颜色
	const borderColor = "\033[1;34m"
	const headerColor = "\033[1;34m"
	const titleColor = "\033[1;37m"
	const tagsColor = "\033[0;33m"
	const sourceColor = "\033[0;37m"
	const categoryColor = "\033[1;33m"
	const cellBorderColor = "\033[0;34m"

	// 按工具对笔记进行分组
	toolMap := make(map[string][]models.Note)
	var tools []string

	for _, note := range cfg.WebNotes.Notes {
		toolName := note.Tool
		if toolName == "" {
			toolName = "未分类"
		}

		if _, exists := toolMap[toolName]; !exists {
			tools = append(tools, toolName)
		}
		toolMap[toolName] = append(toolMap[toolName], note)
	}

	// 对工具进行排序
	sort.Strings(tools)

	// 打印标题
	titleBorder := strings.Repeat("─", TableTotalWidth)
	title := "网页笔记列表"
	titleLen := runeWidth(title)
	titlePadding := (TableTotalWidth - titleLen) / 2
	leftPadding := strings.Repeat(" ", titlePadding)
	rightPadding := strings.Repeat(" ", TableTotalWidth-titleLen-titlePadding)

	fmt.Printf("%s┌%s┐\033[0m\n", borderColor, titleBorder)
	fmt.Printf("%s│%s%s%s│\033[0m\n", borderColor, leftPadding, title, rightPadding)
	fmt.Printf("%s└%s┘\033[0m\n", borderColor, titleBorder)

	// 按工具输出笔记
	for _, toolName := range tools {
		notes := toolMap[toolName]

		// 对每个工具中的笔记按标题排序
		sort.Slice(notes, func(i, j int) bool {
			return notes[i].Title < notes[j].Title
		})

		// 创建分类表格
		categoryTable := Table{
			CategoryTitle: toolName,
			BorderColor:   categoryColor,
			HeaderColor:   headerColor,
			CellColor:     titleColor,
			Columns: []TableColumn{
				{Title: "标题", Width: TitleColWidth, Color: titleColor},
				{Title: "标签", Width: TagsColWidth, Color: tagsColor},
				{Title: "来源", Width: SourceColWidth, Color: sourceColor},
			},
		}

		// 添加数据行
		for _, note := range notes {
			tags := strings.Join(note.Tags, ", ")
			tags = truncateString(cleanString(tags), TagsColWidth)

			source := cleanString(note.Source)
			source = truncateString(source, SourceColWidth)

			title := cleanString(note.Title)
			title = truncateString(title, TitleColWidth)

			categoryTable.Rows = append(categoryTable.Rows, TableRow{
				Columns: []string{title, tags, source},
			})
		}

		// 打印分类表格
		printTable(categoryTable)
	}

	// 输出笔记总数
	fmt.Printf("\n%s总计: %d 个笔记, %d 个分类\033[0m\n", borderColor, len(cfg.WebNotes.Notes), len(tools))
}

// 处理字符串，替换换行符为空格
func cleanString(s string) string {
	// 替换所有换行符为空格
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")

	// 替换多个连续空格为单个空格
	for strings.Contains(s, "  ") {
		s = strings.ReplaceAll(s, "  ", " ")
	}

	return strings.TrimSpace(s)
}

// 打印表格
func printTable(table Table) {
	// 如果没有行数据且没有分类标题，说明是主标题表格
	isMainTitle := len(table.Rows) == 0 && table.CategoryTitle == ""

	// 打印表格标题（仅当有标题时）
	if table.Title != "" {
		titleBorder := strings.Repeat("─", TableTotalWidth-2) // 减去左右边框
		fmt.Printf("%s┌%s┐\033[0m\n", table.TitleColor, titleBorder)

		// 计算标题居中位置
		titleLen := runeWidth(table.Title)
		titlePadding := (TableTotalWidth - 2 - titleLen) / 2 // 减去左右边框
		if titlePadding < 0 {
			titlePadding = 0
		}

		leftPadding := strings.Repeat(" ", titlePadding)
		rightPadding := strings.Repeat(" ", TableTotalWidth-2-titleLen-titlePadding) // 减去左右边框

		fmt.Printf("%s│%s%s%s│\033[0m\n", table.TitleColor, leftPadding, table.Title, rightPadding)

		fmt.Printf("%s└%s┘\033[0m\n", table.TitleColor, titleBorder)
	}

	// 如果是主标题表格，不打印内容
	if isMainTitle {
		return
	}

	// 如果有分类标题，打印分类标题
	if table.CategoryTitle != "" {
		fmt.Printf("\n%s【%s】\033[0m\n", table.BorderColor, table.CategoryTitle)
	}

	// 打印表头边框
	printTableBorder(table, "┌", "┬", "┐")

	// 打印表头
	printTableHeader(table)

	// 打印表头和数据之间的分隔线
	printTableBorder(table, "├", "┼", "┤")

	// 打印数据行
	for _, row := range table.Rows {
		printTableRow(table, row)
	}

	// 打印底部边框
	printTableBorder(table, "└", "┴", "┘")
}

// 打印表格边框
func printTableBorder(table Table, left, middle, right string) {
	fmt.Print(table.BorderColor)
	fmt.Print(left)

	// 如果没有列或者是主标题表格
	if len(table.Columns) == 0 {
		fmt.Print(strings.Repeat("─", TableTotalWidth-2)) // -2 是因为左右边框各占一个字符
	} else {
		// 计算每列宽度和总宽度
		totalWidth := -1 // 初始值为-1是因为最后一列后没有分隔符
		for i, col := range table.Columns {
			totalWidth += col.Width + 2 // 列宽 + 两侧空格
			if i < len(table.Columns)-1 {
				totalWidth += 1 // 加上分隔符宽度
			}
		}

		// 打印每列边框
		for i, col := range table.Columns {
			width := col.Width + 2 // 列宽 + 两侧空格

			// 如果是最后一列，并且总宽度不足，则调整最后一列宽度
			if i == len(table.Columns)-1 && totalWidth < TableTotalWidth-2 {
				width += (TableTotalWidth - 2) - totalWidth
			}

			fmt.Print(strings.Repeat("─", width))
			if i < len(table.Columns)-1 {
				fmt.Print(middle)
			}
		}
	}

	fmt.Print(right)
	fmt.Print("\033[0m\n")
}

// 打印表头
func printTableHeader(table Table) {
	// 如果没有列，直接返回
	if len(table.Columns) == 0 {
		return
	}

	// 计算总宽度
	totalWidth := -1 // 初始值为-1是因为最后一列后没有分隔符
	for i, col := range table.Columns {
		totalWidth += col.Width + 2 // 列宽 + 两侧空格
		if i < len(table.Columns)-1 {
			totalWidth += 1 // 加上分隔符宽度
		}
	}

	fmt.Print(table.BorderColor + "│\033[0m")

	for i, col := range table.Columns {
		width := col.Width

		// 如果是最后一列，并且总宽度不足，则调整最后一列宽度
		if i == len(table.Columns)-1 && totalWidth < TableTotalWidth-2 {
			width += (TableTotalWidth - 2) - totalWidth
		}

		paddedTitle := padString(col.Title, width)
		fmt.Printf(" %s%s\033[0m %s│\033[0m", table.HeaderColor, paddedTitle, table.BorderColor)
	}

	fmt.Println()
}

// 打印表格行
func printTableRow(table Table, row TableRow) {
	// 如果没有列，直接返回
	if len(table.Columns) == 0 {
		return
	}

	// 计算总宽度
	totalWidth := -1 // 初始值为-1是因为最后一列后没有分隔符
	for i, col := range table.Columns {
		totalWidth += col.Width + 2 // 列宽 + 两侧空格
		if i < len(table.Columns)-1 {
			totalWidth += 1 // 加上分隔符宽度
		}
	}

	fmt.Print(table.BorderColor + "│\033[0m")

	for i, col := range row.Columns {
		if i < len(table.Columns) {
			width := table.Columns[i].Width

			// 如果是最后一列，并且总宽度不足，则调整最后一列宽度
			if i == len(table.Columns)-1 && totalWidth < TableTotalWidth-2 {
				width += (TableTotalWidth - 2) - totalWidth
			}

			// 确保列内容不超过列宽
			if runeWidth(col) > width {
				col = truncateString(col, width)
			}

			// 确保填充到正确的宽度，考虑中文字符宽度
			paddedCol := padString(col, width)
			fmt.Printf(" %s%s\033[0m %s│\033[0m", table.Columns[i].Color, paddedCol, table.BorderColor)
		}
	}

	fmt.Println()
}

// 列出所有离线工具
func listOfflineTools() {
	if len(cfg.OfflineTools.Tools) == 0 {
		fmt.Println("没有可用的离线工具")
		return
	}

	// 设置颜色
	const borderColor = "\033[1;36m"
	const headerColor = "\033[1;36m"
	const nameColor = "\033[1;37m"
	const tagsColor = "\033[0;33m"
	const descColor = "\033[0;37m"
	const categoryColor = "\033[1;33m"
	const cellBorderColor = "\033[0;32m"

	// 按分类对工具进行分组
	categoryMap := make(map[string][]models.OfflineTool)
	var categories []string

	for _, tool := range cfg.OfflineTools.Tools {
		if _, exists := categoryMap[tool.Category]; !exists {
			categories = append(categories, tool.Category)
		}
		categoryMap[tool.Category] = append(categoryMap[tool.Category], tool)
	}

	// 对分类进行排序
	sort.Strings(categories)

	// 打印标题
	titleBorder := strings.Repeat("─", TableTotalWidth)
	title := "离线工具列表"
	titleLen := runeWidth(title)
	titlePadding := (TableTotalWidth - titleLen) / 2
	leftPadding := strings.Repeat(" ", titlePadding)
	rightPadding := strings.Repeat(" ", TableTotalWidth-titleLen-titlePadding)

	fmt.Printf("%s┌%s┐\033[0m\n", borderColor, titleBorder)
	fmt.Printf("%s│%s%s%s│\033[0m\n", borderColor, leftPadding, title, rightPadding)
	fmt.Printf("%s└%s┘\033[0m\n", borderColor, titleBorder)

	// 按分类输出工具
	for _, category := range categories {
		tools := categoryMap[category]

		// 对每个分类中的工具按名称排序
		sort.Slice(tools, func(i, j int) bool {
			return tools[i].Name < tools[j].Name
		})

		// 创建分类表格
		categoryTable := Table{
			CategoryTitle: category,
			BorderColor:   categoryColor,
			HeaderColor:   headerColor,
			CellColor:     nameColor,
			Columns: []TableColumn{
				{Title: "名称", Width: NameColWidth, Color: nameColor},
				{Title: "标签", Width: TagsColWidth, Color: tagsColor},
				{Title: "描述", Width: DescColWidth, Color: descColor},
			},
		}

		// 添加数据行
		for _, tool := range tools {
			tags := strings.Join(tool.Tags, ", ")
			tags = truncateString(cleanString(tags), TagsColWidth)

			desc := cleanString(tool.Description)
			desc = truncateString(desc, DescColWidth)

			name := cleanString(tool.Name)
			name = truncateString(name, NameColWidth)

			categoryTable.Rows = append(categoryTable.Rows, TableRow{
				Columns: []string{name, tags, desc},
			})
		}

		// 打印分类表格
		printTable(categoryTable)
	}

	// 输出工具总数
	fmt.Printf("\n%s总计: %d 个工具, %d 个分类\033[0m\n", borderColor, len(cfg.OfflineTools.Tools), len(categories))
}

// 列出所有网页工具
func listWebTools() {
	if len(cfg.WebTools.Tools) == 0 {
		fmt.Println("没有可用的网页工具")
		return
	}

	// 设置颜色
	const borderColor = "\033[1;35m"
	const headerColor = "\033[1;35m"
	const nameColor = "\033[1;37m"
	const tagsColor = "\033[0;33m"
	const descColor = "\033[0;37m"
	const categoryColor = "\033[1;33m"
	const cellBorderColor = "\033[0;35m"

	// 按分类对工具进行分组
	categoryMap := make(map[string][]models.WebTool)
	var categories []string

	for _, tool := range cfg.WebTools.Tools {
		if _, exists := categoryMap[tool.Category]; !exists {
			categories = append(categories, tool.Category)
		}
		categoryMap[tool.Category] = append(categoryMap[tool.Category], tool)
	}

	// 对分类进行排序
	sort.Strings(categories)

	// 打印标题
	titleBorder := strings.Repeat("─", TableTotalWidth)
	title := "网页工具列表"
	titleLen := runeWidth(title)
	titlePadding := (TableTotalWidth - titleLen) / 2
	leftPadding := strings.Repeat(" ", titlePadding)
	rightPadding := strings.Repeat(" ", TableTotalWidth-titleLen-titlePadding)

	fmt.Printf("%s┌%s┐\033[0m\n", borderColor, titleBorder)
	fmt.Printf("%s│%s%s%s│\033[0m\n", borderColor, leftPadding, title, rightPadding)
	fmt.Printf("%s└%s┘\033[0m\n", borderColor, titleBorder)

	// 按分类输出工具
	for _, category := range categories {
		tools := categoryMap[category]

		// 对每个分类中的工具按名称排序
		sort.Slice(tools, func(i, j int) bool {
			return tools[i].Name < tools[j].Name
		})

		// 创建分类表格
		categoryTable := Table{
			CategoryTitle: category,
			BorderColor:   categoryColor,
			HeaderColor:   headerColor,
			CellColor:     nameColor,
			Columns: []TableColumn{
				{Title: "名称", Width: NameColWidth, Color: nameColor},
				{Title: "标签", Width: TagsColWidth, Color: tagsColor},
				{Title: "描述", Width: DescColWidth, Color: descColor},
			},
		}

		// 添加数据行
		for _, tool := range tools {
			tags := strings.Join(tool.Tags, ", ")
			tags = truncateString(cleanString(tags), TagsColWidth)

			desc := cleanString(tool.Description)
			desc = truncateString(desc, DescColWidth)

			name := cleanString(tool.Name)
			name = truncateString(name, NameColWidth)

			categoryTable.Rows = append(categoryTable.Rows, TableRow{
				Columns: []string{name, tags, desc},
			})
		}

		// 打印分类表格
		printTable(categoryTable)
	}

	// 输出工具总数
	fmt.Printf("\n%s总计: %d 个工具, %d 个分类\033[0m\n", borderColor, len(cfg.WebTools.Tools), len(categories))
}

func displayHelp() {
	fmt.Println("Matu7 工具启动器 - 帮助")
	fmt.Println("\n配置命令:")
	fmt.Println("  --add-path <路径>   添加配置文件夹路径")

	fmt.Println("\n功能命令:")
	fmt.Println("  -t [名称]          不加参数显示所有离线工具，加参数搜索并启动离线工具")
	fmt.Println("  -tm <标签>         根据标签搜索离线工具并显示")
	fmt.Println("  -w [名称]          不加参数显示所有网页工具，加参数搜索并打开网页工具")
	fmt.Println("  -wm <标签>         根据标签搜索网页工具并显示")
	fmt.Println("  -n [关键词]        不加参数显示所有笔记，加参数搜索网页笔记")
	fmt.Println("  -nm <标签>         根据标签搜索网页笔记并显示")
	fmt.Println("  help               显示帮助信息")

	fmt.Println("\n示例:")
	fmt.Println("  start --add-path /path/to/config    添加配置路径")
	fmt.Println("  start -t                           显示所有离线工具")
	fmt.Println("  start -t sqlmap                    启动sqlmap工具")
	fmt.Println("  start -t sqlm,数据库                搜索名称包含sqlm且标签或描述包含数据库的工具")
	fmt.Println("  start -tm framework                显示所有标签为framework的工具")
	fmt.Println("  start -w                           显示所有网页工具")
	fmt.Println("  start -n                           显示所有笔记")
	fmt.Println("  start -n Resin                     搜索标题或标签包含Resin的笔记")
	fmt.Println("  start -n Resin,攻击                 搜索标题或标签包含Resin且包含攻击的笔记")
	fmt.Println("  start -nm CauchoResin              显示所有标签为CauchoResin的笔记")
}

// runeWidth 返回字符串的显示宽度（考虑中文等宽字符）
func runeWidth(s string) int {
	width := 0
	for _, r := range s {
		if unicode.Is(unicode.Han, r) || unicode.Is(unicode.Katakana, r) ||
			unicode.Is(unicode.Hiragana, r) || r > 0x3000 {
			width += 2 // 中文和其他全角字符占2个宽度
		} else {
			width += 1 // 英文和其他ASCII字符占1个宽度
		}
	}
	return width
}

// 截断字符串到指定显示宽度
func truncateString(s string, maxWidth int) string {
	if runeWidth(s) <= maxWidth {
		return s
	}

	var result []rune
	width := 0
	for _, r := range s {
		charWidth := 1
		if unicode.Is(unicode.Han, r) || unicode.Is(unicode.Katakana, r) ||
			unicode.Is(unicode.Hiragana, r) || r > 0x3000 {
			charWidth = 2
		}

		if width+charWidth > maxWidth-1 { // 留1个位置给省略号
			break
		}

		width += charWidth
		result = append(result, r)
	}

	return string(result) + "…"
}

// 填充字符串到指定宽度
func padString(s string, width int) string {
	currentWidth := runeWidth(s)
	if currentWidth >= width {
		return s
	}

	padding := width - currentWidth
	return s + strings.Repeat(" ", padding)
}

// handleNoteSearch 处理笔记搜索
func handleNoteSearch(query string) {
	if query == "" {
		// 如果没有提供查询参数，显示所有笔记
		displayNotes()
		return
	}

	// 处理查询条件
	results := search.FuzzySearchNotes(cfg.WebNotes.Notes, query)

	if len(results) == 0 {
		fmt.Println("未找到匹配的笔记")
		// 提供一些建议的标签
		fmt.Println("您可以尝试以下热门标签:")
		displayTopNoteTags(cfg.WebNotes.Notes, 5)
		return
	}

	// 设置颜色
	const borderColor = "\033[1;34m"
	const headerColor = "\033[1;34m"
	const titleColor = "\033[1;37m"
	const tagsColor = "\033[0;33m"
	const sourceColor = "\033[0;37m"
	const categoryColor = "\033[1;33m"

	// 打印标题
	titleBorder := strings.Repeat("─", TableTotalWidth)
	title := "笔记搜索结果"
	titleLen := runeWidth(title)
	titlePadding := (TableTotalWidth - titleLen) / 2
	leftPadding := strings.Repeat(" ", titlePadding)
	rightPadding := strings.Repeat(" ", TableTotalWidth-titleLen-titlePadding)

	fmt.Printf("%s┌%s┐\033[0m\n", borderColor, titleBorder)
	fmt.Printf("%s│%s%s%s│\033[0m\n", borderColor, leftPadding, title, rightPadding)
	fmt.Printf("%s└%s┘\033[0m\n", borderColor, titleBorder)

	// 按工具对笔记进行分组
	toolMap := make(map[string][]models.Note)
	var tools []string

	for _, note := range results {
		toolName := note.Tool
		if toolName == "" {
			toolName = "目录类"
		}

		if _, exists := toolMap[toolName]; !exists {
			tools = append(tools, toolName)
		}
		toolMap[toolName] = append(toolMap[toolName], note)
	}

	// 对工具进行排序
	sort.Strings(tools)

	// 按工具输出笔记
	for _, toolName := range tools {
		notes := toolMap[toolName]

		// 对每个工具中的笔记按标题排序
		sort.Slice(notes, func(i, j int) bool {
			return notes[i].Title < notes[j].Title
		})

		// 创建分类表格
		categoryTable := Table{
			CategoryTitle: toolName,
			BorderColor:   categoryColor,
			HeaderColor:   headerColor,
			CellColor:     titleColor,
			Columns: []TableColumn{
				{Title: "标题", Width: TitleColWidth, Color: titleColor},
				{Title: "标签", Width: TagsColWidth, Color: tagsColor},
				{Title: "来源", Width: SourceColWidth, Color: sourceColor},
			},
		}

		// 添加数据行
		for _, note := range notes {
			tags := strings.Join(note.Tags, ", ")
			tags = truncateString(cleanString(tags), TagsColWidth)

			source := cleanString(note.Source)
			source = truncateString(source, SourceColWidth)

			title := cleanString(note.Title)
			title = truncateString(title, TitleColWidth)

			categoryTable.Rows = append(categoryTable.Rows, TableRow{
				Columns: []string{title, tags, source},
			})
		}

		// 打印分类表格
		printTable(categoryTable)
	}

	// 输出笔记总数
	fmt.Printf("\n%s总计: %d 个笔记, %d 个分类\033[0m\n", borderColor, len(results), len(tools))
}

// displayTopNoteTags 显示最常用的笔记标签
func displayTopNoteTags(notes []models.Note, count int) {
	tagCounts := make(map[string]int)

	// 统计标签出现次数
	for _, note := range notes {
		for _, tag := range note.Tags {
			tagCounts[tag]++
		}
	}

	// 转换为切片并排序
	type tagCount struct {
		Tag   string
		Count int
	}

	var tagList []tagCount
	for tag, count := range tagCounts {
		tagList = append(tagList, tagCount{Tag: tag, Count: count})
	}

	// 按出现次数排序
	sort.Slice(tagList, func(i, j int) bool {
		return tagList[i].Count > tagList[j].Count
	})

	// 显示前N个标签
	fmt.Print("  ")
	for i := 0; i < count && i < len(tagList); i++ {
		fmt.Printf("\033[0;33m%s\033[0m ", tagList[i].Tag)
	}
	fmt.Println("\n使用 -nm <标签> 命令可按标签搜索笔记")
}

// displayAllNoteTags 显示所有笔记标签
func displayAllNoteTags() {
	tagCounts := search.GetAllNotesTagsWithCount(cfg.WebNotes.Notes)

	// 设置颜色
	const borderColor = "\033[1;34m"
	const headerColor = "\033[1;34m"
	const nameColor = "\033[1;37m"
	const countColor = "\033[0;33m"

	// 打印标题
	titleBorder := strings.Repeat("─", TableTotalWidth)
	title := "笔记标签列表"
	titleLen := runeWidth(title)
	titlePadding := (TableTotalWidth - titleLen) / 2
	leftPadding := strings.Repeat(" ", titlePadding)
	rightPadding := strings.Repeat(" ", TableTotalWidth-titleLen-titlePadding)

	fmt.Printf("%s┌%s┐\033[0m\n", borderColor, titleBorder)
	fmt.Printf("%s│%s%s%s│\033[0m\n", borderColor, leftPadding, title, rightPadding)
	fmt.Printf("%s└%s┘\033[0m\n", borderColor, titleBorder)

	// 横排布局，每行显示4个标签
	const tagsPerRow = 4
	const tagWidth = TableTotalWidth/tagsPerRow - 2 // 减去分隔符宽度

	// 创建表格
	fmt.Printf("%s┌", borderColor)
	for i := 0; i < tagsPerRow; i++ {
		fmt.Print(strings.Repeat("─", tagWidth))
		if i < tagsPerRow-1 {
			fmt.Print("┬")
		}
	}
	fmt.Print("┐\033[0m\n")

	// 打印行
	for i := 0; i < len(tagCounts); i += tagsPerRow {
		fmt.Printf("%s│", borderColor)
		for j := 0; j < tagsPerRow; j++ {
			if i+j < len(tagCounts) {
				tag := tagCounts[i+j].Tag
				count := tagCounts[i+j].Count
				displayText := fmt.Sprintf("%s (%d)", tag, count)
				paddedText := padString(displayText, tagWidth-2)
				fmt.Printf(" %s%s %s│", nameColor, paddedText, borderColor)
			} else {
				fmt.Printf(" %s%s %s│", nameColor, padString("", tagWidth-2), borderColor)
			}
		}
		fmt.Print("\033[0m\n")

		// 打印行间分隔符，除了最后一行
		if i+tagsPerRow < len(tagCounts) {
			fmt.Printf("%s├", borderColor)
			for j := 0; j < tagsPerRow; j++ {
				fmt.Print(strings.Repeat("─", tagWidth))
				if j < tagsPerRow-1 {
					fmt.Print("┼")
				}
			}
			fmt.Print("┤\033[0m\n")
		}
	}

	// 打印底部边框
	fmt.Printf("%s└", borderColor)
	for i := 0; i < tagsPerRow; i++ {
		fmt.Print(strings.Repeat("─", tagWidth))
		if i < tagsPerRow-1 {
			fmt.Print("┴")
		}
	}
	fmt.Print("┘\033[0m\n")

	// 输出标签总数
	fmt.Printf("\n%s总计: %d 个标签%s\n", borderColor, len(tagCounts), "\033[0m")
}

// handleNoteByTag 根据标签搜索笔记
func handleNoteByTag(tag string) {
	if tag == "" {
		displayAllNoteTags()
		return
	}

	// 从标签中模糊搜索，而不是完全匹配
	results := search.SearchNotesByTag(cfg.WebNotes.Notes, tag)

	if len(results) == 0 {
		fmt.Println("未找到匹配标签的笔记")
		return
	}

	// 设置颜色
	const borderColor = "\033[1;34m"
	const headerColor = "\033[1;34m"
	const titleColor = "\033[1;37m"
	const tagsColor = "\033[0;33m"
	const sourceColor = "\033[0;37m"
	const categoryColor = "\033[1;33m"

	// 打印标题
	titleBorder := strings.Repeat("─", TableTotalWidth)
	title := fmt.Sprintf("标签 '%s' 下的笔记", tag)
	titleLen := runeWidth(title)
	titlePadding := (TableTotalWidth - titleLen) / 2
	leftPadding := strings.Repeat(" ", titlePadding)
	rightPadding := strings.Repeat(" ", TableTotalWidth-titleLen-titlePadding)

	fmt.Printf("%s┌%s┐\033[0m\n", borderColor, titleBorder)
	fmt.Printf("%s│%s%s%s│\033[0m\n", borderColor, leftPadding, title, rightPadding)
	fmt.Printf("%s└%s┘\033[0m\n", borderColor, titleBorder)

	// 按工具对笔记进行分组
	toolMap := make(map[string][]models.Note)
	var tools []string

	for _, note := range results {
		toolName := note.Tool
		if toolName == "" {
			toolName = "目录类"
		}

		if _, exists := toolMap[toolName]; !exists {
			tools = append(tools, toolName)
		}
		toolMap[toolName] = append(toolMap[toolName], note)
	}

	// 对工具进行排序
	sort.Strings(tools)

	// 按工具输出笔记
	for _, toolName := range tools {
		notes := toolMap[toolName]

		// 对每个工具中的笔记按标题排序
		sort.Slice(notes, func(i, j int) bool {
			return notes[i].Title < notes[j].Title
		})

		// 创建分类表格
		categoryTable := Table{
			CategoryTitle: toolName,
			BorderColor:   categoryColor,
			HeaderColor:   headerColor,
			CellColor:     titleColor,
			Columns: []TableColumn{
				{Title: "标题", Width: TitleColWidth, Color: titleColor},
				{Title: "来源", Width: SourceColWidth + TagsColWidth, Color: sourceColor},
			},
		}

		// 添加数据行
		for _, note := range notes {
			source := cleanString(note.Source)
			source = truncateString(source, SourceColWidth+TagsColWidth)

			title := cleanString(note.Title)
			title = truncateString(title, TitleColWidth)

			categoryTable.Rows = append(categoryTable.Rows, TableRow{
				Columns: []string{title, source},
			})
		}

		// 打印分类表格
		printTable(categoryTable)
	}

	// 输出笔记总数
	fmt.Printf("\n%s总计: %d 个笔记, %d 个分类\033[0m\n", borderColor, len(results), len(tools))
}
