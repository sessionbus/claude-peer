// SPDX-License-Identifier: MIT

package interactive

import (
	"testing"

	kit "github.com/antst/sessionbus/bus/sdk/go"
)

func TestManagedToolRejectsExactNativeDeny(t *testing.T) {
	for _, arguments := range [][]string{
		{"--disallowedTools", PublicTool},
		{"--disallowed-tools", PublicTool},
		{"--disallowedTools=" + PublicTool},
		{"--disallowed-tools=Bash," + PublicTool},
		{"--disallowedTools", "Bash(git status)", PublicTool, "--model", "sonnet"},
		{"--disallowedTools", "Read", "--disallowed-tools", "Bash(git status) " + PublicTool},
		{"--system-prompt", "--", "--disallowedTools", PublicTool},
		{"--disallowedTools", "--", PublicTool},
	} {
		if err := ValidateManagedToolArguments(arguments); err == nil {
			t.Fatalf("exact managed deny accepted: %q", arguments)
		}
		if _, _, err := LaunchPlan(arguments, nil, "/cwd", "/plugin", 1000); err == nil {
			t.Fatalf("interactive launch accepted exact managed deny: %q", arguments)
		}
	}
}

func TestManagedToolPreservesOtherNativePolicyAndBoundaries(t *testing.T) {
	for _, arguments := range [][]string{
		{"--disallowedTools", "Bash(git status),Read"},
		{"--disallowed-tools=OtherTool", "prompt"},
		{"--system-prompt", "--disallowedTools", PublicTool},
		{"--system-prompt=--disallowedTools", PublicTool},
		{"--allowedTools", "--disallowedTools", PublicTool},
		{"--disallowedTools", "Read", "--model", PublicTool},
		{"--", "--disallowedTools", PublicTool},
	} {
		if err := ValidateManagedToolArguments(arguments); err != nil {
			t.Fatalf("native arguments %q rejected: %v", arguments, err)
		}
	}
}

func TestManagedToolRuleSplitterMatchesNativeLists(t *testing.T) {
	if !containsManagedToolRule([]string{"Bash(git status, --short), " + PublicTool}) {
		t.Fatal("managed tool absent from comma-separated native rule list")
	}
	if containsManagedToolRule([]string{"Bash(" + PublicTool + ", status) OtherTool"}) {
		t.Fatal("rule content was treated as a separate tool")
	}
	if containsManagedToolRule([]string{"OtherTool\t" + PublicTool}) {
		t.Fatal("native literal-space splitter incorrectly treated a tab as a separator")
	}
	if !containsManagedToolRule([]string{"\ufeff" + PublicTool + "\ufeff"}) {
		t.Fatal("native ECMAScript trim did not remove byte-order marks")
	}
	if containsManagedToolRule([]string{"\u0085" + PublicTool}) {
		t.Fatal("native ECMAScript trim incorrectly removed NEXT LINE")
	}
}

func TestInteractiveGuardUsesProjectedNativeArguments(t *testing.T) {
	for _, group := range []string{"-g", "--group"} {
		arguments, values, err := LaunchPlan([]string{group, "--disallowedTools", PublicTool}, nil, "/cwd", "/plugin", 1000)
		if err != nil {
			t.Fatalf("%s flag-looking group value rejected: %v", group, err)
		}
		if got := Environment(values)["SESSIONBUS_GROUPS"]; got != `["--disallowedTools"]` {
			t.Fatalf("%s groups = %s", group, got)
		}
		if got := arguments[len(arguments)-1]; got != PublicTool {
			t.Fatalf("%s native positional = %q", group, got)
		}
		if _, _, err := LaunchPlan([]string{group, "team", "--disallowedTools", PublicTool}, nil, "/cwd", "/plugin", 1000); err == nil {
			t.Fatalf("%s genuine native deny accepted", group)
		}
	}
}

func TestManagedToolGuardPreservesCurrentNativeRequiredValues(t *testing.T) {
	for _, option := range []string{"--system-prompt-file", "--permission-prompt-tool", "--managed-settings"} {
		arguments := []string{option, "--disallowedTools", PublicTool}
		if err := ValidateManagedToolArguments(arguments); err != nil {
			t.Fatalf("%s flag-looking value rejected: %v", option, err)
		}
		projected, _, err := LaunchPlan(arguments, nil, "/cwd", "/plugin", 1000)
		if err != nil {
			t.Fatalf("%s interactive value rejected: %v", option, err)
		}
		if got := projected[len(projected)-3:]; got[0] != option || got[1] != "--disallowedTools" || got[2] != PublicTool {
			t.Fatalf("%s native values changed: %q", option, got)
		}
	}
	for _, arguments := range [][]string{{"--channels", "--disallowedTools", PublicTool}} {
		if err := ValidateManagedToolArguments(arguments); err != nil {
			t.Fatalf("native values %q rejected: %v", arguments, err)
		}
	}
	for _, arguments := range [][]string{
		{"--channels", "server", "--disallowedTools", PublicTool},
		{"--resume", "session", "--disallowedTools", PublicTool},
		{"--resume", "--disallowedTools", PublicTool},
	} {
		if err := ValidateManagedToolArguments(arguments); err == nil {
			t.Fatalf("genuine native deny accepted after %q", arguments)
		}
	}
}

func TestTypedArgumentsRejectRawSelectors(t *testing.T) {
	for _, tc := range []struct {
		field     string
		open      kit.OpenOptions
		arguments []string
	}{
		{"model", kit.OpenOptions{Model: "sonnet"}, []string{"--model", "haiku"}},
		{"model", kit.OpenOptions{Model: "sonnet"}, []string{"--model=haiku"}},
		{"reasoning_effort", kit.OpenOptions{ReasoningEffort: "high"}, []string{"--effort", "low"}},
		{"reasoning_effort", kit.OpenOptions{ReasoningEffort: "high"}, []string{"--effort=low"}},
		{"permission_mode", kit.OpenOptions{PermissionMode: "default"}, []string{"--permission-mode", "plan"}},
		{"permission_mode", kit.OpenOptions{PermissionMode: "default"}, []string{"--permission-mode=plan"}},
		{"permission_mode", kit.OpenOptions{PermissionMode: "default"}, []string{"--dangerously-skip-permissions"}},
		{"permission_mode", kit.OpenOptions{PermissionMode: "bypassPermissions"}, []string{"--dangerously-skip-permissions"}},
	} {
		tc.open.Arguments = append([]string{"--verbose"}, tc.arguments...)
		err := ValidateTypedArguments(tc.open)
		if err == nil || err.Error() != "argument conflicts with typed field "+tc.field {
			t.Fatalf("%q with typed %s: %v", tc.arguments, tc.field, err)
		}
		// The same selector stays native passthrough when its field is untyped.
		untyped := kit.OpenOptions{Model: "sonnet", ReasoningEffort: "high", PermissionMode: "default", Arguments: tc.open.Arguments}
		switch tc.field {
		case "model":
			untyped.Model = ""
		case "reasoning_effort":
			untyped.ReasoningEffort = ""
		case "permission_mode":
			untyped.PermissionMode = ""
		}
		if err := ValidateTypedArguments(untyped); err != nil {
			t.Fatalf("%q with untyped %s rejected: %v", tc.arguments, tc.field, err)
		}
		if err := ValidateTypedArguments(kit.OpenOptions{Arguments: tc.open.Arguments}); err != nil {
			t.Fatalf("%q without typed fields rejected: %v", tc.arguments, err)
		}
	}
}

func TestTypedArgumentsKeepOperandsAndRequiredValues(t *testing.T) {
	typed := kit.OpenOptions{Model: "sonnet", ReasoningEffort: "high", PermissionMode: "bypassPermissions"}
	for _, arguments := range [][]string{
		{"--", "--model", "haiku", "--effort=low", "--permission-mode", "plan", "--dangerously-skip-permissions"},
		{"--append-system-prompt", "--model"},
		{"--append-system-prompt", "--", "--model", "haiku"},
	} {
		typed.Arguments = arguments
		if err := ValidateTypedArguments(typed); err != nil {
			t.Fatalf("native arguments %q rejected: %v", arguments, err)
		}
	}
}
