#!/usr/bin/env bash
#
# hyperpixel-check.sh (Ubuntu Pi + Raspberry Pi OS tolerant)
#

set -u

PASS=0
WARN=0
FAIL=0

green()  { printf "\033[32m%s\033[0m\n" "$*"; }
yellow() { printf "\033[33m%s\033[0m\n" "$*"; }
red()    { printf "\033[31m%s\033[0m\n" "$*"; }

pass() { green "[PASS] $1"; ((PASS++)); }
warn() { yellow "[WARN] $1"; ((WARN++)); }
fail() { red "[FAIL] $1"; ((FAIL++)); }

echo "=== HyperPixel Validation ==="
echo

# OS / kernel
if [[ -f /etc/os-release ]]; then
    . /etc/os-release
    echo "OS        : ${PRETTY_NAME:-unknown}"
fi

echo "Kernel    : $(uname -r)"
echo

# config.txt location
BOOTCFG=""
for f in /boot/firmware/config.txt /boot/config.txt; do
    [[ -f "$f" ]] && BOOTCFG="$f"
done

if [[ -n "$BOOTCFG" ]]; then
    pass "Found boot config: $BOOTCFG"
else
    fail "config.txt not found"
fi

# Collect overlay directories (Ubuntu + Pi OS variants)
OVERLAY_DIRS=()

while IFS= read -r d; do
    [[ -d "$d" ]] && OVERLAY_DIRS+=("$d")
done < <(
    find /boot /usr/lib /lib -type d -name overlays 2>/dev/null
)

# Deduplicate
declare -A seen
UNIQUE_DIRS=()

for d in "${OVERLAY_DIRS[@]}"; do
    if [[ -z "${seen[$d]+x}" ]]; then
        seen["$d"]=1
        UNIQUE_DIRS+=("$d")
    fi
done

if [[ ${#UNIQUE_DIRS[@]} -eq 0 ]]; then
    fail "No overlay directories found"
else
    pass "Overlay directories found:"
    for d in "${UNIQUE_DIRS[@]}"; do
        echo "  - $d"
    done
fi

echo

# Search for HyperPixel DT overlays
HYPERPIXEL_FOUND=0

for d in "${UNIQUE_DIRS[@]}"; do
    while IFS= read -r f; do
        pass "HyperPixel overlay found: $(basename "$f")"
        HYPERPIXEL_FOUND=1
    done < <(find "$d" -maxdepth 1 -type f -iname "*hyperpixel*.dtbo" 2>/dev/null)
done

if [[ $HYPERPIXEL_FOUND -eq 0 ]]; then
    fail "No HyperPixel DT overlays found"
fi

echo

# Check config for dtoverlay usage
if [[ -n "$BOOTCFG" ]]; then
    if grep -Eq 'dtoverlay=.*hyperpixel' "$BOOTCFG"; then
        pass "HyperPixel dtoverlay configured in config.txt"
    else
        warn "HyperPixel dtoverlay NOT configured in config.txt"
    fi

    if grep -Eq 'dtoverlay=vc4-kms-v3d' "$BOOTCFG"; then
        pass "vc4-kms-v3d enabled"
    else
        warn "vc4-kms-v3d not found"
    fi
fi

echo

# DRM sanity
if [[ -d /sys/class/drm ]]; then
    pass "/sys/class/drm present"

    echo "DRM devices:"
    ls /sys/class/drm
else
    fail "DRM subsystem missing"
fi

echo

# Kernel modules (informational only)
for m in vc4; do
    if lsmod | awk '{print $1}' | grep -qx "$m"; then
        pass "Kernel module loaded: $m"
    else
        warn "Kernel module not loaded: $m"
    fi
done

echo

# Tools
for p in kmsprint modetest fbset i2cdetect; do
    if command -v "$p" >/dev/null 2>&1; then
        pass "$p available"
    else
        warn "$p missing"
    fi
done

echo
echo "===== Summary ====="
echo "PASS: $PASS"
echo "WARN: $WARN"
echo "FAIL: $FAIL"

echo

exit $([[ $FAIL -gt 0 ]] && echo 1 || echo 0)

# dtoverlay=vc4-kms-dpi-hyperpixel4
# dtoverlay=vc4-kms-v3d
# sudo apt install -y \
  #  libdrm-tests \
  #  mesa-utils \
  #  mesa-utils-extra

# dtoverlay -l 2>/dev/null
# ls /sys/class/drm/
