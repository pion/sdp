package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSsrcGroup(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrSsrcGroup1, exampleAttrSsrcGroup1Line},
		{exampleAttrSsrcGroup2, exampleAttrSsrcGroup2Line},
	}

	for i, u := range tests {
		actual := SsrcGroup{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}

func TestSsrcGroup_LifeCycle(t *testing.T) {
	group := SsrcGroup{}
	assert.Equal(t, 0, len(group.SSRCs))
	assert.Nil(t, group.FirstSSRC())
	assert.Equal(t, -1, group.HasSSRC("0"))
	assert.NotEqual(t, -1, group.AddSSRC("0"))
	assert.NotEqual(t, -1, group.AddSSRC("1"))
	assert.Equal(t, 2, len(group.SSRCs))
	assert.Equal(t, 0, group.HasSSRC("0"))
	assert.Equal(t, 1, group.HasSSRC("1"))
	assert.Equal(t, 1, group.AddSSRC("1"))
	assert.NotNil(t, group.FirstSSRC())
	assert.Equal(t, 2, len(group.SSRCs))
	assert.Equal(t, false, group.RemoveSSRC("2"))
	assert.Equal(t, true, group.RemoveSSRC("0"))
	assert.Equal(t, 1, len(group.SSRCs))
	assert.Equal(t, 0, group.HasSSRC("1"))
	assert.NotNil(t, group.FirstSSRC())
}
