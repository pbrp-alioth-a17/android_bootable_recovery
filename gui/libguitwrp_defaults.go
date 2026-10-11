package twrp

import (
	"android/soong/android"
	"android/soong/cc"
	"fmt"
	"strconv"
	"strings"
)

func printThemeWarning(theme string) {
	if theme == "" {
		theme = "not set"
	}
	themeWarning := "***************************************************************************\n"
	themeWarning += "Could not find ui.xml for TW_THEME: "
	themeWarning += theme
	themeWarning += "\nSet TARGET_SCREEN_WIDTH and TARGET_SCREEN_HEIGHT to automatically select\n"
	themeWarning += "an appropriate theme, or set TW_THEME to one of the following:\n"
	themeWarning += "landscape_hdpi landscape_mdpi portrait_hdpi portrait_mdpi watch_mdpi\n"
	themeWarning += "****************************************************************************\n"
	themeWarning += "(theme selection failed; exiting)\n"

	fmt.Printf(themeWarning)
}

func printCustomThemeWarning(theme string, location string) {
	customThemeWarning := "****************************************************************************\n"
	customThemeWarning += "Could not find ui.xml for TW_CUSTOM_THEME: "
	customThemeWarning += theme + "\n"
	customThemeWarning += "Expected to find custom theme's ui.xml at: "
	customThemeWarning += location
	customThemeWarning += "Please fix this or set TW_THEME to one of the following:\n"
	customThemeWarning += "landscape_hdpi landscape_mdpi portrait_hdpi portrait_mdpi watch_mdpi\n"
	customThemeWarning += "****************************************************************************\n"
	customThemeWarning += "(theme selection failed; exiting)\n"
	fmt.Printf(customThemeWarning)
}

func determineTheme(ctx android.BaseContext) string {
	guiWidth := 0
	guiHeight := 0
	if getMakeVars(ctx, "TW_CUSTOM_THEME") == "" {
		if getMakeVars(ctx, "TW_THEME") == "" {
			if getMakeVars(ctx, "DEVICE_RESOLUTION") == "" {
				width, err := strconv.Atoi(getMakeVars(ctx, "TARGET_SCREEN_WIDTH"))
				if err == nil {
					guiWidth = width
				}
				height, err := strconv.Atoi(getMakeVars(ctx, "TARGET_SCREEN_HEIGHT"))
				if err == nil {
					guiHeight = height
				}
			} else {
				deviceRes := getMakeVars(ctx, "DEVICE_RESOLUTION")
				width, err := strconv.Atoi(strings.Split(deviceRes, "x")[0])
				if err == nil {
					guiWidth = width
				}
				height, err := strconv.Atoi(strings.Split(deviceRes, "x")[1])
				if err == nil {
					guiHeight = height
				}
			}
		}
		if guiWidth > 100 {
			if guiHeight > 100 {
				if guiWidth > guiHeight {
					if guiWidth > 1280 {
						return "landscape_hdpi"
					} else {
						return "landscape_mdpi"
					}
				} else if guiWidth < guiHeight {
					if guiWidth > 720 {
						return "portrait_hdpi"
					} else {
						return "portrait_mdpi"
					}
				} else if guiWidth == guiHeight {
					return "watch_mdpi"
				}
			}
		}
	}

	return getMakeVars(ctx, "TW_THEME")
}

// Resource copying belongs to gui/Android.mk, not Soong graph generation.
func validateTheme(ctx android.BaseContext) bool {
	customTheme := getMakeVars(ctx, "TW_CUSTOM_THEME")
	if customTheme != "" {
		if android.ExistentPathForSource(ctx, customTheme+"/ui.xml").Valid() {
			return true
		}
		printCustomThemeWarning(customTheme, customTheme+"/ui.xml")
		return false
	}
	theme := determineTheme(ctx)
	if android.ExistentPathForSource(ctx, fmt.Sprintf("%s/theme/%s/ui.xml", ctx.ModuleDir(), theme)).Valid() {
		return true
	}
	printThemeWarning(theme)
	return false
}

func globalFlags(ctx android.BaseContext) []string {
	var cflags []string

	if framerate := ctx.Config().VendorConfig("pbrpGlobalVars").String("tw_framerate"); framerate != "" {
		cflags = append(cflags, "-DTW_FRAMERATE="+framerate)
	}

	if getMakeVars(ctx, "AB_OTA_UPDATER") == "true" {
		cflags = append(cflags, "-DAB_OTA_UPDATER=1")
	}

	if getMakeVars(ctx, "PB_TORCH_PATH") != "" {
		cflags = append(cflags, "-DPB_TORCH_PATH="+getMakeVars(ctx, "PB_TORCH_PATH"))
	}

	if getMakeVars(ctx, "PB_TORCH_MAX_BRIGHTNESS") != "" {
		cflags = append(cflags, "-DPB_MAX_BRIGHT_VALUE="+getMakeVars(ctx, "PB_TORCH_MAX_BRIGHTNESS"))
	}
	return cflags
}

func globalSrcs(ctx android.BaseContext) []string {
	var srcs []string

	if getMakeVars(ctx, "TWRP_CUSTOM_KEYBOARD") != "" {
		srcs = append(srcs, getMakeVars(ctx, "TWRP_CUSTOM_KEYBOARD"))
	} else {
		srcs = append(srcs, "hardwarekeyboard.cpp")
	}
	return srcs
}

func libGuiDefaults(ctx android.LoadHookContext) {
	type props struct {
		Target struct {
			Android struct {
				Cflags  []string
				Enabled *bool
			}
		}
		Cflags       []string
		Srcs         []string
		Include_dirs []string
	}

	p := &props{}
	p.Cflags = globalFlags(ctx)
	s := globalSrcs(ctx)
	p.Srcs = s
	ctx.AppendProperties(p)
	if !validateTheme(ctx) {
		ctx.ModuleErrorf("invalid recovery theme")
	}
}

func init() {
	android.RegisterModuleType("libguitwrp_defaults", libGuiDefaultsFactory)
}

func libGuiDefaultsFactory() android.Module {
	module := cc.DefaultsFactory()
	android.AddLoadHook(module, libGuiDefaults)

	return module
}
