package jvalidator

type invalids map[string][]string

func (i invalids) Add(name, message string) {
	i[name] = append(i[name], message)
}

func (i invalids) Get(name string) string {
	es := i[name]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
