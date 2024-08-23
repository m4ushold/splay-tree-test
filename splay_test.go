package splay_test

import (
	gojson "encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	splay "splay-test/splayTree"
	"testing"
)

type editTrace struct {
	Edits     [][]interface{} `json:"edits"`
	FinalText string          `json:"finalText"`
}

// func TestSplayTree(t *testing.T) {

// readEditingTraceFromFile reads trace from editing-trace.json.
func readEditingTraceFromFile(b *testing.B) (*editTrace, error) {
	var trace editTrace

	file, err := os.Open("./editing-trace.json")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = file.Close(); err != nil {
			b.Fatal(err)
		}
	}()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if err = gojson.Unmarshal(byteValue, &trace); err != nil {
		return nil, err
	}

	return &trace, err
}

type stringValue struct {
	content string
	removed bool
}

func newSplayNode(content string) *splay.Node[*stringValue] {
	return splay.NewNode(&stringValue{
		content: content,
	})
}

func (v *stringValue) Len() int {
	if v.removed {
		return 0
	}
	return len(v.content)
}

func (v *stringValue) String() string {
	return v.content
}

type splayTrees interface {
	Insert(*splay.Node[*stringValue]) *splay.Node[*stringValue]
	Find(int) (*splay.Node[*stringValue], int, error)
	Delete(*splay.Node[*stringValue])
	Height() int
	RotateCount() int
}

func Execute(tree splayTrees, editingTrace *editTrace) (int, int, int) {
	maxHeight := 0
	maxRotate := 0
	prevRotate := 0
	for _, edit := range editingTrace.Edits {
		cursor := int(edit[0].(float64))
		mode := int(edit[1].(float64))

		if mode == 0 {
			strValue, ok := edit[2].(string)
			if ok {
				tree.Insert(newSplayNode(strValue))
			}
		} else {
			node, _, err := tree.Find(cursor)
			if err != nil {
				tree.Delete(node)
			}
		}
		if tree.RotateCount()-prevRotate > maxRotate {
			maxRotate = tree.RotateCount() - prevRotate
		}
		if maxHeight < tree.Height() {
			maxHeight = tree.Height()
		}
		prevRotate = tree.RotateCount()
	}
	sumRotate := tree.RotateCount()

	res := fmt.Sprintf("\nresult(%d,%d,%d)", maxHeight, maxRotate, sumRotate)
	fmt.Println(res)

	return maxHeight, maxRotate, tree.RotateCount()
}

func BenchmarkBasicSplayTree(b *testing.B) {
	b.StopTimer()
	editingTrace, err := readEditingTraceFromFile(b)
	if err != nil {
		b.Fatal(err)
	}
	b.StartTimer()

	tree := splay.NewBasicSplayTree[*stringValue](nil)
	Execute(tree, editingTrace)
	b.StopTimer()
}

func BenchmarkRandom(b *testing.B) {
	b.StopTimer()
	editingTrace, err := readEditingTraceFromFile(b)
	if err != nil {
		b.Fatal(err)
	}
	b.StartTimer()

	functions := []func(i int) int{
		func(i int) int { return int(math.Sqrt(float64(i))) },
		func(i int) int { return int(math.Log2(float64(i))) },
		func(i int) int {
			x := int(math.Log2(float64(i)))
			return x * x
		},
		func(i int) int {
			k := 3
			return k * i / (int(math.Log2(float64(i))) + 1)
		},
	}
	functionNames := []string{"sqrt", "log", "log^2", "k(n/ln{n})"}
	splayCounts := []int{1, 2, 3}
	ratios := []int{10, 20, 30, 40, 50}
	bounds := []int{10000, 50000, 100000}

	b.Run("random splay by counting", func(b *testing.B) {
		for fc := 0; fc < 4; fc++ {
			for spCnt := 0; spCnt < 3; spCnt++ {
				b.Run(fmt.Sprintf("function: %s, splayCount: %d", functionNames[fc], splayCounts[spCnt]), func(b *testing.B) {
					tree := splay.NewRandomByCountSplayTree[*stringValue](nil, functions[fc], splayCounts[spCnt])
					Execute(tree, editingTrace)
				})
			}
		}
	})

	b.Run("random splay by k-height", func(b *testing.B) {
		for r := 0; r < 5; r++ {
			for spCnt := 0; spCnt < 3; spCnt++ {
				b.Run(fmt.Sprintf("ratio: %d, splayCount: %d", ratios[r], splayCounts[spCnt]), func(b *testing.B) {
					tree := splay.NewRandomKSplayTree[*stringValue](nil, ratios[r], splayCounts[spCnt])
					Execute(tree, editingTrace)
				})
			}
		}
	})

	b.Run("random splay by bounded height", func(b *testing.B) {
		for bd := 0; bd < 3; bd++ {
			for spCnt := 0; spCnt < 3; spCnt++ {
				b.Run(fmt.Sprintf("bound: %d, splayCount: %d", bounds[bd], splayCounts[spCnt]), func(b *testing.B) {
					tree := splay.NewRandomBoundSplayTree[*stringValue](nil, bounds[bd], splayCounts[spCnt])
					Execute(tree, editingTrace)
				})
			}
		}
	})

	b.StopTimer()
}

func BenchmarkMaxHeightSplay(b *testing.B) {
	b.StopTimer()
	editingTrace, err := readEditingTraceFromFile(b)
	if err != nil {
		b.Fatal(err)
	}
	b.StartTimer()

	functions := []func(i int) int{
		func(i int) int { return int(math.Sqrt(float64(i))) },
		func(i int) int { return int(math.Log2(float64(i))) },
		func(i int) int {
			x := int(math.Log2(float64(i)))
			return x * x
		},
		func(i int) int {
			k := 3
			return k * i / (int(math.Log2(float64(i))) + 1)
		},
	}
	functionNames := []string{"sqrt", "log", "log^2", "k(n/ln{n})"}
	splayCounts := []int{1, 2, 3}
	ratios := []int{10, 20, 30, 40, 50}
	bounds := []int{10000, 50000, 100000}

	b.Run("max height splay by counting", func(b *testing.B) {
		for fc := 0; fc < 4; fc++ {
			for spCnt := 0; spCnt < 3; spCnt++ {
				b.Run(fmt.Sprintf("function: %s, splayCount: %d", functionNames[fc], splayCounts[spCnt]), func(b *testing.B) {
					tree := splay.NewMaxHeightByCountSplayTree[*stringValue](nil, functions[fc], splayCounts[spCnt])
					Execute(tree, editingTrace)
				})
			}
		}
	})

	b.Run("max height splay by k-height", func(b *testing.B) {
		for r := 0; r < 5; r++ {
			for spCnt := 0; spCnt < 3; spCnt++ {
				b.Run(fmt.Sprintf("ratio: %d, splayCount: %d", ratios[r], splayCounts[spCnt]), func(b *testing.B) {
					tree := splay.NewMaxHeightKSplayTree[*stringValue](nil, ratios[r], splayCounts[spCnt])
					Execute(tree, editingTrace)
				})
			}
		}
	})

	b.Run("max height splay by bounded height", func(b *testing.B) {
		for bd := 0; bd < 3; bd++ {
			for spCnt := 0; spCnt < 3; spCnt++ {
				b.Run(fmt.Sprintf("bound: %d, splayCount: %d", bounds[bd], splayCounts[spCnt]), func(b *testing.B) {
					tree := splay.NewMaxHeightBoundSplayTree[*stringValue](nil, bounds[bd], splayCounts[spCnt])
					Execute(tree, editingTrace)
				})
			}
		}
	})

	b.StopTimer()
}

func BenchmarkSTLB(b *testing.B) {
	b.StopTimer()
	editingTrace, err := readEditingTraceFromFile(b)
	if err != nil {
		b.Fatal(err)
	}
	b.StartTimer()

	thresholds := []int{1000, 2000, 3000, 4000, 5000, 6000, 7000, 8000, 9000, 10000}

	b.Run("max height splay by counting", func(b *testing.B) {
		for i := 0; i < 10; i++ {
			b.Run(fmt.Sprintf("ration: %d", thresholds[i]), func(b *testing.B) {
				tree := splay.NewSTLB[*stringValue](nil, thresholds[i])
				Execute(tree, editingTrace)
			})
		}
	})

	b.StopTimer()
}
