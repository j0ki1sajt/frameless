package iterators_test

import (
	"testing"

	"github.com/adamluzsi/frameless/iterators"
	"github.com/stretchr/testify/require"
)

func TestNextDecoder_IteratorGiven_ValidDecoderReturnedThanCanDecodeTheFirstValueFromTheIterator(t *testing.T) {
	t.Parallel()

	var expected int = 42
	var actually int

	i := iterators.NewForSlice([]int{expected, 4, 2})
	defer i.Close()

	if err := iterators.NextDecoder(i).Decode(&actually); err != nil {
		t.Fatal(err)
	}

	require.Equal(t, expected, actually)
}
