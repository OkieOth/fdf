package progressbar

import (
	"github.com/schollz/progressbar/v3"
)

var bar *progressbar.ProgressBar

func Init(maxItems int64, description string) {
	bar = progressbar.Default(maxItems)
	bar.Describe(description)
}

func ProgressOne() {
	bar.Add(1)
}

func Description(description string) {
	bar.Describe(description)
}
