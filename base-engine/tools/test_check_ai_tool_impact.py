import unittest

from check_ai_tool_impact import validate_impact


VALID_NO_IMPACT = '## AI_TOOL_IMPACT\n```json\n{"outcome":"NO_IMPACT","operationIds":[],"toolIds":[],"reason":"Only visual spacing changed","evidence":""}\n```'
VALID_UPDATED = '## AI_TOOL_IMPACT\n```json\n{"outcome":"UPDATED","operationIds":["graphql.query.stores"],"toolIds":["FranchiseStores"],"reason":"","evidence":"go test ./src/integration"}\n```'


class ImpactGateTest(unittest.TestCase):
    def test_mixed_callable_and_non_callable_operations(self):
        body = '''## AI_TOOL_IMPACT
```json
{"impacts":[{"outcome":"UPDATED","operationIds":["graphql.query.stores"],"toolIds":["FranchiseStores"],"reason":"","evidence":"go test ./src/integration"},{"outcome":"NON_CALLABLE","operationIds":["graphql.mutation.login"],"toolIds":[],"reason":"Session operations are prohibited","evidence":""}]}
```'''
        validate_impact(body, ["src/services/store/store.go"], "engine")

    def test_same_operation_cannot_have_conflicting_outcomes(self):
        body = '''## AI_TOOL_IMPACT
```json
{"impacts":[{"outcome":"UPDATED","operationIds":["graphql.query.stores"],"toolIds":["FranchiseStores"],"reason":"","evidence":"go test ./src/integration"},{"outcome":"NON_CALLABLE","operationIds":["graphql.query.stores"],"toolIds":[],"reason":"Session operations are prohibited","evidence":""}]}
```'''
        with self.assertRaisesRegex(ValueError, "duplicate operation ID"):
            validate_impact(body, ["src/services/store/store.go"], "engine")

    def test_array_and_legacy_outcome_cannot_conflict(self):
        body = '''## AI_TOOL_IMPACT
```json
{"outcome":"NO_IMPACT","impacts":[{"outcome":"UPDATED","operationIds":["graphql.query.stores"],"toolIds":["FranchiseStores"],"reason":"","evidence":"go test ./src/integration"}]}
```'''
        with self.assertRaisesRegex(ValueError, "ambiguous"):
            validate_impact(body, ["src/services/store/store.go"], "engine")

    def test_source_change_requires_declaration(self):
        with self.assertRaises(ValueError):
            validate_impact("", ["src/services/store/store.go"], "engine")

    def test_wrong_outcome_is_rejected(self):
        body = VALID_NO_IMPACT.replace("NO_IMPACT", "IGNORED")
        with self.assertRaises(ValueError):
            validate_impact(body, ["src/a.tsx"], "web")

    def test_unchanged_tool_needs_specific_reason(self):
        body = VALID_NO_IMPACT.replace("Only visual spacing changed", "")
        with self.assertRaises(ValueError):
            validate_impact(body, ["lib/a.dart"], "app")

    def test_changed_operation_needs_tool_ids(self):
        body = VALID_UPDATED.replace('"toolIds":["FranchiseStores"]', '"toolIds":[]')
        with self.assertRaises(ValueError):
            validate_impact(body, ["src/a.tsx"], "web")

    def test_verified_no_change_needs_evidence(self):
        body = VALID_UPDATED.replace("UPDATED", "VERIFIED_NO_CHANGE").replace("go test ./src/integration", "")
        with self.assertRaises(ValueError):
            validate_impact(body, ["src/a.tsx"], "web")

    def test_non_callable_requires_operation_and_reason(self):
        body = '## AI_TOOL_IMPACT\n```json\n{"outcome":"NON_CALLABLE","operationIds":["http.ai.status"],"toolIds":[],"reason":"AI must not manage its own configuration","evidence":""}\n```'
        validate_impact(body, ["src/a.tsx"], "web")
        with self.assertRaises(ValueError):
            validate_impact(body.replace('"operationIds":["http.ai.status"]', '"operationIds":[]'), ["src/a.tsx"], "web")

    def test_duplicate_declarations_are_rejected(self):
        with self.assertRaises(ValueError):
            validate_impact(VALID_NO_IMPACT + VALID_NO_IMPACT, ["lib/a.dart"], "app")

    def test_valid_declarations_and_docs_only(self):
        validate_impact(VALID_NO_IMPACT, ["src/a.tsx"], "web")
        validate_impact(VALID_UPDATED, ["lib/a.dart"], "app")
        validate_impact("", ["README.md"], "engine")

    def test_reason_can_contain_braces(self):
        body = VALID_NO_IMPACT.replace("Only visual spacing changed", "Only visual spacing changed in {drawer}")
        validate_impact(body, ["src/a.tsx"], "web")

    def test_visible_json_block_is_accepted(self):
        body = '## AI_TOOL_IMPACT\n```json\n{"outcome":"NO_IMPACT","operationIds":[],"toolIds":[],"reason":"Only visual spacing changed","evidence":""}\n```'
        validate_impact(body, ["src/a.tsx"], "web")

    def test_repository_template_is_parseable_after_editing(self):
        with open('.github/PULL_REQUEST_TEMPLATE.md', encoding='utf-8') as source:
            body = source.read().replace('"reason":""', '"reason":"Only visual spacing changed"')
        self.assertIn('"impacts":[', body)
        validate_impact(body, ["src/a.tsx"], "web")


if __name__ == "__main__":
    unittest.main()
