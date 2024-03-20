# gutil
go 语言统一工具类
## Installation

When used with Go modules, use the following import path:
```go
go get github.com/HC74/gutil@latest
```

## 工具详情

# 字符串


### `StringIsEmpty`

#### 描述
判断字符串是否为空。

#### 参数
- `str`: 要判断的字符串。

#### 返回值
布尔值，如果为空就返回true。

#### 示例
```go

import (
	"github.com/HC74/gutil"
)
str := ""
isnull := gutil.StringIsEmpty(str) // true
```