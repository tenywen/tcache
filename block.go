package cache

type block struct {
	total uint16
	kl    int16
	vl    int16
	s     int32
}

type sortBlocks []block

func (sb sortBlocks) Len() int {
	return len(sb)
}

func (sb sortBlocks) Less(i, j int) bool {
	return sb[i].total < sb[j].total
}

func (sb sortBlocks) Swap(i, j int) {
	sb[i], sb[j] = sb[j], sb[i]
}
