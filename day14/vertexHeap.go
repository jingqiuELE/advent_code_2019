package main

type VertexHeap []*Vertex

func (v VertexHeap) Len() int {
	return len(v)
}

func (v VertexHeap) Less(i, j int) bool {
	return v[i].Weight > v[j].Weight
}

func (v VertexHeap) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v *VertexHeap) Push(x interface{}) {
	*v = append(*v, x.(*Vertex))
	x.(*Vertex).InHeap = true
}

func (v *VertexHeap) Pop() interface{} {
	old := *v
	n := len(old)
	x := old[n-1]
	*v = old[0 : n-1]
	x.InHeap = false
	return x
}
