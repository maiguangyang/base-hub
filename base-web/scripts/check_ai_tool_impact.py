"""Require an explicit AI tool impact declaration for source changes in a PR."""

import argparse
import json
import os
import re
import subprocess
import sys


IMPACT_BLOCK = re.compile(r"^## AI_TOOL_IMPACT[ \t]*\r?\n(?:[ \t]*\r?\n)*```json[ \t]*\r?\n(.*?)\r?\n```[ \t]*(?:\r?\n|$)", re.MULTILINE | re.DOTALL)
OUTCOMES = {"UPDATED", "VERIFIED_NO_CHANGE", "NON_CALLABLE", "NO_IMPACT"}


def touches_function(path: str, scope: str) -> bool:
    if scope == "web":
        return path.startswith("src/")
    if scope == "app":
        return path.startswith("lib/")
    if scope == "engine":
        return ((path.endswith(".go") and not path.endswith("_test.go") and not path.startswith("gen/"))
                or (path.startswith("model/") and path.endswith(".graphql"))
                or path == "tools/contract_inventory.json")
    raise ValueError(f"unknown repository scope: {scope}")


def checked_ids(value: object, label: str) -> list[str]:
    if not isinstance(value, list) or any(not isinstance(item, str) or not item.strip() for item in value):
        raise ValueError(f"{label} must be an array of nonempty IDs")
    return value


def validate_impact(body: str, changed_paths: list[str], scope: str) -> None:
    if not any(touches_function(path, scope) for path in changed_paths):
        return
    blocks = IMPACT_BLOCK.findall(body)
    if len(blocks) != 1:
        raise ValueError("source changes require exactly one visible AI_TOOL_IMPACT JSON block in the PR body")
    try:
        impact = json.loads(blocks[0])
    except json.JSONDecodeError as error:
        raise ValueError(f"invalid AI_TOOL_IMPACT JSON: {error}") from error
    if not isinstance(impact, dict):
        raise ValueError("AI tool impact declaration must be a JSON object")
    if "impacts" in impact and "outcome" in impact:
        raise ValueError("ambiguous AI tool impact declaration")
    entries = impact.get("impacts", [impact])
    if not isinstance(entries, list) or not entries:
        raise ValueError("AI tool impacts must be a nonempty array")
    seen_operations: set[str] = set()
    for entry in entries:
        validate_entry(entry)
        for operation in entry["operationIds"]:
            if operation in seen_operations:
                raise ValueError(f"duplicate operation ID: {operation}")
            seen_operations.add(operation)


def validate_entry(impact: object) -> None:
    if not isinstance(impact, dict) or impact.get("outcome") not in OUTCOMES:
        raise ValueError("AI tool impact outcome is invalid")
    operations = checked_ids(impact.get("operationIds"), "operationIds")
    tools = checked_ids(impact.get("toolIds"), "toolIds")
    reason = impact.get("reason")
    evidence = impact.get("evidence")
    outcome = impact["outcome"]
    if outcome in {"UPDATED", "VERIFIED_NO_CHANGE"}:
        if not operations or not tools or not isinstance(evidence, str) or not evidence.strip():
            raise ValueError(f"{outcome} requires operation IDs, tool IDs, and test/review evidence")
    elif outcome == "NON_CALLABLE":
        if not operations or tools or not isinstance(reason, str) or len(reason.strip()) < 12:
            raise ValueError("NON_CALLABLE requires operation IDs and a concrete reason, without tool IDs")
    elif operations or tools or not isinstance(reason, str) or len(reason.strip()) < 12:
        raise ValueError("NO_IMPACT requires empty IDs and a concrete reason")


def changed_since(base: str) -> list[str]:
    if re.fullmatch(r"[0-9a-f]{40}", base) is None:
        raise ValueError("base must be a 40-character git commit SHA")
    result = subprocess.run(["git", "diff", "--name-only", f"{base}...HEAD", "--"],
                            check=True, capture_output=True, text=True)
    return result.stdout.splitlines()


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--base", required=True)
    parser.add_argument("--scope", choices=("engine", "web", "app"), required=True)
    args = parser.parse_args()
    try:
        validate_impact(os.environ.get("AI_TOOL_IMPACT_BODY", ""), changed_since(args.base), args.scope)
    except (ValueError, subprocess.CalledProcessError) as error:
        print(f"AI tool impact gate: {error}", file=sys.stderr)
        return 1
    print("AI tool impact gate passed")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
