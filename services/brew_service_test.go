package services

import "testing"

func TestParseNameLines(t *testing.T) {
	input := "\nfoo\nbar  \n\n baz\n"
	got := parseNameLines(input)
	if len(got) != 3 {
		t.Fatalf("expected 3 items, got %d", len(got))
	}
	if got[0] != "foo" || got[1] != "bar" || got[2] != "baz" {
		t.Fatalf("unexpected parsed lines: %#v", got)
	}
}

func TestCountNonEmptyLines(t *testing.T) {
	input := "\na\n\n b \n"
	if got := countNonEmptyLines(input); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestIsPinnedUpgradeNotice(t *testing.T) {
	input := "Error: Not upgrading 1 pinned package:\nnode 23.0.0\n"
	if !isPinnedUpgradeNotice(input) {
		t.Fatalf("expected pinned upgrade notice to be detected")
	}
	if isPinnedUpgradeNotice("Error: Failed to download bottle") {
		t.Fatalf("expected unrelated errors to stay failures")
	}
}

func TestPinCommandArgs(t *testing.T) {
	pinArgs, err := pinCommandArgs("node", true)
	if err != nil {
		t.Fatalf("pinCommandArgs returned error: %v", err)
	}
	if len(pinArgs) != 2 || pinArgs[0] != "pin" || pinArgs[1] != "node" {
		t.Fatalf("unexpected pin args: %#v", pinArgs)
	}

	unpinArgs, err := pinCommandArgs("node", false)
	if err != nil {
		t.Fatalf("pinCommandArgs returned error: %v", err)
	}
	if len(unpinArgs) != 2 || unpinArgs[0] != "unpin" || unpinArgs[1] != "node" {
		t.Fatalf("unexpected unpin args: %#v", unpinArgs)
	}

	if _, err := pinCommandArgs(" ", true); err == nil {
		t.Fatalf("expected empty package name to fail")
	}
}

func TestParsePackageInfoJSONIncludesCask(t *testing.T) {
	raw := `{"formulae":[],"casks":[{"name":"visual-studio-code","full_name":"visual-studio-code","tap":"homebrew/cask","desc":"Open-source code editor","homepage":"https://code.visualstudio.com/","version":"1.99.0","installed":"1.98.2","auto_updates":true,"token":"visual-studio-code"}]}`
	info, err := parsePackageInfoJSON(raw, "visual-studio-code", "cask")
	if err != nil {
		t.Fatalf("parsePackageInfoJSON returned error: %v", err)
	}
	if info.Type != "cask" || info.Name != "visual-studio-code" || info.CurrentVersion != "1.99.0" || info.InstalledVersion != "1.98.2" {
		t.Fatalf("unexpected cask info: %+v", info)
	}
	if !info.AutoUpdates {
		t.Fatalf("expected auto_updates to be true")
	}
}

func TestParsePackageInfoJSONIncludesFormula(t *testing.T) {
	raw := `{"formulae":[{"name":"wget","full_name":"wget","tap":"homebrew/core","desc":"Internet file retriever","homepage":"https://www.gnu.org/software/wget/","license":"GPL-3.0","linked_keg":"1.24.5","pinned":false,"installed":[{"version":"1.24.5","installed_as_dependency":false,"installed_on_request":true}],"versions":{"stable":"1.24.5","head":"HEAD"},"dependencies":["openssl@3","zlib"]}],"casks":[]}`
	info, err := parsePackageInfoJSON(raw, "wget", "formula")
	if err != nil {
		t.Fatalf("parsePackageInfoJSON returned error: %v", err)
	}
	if info.Type != "formula" || info.Name != "wget" || info.CurrentVersion != "1.24.5" || info.InstalledVersion != "1.24.5" {
		t.Fatalf("unexpected formula info: %+v", info)
	}
	if len(info.Dependencies) != 2 || info.Dependencies[0] != "openssl@3" {
		t.Fatalf("unexpected dependencies: %+v", info.Dependencies)
	}
}
