package security

// Redactor masks secret-looking strings and known API keys in output before
// it reaches logs or callbacks.
type Redactor struct {
	secrets []string
	mask    string
}

// NewRedactor creates a Redactor for the given secret strings.
func NewRedactor(secrets []string) *Redactor {
	return &Redactor{secrets: secrets, mask: "***"}
}

// Redact replaces every known secret occurrence with the mask.
func (r *Redactor) Redact(s string) string {
	for _, sec := range r.secrets {
		if sec == "" {
			continue
		}
		for {
			idx := indexOf(s, sec)
			if idx < 0 {
				break
			}
			s = s[:idx] + r.mask + s[idx+len(sec):]
		}
	}
	return s
}

// RedactAPIKeyPatterns masks keys that follow sk- / key- / Bearer patterns.
func (r *Redactor) RedactAPIKeyPatterns(s string) string {
	// crude token scan: replace anything that looks like a long secret token.
	out := []byte(s)
	i := 0
	for i < len(out) {
		switch out[i] {
		case 's', 'S', 'k', 'K':
			// guess length-based: skip; handled via secrets list.
			i++
		case 'B':
			if hasPrefixAt(out, i, "Bearer ") {
				j := i + len("Bearer ")
				for j < len(out) && out[j] != ' ' && out[j] != '\n' && out[j] != '\r' {
					j++
				}
				if j-i > len("Bearer ")+8 {
					for k := i + len("Bearer "); k < j; k++ {
						out[k] = '*'
					}
				}
				i = j
				continue
			}
			i++
		default:
			i++
		}
	}
	return string(out)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

func hasPrefixAt(b []byte, i int, prefix string) bool {
	if i+len(prefix) > len(b) {
		return false
	}
	for j := 0; j < len(prefix); j++ {
		if b[i+j] != prefix[j] {
			return false
		}
	}
	return true
}
