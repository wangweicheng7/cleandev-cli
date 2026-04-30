# cleandev-cli

一个面向 macOS 开发者的命令行清理工具，默认安全优先（先预览，后清理）。

## 功能

- `scan` / `plan`：扫描可清理项并输出预览。
- `clean --confirm`：显式确认后执行清理。
- `profile` 分层：`safe` / `dev` / `aggressive`。
- 高风险项自动降级为 `report_only`。
- 保护路径硬规则，避免误删关键目录。
- 审计日志记录清理行为。

## 快速开始

```bash
go run ./cmd/cleaner plan --profile dev
go run ./cmd/cleaner scan --profile dev --json
go run ./cmd/cleaner clean --profile dev --confirm
```

## 命令

- `cleaner scan [--profile dev] [--json] [--category cache,logs]`
- `cleaner plan [--profile dev] [--json] [--category cache,logs]`
- `cleaner clean --confirm [--profile dev] [--category cache,logs]`
- `cleaner doctor`
- `cleaner config init`

## 配置文件

可在项目目录执行以下命令生成模板：

```bash
go run ./cmd/cleaner config init
```

默认配置文件名：`.cleandevrc.json`。

优先级：CLI 参数 > 显式 `--config` > 当前目录配置文件 > 默认配置。

## 开发

```bash
make fmt
make test
make run ARGS="plan --profile dev"
make build
```

## 安装

- 安装到用户目录（推荐，无需 sudo）：

```bash
make install-user
echo 'export PATH="$HOME/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

- 安装到系统目录：

```bash
make install
```

- 卸载：

```bash
make uninstall
```

## Homebrew Tap 骨架

仓库内已提供公式模板：`homebrew-tap/Formula/cleaner.rb`。

- 你后续只需替换：
  - `url`
  - `sha256`
  - `version`
- 然后把该目录推到独立 tap 仓库（例如 `homebrew-tap`）即可给他人 `brew install` 使用。
