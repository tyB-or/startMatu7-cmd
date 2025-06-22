# Matu7 工具启动器

Matu7是一个强大的命令行工具启动器，专为安全工具管理和快速访问设计。它能够帮助用户组织、搜索和启动各种离线工具和网页资源，同时提供标签管理和笔记功能。

## 功能特性

- **离线工具管理**：快速搜索和启动本地安全工具
- **网页工具访问**：打开常用的网页工具
- **笔记管理**：检索相关的笔记
- **标签系统**：通过标签对工具和笔记进行分类和搜索
- **模糊搜索**：支持按名称、描述和标签进行模糊搜索
- **多条件搜索**：支持逗号分隔的多关键词搜索

## 安装方法

### 前置条件

- Go 1.16 或更高版本

### 编译安装

1.使用matu7，并维护好自己日常使用和收集的工具

matu7项目：https://github.com/tyB-or/matu7

2. 编译本项目：

```bash
go build -o start cmd/start/main.go
```

3. 将编译好的二进制文件移动到系统路径（可选）：

```bash
sudo mv start /usr/local/bin/
```

## 配置

首次使用前，需要设置配置文件路径：（matu7工具的config目录）

```bash
./start --add-path /path/to/config  
```

配置文件夹应包含以下JSON文件：

- `offline_tools_xxx.json`：离线工具配置
- `web_tools.json`：网页工具配置
- `web_notes.json`：网页笔记配置

## 使用说明


https://github.com/user-attachments/assets/60031612-8462-4afa-bf98-d3c00eac01f5


### 基本命令

- **离线工具**：
  - `-t [关键词]`：不加参数显示所有离线工具，加参数搜索并启动离线工具(模糊搜索，不区分大小写),搜索逻辑：从名称、标签、描述中查询
  - `-tm <标签>`：根据标签搜索离线工具，支持模糊搜索，不区分大小写

- **网页工具**：
  - `-w [关键词]`：不加参数显示所有网页工具，加参数搜索并打开网页工具(模糊搜索，不区分大小写),搜索逻辑：从名称、标签、描述中查询
  - `-wm <标签>`：根据标签搜索网页工具,支持模糊搜索，不区分大小写

- **笔记管理**：
  - `-n [关键词]`：不加参数显示所有笔记，加参数搜索网页笔记，(模糊搜索，不区分大小写),搜索逻辑：从名称、标签中查询
  - `-nm <标签>`：根据标签搜索网页笔记,支持模糊搜索，不区分大小写

- **帮助**：
  - `help`：显示帮助信息

- **启动说明**：
  - `离线工具`：优先根据配置文件中的命令进行启动，没有的话检测常见的可执行程序后缀并启动，还没有就进入工具目录的终端路径
  - `网页工具`：调用默认浏览器打开指定网址

### 搜索语法

- 单个关键词：`-t sqlmap`
- 多条件搜索：`-t sqlm,数据库`（搜索包含sqlm和数据库关键词的工具）
- 标签搜索：`-tm framework`（显示所有标签为framework的工具）


```bash
./start
```



## 使用示例

### 显示所有离线工具

```bash
./start -t
```

### 搜索并启动特定工具

```bash
./start -t sqlmap
```

### 按标签搜索工具

```bash
./start -tm framework
```

### 显示所有网页工具

```bash
./start -w
```

### 搜索笔记

```bash
./start -n Resin
```

### 多条件搜索笔记

```bash
./start -n Resin,攻击
```
### 使用截图
![t](https://github.com/user-attachments/assets/7562cadf-c223-496a-84da-67f7f46e48ca)
![w](https://github.com/user-attachments/assets/4c1b6aae-67a5-4e96-a09f-42e43b7e5d4c)
![n](https://github.com/user-attachments/assets/72240468-3d19-4e83-bd12-fb027c89ff06)
![t-jar](https://github.com/user-attachments/assets/76674273-099f-47ce-b7fb-f1c41896c92a)
![tm](https://github.com/user-attachments/assets/67473abc-9756-438b-8050-6008a80c34a0)


## 配置文件格式

### offline_tools.json

```json
{
  "scan_path": "/path/to/tools",
  "auto_refresh": false,
  "tools": [
    {
      "id": "1",
      "name": "工具名称",
      "category": "分类",
      "path": "/path/to/tool",
      "description": "工具描述",
      "tags": ["tag1", "tag2"],
      "command": "可选的启动命令"
    }
  ]
}
```

### web_tools.json

```json
{
  "tools": [
    {
      "id": "1",
      "name": "网页工具名称",
      "url": "https://example.com",
      "description": "网页工具描述",
      "category": "分类",
      "tags": ["tag1", "tag2"]
    }
  ]
}
```

### web_notes.json

```json
{
  "notes": [
    {
      "id": "1",
      "title": "笔记标题",
      "url": "https://example.com",
      "source": "来源",
      "tool": "相关工具",
      "tags": ["tag1", "tag2"],
      "note": "笔记内容"
    }
  ]
}
```

## 项目结构

```
matu7/
  ├── cmd/
  │   └── start/
  │       └── main.go        # 主程序入口
  ├── internal/              # 内部包
  │   ├── config/            # 配置管理
  │   ├── launcher/          # 工具启动逻辑
  │   └── search/            # 搜索功能
  ├── pkg/                   # 公共包
  │   └── models/            # 数据模型
  └── config/                # 配置文件示例
```
