package render

import (
	"github.com/thomaso-mirodin/intmath/intgr"
)

// A cell whose size has been calculated to fit
type SizedCell struct {
	Cell  *Cell
	Width int
}

// Allocate a given width to a list of cells. The process is as follows:
// 1. Start out by satisfying every cell's minimum
// 2. If the maximum size is exceeded, start trimming shy cells, tiebreaking by
//    removing cells further to the right and obeying glue constraints
// 3. Divide any free space amongst greedy cells up to their maximums
// 4. Divide any free space amongst nongreedy cells up to their maximums
func AllocateCellSizes(cells []*Cell, width int) []*SizedCell {
	if len(cells) == 0 {
		return []*SizedCell{}
	}

	head := makeSizedCellList(cells)

	head.applyMinimumSizes(width)
	head = head.trimShyCells(width)

	head.expandCells(width, true)
	head.expandCells(width, false)

	return head.toSlice()
}

func makeSizedCellList(cells []*Cell) *scNode {
	// Create sized cell list
	head := &scNode{
		val: &SizedCell{Cell: cells[0]},
	}

	lastNode := head
	for _, c := range cells[1:] {
		node := &scNode{
			val: &SizedCell{Cell: c},
		}
		lastNode.next = node
		node.prev = lastNode
		lastNode = node
	}
	return head
}

// Doubly linked list used for allocation
type scNode struct {
	val  *SizedCell
	next *scNode
	prev *scNode
}

func (n *scNode) each(cb func(*scNode)) {
	for node := n; node != nil; node = node.next {
		cb(node)
	}
}

func (n *scNode) tail() *scNode {
	var lastNode *scNode
	n.each(func(node *scNode) {
		lastNode = node
	})
	return lastNode
}

func (n *scNode) eachReverse(cb func(*scNode)) {
	for node := n.tail(); node != nil; node = node.prev {
		cb(node)
	}
}

func (n *scNode) toSlice() []*SizedCell {
	result := []*SizedCell{}
	n.each(func(node *scNode) {
		result = append(result, node.val)
	})

	return result
}

func (n *scNode) delete() {
	if n.prev != nil {
		n.prev.next = n.next
	}
	if n.next != nil {
		n.next.prev = n.prev
	}
}

func (n *scNode) applyMinimumSizes(width int) {
	n.each(func(node *scNode) {
		node.val.Width = node.val.Cell.MinWidth
	})
}

// Remove shy cells and return the new list head (since the existing head may be
// trimmed)
func (n *scNode) trimShyCells(width int) *scNode {
	total := n.totalWidth()
	for total > width {
		toDelete := n.shyestNode()

		// Walk back to first glued node
		for {
			if toDelete.val.Cell.Glue && toDelete.prev != nil {
				toDelete = toDelete.prev
			} else {
				break
			}
		}

		// Start deleting nodes until we no longer see one that is glued
		for {
			total -= toDelete.val.Width
			// Unless deleting the last node, we also remove the padding
			if toDelete.next != nil || toDelete.prev != nil {
				total -= 1
			}

			// Special case where we're deleting the head
			if toDelete == n {
				n = n.next
			}

			toDelete.delete()
			if toDelete.next != nil && toDelete.next.val.Cell.Glue {
				toDelete = toDelete.next
			} else {
				break
			}
		}
	}

	return n
}

func (n *scNode) expandCells(width int, greedy bool) {
	freeSpace := width - n.totalWidth()
	for freeSpace > 0 {
		// Count cells that can expand
		expandableCount := 0
		n.each(func(node *scNode) {
			if val := node.val; val.Cell.Greedy == greedy && val.Width < val.Cell.MaxWidth {
				expandableCount++
			}
		})

		// Done expanding this type of cell
		if expandableCount == 0 {
			return
		}

		// Create dither generator to spread width out fairly
		allocator := dither(freeSpace, expandableCount)
		n.each(func(node *scNode) {
			val := node.val
			if val.Cell.Greedy == greedy && val.Width < val.Cell.MaxWidth {
				// Give the cell a fair allocation, or as much as it will take
				origWidth := val.Width
				val.Width += allocator()
				val.Width = intgr.Min(val.Cell.MaxWidth, val.Width)

				// Count off from the free space
				freeSpace -= (val.Width - origWidth)
			}
		})
	}
}

func (n *scNode) shyestNode() *scNode {
	shyest := n
	n.each(func(node *scNode) {
		// Using >= instead of > means we tiebreak to the right
		if node.val.Cell.Shyness >= shyest.val.Cell.Shyness {
			shyest = node
		}
	})
	return shyest
}

func (n *scNode) totalWidth() int {
	total := 0
	n.each(func(node *scNode) {
		total += node.val.Width + 1
	})
	// Delete padding from the last cell
	total -= 1
	return total
}
