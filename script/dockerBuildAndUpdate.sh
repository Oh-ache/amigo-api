#!/bin/sh
set -e

REGISTRY="registry.cn-hangzhou.aliyuncs.com/theacheapp"
COMPOSE="deploy/docker-compose.yml"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

# service 名 -> (子脚本, 模块名)  为空时表示参数非法
to_build_cmd() {
  case "$1" in
    sdk-rpc)      echo "script/dockerRpc sdk" ;;
    sdk-api)      echo "script/dockerApi sdk" ;;
    basecode-rpc) echo "script/dockerRpc baseCode" ;;
    basecode-api) echo "script/dockerApi baseCode" ;;
    user-rpc)     echo "script/dockerRpc user" ;;
    user-api)     echo "script/dockerApi user" ;;
    device-rpc)   echo "script/dockerRpc device" ;;
    device-api)   echo "script/dockerApi device" ;;
    mqueue-job)   echo "script/dockerMqueue" ;;
    gateway)      echo "script/dockerGateway" ;;
    *) echo "" ;;
  esac
}

usage() {
  cat <<EOF
用法:
  $0 <service> [<service> ...] [<tag>]
  $0 --auto [<tag>]

service 列表 (与 docker-compose.yml 中的 service 名一致):
  sdk-rpc  sdk-api  basecode-rpc  basecode-api
  user-rpc  user-api  device-rpc  device-api
  mqueue-job  gateway

示例:
  $0 sdk-rpc
  $0 sdk-rpc sdk-api
  $0 sdk-rpc basecode-api 20260920220000
  $0 --auto

参数:
  --auto     根据 git diff HEAD~1 HEAD 自动识别改动了哪些 service
             （app/<m>/<rpc|api>/ 或 common/utils/ 视为影响 sdk-rpc）
  --dry-run  跳过 docker 构建/推送，仅跑参数解析+备份+compose 修改
EOF
}

# ---------- 入参解析 ----------
tag=""
auto=0
dryrun=0
services=""
for arg in "$@"; do
  case "$arg" in
    --auto) auto=1 ;;
    --dry-run) dryrun=1 ;;
    -h|--help) usage; exit 0 ;;
    -*) echo "未知参数: $arg"; usage; exit 1 ;;
    *)
      if [ -z "$tag" ] && echo "$arg" | grep -qE '^[0-9]+$'; then
        tag="$arg"
      else
        services="$services $arg"
      fi
      ;;
  esac
done

if [ "$auto" = 1 ]; then
  services=$(cd "$ROOT_DIR" && git diff --name-only HEAD~1 HEAD 2>/dev/null \
    | python3 -c '
import sys, re
seen = set()
for line in sys.stdin:
    p = line.strip()
    if not p: continue
    m = re.match(r"^app/([^/]+)/(rpc|api)/", p)
    if m:
        seen.add(m.group(1) + "-" + m.group(2))
        continue
    if p.startswith("common/utils/"):
        seen.add("sdk-rpc")
print("\n".join(sorted(seen)))
') || true
  if [ -z "$services" ]; then
    echo "--auto 未检测到任何改动。请检查 git 历史或改用手动指定 service。"
    exit 1
  fi
  echo "[--auto] 检测到改动 service:$services"
fi

[ -z "$services" ] && { usage; exit 1; }

# ---------- 生成 tag ----------
[ -n "$tag" ] || tag=$(date "+%Y%m%d%H%M%S")

# ---------- 校验所有 service 名 ----------
for svc in $services; do
  cmd=$(to_build_cmd "$svc")
  if [ -z "$cmd" ]; then
    echo "错误: 未知 service '$svc'"
    exit 1
  fi
done

echo ""
echo "========================================"
echo " 将要构建的镜像 (tag=$tag):"
for svc in $services; do
  printf "  - %-15s -> sh %s %s\n" "$svc" "$(to_build_cmd "$svc")" "$tag"
done
echo "========================================"
echo ""

# ---------- 1. 备份 ----------
cd "$ROOT_DIR"
cp "$COMPOSE" "${COMPOSE}.bak.${tag}"
echo "[1/3] 已备份 ${COMPOSE} -> ${COMPOSE}.bak.${tag}"

# ---------- 2. 构建 + 推送 ----------
if [ "$dryrun" = 1 ]; then
  echo "[2/3] --dry-run 跳过 docker 构建/推送"
else
  echo "[2/3] 开始构建并推送镜像..."
  for svc in $services; do
    cmd=$(to_build_cmd "$svc")
    echo "----> $svc  (sh $cmd $tag)"
    sh $cmd $tag
  done
fi

# ---------- 3. 修改 docker-compose.yml ----------
echo "[3/3] 更新 $COMPOSE 中的 image tag..."
python3 - "$COMPOSE" "$tag" $services <<'PY'
import re, sys
path, tag, svcs = sys.argv[1], sys.argv[2], sys.argv[3:]
with open(path) as f:
    content = f.read()
total = 0
for svc in svcs:
    pat = re.compile(
        r"(image:\s*registry\.cn-hangzhou\.aliyuncs\.com/theacheapp/"
        + re.escape(svc) + r":)\d+"
    )
    new, n = pat.subn(r"\g<1>" + tag, content)
    if n:
        content = new
        total += n
        print(f"  {svc}: 替换 {n} 处")
    else:
        print(f"  {svc}: 未匹配到（该行可能已被注释）")
with open(path, "w") as f:
    f.write(content)
sys.exit(0 if total > 0 else 2)
PY

echo ""
echo "✅ 完成。新 tag = $tag"
echo ""
echo "--- git diff ---"
git --no-pager diff -- "$COMPOSE" || true
echo ""
echo "下一步：审阅 diff 后手动提交："
echo "  git add $COMPOSE"
echo "  git commit -m \"chore: 升级镜像 tag=$tag\""
echo ""
if [ "$dryrun" = 1 ]; then
  echo "⚠️  本次为 --dry-run，未实际构建/推送镜像。"
  echo "    回滚 compose 改动:  mv ${COMPOSE}.bak.${tag} ${COMPOSE}"
fi
