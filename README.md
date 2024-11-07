# Tagger

自动化版本标签管理工具，用于自动查找仓库最新标签并递增版本号。 用来和 CI/CD 配合使用更佳。

> 我会用它来自动打标签来触发流水线。

![](./tagger_pic.gif)

## 安装

```bash
curl -s https://cdn.mereith.com/tagger/tagger.sh | sh
```

> **Mac 用户注意：首次运行请在"系统设置 -> 隐私与安全"中允许应用运行**

## 使用说明

基本格式：

```bash
tagger [版本类型] [-p <前缀>] [-s <后缀>]
```

参数说明：

- 版本类型：patch(修订版本，默认) | minor(次版本) | major(主版本)
- -p：标签前缀，默认为"v"，如：v0.0.1、prod-0.0.1
- -s：标签后缀，默认为空，如：v0.0.1-dev

相关命令：

```bash
# 设置默认前缀
tagger set-default-prefix <前缀>

# 设置默认后缀
tagger set-default-suffix <后缀>

# 查看当前默认配置
tagger info
```

使用示例：

```bash
# 打补丁版本（v0.0.1 -> v0.0.2）
tagger

# 打次版本（v0.0.1 -> v0.1.0）
tagger minor

# 打主版本（v0.0.1 -> v1.0.0）
tagger major

# 使用自定义前缀（prod-0.0.1 -> prod-0.0.2）
tagger -p prod-

# 使用自定义后缀（v0.0.1-dev -> v0.0.2-dev）
tagger -s -dev

# 同时使用前缀和后缀（rc-0.0.1-beta -> rc-0.0.2-beta）
tagger -p rc- -s -beta

# 打次版本并使用前后缀（test-0.1.0-alpha -> test-0.2.0-alpha）
tagger minor -p test- -s -alpha

# 设置默认前缀为 release-
tagger set-default-prefix release-

# 设置默认后缀为 -stable
tagger set-default-suffix -stable

# 使用默认前后缀打标签（release-0.0.1-stable -> release-0.0.2-stable）
tagger
```
