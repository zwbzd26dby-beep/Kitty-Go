package compute

import (
	"sort"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/compute"
)

// candidate is a device with its computed score.
type candidate struct {
	device compute.Device
	score  float64
}

// candidates filters the registry by hard constraints and scores survivors.
func (r *Router) candidates(req Request, excluded map[string]bool) []candidate {
	var cands []candidate
	for _, d := range r.registry.List() {
		if excluded[d.ID] {
			continue
		}
		if !filter(d, req) {
			continue
		}
		cands = append(cands, candidate{device: d, score: r.score(d)})
	}
	sortByScore(cands)
	return cands
}

// filter enforces hard resource/model constraints.
func filter(d compute.Device, req Request) bool {
	if d.Health != compute.HealthHealthy && d.Health != compute.HealthUnknown {
		return false
	}
	if req.NeedVRAM > 0 && d.Resources.VRAMMB < req.NeedVRAM {
		return false
	}
	if req.NeedRAM > 0 && d.Resources.RAMMB < req.NeedRAM {
		return false
	}
	if req.PreferredModel != "" {
		ok := false
		for _, m := range d.SupportedModels {
			if m == req.PreferredModel {
				ok = true
				break
			}
		}
		if !ok {
			return false
		}
	}
	return true
}

// score ranks a device 0..1. Load/latency dominate when available; otherwise
// it favours larger free resources.
func (r *Router) score(d compute.Device) float64 {
	if r.metrics != nil {
		load := r.metrics.Load(d.ID)
		lat := r.metrics.LatencyMS(d.ID)
		// Lower load and latency are better; frees up to +1 weight.
		score := (1.0 - load) + (1.0 / (1.0 + lat))
		return score
	}
	// Resource heuristic: normalize RAM as the dominant free-resource signal.
	ramScore := 0.0
	if d.Resources.RAMMB > 0 {
		ramScore = float64(d.Resources.RAMMB) / 1024.0 // per 1GB
	}
	return 0.5 + ramScore
}

// sortByScore sorts candidates by descending score, breaking ties by device ID
// so ordering is deterministic regardless of map iteration order.
func sortByScore(cands []candidate) {
	sort.SliceStable(cands, func(i, j int) bool {
		if cands[i].score != cands[j].score {
			return cands[i].score > cands[j].score
		}
		return cands[i].device.ID < cands[j].device.ID
	})
}
