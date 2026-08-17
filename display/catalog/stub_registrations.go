package catalog

func init() {
	RegisterSnapshotter("testpattern", StubSnapshotter{})
	RegisterSnapshotter("testfonts", StubSnapshotter{})
	RegisterSnapshotter("testicons", StubSnapshotter{})
	RegisterSnapshotter("testwidgets", StubSnapshotter{})
}
