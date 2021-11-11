package cache

import "sort"

type block struct {
	start int
	size  int
}

type sortBlocks []block

type remove struct {
	blocks sortBlocks
}

func (sb sortBlocks) Len() int {
	return len(sb)
}

func (sb sortBlocks) Swap(i, j int) {
	sb[i], sb[j] = sb[j], sb[i]
}

func (sb sortBlocks) Less(i, j int) bool {
	return sb[i].size < sb[j].size
}

func (sb *sortBlocks) add(start int, size int) {
	*sb = append(*sb, block{start: start, size: size})
}

func (sb *sortBlocks) getBlock(size int) int {
	sort.Sort(*sb)
	length := sb.Len()
	index := sort.Search(length, func(i int) bool {
		return (*sb)[i].size >= size
	})

	if index >= length || (*sb)[index].size < size {
		return undefined
	}

	(*sb)[index].size -= size
	if (*sb)[index].size == 0 && length != 1 {
		(*sb)[index] = (*sb)[length-1]
		(*sb) = (*sb)[:length-1]
	}

	return index
}

func (remove *remove) add(start, size int) {
	remove.blocks.add(start, size)
}

func (remove *remove) getBlock(size int) int {
	if len(remove.blocks) == 0 {
		return undefined
	}

	return remove.blocks.getBlock(size)
}
