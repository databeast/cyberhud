// Package image provides a display mode for displaying externally-supplied images.
package image

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
	"sync"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/borderframe"
	"github.com/databeast/cyberhud/runtime/action"
)

// Allowed style values for the image mode.
const (
	StyleDefault  = "default"
	StyleBordered = "bordered"
)

var allowedStyles = []string{StyleDefault, StyleBordered}

func init() {
	catalog.Register(catalog.Definition{
		ID:      "image",
		Title:   "Image",
		Scope:   "any",
		Summary: "Externally supplied image displayed full-panel with configurable scaling policy.",
		Order:   80,
		Options: []catalog.OptionDefinition{
			{Key: "fit", Type: "string", Summary: "How the image should fit the panel bounds.", Default: FitScale, Allowed: []string{FitTruncate, FitScale, FitStretch}},
			{Key: "style", Type: "string", Summary: "Visual presentation style (default or bordered).", Default: "", Allowed: allowedStyles},
		},
	})
	catalog.RegisterCommand(catalog.CommandDefinition{
		Verb:    "image",
		Summary: "Manage externally supplied images and fit policy.",
		Usage:   "image <set|policy|style|clear> ...",
		Handle:  HandleConsoleCommand,
	})
}

const (
	FitTruncate = "truncate"
	FitScale    = "scale"
	FitStretch  = "stretch"
)

var imageState = struct {
	sync.RWMutex
	current image.Image
	policy  Policy
	info    string
}{
	current: nil,
	policy:  DefaultPolicy(),
	info:    "(no image)",
}

// Policy controls how images are scaled to fit the panel.
type Policy struct {
	Fit   string // truncate, scale, stretch
	Style string // default, bordered
}

// DefaultPolicy returns the baseline image scaling policy.
func DefaultPolicy() Policy {
	return Policy{Fit: FitScale, Style: StyleDefault}
}

// NormalizePolicy ensures a policy has valid values.
func NormalizePolicy(p Policy) Policy {
	p.Fit = strings.ToLower(strings.TrimSpace(p.Fit))
	switch p.Fit {
	case FitTruncate, FitScale, FitStretch:
		// valid
	default:
		p.Fit = FitScale
	}

	p.Style = strings.ToLower(strings.TrimSpace(p.Style))
	validStyle := false
	for _, s := range allowedStyles {
		if p.Style == s {
			validStyle = true
			break
		}
	}
	if !validStyle {
		p.Style = StyleDefault
	}
	return p
}

// Handler implements action.Handler for image mode.
type Handler struct{}

func (Handler) HandleAction(act action.Action, cursor, itemCount int) action.Result {
	switch act {
	case action.Primary, action.Secondary:
		return action.Result{Navigate: "menu"}
	}
	return action.Result{}
}

// SetImage replaces the displayed image with a new one.
func SetImage(img image.Image, policy Policy) {
	imageState.Lock()
	defer imageState.Unlock()
	imageState.current = img
	imageState.policy = NormalizePolicy(policy)
	if img == nil {
		imageState.info = "(no image)"
	} else {
		bounds := img.Bounds()
		imageState.info = fmt.Sprintf("%dx%d (%s)", bounds.Dx(), bounds.Dy(), imageState.policy.Fit)
	}
}

// Current returns the currently displayed image and its scaling policy.
func Current() (image.Image, Policy) {
	imageState.RLock()
	defer imageState.RUnlock()
	return imageState.current, imageState.policy
}

// Info returns a brief text description of the current image.
func Info() string {
	imageState.RLock()
	defer imageState.RUnlock()
	return imageState.info
}

// Clear removes the current image.
func Clear() {
	imageState.Lock()
	defer imageState.Unlock()
	imageState.current = nil
	imageState.policy = DefaultPolicy()
	imageState.info = "(no image)"
}

// SetPolicy updates the scaling policy for the current image.
func SetPolicy(policy Policy) {
	imageState.Lock()
	defer imageState.Unlock()
	imageState.policy = NormalizePolicy(policy)
	if imageState.current != nil {
		bounds := imageState.current.Bounds()
		imageState.info = fmt.Sprintf("%dx%d (%s)", bounds.Dx(), bounds.Dy(), imageState.policy.Fit)
	}
}

// PolicySnapshot returns the current policy.
func PolicySnapshot() Policy {
	imageState.RLock()
	defer imageState.RUnlock()
	return imageState.policy
}

// BuildItems returns text items for the image mode panel.
func BuildItems() []string {
	info := Info()
	if info == "(no image)" {
		return []string{info}
	}
	return []string{
		fmt.Sprintf("Image: %s", info),
		"",
		"(full-panel rendering)",
	}
}

// Signature returns a stable identifier for change detection.
func Signature() string {
	imageState.RLock()
	defer imageState.RUnlock()
	if imageState.current == nil {
		return "image:none"
	}
	bounds := imageState.current.Bounds()
	return fmt.Sprintf("image:%dx%d:%s:%s", bounds.Dx(), bounds.Dy(), imageState.policy.Fit, imageState.policy.Style)
}

// HandleConsoleCommand handles the top-level "image" console verb.
func HandleConsoleCommand(args []string) string {
	if len(args) < 1 {
		return "ERR usage: image <set|policy|style|clear> ..."
	}
	sub := strings.ToLower(strings.TrimSpace(args[0]))
	switch sub {
	case "set":
		if len(args) < 3 {
			return "ERR usage: image set <file|base64> <data>"
		}
		source := strings.ToLower(strings.TrimSpace(args[1]))
		data := strings.TrimSpace(strings.Join(args[2:], " "))
		switch source {
		case "file":
			return setFromFile(data)
		case "base64":
			return setFromBase64(data)
		default:
			return "ERR source must be 'file' or 'base64'"
		}
	case "policy":
		if len(args) == 1 {
			return formatPolicyResponse(PolicySnapshot())
		}
		p := PolicySnapshot()
		for _, token := range args[1:] {
			kv := strings.SplitN(token, "=", 2)
			if len(kv) != 2 {
				return "ERR usage: image policy [fit=<truncate|scale|stretch>] [style=<default|bordered>]"
			}
			key := strings.ToLower(strings.TrimSpace(kv[0]))
			val := strings.TrimSpace(kv[1])
			switch key {
			case "fit":
				p.Fit = val
			case "style":
				lower := strings.ToLower(val)
				found := false
				for _, a := range allowedStyles {
					if lower == a {
						found = true
						break
					}
				}
				if !found {
					return fmt.Sprintf("ERR style: must be one of [%s]", strings.Join(allowedStyles, ", "))
				}
				p.Style = lower
			default:
				return fmt.Sprintf("ERR unknown image policy key %q", key)
			}
		}
		SetPolicy(p)
		return formatPolicyResponse(PolicySnapshot())
	case "style":
		if len(args) == 1 {
			// Query current style.
			p := PolicySnapshot()
			return fmt.Sprintf("OK image style=%s", p.Style)
		}
		val := strings.ToLower(strings.TrimSpace(args[1]))
		found := false
		for _, a := range allowedStyles {
			if val == a {
				found = true
				break
			}
		}
		if !found {
			return fmt.Sprintf("ERR style: must be one of [%s]", strings.Join(allowedStyles, ", "))
		}
		p := PolicySnapshot()
		p.Style = val
		SetPolicy(p)
		return fmt.Sprintf("OK image style=%s", val)
	case "clear":
		Clear()
		return "OK image cleared"
	default:
		return "ERR usage: image <set|policy|style|clear> ..."
	}
}

func setFromFile(path string) string {
	if path == "" {
		return "ERR image file path must not be empty"
	}
	f, err := os.Open(path)
	if err != nil {
		return fmt.Sprintf("ERR failed to open image file: %v", err)
	}
	defer f.Close()

	var img image.Image
	switch {
	case strings.HasSuffix(strings.ToLower(path), ".png"):
		decoded, err := png.Decode(f)
		if err != nil {
			return fmt.Sprintf("ERR failed to decode PNG: %v", err)
		}
		img = decoded
	case strings.HasSuffix(strings.ToLower(path), ".jpg"), strings.HasSuffix(strings.ToLower(path), ".jpeg"):
		decoded, err := jpeg.Decode(f)
		if err != nil {
			return fmt.Sprintf("ERR failed to decode JPEG: %v", err)
		}
		img = decoded
	default:
		return "ERR unsupported image format (PNG or JPEG required)"
	}

	SetImage(img, PolicySnapshot())
	bounds := img.Bounds()
	return fmt.Sprintf("OK image loaded %dx%d", bounds.Dx(), bounds.Dy())
}

func setFromBase64(data string) string {
	if data == "" {
		return "ERR base64 image data must not be empty"
	}
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return fmt.Sprintf("ERR failed to decode base64: %v", err)
	}

	var img image.Image
	pngImg, pngErr := png.Decode(bytes.NewReader(decoded))
	if pngErr == nil {
		img = pngImg
	} else {
		jpgImg, jpgErr := jpeg.Decode(bytes.NewReader(decoded))
		if jpgErr == nil {
			img = jpgImg
		}
	}

	if img == nil {
		return "ERR failed to decode image (PNG or JPEG required)"
	}

	SetImage(img, PolicySnapshot())
	bounds := img.Bounds()
	return fmt.Sprintf("OK image loaded %dx%d", bounds.Dx(), bounds.Dy())
}

func formatPolicyResponse(p Policy) string {
	return fmt.Sprintf("OK image policy fit=%s style=%s", p.Fit, p.Style)
}

// ViewData holds the output of BuildView for the registry to convert to State.
type ViewData struct {
	Title  string
	Items  []string
	Hint   string
	Static bool

	Sprites []widgets.Sprite
	// ScaledBounds is the computed scaling target rectangle for the image.
	// When the image should be drawn scaled, this provides the destination bounds.
	ScaledBounds image.Rectangle
}

// BorderBuilder is a function that builds border frame sprites for a given pixel bounds.
// Retained for backward compatibility with callers that supply a builder callback.
type BorderBuilder func(bounds image.Rectangle) []widgets.Sprite

// spriteRenderable wraps a pre-built *Sprite into a Renderable for use with the Compositor.
type spriteRenderable struct {
	sprite *widgets.Sprite
}

func (s *spriteRenderable) RenderFrame() *widgets.Sprite {
	return s.sprite
}

func BuildView(hints textlayout.TextHints) ViewData {
	p := PolicySnapshot()
	img, _ := Current()

	borderInset := 0
	hasBorder := p.Style == StyleBordered && hints.PixelWidth >= 16 && hints.PixelHeight >= 16
	if hasBorder {
		borderInset = 8
	}

	// Construct SuppressionContext from panel dimensions.
	ctx := widgets.SuppressionContext{
		AvailableWidth:  hints.PixelWidth,
		AvailableHeight: hints.PixelHeight,
	}

	// Border Compositor: border renders behind content (prepended).
	borderComp := widgets.NewCompositor(ctx)
	borderComp.AddIf(hasBorder, borderframe.New(borderframe.Config{
		Bounds: image.Rect(0, 0, hints.PixelWidth, hints.PixelHeight),
	}))

	// Content Compositor: image sprite renders on top of border.
	contentComp := widgets.NewCompositor(ctx)

	// Compute scaling target dimensions and add the image as a sprite.
	var scaledBounds image.Rectangle
	if img != nil {
		srcBounds := img.Bounds()
		srcW := srcBounds.Dx()
		srcH := srcBounds.Dy()

		// Determine target viewport (reduce by border inset if bordered).
		targetW := hints.PixelWidth - borderInset*2
		targetH := hints.PixelHeight - borderInset*2

		// Fall back to native dimensions if TextHints are zero/unavailable.
		if targetW <= 0 || targetH <= 0 {
			targetW = srcW
			targetH = srcH
		}

		switch p.Fit {
		case FitScale:
			scaled := ScaleToFit(srcW, srcH, targetW, targetH)
			scaledBounds = image.Rect(
				borderInset, borderInset,
				borderInset+scaled.Width, borderInset+scaled.Height,
			)
		case FitStretch:
			scaledBounds = image.Rect(
				borderInset, borderInset,
				borderInset+targetW, borderInset+targetH,
			)
		default: // truncate — draw at native size from origin
			scaledBounds = image.Rect(
				borderInset, borderInset,
				borderInset+srcW, borderInset+srcH,
			)
		}

		contentComp.Add(&spriteRenderable{sprite: &widgets.Sprite{
			Image:    img,
			Position: image.Pt(scaledBounds.Min.X, scaledBounds.Min.Y),
			Bounds:   scaledBounds,
			Label:    "image",
		}})
	}

	// Assemble final sprites: border behind, content on top.
	sprites := borderComp.Sprites()
	sprites = append(sprites, contentComp.Sprites()...)

	// Build items for text display fallback (shown when no image is loaded).
	info := Info()
	var items []string
	if info == "(no image)" {
		items = []string{info}
	}

	return ViewData{
		Title:        "",
		Items:        items,
		Hint:         "[K1]menu",
		Static:       true,
		Sprites:      sprites,
		ScaledBounds: scaledBounds,
	}
}

// ScaledDimensions holds the result of aspect-ratio-preserving scaling.
type ScaledDimensions struct {
	Width  int
	Height int
}

// ScaleToFit computes the largest dimensions that fit within targetW×targetH
// while preserving the aspect ratio of srcW×srcH. Returns the source dimensions
// if either source or target dimension is zero or negative.
func ScaleToFit(srcW, srcH, targetW, targetH int) ScaledDimensions {
	if srcW <= 0 || srcH <= 0 || targetW <= 0 || targetH <= 0 {
		if srcW <= 0 {
			srcW = 1
		}
		if srcH <= 0 {
			srcH = 1
		}
		return ScaledDimensions{Width: srcW, Height: srcH}
	}

	// Scale by the dimension that is the limiting factor.
	ratioW := float64(targetW) / float64(srcW)
	ratioH := float64(targetH) / float64(srcH)

	ratio := ratioW
	if ratioH < ratioW {
		ratio = ratioH
	}

	w := int(float64(srcW) * ratio)
	h := int(float64(srcH) * ratio)

	// Ensure at least 1×1.
	if w < 1 {
		w = 1
	}
	if h < 1 {
		h = 1
	}

	return ScaledDimensions{Width: w, Height: h}
}

// BorderInset returns the pixel inset for the current style (8 if bordered, 0 if default).
func BorderInset() int {
	p := PolicySnapshot()
	if p.Style == StyleBordered {
		return 8
	}
	return 0
}
