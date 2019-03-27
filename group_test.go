package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroup(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrGroup1, exampleAttrGroup1Line},
		{exampleAttrGroup2, exampleAttrGroup2Line},
	}

	for i, u := range tests {
		actual := Group{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}

func TestGroup_LifeCycle(t *testing.T) {
	group := Group{}
	assert.Equal(t, 0, len(group.MIDs))
	assert.Nil(t, group.FirstMID())
	assert.Equal(t, -1, group.FindMID("0"))
	assert.Equal(t, false, group.HasMID("0"))
	assert.NotEqual(t, -1, group.AddMID("0"))
	assert.NotEqual(t, -1, group.AddMID("1"))
	assert.Equal(t, 2, len(group.MIDs))
	assert.Equal(t, 0, group.FindMID("0"))
	assert.Equal(t, true, group.HasMID("0"))
	assert.Equal(t, 1, group.FindMID("1"))
	assert.Equal(t, true, group.HasMID("1"))
	assert.Equal(t, 1, group.AddMID("1"))
	assert.NotNil(t, group.FirstMID())
	assert.Equal(t, 2, len(group.MIDs))
	assert.Equal(t, false, group.RemoveMID("2"))
	assert.Equal(t, true, group.RemoveMID("0"))
	assert.Equal(t, 1, len(group.MIDs))
	assert.Equal(t, 0, group.FindMID("1"))
	assert.Equal(t, true, group.HasMID("1"))
	assert.NotNil(t, group.FirstMID())
}
