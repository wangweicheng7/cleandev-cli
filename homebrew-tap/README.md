# Homebrew Tap Skeleton

这个目录是 tap 仓库内容模板。

## 使用方式

1. 将此目录内容推送到独立仓库（建议名：`homebrew-tap`）。
2. 在 `Formula/cleaner.rb` 中替换：
   - `homepage`
   - `url`
   - `sha256`
3. 发布后用户可执行：

```bash
brew tap your-org/tap-repo
brew install cleaner
```
