package ffmpeg

// Progress ...
type Progress struct {
	FramesProcessed string
	CurrentTime     string
	CurrentBitrate  string
	Progress        float64
	Speed           string
}

// GetPercent ...
func (p *Progress) GetPercent() float64 {
	return p.Progress
}
