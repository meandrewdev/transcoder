package ffmpeg

import (
	"reflect"
	"testing"
)

func TestOptions_GetInputStrArguments(t *testing.T) {
	test := "test"
	zero := "0"
	tests := []struct {
		name string
		opts Options
		want []string
	}{
		{
			opts: Options{
				OutputFormat: &test,
				Aspect:       &test,
				InputFormat:  &test,
				InputSafe:    &zero,
				InputExtraArgs: map[string]interface{}{
					"-test": "test",
				},
			},
			want: []string{"-f", "test", "-safe", "0", "-test", "test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.opts.GetInputStrArguments(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Options.GetInputStrArguments() = %v, want %v", got, tt.want)
			}
		})
	}
}
