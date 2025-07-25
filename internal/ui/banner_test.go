package ui

import (
	"strings"
	"testing"
)

func TestRenderBanner(t *testing.T) {
	config := BannerConfig{
		Width:       100,
		ShowVersion: true,
		ColorScheme: "default",
	}

	banner := RenderBanner(config)

	if banner == "" {
		t.Error("Expected banner to return non-empty string")
	}

	if !strings.Contains(banner, "___ ___") {
		t.Error("Expected banner to contain ASCII art")
	}

	if !strings.Contains(banner, "v0.1.0") {
		t.Error("Expected banner to contain version when ShowVersion is true")
	}
}

func TestRenderBannerCompact(t *testing.T) {
	config := BannerConfig{
		Width:       50, // Very narrow terminal to trigger compact mode
		ShowVersion: false,
		ColorScheme: "default",
	}

	banner := RenderBanner(config)

	if banner == "" {
		t.Error("Expected compact banner to return non-empty string")
	}

	// Should use compact layout for narrow terminals
	if !strings.Contains(banner, "Context Engine") {
		t.Error("Expected compact banner to contain title")
	}
}

func TestRenderBannerDefault(t *testing.T) {
	banner := RenderBannerDefault()

	if banner == "" {
		t.Error("Expected default banner to return non-empty string")
	}
}

func TestCenterText(t *testing.T) {
	text := "Hello"
	width := 20
	centered := centerText(text, width)

	// Should have padding on the left
	expectedPadding := (width - len(text)) / 2
	actualPadding := 0
	for _, char := range centered {
		if char == ' ' {
			actualPadding++
		} else {
			break
		}
	}

	if actualPadding != expectedPadding {
		t.Errorf("Expected padding %d, got %d", expectedPadding, actualPadding)
	}
}

func TestGetTerminalWidth(t *testing.T) {
	width := GetTerminalWidth()

	// Should return at least the minimum width
	if width < 80 {
		t.Errorf("Expected minimum width of 80, got %d", width)
	}
}