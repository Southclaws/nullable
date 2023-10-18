package nullable

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Data struct {
	Value        string
	CanBeNull    Nullable[string]
	CannotBeNull *string
}

var pstring = "ptr"

// const example = `{"Value": "value","CanBeNull":null}`

func TestNullableMarshalJSON(t *testing.T) {
	for i, v := range []struct {
		want  string
		input Data
	}{
		{
			input: Data{},
			want:  `{"Value":"","CanBeNull":null,"CannotBeNull":null}`,
		},
		{
			input: Data{
				Value: "value",
			},
			want: `{"Value":"value","CanBeNull":null,"CannotBeNull":null}`,
		},
		{
			input: Data{
				Value:        "value",
				CanBeNull:    New[string]("value"),
				CannotBeNull: &pstring,
			},
			want: `{"Value":"value","CanBeNull":"value","CannotBeNull":"ptr"}`,
		},
		{
			input: Data{
				Value:     "value",
				CanBeNull: NewNull[string](),
			},
			want: `{"Value":"value","CanBeNull":null,"CannotBeNull":null}`,
		},
	} {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			r := require.New(t)

			got, err := json.Marshal(v.input)
			r.NoError(err)
			r.Equal(v.want, string(got))
		})
	}
}

func TestNullableUnmarshalJSON(t *testing.T) {
	for i, v := range []struct {
		want  Data
		input string
	}{
		{
			input: `{"Value":"","CanBeNull":null,"CannotBeNull":null}`,
			want: Data{
				CanBeNull: NewNull[string](),
			},
		},
		{
			input: `{"Value":"value","CanBeNull":null,"CannotBeNull":null}`,
			want: Data{
				Value:     "value",
				CanBeNull: NewNull[string](),
			},
		},
		{
			input: `{"Value":"value","CanBeNull":"value","CannotBeNull":"ptr"}`,
			want: Data{
				Value:        "value",
				CanBeNull:    New[string]("value"),
				CannotBeNull: &pstring,
			},
		},
		{
			input: `{"Value":"value","CanBeNull":null,"CannotBeNull":null}`,
			want: Data{
				Value:     "value",
				CanBeNull: NewNull[string](),
			},
		},
	} {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			r := require.New(t)

			var got Data
			err := json.Unmarshal([]byte(v.input), &got)
			r.NoError(err)
			r.Equal(v.want, got)
		})
	}
}

func TestCheckFunctions(t *testing.T) {
	a := assert.New(t)

	var (
		got1 Data
		err1 error
	)

	err1 = json.Unmarshal([]byte(`{"Value":"value","CanBeNull":null,"CannotBeNull":null}`), &got1)
	a.NoError(err1)
	a.True(got1.CanBeNull.IsNull())
	a.False(got1.CanBeNull.IsUnset())

	var (
		got2 Data
		err2 error
	)

	err2 = json.Unmarshal([]byte(`{"Value":"value","CannotBeNull":null}`), &got2)
	a.NoError(err2)
	a.False(got2.CanBeNull.IsNull())
	a.True(got2.CanBeNull.IsUnset())

	var (
		got3 Data
		err3 error
	)

	err3 = json.Unmarshal([]byte(`{"Value":"value","CanBeNull":"hello!","CannotBeNull":null}`), &got3)
	a.NoError(err3)
	a.False(got3.CanBeNull.IsNull())
	a.False(got3.CanBeNull.IsUnset())
	v, ok := got3.CanBeNull.Value()
	a.True(ok)
	a.Equal("hello!", v)
}
