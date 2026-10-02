package tools

import (
	"fmt"

	"base-engine/auth"
	"base-engine/src/services/ai"
	"base-engine/system_prompt"
)

func toolMetadata(id string, mode ai.ToolMode, workspace auth.WorkspaceType) (string, string) {
	title, _ := system_prompt.ToolDescription(id)
	area := "总部后台"
	if workspace == auth.WorkspaceTypeFranchise {
		area = "加盟商后台"
	}
	if mode == ai.ModeReadOnly {
		return title, fmt.Sprintf("在%s执行“%s”。此工具只读取当前账号有权限查看的业务数据，不修改后台。请按参数说明提供目标、分页或筛选条件；需要多页结果时逐页调用，不猜测未返回的数据。", area, title)
	}
	return title, fmt.Sprintf("在%s执行“%s”。此工具会修改业务数据，只能在用户明确提出该操作、服务端生成完整写入计划并由用户确认后调用。目标或创建范围、完整参数及调用次数必须与批准计划一致；先用只读工具核实不确定的对象。", area, title)
}
