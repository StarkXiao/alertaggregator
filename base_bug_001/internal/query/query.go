package query

import "alertaggregator/internal/model"

func ByService(in []model.Alert, service string) []model.Alert {
	out := []model.Alert{}
	for _, a := range in {
		if service == "" || a.Service == service {
			out = append(out, a)
		}
	}
	return out
}
