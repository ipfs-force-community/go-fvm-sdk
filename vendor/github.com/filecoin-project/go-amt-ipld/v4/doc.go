/*
Package amt provides a reference implementation of the IPLD AMT (Array Mapped
Trie) used in the Filecoin blockchain.

The AMT algorithm is similar to a HAMT
https://en.wikipedia.org/wiki/Hash_array_mapped_trie but instead presents an
array-like interface where the indexes themselves form the mapping to nodes in
the trie structure. An AMT is suitable for storing sparse array data as a
minimum amount of intermediate nodes are required to address a small number of
entries even when their indexes span a large distance. AMT is also a suitable
means of storing non-sparse array data as required, with a small amount of
storage and algorithmic overhead required to handle mapping that assumes that
some elements within any range of data may not be present.

Algorithm Overview

The AMT algorithm produces a tree-like graph, with a single root node
addressing a collection of child nodes which connect downward toward leaf nodes
which store the actual entries. No terminal entries are stored in intermediate
elements of the tree, unlike in a HAMT. We can divide up the AMT tree structure
into "levels" or "heights", where a height of zero contains the terminal
elements, and the maximum height of the tree contains the single root node.
Intermediate nodes are used to span across the range of indexes.

Any AMT instance uses a fixed "width" that is consistent across the tree's
nodes. An AMT's "bitWidth" dictates the width, or maximum-brancing factor
(arity) of the AMT's nodes by determining how many bits of the original index
are used to determine the index at any given level. A bitWidth of 3 (the
default for this implementation) can generate indexes in the range of 0 to
(2^3)-1=7, i.e. a "width" of 8. In practice, this means that an AMT with a
bitWidth of 3 has a branching factor of _between 1 and 8_ for any node in the
structure.

Considering the minimal case: a minimal AMT contains a single node which serves
as both the root and the leaf node and can hold zero or more elements
(an empty AMT is possible, although a special-case, and consists of a
zero-length root). This minimal AMT can store array indexes from 0 to width-1
(8 for the default bitWidth of 3) without requiring the addition of additional
nodes. Attempts to add additional indexes beyond width-1 will result in
additional nodes being added and a tree structure in order to address the new
elements. The minimal AMT node is said to have a height of 0. Every node in an
AMT has a height that indicates its distance from the leaf nodes. All leaf
nodes have a height of 0. The height of the root node dictates the overall
height of the entire AMT. In the case of the minimal AMT, this is 0.

Elements are stored in a compacted form within nodes, they are
"position-mapped" by a bitmap field that is stored with the node. The bitmap is
a simple byte array, where each bit represents an element of the data that can
be stored in the node. With a width of 8, the bitmap is a single byte and up to
8 elements can be stored in the node. The data array of a node _only stores
elements that are present in that node_, so the array is commonly shorter than
the maximum width. An empty AMT is a special-case where the single node can
have zero elements, therefore a zero-length data array and a bitmap of `0x00`.
In all other cases, the data array must have between 1 and width elements.

Determining the position of an index within the data array requires counting
the number of set bits within the bitmap up to the element we are concerned
with. If the bitmap has bits 2, 4 and 6 set, we can see that only 3 of the bits
are set so our data array should hold 3 elements. To address index 4, we know
that the first element will be index 2 and therefore the second will hold index
4. This format allows us to store only the elements that are set in the node.

Overflow beyond the single node AMT by adding an index beyond width-1 requires
an increase in height in order to address all elements. If an element in the
range of width to (width*2)-1 is added, a single additional height is required
which will result in a new root node which is used to address two consecutive
leaf nodes. Because we have an arity of up to width at any node, the addition
of indexes in the range of 0 to (width^2)-1 will still require only the
addition of a single additional height above the leaf nodes, i.e. height 1.

From the width of an AMT we can derive the maximum range of indexes that can be
contained by an AMT at any given `height` with the formula width^(height+1)-1.
e.g. an AMT with a width of 8 and a height of 2 can address indexes 0 to
8^(2+1)-1=511. Incrementing the height doubles the range of indexes that can be
contained within that structure.

Nodes above height 0 (non-leaf nodes) do not contain terminal elements, but
instead, their data array contains links to child nodes. The index compaction
using the bitmap is the same as for leaf nodes, so each non-leaf node only
stores as many links as it has child nodes.

Because additional height is required to address larger indexes, even a
single-element AMT will require more than one node where the index is greater
than the width of the AMT. For a width of 8, indexes 8 to 63 require a height
of 1, indexes 64 to 511 require a height of 2, indexes 512 to 4095 require a
height of 3, etc.

Retrieving elements from the AMT requires extracting only the portion of the
requested index that is required at each height to determine the position in
the data array to navigate into. When traversing through the tree, we only need
to select from indexes 0 to width-1. To do this, we take log2(width) bits from
the index to form a number that is between 0 and width-1. e.g. for a width of
8, we only need 3 bits to form a number between 0 and 7, so we only consume
3 bits per level of the AMT as we traverse. A simple method to calculate this
at any height in the AMT (assuming bitWidth of 3, i.e. a width of 8) is:

1. Calculate the maximum number of nodes (not entries) that may be present in
an sub-tree rooted at the current height. width^height provides this number.
e.g. at height 0, only 1 node can be present, but at height 3, we may have a
tree of up to 512 nodes (storing up to 8^(3+1)=4096 entries).

2. Divide the index by this number to find the index for this height. e.g. an
index of 3 at height 0 will be 3/1=3, or an index of 20 at height 1 will be
20/8=2.

3. If we are at height 0, the element we want is at the data index,
position-mapped via the bitmap.

4. If we are above height 0, we need to navigate to the child element at the
index we calculated, position-mapped via the bitmap. When traversing to the
child, we discard the upper portion of the index that we no longer need.
This can be achieved by a mod operation against the number-of-nodes value.
e.g. an index of 20 at height 1 requires navigation to the element at
position 2, when moving to that element (which is height 0), we truncate the
index with 20%8=4, at height 0 this index will be the index in our data
array (position-mapped via the bitmap).

In this way, each sub-tree root consumes a small slice, log2(width) bits long,
of the original index.

Adding new elements to an AMT may require up to 3 steps:

1. Increasing the height to accommodate a new index if the current height is
not sufficient to address the new index. Increasing the height requires turning
the current root node into an intermediate and adding a new root which
links to the old (repeated until the required height is reached).

2. Adding any missing intermediate and leaf nodes that are required to address
the new index. Depending on the density of existing indexes, this may require
the addition of up to height-1 new nodes to connect the root to the required
leaf. Sparse indexes will mean large gaps in the tree that will need filling to
address new, equally sparse, indexes.

3. Setting the element at the leaf node in the appropriate position in the data
array and setting the appropriate bit in the bitmap.

Removing elements requires a reversal of this process. Any empty node (other
than the case of a completely empty AMT) must be removed and its parent should
have its child link removed. This removal may recurse up the tree to remove
many unnecessary intermediate nodes. The root node may also be removed if the
current height is no longer necessary to contain the range of indexes still in
the AMT. This can be easily determined if _only_ the first bit of the root's
bitmap is set, meaning only the left-most is present, which will become the
new root node (repeated until the new root has more than the first bit set or
height of 0, the single-node case).

Further Reading

See https://github.com/ipld/specs/blob/master/data-structures/hashmap.md for a
description of a HAMT algorithm. And
https://github.com/ipld/specs/blob/master/data-structures/vector.md for a
description of a similar algorithm to an AMT that doesn't support internal node
compression and therefore doesn't support sparse arrays.

Usage Considerations

Unlike a HAMT, the AMT algorithm doesn't benefit from randomness introduced by
a hash algorithm. Therefore an AMT used in cases where user-input can
influence indexes, larger-than-necessary tree structures may present risks as
well as the challenge imposed by having a strict upper-limit on the indexes
addressable by the AMT. A width of 8, using 64-bit integers for indexing,
allows for a tree height of up to 64/log2(8)=21 (i.e. a width of 8 has a
bitWidth of 3, dividing the 64 bits of the uint into 21 separate per-height
indexes). Careful placement of indexes could create extremely sub-optimal forms
with large heights connecting leaf nodes that are sparsely packed. The overhead
of the large number of intermediate nodes required to connect leaf nodes in
AMTs that contain high indexes can be abused to create perverse forms that
contain large numbers of nodes to store a minimal number of elements.

Minimal nodes will be created where indexes are all in the lower-range. The
optimal case for an AMT is contiguous index values starting from zero. As
larger indexes are introduced that span beyond the current maximum, more nodes
are required to address the new nodes _and_ the existing lower index nodes.
Consider a case where a width=8 AMT is only addressing indexes less than 8 and
requiring a single height. The introduction of a single index within 8 of the
maximum 64-bit unsigned integer range will require the new root to have a
height of 21 and have enough connecting nodes between it and both the existing
elements and the new upper index. This pattern of behavior may be acceptable if
there is significant density of entries under a particular maximum index.

There is a direct relationship between the sparseness of index values and the
number of nodes required to address the entries. This should be the key
consideration when determining whether an AMT is a suitable data-structure for
a given application.

*/
package amt
