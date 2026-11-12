# 面料库存 / Fabric inventory

仅使用 supplier-fabric 当前 tools/list 暴露的 `search_fabric_inventory`、`preview_fabric_inventory_import`、`preview_fabric_inventory_adjustments`、`commit_fabric_inventory_import`、`commit_fabric_inventory_adjustments`。五项 Capability 分别授权；只有预检授权不能提交，缺失对应工具时停止该动作。

本地 Python 3 标准库脚本（路径相对已安装 Skill）：

```sh
python3 scripts/fabric_inventory_excel.py template --output template.xlsx
python3 scripts/fabric_inventory_excel.py parse --input filled.xlsx --output local.json --errors-output local-errors.xlsx
python3 scripts/fabric_inventory_excel.py export --input query.json --output inventory.xlsx
python3 scripts/fabric_inventory_excel.py errors --input local.json --result preview-errors.json --output mes-errors.xlsx
```

固定四列：面料编码、面料名称、色号、公司专属成品库存。编码/名称采用供应商字段，不是公司字段或花型编号。编码和色号应输入为文本，保留前导零；名称空白不参与匹配，非空保持原文精确匹配。库存为非负 int32 目标值；0 有效，空值不等于0。只支持 XLSX，不支持 XLS；公式须先粘贴为值。文本形式的 =、+、-、@ 不作为公式执行。

parse 输出是本地 envelope：items、row_mapping、raw_rows、errors。任何本地错误均阻止整批预检；只把 `{"items": local.items}` 发送给预检 MCP，不发送 envelope、Excel、映射或本地路径。每批1–500条，超过时向用户说明并确定新批次范围，不能自动拆分写入。

查询按页读取（默认50，上限100），汇集返回 items 到 query.json 后本地导出；实时分页不保证快照一致。服务端不接收或传输文件。不得使用机器凭证、浏览器登录或其他领域绕过权限。

预检仅是基于当时库存的建议性结果，不保存、不提供 preview_id 或过期时间。导入用 `preview_fabric_inventory_import`，调整用 `preview_fabric_inventory_adjustments`，分别要求相应授权。展示绑定供应商、来源、目标字段、总行数、修改/不变行数及前后值摘要；changes_truncated 时明确摘要不完整。取得用户针对该整批目标值的明确确认后，导入调用 `commit_fabric_inventory_import`，调整调用 `commit_fabric_inventory_adjustments`；参数为 `{"idempotency_key":"本次确认生成的稳定键","items":local.items}`。items 必须与预检的1–500条四字段数据一致；不传本地 envelope、preview_id、source、附件或身份覆盖参数。服务端重新鉴权、解析并校验整批，按提交时库存计算变更；提交结果可能与预检摘要不同。成功仅返回来源、总/变更/不变行数和 request_id，不包含业务明细。事务竞争返回 `FABRIC_INVENTORY_STATE_CHANGED`（409）时重新预检并请用户再次确认；提交时整批行数据无效返回 `INVALID_ARGUMENT`（400），修正后重新预检并确认。

每次经过用户确认的**新提交**生成新幂等键；同次网络重试沿用原键和完全相同的 items。统一 WRITE Executor 在结果保留期内对同键同内容返回首次成功摘要，同键不同内容报幂等冲突，同键仍执行报处理中。用户要求再次执行相同目标值时重新预检、确认并用新键；当前库存相同可能全部计为未变更。提交超时、断线或结果不明确时，不自动重试、不宣称成功或失败；先用 `search_fabric_inventory` 查询当前目标库存，说明查询结果不能证明上次请求是否执行，再请用户决定是否重新确认。原键仅在平台明确可安全重放时用于完全相同请求的重试；不能自行改用新键掩盖不确定结果。授权、网络或事务错误不得伪造为行错误。

本地或MES行错误自动生成独立错误文件并告诉用户“整批未提交”和文件路径；保留整批原四列和空白行，第五列追加错误原因，按1-based items row_index映射原行。五列修正版忽略旧错误列，重新解析、预检、确认。errors_truncated 必须在文件和用户回复说明“错误详情未完整返回”，不能声称未标记行通过。认证、权限、网络或模板损坏不伪造行错误文件。输出文件已存在即拒绝覆盖。

## English

Use only the five inventory tools exposed by the fabric connection; preview and commit have separate grants. Run the local Python 3 commands above; XLSX only, no remote file transfer. The four columns represent supplier material code/name, color number and nonnegative int32 target stock. Preserve code text and leading zeros; blank names omit name matching; zero is valid but blank stock is invalid.

The parse envelope stays local. Send only `{"items": local.items}` to the appropriate preview tool after all local rows pass. Never automatically split more than 500 write rows. Import preview and adjustment preview have separate grants. Preview is advisory, with no preview ID or expiry. Show the bound supplier, source, target field, whole-batch counts and before/after summary; disclose truncated changes. After explicit approval, use `commit_fabric_inventory_import` or `commit_fabric_inventory_adjustments`, respectively, with `{"idempotency_key":"stable-key-for-this-confirmation","items":local.items}` and exactly the same 1-500 four-field rows. Never send the local envelope, preview ID, source, attachment or identity override. The server reauthorizes and revalidates against current stock, so commit counts may differ. Success contains only source, total/changed/unchanged counts and request ID. For transaction races, `FABRIC_INVENTORY_STATE_CHANGED` (409) means preview and confirm again. Invalid batch rows at commit must be corrected, then previewed and confirmed again.

Generate a new idempotency key for each separately confirmed submission. For a network retry of the same request, reuse the original key and identical items; the shared WRITE Executor replays a completed result, rejects changed input for the same key, and reports in-progress requests. The same target values can be committed again with a new key only after a new confirmation; that does not guarantee exactly-once execution beyond the retention period. If a commit times out or its result is unclear, do not retry automatically or claim success/failure. Query current inventory, explain that even matching stock does not prove whether that request executed, and ask the user before another submission. Do not turn authorization, network or transaction failures into fabricated row errors.

Generate a distinct local error workbook for row errors, preserving all original rows and four raw columns. MES row_index is a 1-based item index mapped through row_mapping. Corrected five-column files ignore old error text and repeat the complete workflow. Disclose truncated error detail in both file and reply; unmarked rows are not declared valid. Authentication, permission and network failures do not create fabricated row errors. Query pagination is live rather than a consistent snapshot; collect authorized pages locally before export.
