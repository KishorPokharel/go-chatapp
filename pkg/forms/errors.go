package forms

type FormErrors map[string][]string

func (e FormErrors) Add(field, message string) {
	e[field] = append(e[field], message)
}

func (e FormErrors) Get(field string) string {
	es, ok := e[field]
	if !ok || len(es) == 0 {
		return ""
	}
	return es[0]
}
