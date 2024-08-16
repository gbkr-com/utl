package utl

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWait(t *testing.T) {

	b := WithTimeout(time.Second)

	b.Done()
	err := b.Wait()
	assert.Nil(t, err, "Done before timeout")

	b = WithTimeout(100 * time.Millisecond)

	err = b.Wait()
	assert.NotNil(t, err, "Timeout before Done")
	assert.True(t, errors.Is(err, ErrTimeout))

	b = WithTimeout(time.Hour)

	b.Done()
	b.Done()
	err = b.Wait()
	assert.Nil(t, err, "More than one Done has no affect")

	err = b.Wait()
	assert.NotNil(t, err, "More than one Wait returns a timeout error")
	assert.True(t, errors.Is(err, ErrTimeout))

}
