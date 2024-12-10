package sfinder

type SerialType string

const (
	CH34x SerialType = "ch34x"
	FT2x SerialType = "ft23x"
	CP2x SerialType = "cp21x"
	PL23X SerialType = "pl23x"
)

var vendor = make(map[SerialType][]string)

func init() {
	vendor[CH34x] = []string{"1a86"}
	vendor[FT2x] = []string{"0403", "165c"}
	vendor[CP2x] = []string{"10c4"}
	vendor[PL23X] = []string{"067b"}
}

