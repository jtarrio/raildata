package util

import (
	"github.com/fatih/color"
	"github.com/jtarrio/raildata"
)

func FindStation(codeOrName string) (*raildata.StationCode, bool) {
	station, found := raildata.FindStation().WithCode(raildata.StationCode(codeOrName)).WithName(codeOrName).Search()
	if found {
		return &station.Code, true
	}
	return nil, false
}

func FindLine(codeOrName string) (*raildata.LineCode, bool) {
	line, found := raildata.FindLine().WithCode(raildata.LineCode(codeOrName)).WithName(codeOrName).Search()
	if found {
		return &line.Code, true
	}
	return nil, false
}

func HtmlColors(fg *raildata.Color, bg *raildata.Color) *color.Color {
	if fg != nil && bg != nil {
		fR, fG, fB := fg.RGB()
		bR, bG, bB := bg.RGB()
		return color.RGB(fR, fG, fB).AddBgRGB(bR, bG, bB)
	} else if fg != nil {
		fR, fG, fB := fg.RGB()
		return color.RGB(fR, fG, fB)
	} else if bg != nil {
		bR, bG, bB := bg.RGB()
		return color.BgRGB(bR, bG, bB)
	} else {
		return color.New(color.Reset)
	}
}
