package memory

import "math"

// Embedder produces vector embeddings for text. The default implementation
// uses a deterministic character bigram hashing so it works offline.
type Embedder interface {
	Embed(text string) ([]float64, error)
}

// HashEmbedder is a deterministic, dependency-free Embedder.
type HashEmbedder struct {
	dims int
}

// NewHashEmbedder creates a HashEmbedder with the given vector dimensions.
func NewHashEmbedder(dims int) *HashEmbedder {
	if dims <= 0 {
		dims = 64
	}
	return &HashEmbedder{dims: dims}
}

// Embed hashes character bigrams into a normalized vector.
func (h *HashEmbedder) Embed(text string) ([]float64, error) {
	vec := make([]float64, h.dims)
	if text == "" {
		return vec, nil
	}
	buf := []rune(text)
	var weight float64
	push := func(i int, v float64) {
		if i < 0 {
			i = (i%h.dims + h.dims) % h.dims
		}
		vec[i%h.dims] += v
		weight += v
	}
	for i := 0; i < len(buf); i++ {
		push(2497527*int(buf[i]), 1)
		if i > 0 {
			push(4671397*int(buf[i-1])+886609*int(buf[i]), 0.5)
		}
	}
	if weight > 0 {
		for i := range vec {
			vec[i] /= weight
		}
	}
	return vec, nil
}

// cosine returns the cosine similarity between two vectors.
func cosine(a, b []float64) float64 {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}
	var dot, na, nb float64
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		dot += a[i] * b[i]
		na += a[i] * a[i]
		nb += b[i] * b[i]
	}
	if na == 0 || nb == 0 {
		return 0
	}
	return dot / (math.Sqrt(na) * math.Sqrt(nb))
}
