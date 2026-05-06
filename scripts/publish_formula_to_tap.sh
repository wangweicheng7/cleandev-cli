#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
TAP_DIR="${ROOT_DIR}/homebrew-tap"
FORMULA_REL="Formula/devclean-cli.rb"
FORMULA_SRC="${TAP_DIR}/${FORMULA_REL}"

if [[ ! -d "${TAP_DIR}/.git" ]]; then
  echo "homebrew-tap is not a git repository: ${TAP_DIR}" >&2
  exit 1
fi
if [[ ! -f "${FORMULA_SRC}" ]]; then
  echo "formula file not found: ${FORMULA_SRC}" >&2
  exit 1
fi

cd "${TAP_DIR}"
git add "${FORMULA_REL}"

if git diff --cached --quiet; then
  echo "no formula changes to publish"
  exit 0
fi

MSG="${MSG:-chore: publish devclean-cli formula update}"
git commit -m "${MSG}"
git push
echo "published ${FORMULA_REL} to tap remote"

