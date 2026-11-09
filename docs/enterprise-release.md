# 构建与发布

`make check` 执行 Go 测试与 vet；`make build` 构建本机 CLI；`make dist` 交叉构建五个平台包；`bash scripts/verify-dist.sh` 验证 SHA-256，包内容由自动测试检查。安装包中的 CLI 内嵌 Skill、参考和元数据，无需从 ERP 仓库加载。未提供 C.UTF-8 的系统可使用 `LC_ALL=C LANG=C`。

源码版本、manifest 与 Connector 元数据当前为 `0.1.1`。远程为 `https://github.com/xsy-tech/supplier-agent-skill.git`，私有仓库发布默认目标为 `xsy-tech/supplier-agent-skill`；保留版权并遵守 [LICENSE](../LICENSE)。每个提交包含 GYXT-35342 和 AGENTS.md 内的完整工单链接，先运行 `make install-hooks`。

验证并提交源码、推送到目标仓库且发布已获授权后，执行 `make release VERSION=0.1.1`。脚本先验证、构建，将五个平台包与校验清单上传到草稿，全部成功后发布；不会覆盖已存在标签。其他托管平台使用同一批校验过的安装包按该平台流程分发。文档只填写真实存在的下载地址。
