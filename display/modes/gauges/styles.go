package gauges

import (
	"github.com/databeast/cyberhud/display/modes/gauges/source"
	"github.com/databeast/cyberhud/display/modes/gauges/styles"
	"github.com/databeast/cyberhud/display/style"
)

var gaugesRegistry = style.NewRegistry[source.GaugeSet, source.Policy](
	styles.GaugesDefaultStyle,
	styles.GaugesMonoSlowStyle,
	styles.GaugesMonoFastStyle,
	styles.GaugesGrayscaleSlowStyle,
	styles.GaugesGrayscaleFastStyle,
	styles.GaugesColorSlowStyle,
	styles.GaugesColorFastStyle,
)
