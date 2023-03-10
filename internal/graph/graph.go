package graph

// Trusted - verified device
// Flag - device seen by trusted
// Pawn - device that is not trusted and not seen by trusted
const (
	TrustedScore = -1000
	FlagScore    = -2000
)

type Directed struct {
	trusted   map[string]*Vertex
	untrusted map[string]*Vertex
}

type Vertex struct {
	id    string
	Score int
	edges []*Edge
}

type Edge struct {
	from *Vertex
	to   *Vertex
}

func NewDirected() *Directed {

	return &Directed{
		trusted:   make(map[string]*Vertex),
		untrusted: make(map[string]*Vertex),
	}
}

func (d *Directed) Get(id string) *Vertex {

	return d.untrusted[id]
}

func (d *Directed) AddEdge(fromVertexId string, fromTrusted bool, toVertexId string, toTrusted bool) {

	var from, to *Vertex
	var ok bool

	// from
	if fromTrusted {
		from, ok = d.trusted[fromVertexId]
		if !ok {
			from = &Vertex{id: fromVertexId}
			d.trusted[fromVertexId] = from
		}
	} else {
		from, ok = d.untrusted[fromVertexId]
		if !ok {
			from = &Vertex{id: fromVertexId}
			d.untrusted[fromVertexId] = from
		}
	}

	// to
	if toTrusted {
		to, ok = d.trusted[toVertexId]
		if !ok {
			to = &Vertex{id: toVertexId}
			d.trusted[toVertexId] = to
		}
	} else {
		to, ok = d.untrusted[toVertexId]
		if !ok {
			to = &Vertex{id: toVertexId}
			d.untrusted[toVertexId] = to
		}
	}

	from.edges = append(from.edges, &Edge{from: from, to: to})
}

func (d *Directed) Score() {

	flags := d.scoreFlags()
	for _, v := range d.untrusted {
		if isPawn(v) {
			v.Score = scorePawn(v, flags)
		}
	}
}

// scorePawn returns number of flags that have path to the pawn
func scorePawn(v *Vertex, flags []*Vertex) int {

	score := 0
	for _, currFlag := range flags {
		if findPath(currFlag, v.id) {
			score++
		}
	}

	return score
}

func findPath(from *Vertex, to string) bool {

	for _, currEdge := range from.edges {
		if currEdge.to.id == to {
			return true
		}
		if findPath(currEdge.to, to) {
			return true
		}
	}

	return false
}

func (d *Directed) scoreFlags() []*Vertex {

	flags := []*Vertex{}
	for _, v := range d.trusted {
		v.Score = TrustedScore
		for _, currEdge := range v.edges {
			if _, ok := d.untrusted[currEdge.to.id]; ok {
				currEdge.to.Score = FlagScore
				flags = append(flags, currEdge.to)
			}
		}
	}

	return flags
}

func isPawn(v *Vertex) bool {

	return v.Score != FlagScore
}
