package main

import "testing"

func TestRelativeName(t *testing.T) {
	tests := []struct {
		name, fqdn, zone, want string
	}{
		{"subdomain", "_acme-challenge.web.yuchat.ai.", "yuchat.ai.", "_acme-challenge.web"},
		{"apex", "yuchat.ai.", "yuchat.ai.", "@"},
		{"apex challenge", "_acme-challenge.yuchat.ai.", "yuchat.ai.", "_acme-challenge"},
		{"no trailing dots", "_acme-challenge.web.yuchat.ai", "yuchat.ai", "_acme-challenge.web"},
		{"nested", "_acme-challenge.a.b.yuchat.ai.", "yuchat.ai.", "_acme-challenge.a.b"},
		{"suffix without dot boundary", "_acme-challenge.notyuchat.ai.", "yuchat.ai.", "_acme-challenge.notyuchat.ai"},
		{"case", "_acme-challenge.Web.YuChat.ai.", "yuchat.ai.", "_acme-challenge.Web"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := relativeName(tt.fqdn, tt.zone); got != tt.want {
				t.Errorf("relativeName(%q, %q) = %q, want %q", tt.fqdn, tt.zone, got, tt.want)
			}
		})
	}
}

func TestRecordMatches(t *testing.T) {
	const (
		fqdn = "_acme-challenge.web.yuchat.ai."
		zone = "yuchat.ai."
		key  = "abc123"
	)
	tests := []struct {
		name, recName, content string
		want                   bool
	}{
		{"relative name", "_acme-challenge.web", "abc123", true},
		{"fqdn with dot", "_acme-challenge.web.yuchat.ai.", "abc123", true},
		{"fqdn without dot", "_acme-challenge.web.yuchat.ai", "abc123", true},
		{"quoted content", "_acme-challenge.web", `"abc123"`, true},
		{"other key (another SAN)", "_acme-challenge.web", "other", false},
		{"other name", "_acme-challenge.api", "abc123", false},
		{"apex record", "@", "abc123", false},
		{"empty content", "_acme-challenge.web", "", false},
		{"single quote char", "_acme-challenge.web", `"`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := recordMatches(tt.recName, tt.content, fqdn, zone, key); got != tt.want {
				t.Errorf("recordMatches(%q, %q) = %v, want %v", tt.recName, tt.content, got, tt.want)
			}
		})
	}
	if !recordMatches("@", "abc123", "yuchat.ai.", zone, key) {
		t.Error("apex challenge should match record named @")
	}
}
