package progressbar

import (
	"github.com/schollz/progressbar/v3"
	"time"
)

var bar *progressbar.ProgressBar

func Init(maxItems int64, description string) {
	if maxItems == 0 {
		bar = progressbar.NewOptions64(
			-1,
			progressbar.OptionSetDescription(description),
			progressbar.OptionSetWidth(10),
			progressbar.OptionThrottle(65*time.Millisecond),
			progressbar.OptionShowCount(),
			progressbar.OptionSpinnerType(14),
			progressbar.OptionFullWidth(),
			progressbar.OptionSetRenderBlankState(true),
		)
	} else {
		bar = progressbar.Default(maxItems, description)
	}
}

func ProgressOne() {
	if bar != nil {
		bar.Add(1)
	}
}

func Description(description string) {
	if bar != nil {
		bar.Describe(description)
	}
}
