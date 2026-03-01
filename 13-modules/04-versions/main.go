package main

import (
	"fmt"
)

// 版本管理示例

func main() {
	fmt.Println("=== 版本管理示例 ===")
	fmt.Println()

	fmt.Println("语义化版本 (SemVer):")
	fmt.Println("v1.2.3")
	fmt.Println(" │ │ └─ Patch: Bug 修复（向后兼容）")
	fmt.Println(" │ └─── Minor: 新功能（向后兼容）")
	fmt.Println(" └───── Major: 破坏性变更（不兼容）")
	fmt.Println()

	fmt.Println("版本升级示例:")
	fmt.Println("  v1.0.0 → v1.0.1  # Bug 修复")
	fmt.Println("  v1.0.0 → v1.1.0  # 新功能")
	fmt.Println("  v1.0.0 → v2.0.0  # 破坏性变更")
	fmt.Println()

	fmt.Println("主版本升级 (v2+):")
	fmt.Println("  // v1 版本")
	fmt.Println("  import \"github.com/some/package\"")
	fmt.Println()
	fmt.Println("  // v2 版本（路径变化！）")
	fmt.Println("  import \"github.com/some/package/v2\"")
	fmt.Println()

	fmt.Println("常用命令:")
	fmt.Println("  go get github.com/pkg/name@latest      # 最新版本")
	fmt.Println("  go get github.com/pkg/name@v1.2.3      # 特定版本")
	fmt.Println("  go get github.com/pkg/name@v1.2        # 最新小版本")
	fmt.Println("  go get -u=patch ./...                  # 升级补丁版本")
	fmt.Println("  go list -m -versions github.com/pkg    # 查看所有版本")
}
