package main

/*
#include <stdlib.h>
#include <gtk/gtk.h>
#cgo pkg-config: gtk+-2.0

int iterate_array(gint *array) {
	return array++;
}
*/
import "C"
import "unsafe"

// GTKTheme is a simple wrapper around the C GTKIconTheme struct.
type GTKTheme struct {
	theme *C.GtkIconTheme
}

// GTKIconProperties represents the properties of an icon.
type GTKIconProperties struct {
	// The sizes in which this icon is available.
	Sizes []uint
	// Whether or not the icon has a scalable form.
	Scalable bool
	// The name of the icon.
	Name string
	// The category of the icon ('app', 'emblem', etc.)
	Category string
}

// GTKInit call gtk_init with the given program arguments.
func GTKInit(args []string) {
	// Convert the args to C args.
	argc := C.int(len(args))
	argv := make([]*C.char, argc)

	// Convert each argument.
	for i := 0; i < len(args); i++ {
		argv[i] = C.CString(args[i])
		defer C.free(unsafe.Pointer(argv[i]))
	}

	// Init GTK.
	argvPtr := (**C.char)(unsafe.Pointer(&argv[0]))
	C.gtk_init(&argc, &argvPtr)
}

// CreateTheme creates a new GtkIconTheme and sets it to the given theme name.
func CreateTheme(name string) GTKTheme {
	// Create a new blank theme.
	var t GTKTheme
	t.theme = C.gtk_icon_theme_new()

	// Load the given theme name.
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	C.gtk_icon_theme_set_custom_theme(t.theme, (*C.gchar)(cName))
	return t
}

// GetIcon takes the name of an icon and a size as arguments. It returns the closest
// matching icon file from the theme.
func (t GTKTheme) GetIcon(name string, size int) string {
	// TODO: User-specified flags.
	var flags C.GtkIconLookupFlags

	// Lookup the icon.
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	iconInfo := C.gtk_icon_theme_lookup_icon(t.theme, (*C.gchar)(cName), C.gint(size), flags)

	// Check the lookup was successful.
	if iconInfo == nil {
		return ""
	}
	// gtk_icon_info_free is deprecated, but we seem to have issues using g_object_unref.
	defer C.gtk_icon_info_free(iconInfo)
	//defer C.g_object_unref(C.gpointer(iconInfo))

	// Check the icon's filename.
	cFilename := C.gtk_icon_info_get_filename(iconInfo)
	if cFilename == nil {
		return ""
	}

	filename := C.GoString((*C.char)(cFilename))
	return filename
}

// Returns all the icons in the theme, including inherited and hicolor icons.
func (t GTKTheme) GetAllIcons() []string {
	// Get the list of all icons in this theme.
	out := make([]string, 0)
	list := C.gtk_icon_theme_list_icons(t.theme, nil)
	defer C.g_list_free(list)

	// Convert the list into a slice, freeing used elements as we go.
	for ptr := list; ptr != nil; ptr = ptr.next {
		out = append(out, C.GoString((*C.char)(ptr.data)))
		C.g_free(ptr.data)
	}

	return out
}

// GetIconSizes gets the sizes an icon is available in and whether or not the icon has
// a scalable version.
func (t GTKTheme) GetIconSizes(icon string) (sizes []uint, isScalable bool) {
	// Get the array of sizes.
	cIcon := C.CString(icon)
	defer C.free(unsafe.Pointer(cIcon))
	cSizes := C.gtk_icon_theme_get_icon_sizes(t.theme, (*C.gchar)(cIcon))
	defer C.g_free(C.gpointer(cSizes))

	// Convert the data into Go types.
	for i := C.iter_array(cSizes); *i != 0; i = C.iter_array(i) {
		size := int(*i)
		if size == -1 {
			isScalable = true
		} else {
			sizes = append(sizes, size)
		}
	}

	return
}

// GetIconProperties attempts to determine all the properties of the given icon name.
func (t GTKTheme) GetIconProperties(icon string) GTKIconProperties {
	// Record the name of the icon.
	var properties GTKIconProperties
	properties.Name = icon

	return properties
}
