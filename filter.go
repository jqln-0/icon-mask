package main

import "strings"

// CreatePropertyFiller creates a pipeline worker that creates a GTKIconProperties struct
// from each input string representing an icon name.
func CreatePropertyFiller(theme GTKTheme, in <-chan string) <-chan GTKIconProperties {
	out := make(chan GTKIconProperties)
	go func() {
		for val := range in {
			properties, err := theme.GetIconProperties(val)
			if err == nil {
				out <- properties
			}
		}
		close(out)
	}()
	return out
}

// CreateThemeFilter creates a pipeline worker which filters out any icon names which are
// provided by the given icon theme.
func CreateThemeFilter(theme GTKTheme, in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		for val := range in {
			// We determine whether or not the icon is in the theme simple by looking
			// at the filepath. TODO: Find a better way.
			path := theme.GetIcon(val, 64)
			if strings.Contains(path, theme.name) {
			} else {
				out <- val
			}
		}
		close(out)
	}()
	return out
}

// Generator for the icon name filter pipeline. Places into the output channel the
// names of all icons provided by the given theme.
func GenerateIconNames(theme GTKTheme) <-chan string {
	out := make(chan string)
	go func() {
		names := theme.GetAllIcons()
		for _, val := range names {
			out <- val
		}
		close(out)
	}()
	return out
}

// isInCategories checks whether or not the given candidate string is included in
// the given list of category names.
func isInCategories(categories []string, candidate string) bool {
	for _, cat := range categories {
		if candidate == cat {
			return true
		}
	}
	return false
}

// CreateCategoryFilter creates a pipeline filter that discards any icon not of one of the
// given categories.
func CreateCategoryFilter(categories []string, in <-chan GTKIconProperties) <-chan GTKIconProperties {
	out := make(chan GTKIconProperties)
	go func() {
		for val := range in {
			if isInCategories(categories, val.Category) {
				out <- val
			}
		}
		close(out)
	}()
	return out
}
