package ggrender

import "image/color"

// FillRect draws a filled rectangle at (x, y) with the given width, height, and color.
func (c *Canvas) FillRect(x, y, w, h float64, col color.RGBA) {
	c.ctx.DrawRectangle(x, y, w, h)
	c.ctx.SetColor(col)
	c.ctx.Fill()
}

// FillRoundedRect draws a filled rounded rectangle at (x, y) with the given
// width, height, corner radius, and color.
func (c *Canvas) FillRoundedRect(x, y, w, h, radius float64, col color.RGBA) {
	c.ctx.DrawRoundedRectangle(x, y, w, h, radius)
	c.ctx.SetColor(col)
	c.ctx.Fill()
}

// FillCircle draws a filled circle at center (cx, cy) with the given radius and color.
func (c *Canvas) FillCircle(cx, cy, radius float64, col color.RGBA) {
	c.ctx.DrawCircle(cx, cy, radius)
	c.ctx.SetColor(col)
	c.ctx.Fill()
}

// StrokeLine draws a stroked line from (x1, y1) to (x2, y2) with the given
// stroke width and color.
func (c *Canvas) StrokeLine(x1, y1, x2, y2, strokeWidth float64, col color.RGBA) {
	c.ctx.DrawLine(x1, y1, x2, y2)
	c.ctx.SetLineWidth(strokeWidth)
	c.ctx.SetColor(col)
	c.ctx.Stroke()
}

// FillArc draws a filled arc (pie wedge) centered at (cx, cy) with the given
// radius, start angle, and end angle (in radians), filled with the specified color.
func (c *Canvas) FillArc(cx, cy, radius, startAngle, endAngle float64, col color.RGBA) {
	c.ctx.DrawArc(cx, cy, radius, startAngle, endAngle)
	c.ctx.LineTo(cx, cy)
	c.ctx.ClosePath()
	c.ctx.SetColor(col)
	c.ctx.Fill()
}
