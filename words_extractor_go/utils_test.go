package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_removeCharacters(t *testing.T) {
	type args struct {
		input      string
		characters string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Characters removal: suffix",
			args: args{
				input:      "Załoenie;!",
				characters: ";!",
			},
			want: "Załoenie",
		},
		{
			name: "Characters removal: prefix",
			args: args{
				input:      ",#Załoenie",
				characters: ";!-,#",
			},
			want: "Załoenie",
		},
		{
			name: "Characters removal: both",
			args: args{
				input:      "-!Załoenie;!",
				characters: ";!-",
			},
			want: "Załoenie",
		},
	}
	assert := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeCharacters(tt.args.input, tt.args.characters)
			assert.Equal(tt.want, got, "Unexpected result in test: "+tt.name)
		})
	}
}

func Test_resultsArray_extractWords(t *testing.T) {
	type fields struct {
		Results []string
	}
	type args struct {
		s string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wants  []string
	}{
		{
			name: "Simple sentence",
			args: args{
				s: "Within this tutorial, we are going to look at how you can effectively read and write to files within your filesystem using the go programming language.",
			},
			wants: []string{"within", "this", "tutorial", "we", "are", "going", "to", "look", "at", "how", "you", "can", "effectively", "read", "and", "write", "to", "files", "within", "your", "filesystem", "using", "the", "go", "programming", "language"},
		},
		{
			name: "Multiline sentence",
			args: args{
				s: `The UK has recorded another five COVID deaths and 2,047 more cases in the latest daily figures.

				It compares with seven deaths and 1,907 cases this time last week, while the latest seven-day rolling average is 11.3 and 2,080.`,
			},
			wants: []string{"the", "uk", "has", "recorded", "another", "five", "covid", "deaths", "and", "2047", "more", "cases", "in", "the", "latest", "daily", "figures", "it", "compares", "with", "seven", "deaths", "and", "1907", "cases", "this", "time", "last", "week", "while", "the", "latest", "sevenday", "rolling", "average", "is", "113", "and", "2080"},
		},
	}
	assert := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &resultsArray{
				Results: tt.fields.Results,
			}
			r.extractWords(tt.args.s)
			assert.Equal(tt.wants, r.Results, "Unexpected result in test: "+tt.name)
		})
	}
}
