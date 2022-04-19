package utils

import (
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl"
	"log"
)

// Block is a wrapper for fbdl.Block. It is needed, because in some languages
// or implementations the hierarchy must be flattened. In such case, there
// is a need to resolve name conflicts.
type Block struct {
	Name      string
	NameLevel int
	Path      []string
	Block     *fbdl.Block
}

// Rename renames Block based on the NameLevel and Path.
// NameLevel indicates how many Path elements should be included in the Name.
// NameLevel field is incremented internally. Path elements are taken from the end.
func (b *Block) Rename() {
	if b.NameLevel < len(b.Path) {
		b.NameLevel += 1
	}

	name := b.Path[len(b.Path)-1]
	for i := 1; i < b.NameLevel; i++ {
		name = b.Path[len(b.Path)-i-1] + "_" + name
	}

	b.Name = name
}

func CollectBlocks(blk *fbdl.Block, blocks []Block, path []string) []Block {
	if blocks == nil {
		blocks = []Block{Block{
			Name: "Main", NameLevel: 1, Path: []string{"Main"}, Block: blk},
		}
		path = append(path, "Main")
	} else {
		p := make([]string, len(path))
		n := copy(p, path)
		if n != len(path) {
			log.Fatalf("utils: colllect blocks: copying block path failed, copied %d, expected %d", n, len(path))
		}

		ent := Block{Name: blk.Name, Path: p, Block: blk}
		blocks = append(blocks, ent)
	}

	for _, b := range blk.Subblocks {
		path = append(path, b.Name)
		blocks = CollectBlocks(b, blocks, path)
		path = path[:len(path)-1]
	}

	return blocks
}

func ResolveBlockNameConflicts(blocks []Block) {
	for i, _ := range blocks[:len(blocks)-1] {
		conflicts := []*Block{&blocks[i]}
		for j, _ := range blocks[i+1:] {
			if blocks[i].Name == blocks[i+j+1].Name {
				conflicts = append(conflicts, &blocks[i+j+1])
			}
		}
		if len(conflicts) == 1 {
			continue
		}

		for {
			for _, be := range conflicts {
				be.Rename()
			}

			foundConflict := false
			newNames := map[string]bool{}
			for _, be := range conflicts {
				if _, exist := newNames[be.Name]; exist {
					foundConflict = true
					break
				} else {
					newNames[be.Name] = true
				}
			}

			if !foundConflict {
				break
			}
		}

		ResolveBlockNameConflicts(blocks)
	}
}
