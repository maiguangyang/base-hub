# agentkit

可复制的 Go ADK 运行包。它提供独立会话运行、流式文本、Chat Completions 模型连接器、一次性审批令牌和批准步骤约束。当前基线：Go 1.27.1、`google.golang.org/adk/v2 v2.4.0`。

## 复制到其他项目

把整个 `agentkit/` 目录复制到目标 Go module 的根目录，保留 `chatmodel/` 子目录。然后在目标项目运行：

```sh
go get google.golang.org/adk/v2@v2.4.0
go mod tidy
go test ./agentkit/...
```

导入路径随目标项目的 module 名称变化，例如 `your.module/agentkit` 和 `your.module/agentkit/chatmodel`。功能包内部没有 Engine 项目的导入路径。目标项目需要提供模型凭据、可信身份解析、业务工具及其权限校验。

## 最小运行示例

```go
ctx := context.Background()
llm, err := chatmodel.NewModel(ctx, chatmodel.ModelConfig{
    Name: "your-model", BaseURL: "https://model.example/v1", APIKey: apiKey,
})
if err != nil { return err }

_, err = agentkit.Run(ctx, agentkit.Config{
    Name: "assistant", Instruction: instruction, Model: llm,
    Tools: approvedTools,
}, trustedUserID, userPrompt, func(delta string) error {
    return sendTextToClient(delta)
})
return err
```

`approvedTools` 由目标项目按当前用户和当前阶段生成。运行器会把工具错误作为本次运行错误返回，即使模型在工具失败后继续输出文字。

## 写入审批

预览阶段由目标项目读取数据、校验完整工具参数并构造计划；`ApprovalStore` 保存计划和提示词、签发一次性令牌。执行阶段使用从**服务端会话**导出的身份绑定消耗令牌，再在每个业务写调用前核对工具名称和完整 JSON 参数：

```go
store, err := agentkit.NewApprovalStore[MyPlan](time.Now)
if err != nil { return err }
steps := []agentkit.ApprovedStep{{ToolID: "create", Arguments: approvedJSON, MaxCalls: 1}}
if err := store.Register(plan.ID, trustedBinding, plan.ExpiresAt, steps, plan, prompt); err != nil { return err }
token, err := store.Issue(plan.ID, plan)
if err != nil { return err }

// 用户批准后，在新的请求中执行：
approved, err := store.Consume(token, trustedBinding, func(p MyPlan) bool {
    return permissionsStillValid(p)
})
if err != nil { return err }
defer store.Cancel(approved.ID)
if err := store.AuthorizeStep(approved.ID, "create", approvedJSON); err != nil { return err }
if err := performBusinessWrite(); err != nil { return err }
if err := store.RecordSuccess(approved.ID, "create"); err != nil { return err }
return store.Complete(approved.ID)
```

接入项目必须在每次业务调用前重新检查身份、权限和模型配置，并在写入结果不确定时停止后续步骤。`RecordSuccess` 只用于确定成功的写入。`ApprovalStore` 保存在单个进程内；预览与执行必须命中同一进程。多实例部署需要另行提供共享存储实现。

Engine 的接入方式见 `src/services/ai/`：业务工具目录、GraphQL/HTTP 固定调用、模型配置持久化、审计与 SSE 对外事件仍属于接入项目。
