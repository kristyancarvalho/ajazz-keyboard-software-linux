package keyboard

const (
	VendorID         uint16 = 0x258A
	ProductID        uint16 = 0x00C7
	LegacyVendorID   uint16 = 0x0C45
	LegacyProductID  uint16 = 0x8009
	ReportSize              = 64
	ReportID         byte   = 0x04
	NumChunks               = 9
	ImageChunkSize          = 4123
	CmdSetMode       byte   = 0x02
	CmdCommit        byte   = 0xF5
	CmdSetPerKey     byte   = 0x18
	CmdImageChunk    byte   = 0x13
	CmdImageFinal    byte   = 0xF0
	CmdSetMacro      byte   = 0x30
	TFTWidth                = 128
	TFTHeight               = 128
	TFTBytesPerPixel        = 2
	ThemeBytesPerKey        = 3
	KeyCount                = 81
	MaxMacroSequence        = ReportSize - 4
)
