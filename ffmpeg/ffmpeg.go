package ffmpeg

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/meandrewdev/transcoder"
)

// Transcoder ...
type Transcoder struct {
	config           *Config
	output           string
	options          []string
	opts             transcoder.Options
	metadata         *transcoder.Metadata
	inputPipeReader  *io.ReadCloser
	outputPipeReader *io.ReadCloser
	inputPipeWriter  *io.WriteCloser
	outputPipeWriter *io.WriteCloser
}

// New ...
func New(cfg *Config) transcoder.Transcoder {
	return &Transcoder{config: cfg}
}

// Start ...
func (t *Transcoder) Start(opts transcoder.Options) ([]byte, error) {
	t.opts = opts

	defer t.closePipes()

	// Validates config
	if err := t.validate(); err != nil {
		return nil, err
	}

	/*
		// Get file metadata
		_, err := t.getMetadata()
		if err != nil {
			return nil, err
		}
	*/

	var err error

	args := opts.GetStrArguments()

	// Append output flag
	args = append(args, []string{t.output}...)

	// Initialize command
	cmd := exec.Command(t.config.FfmpegBinPath, args...)

	if t.config.Verbose {
		cmd.Stderr = os.Stdout
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return out, fmt.Errorf("Failed starting transcoding (%s) with args (%s) with error %s", t.config.FfmpegBinPath, args, err)
	}

	return out, err
}

func (t *Transcoder) Probe(input string) (result transcoder.Metadata, err error) {
	_, err = t.getMetadata(input)
	if err != nil {
		return
	}
	result = *t.metadata
	return
}

func (t *Transcoder) getMetadata(input string) (metadata *transcoder.Metadata, err error) {

	if t.config.FfprobeBinPath != "" {
		var outb, errb bytes.Buffer

		args := []string{"-i", input, "-print_format", "json", "-show_format", "-show_streams", "-show_error"}

		cmd := exec.Command(t.config.FfprobeBinPath, args...)
		cmd.Stdout = &outb
		cmd.Stderr = &errb

		err := cmd.Run()
		if err != nil {
			return nil, fmt.Errorf("error executing (%s) with args (%s) | error: %s | message: %s %s", t.config.FfprobeBinPath, args, err, outb.String(), errb.String())
		}

		if err = json.Unmarshal([]byte(outb.String()), &metadata); err != nil {
			return nil, err
		}

		t.metadata = metadata

		return metadata, nil
	}

	return nil, errors.New("ffprobe binary not found")
}

// Output ...
func (t *Transcoder) Output(arg string) transcoder.Transcoder {
	t.output = arg
	return t
}

// OutputPipe ...
func (t *Transcoder) OutputPipe(w *io.WriteCloser, r *io.ReadCloser) transcoder.Transcoder {
	if &t.output == nil {
		t.outputPipeWriter = w
		t.outputPipeReader = r
	}
	return t
}

// validate ...
func (t *Transcoder) validate() error {
	if t.config.FfmpegBinPath == "" {
		return errors.New("ffmpeg binary path not found")
	}

	if len(t.opts.GetInputs()) == 0 {
		return errors.New("missing input option")
	}

	if t.output == "" {
		return errors.New("missing output option")
	}

	return nil
}

// closePipes Closes pipes if opened
func (t *Transcoder) closePipes() {
	if t.inputPipeReader != nil {
		ipr := *t.inputPipeReader
		ipr.Close()
	}

	if t.outputPipeWriter != nil {
		opr := *t.outputPipeWriter
		opr.Close()
	}
}
