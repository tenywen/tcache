package cache

import (
	"testing"
)

func TestAddBlock(t *testing.T) {
	sb := sortBlocks{}
	sb.add(block{
		si:    1,
		total: 10,
	})
	sb.add(block{
		si:    11,
		total: 20,
	})
	sb.add(block{
		si:    32,
		total: 2,
	})
	sb.add(block{
		si:    35,
		total: 1,
	})
}

func TestGetBlock(t *testing.T) {
	sb := sortBlocks{}
	sb.add(block{
		si:    1,
		total: 10,
	})
	sb.add(block{
		si:    11,
		total: 20,
	})
	sb.add(block{
		si:    32,
		total: 2,
	})
	sb.add(block{
		si:    35,
		total: 1,
	})
	b, ok := sb.getBlock(1)
	t.Log(b, ok)
	b, ok = sb.getBlock(1)
	t.Log(b, ok)
	b, ok = sb.getBlock(1)
	t.Log(b, ok)

}
