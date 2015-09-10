package main

import (
	"reflect"
	"sort"
	"testing"
)

func TestEnvironmentStrings(t *testing.T) {
	pairs := []struct {
		input    environment
		expected []string
	}{
		{
			environment{},
			[]string{},
		},
		{
			environment{"FOO": "foo", "BAR": "bar"},
			[]string{"FOO=foo", "BAR=bar"},
		},
		{
			environment{"BAZ": "baz=baz", "quux": "quux"},
			[]string{"BAZ=baz=baz", "quux=quux"},
		},
	}
	for _, pair := range pairs {
		actual := pair.input.Strings()
		sort.Strings(actual)
		sort.Strings(pair.expected)
		equal := reflect.DeepEqual(actual, pair.expected)
		if !equal {
			t.Errorf("Expected %#v Got %#v", pair.expected, actual)
		}
	}
}

func TestTrimLastNewline(t *testing.T) {
	pairs := []struct {
		input, expected string
	}{
		{"", ""},
		{"\n", ""},
		{"\r\n", ""},
		{"\r", ""},
		{"hello ", "hello "},
		{"hello", "hello"},
		{"world\n", "world"},
		{"world\r\n", "world"},
		{"hello世界\n", "hello世界"},
		{"hello世界\r\n", "hello世界"},
		{"world\r", "world"},
		{"\nhello\nworld\n", "\nhello\nworld"},
		{"\r\nhello\nworld\n", "\r\nhello\nworld"},
		{"\nhello\nworld\r\n", "\nhello\nworld"},
	}
	for _, pair := range pairs {
		actual := trimLastNewline(pair.input)
		if actual != pair.expected {
			t.Errorf("Expected %v Got %v", pair.expected, actual)
		}
	}
}

func TestMakeEnvironmentMap(t *testing.T) {
	pairs := []struct {
		expected environment
		input    []string
	}{
		{
			environment{},
			[]string{},
		},
		{
			environment{"HELLO": "hello=hi", "WORLD": "world"},
			[]string{"HELLO=hello=hi", "WORLD=world"},
		},
	}
	for _, pair := range pairs {
		actual := makeEnvironmentMap(pair.input)
		if !reflect.DeepEqual(actual, pair.expected) {
			t.Errorf("Expected %v Got %v", actual, pair.expected)
		}
	}
}
