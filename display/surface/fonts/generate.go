package font

//go:generate go run ../../../buildtools/fontgen -clone spleen -bdf spleen-5x8.bdf -out gen_spleen_5x8.go -id spleen-5x8 -pkg font -width 5 -height 8 -advance 6 -rowheight 10
//go:generate go run ../../../buildtools/fontgen -clone spleen -bdf spleen-6x12.bdf -out gen_spleen_6x12.go -id spleen-6x12 -pkg font -width 6 -height 12 -advance 7 -rowheight 14
//go:generate go run ../../../buildtools/fontgen -clone spleen -bdf spleen-8x16.bdf -out gen_spleen_8x16.go -id spleen-8x16 -pkg font -width 8 -height 16 -advance 9 -rowheight 18
//go:generate go run ../../../buildtools/fontgen -clone spleen -bdf spleen-12x24.bdf -out gen_spleen_12x24.go -id spleen-12x24 -pkg font -width 12 -height 24 -advance 13 -rowheight 26
//go:generate go run ../../../buildtools/fontgen -clone spleen -bdf spleen-16x32.bdf -out gen_spleen_16x32.go -id spleen-16x32 -pkg font -width 16 -height 32 -advance 17 -rowheight 34
//go:generate go run ../../../buildtools/fontgen -clone spleen -bdf spleen-32x64.bdf -out gen_spleen_32x64.go -id spleen-32x64 -pkg font -width 32 -height 64 -advance 33 -rowheight 66
//go:generate go run ../../../buildtools/fontgen -download https://github.com/the-moonwitch/Cozette/releases/download/v.1.30.0/cozette.bdf -bdf cozette.bdf -out gen_cozette_6x13.go -id cozette-6x13 -pkg font -width 6 -height 13 -advance 7 -rowheight 15
//go:generate go run ../../../buildtools/fontgen -clone terminus -bdf ter-u12n.bdf -out gen_terminus_6x12.go -id terminus-6x12 -pkg font -width 6 -height 12 -advance 7 -rowheight 14
//go:generate go run ../../../buildtools/fontgen -clone terminus -bdf ter-u14n.bdf -out gen_terminus_8x14.go -id terminus-8x14 -pkg font -width 8 -height 14 -advance 9 -rowheight 16
//go:generate go run ../../../buildtools/fontgen -clone terminus -bdf ter-u16n.bdf -out gen_terminus_8x16.go -id terminus-8x16 -pkg font -width 8 -height 16 -advance 9 -rowheight 18
//go:generate go run ../../../buildtools/fontgen -clone terminus -bdf ter-u18n.bdf -out gen_terminus_10x18.go -id terminus-10x18 -pkg font -width 10 -height 18 -advance 11 -rowheight 20
//go:generate go run ../../../buildtools/fontgen -clone terminus -bdf ter-u20n.bdf -out gen_terminus_10x20.go -id terminus-10x20 -pkg font -width 10 -height 20 -advance 11 -rowheight 22
//go:generate go run ../../../buildtools/fontgen -clone terminus -bdf ter-u22n.bdf -out gen_terminus_11x22.go -id terminus-11x22 -pkg font -width 11 -height 22 -advance 12 -rowheight 24
//go:generate go run ../../../buildtools/fontgen -clone terminus -bdf ter-u24n.bdf -out gen_terminus_12x24.go -id terminus-12x24 -pkg font -width 12 -height 24 -advance 13 -rowheight 26
//go:generate go run ../../../buildtools/fontgen -clone terminus -bdf ter-u28n.bdf -out gen_terminus_14x28.go -id terminus-14x28 -pkg font -width 14 -height 28 -advance 15 -rowheight 30
//go:generate go run ../../../buildtools/fontgen -clone terminus -bdf ter-u32n.bdf -out gen_terminus_16x32.go -id terminus-16x32 -pkg font -width 16 -height 32 -advance 17 -rowheight 34
//go:generate go run ../../../buildtools/fontgen -ttf "../../../buildtools/fontgen/vendor/matrix-code-font/Matrix Code Font.ttf" -out gen_matrix_10x10.go -id matrix-10x10 -pkg font -width 10 -height 10 -advance 11 -rowheight 12 -targetheight 10 -ranges "33-126,65382-65437"
//go:generate go run ../../../buildtools/fontgen -ttf "../../../buildtools/fontgen/vendor/matrix-code-font/Matrix Code Font.ttf" -out gen_matrix_code.go -id matrix-code -pkg font -width 13 -height 12 -advance 14 -rowheight 14 -targetheight 12 -ranges "33-126,65382-65437"
//go:generate go run ../../../buildtools/gen-icons -codepoints ../../../buildtools/gen-icons/.cache/codepoints -ttf ../../../buildtools/gen-icons/.cache/MaterialSymbolsOutlined.ttf -constout gen_material_icons_constants.go -faceout gen_material_icons_24.go -pkg font -targetheight 24 -faceid material-icons-24
