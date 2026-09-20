#!/bin/sh
# 把 script/hooks/ 下的钩子安装到 .git/hooks/。
# 已存在的同名钩子（非示例、非指向 script/hooks 的符号链接）会备份为 .bak。

set -e

REPO_ROOT="$(git rev-parse --show-toplevel)"
cd "$REPO_ROOT"

hooks_dir=".git/hooks"
src_dir="script/hooks"

if [ ! -d "$src_dir" ]; then
    echo "ERROR: $src_dir 目录不存在"
    exit 1
fi

installed=0
for hook in "$src_dir"/*; do
    [ -f "$hook" ] || continue
    name=$(basename "$hook")
    target="$hooks_dir/$name"

    # 已存在且不是指向源文件的符号链接 → 备份
    if [ -e "$target" ] && { [ ! -L "$target" ] || [ "$(readlink "$target" 2>/dev/null)" != "$hook" ]; }; then
        if [ -e "$target.bak" ]; then
            echo "WARN: $target.bak 已存在，跳过安装 $name"
            continue
        fi
        mv "$target" "$target.bak"
        echo "已备份: $target -> $target.bak"
    fi

    cp "$hook" "$target"
    chmod +x "$target"
    echo "已安装: $target"
    installed=$((installed + 1))
done

if [ "$installed" -eq 0 ]; then
    echo "ERROR: $src_dir 下未找到任何钩子文件"
    exit 1
fi

echo ""
echo "✅ 完成，共安装 $installed 个钩子。"
echo "   跳过钩子: git commit --no-verify"