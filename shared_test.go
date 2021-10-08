package cache

import "testing"

func TestShrink(t *testing.T) {
	shared := newShared()
	hasher := newDefaultHash()
	shared.set("tenywen1", []byte("tenywen1value"), hasher)
	shared.set("tenywen2", []byte("tenywen2value"), hasher)
	shared.set("tenywen3", []byte("tenywen3value"), hasher)
	shared.set("tenywen4", []byte("tenywen4value"), hasher)
	shared.set("tenywen5", []byte("tenywen5value"), hasher)
	shared.set("tenywen6", []byte("tenywen6value"), hasher)
	shared.set("tenywen7", []byte("tenywen7value"), hasher)
	shared.set("tenywen8", []byte("tenywen8value"), hasher)
	shared.set("tenywen9", []byte("tenywen9value"), hasher)
	shared.set("tenywen10", []byte("tenywen10value"), hasher)
	shared.set("tenywen11", []byte("tenywen11value"), hasher)
	t.Log(shared.keys)
	shared.del("tenywen2", hasher)
	shared.del("tenywen3", hasher)
	shared.del("tenywen4", hasher)
	t.Log(shared.removes)
	t.Log(shared.keys)
	shared.shrink(hasher)
	b, err := shared.get("tenywen7", hasher)
	t.Log(string(b), err)
}
