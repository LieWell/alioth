// Package static 把前端静态资源以 embed.FS 形式暴露给其他包使用。
// 由于 Go embed.FS 只能嵌入当前包内的文件，而 web 包需要挂载这些资源，
// 所以把静态文件放在独立的 static 包下，web 包通过 import 引用。
package static

import (
	"embed"
	"io/fs"
)

//go:embed index.html
//go:embed about.html
//go:embed favicon.ico
//go:embed robots.txt
//go:embed css
//go:embed js
var FS embed.FS

// Sub 返回 embed.FS（文件已在 static/ 包根目录下，无需再 Sub）
func Sub() fs.FS {
	return FS
}
