package iterators_test

import (
	"errors"
	"testing"
	"time"

	"github.com/adamluzsi/frameless"

	"github.com/stretchr/testify/require"

	"github.com/adamluzsi/frameless/iterators"
	"github.com/adamluzsi/frameless/iterators/iterateover"
)

func ExampleNewPipe() (frameless.Iterator, *iterators.PipeSender) {
	return iterators.NewPipe()
}

func TestNewPipe_SimpleFeedScenario(t *testing.T) {
	t.Parallel()

	r, w := iterators.NewPipe()

	var expected Entity = Entity{Text: "hitchhiker's guide to the galaxy"}
	var actually Entity

	go func() {
		defer w.Close()
		require.Nil(t, w.Send(&expected))
	}()

	require.True(t, r.Next())            // first next should return the value mean to be sent
	require.Nil(t, r.Decode(&actually))  // and decode it
	require.Equal(t, expected, actually) // the exactly same value passed in
	require.False(t, r.Next())           // no more values left, sender done with its work
	require.Nil(t, r.Err())              // No error sent so there must be no err received
	require.Nil(t, r.Close())            // Than I release this resource too
}

func TestNewPipe_FetchWithCollectAll(t *testing.T) {
	t.Parallel()

	r, w := iterators.NewPipe()

	var actually []*Entity
	var expected []*Entity = []*Entity{
		&Entity{Text: "hitchhiker's guide to the galaxy"},
		&Entity{Text: "The 5 Elements of Effective Thinking"},
		&Entity{Text: "The Art of Agile Development"},
		&Entity{Text: "The Phoenix Project"},
	}

	go func() {
		defer w.Close()

		for _, e := range expected {
			w.Send(e)
		}
	}()

	require.Nil(t, iterateover.AndCollectAll(r, &actually)) // When I collect everything with Collect All and close the resource
	require.True(t, len(actually) > 0)                      // the collection includes all the sent values
	require.Equal(t, expected, actually)                    // which is exactly the same that mean to be sent.
}

func TestNewPipe_ReceiverCloseResourceEarly_FeederNoted(t *testing.T) {
	t.Parallel()

	// skip when only short test expected
	// this test is slow because it has sleep in it
	//
	// This could be fixed by implementing extra logic in the Pipe iterator,
	// but that would be overengineering because after an iterator is closed,
	// it is highly unlikely that next value and decode will be called.
	// So this is only for the sake of describing the iterator behavior in this edge case
	if testing.Short() {
		t.Skip()
	}

	r, w := iterators.NewPipe()

	go func() {
		defer w.Close()

		require.Equal(t, iterators.ErrClosed, w.Send(&Entity{Text: "hitchhiker's guide to the galaxy"}))
	}()

	// normally next should not be called after a Close, but in the test I have to define the behavior
	// so in order to prevent overengineering in sender Send method, I place a sleep here to force thick the scheduler in favor of done channel
	time.Sleep(10 * time.Millisecond)

	require.Nil(t, r.Close()) // I release the resource,
	// for example something went wrong during the processing on my side (receiver) and I can't continue work,
	// but I want to note this to the sender as well
	require.Nil(t, r.Close()) // multiple times because defer ensure and other reasons

	require.False(t, r.Next())            // the sender is notified about this and stopped sending messages
	require.Error(t, r.Decode(&Entity{})) // and for some reason when I want to decode, it tells me the iterator closed. It was the sender who close it
}

func TestNewPipe_SenderSendErrorAboutProcessingToReceiver_ReceiverNotified(t *testing.T) {
	t.Parallel()

	expected := errors.New("Boom!")

	r, w := iterators.NewPipe()

	go func() {
		require.Nil(t, w.Send(&Entity{Text: "hitchhiker's guide to the galaxy"}))
		require.Nil(t, w.Error(expected))
		require.Nil(t, w.Close())
	}()

	require.True(t, r.Next())           // everything goes smoothly, I'm notified about next value
	require.Nil(t, r.Decode(&Entity{})) // I even able to decode it as well
	require.False(t, r.Next())          // Than the sender is notify me that I will not receive any more value
	require.Equal(t, expected, r.Err()) // Also tells me that something went wrong during the processing
	require.Nil(t, r.Close())           // I release the resource because than and go on
	require.Equal(t, expected, r.Err()) // The last error should be available later
}
