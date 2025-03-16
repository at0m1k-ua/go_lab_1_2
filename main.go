package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type MeasurementType struct {
	Name  string
	Label string
	Units string
	Value string
}

type TemplateModel struct {
	Measurements []MeasurementType
	CalcResult   []MeasurementType
	Error        error
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/index.html")

	m := TemplateModel{}
	m.Measurements = []MeasurementType{
		{Name: "hg", Label: "Hг", Units: "%"},
		{Name: "cg", Label: "Cг", Units: "%"},
		{Name: "sg", Label: "Sг", Units: "%"},
		{Name: "og", Label: "Oг", Units: "%"},
		{Name: "vg", Label: "Vг", Units: "мг/кг"},
		{Name: "wg", Label: "Wг", Units: "%"},
		{Name: "ag", Label: "Aг", Units: "%"},
	}

	if r.Method != "POST" {
		_ = tmpl.Execute(w, m)
		return
	}

	err := r.ParseForm()
	if err != nil {
		return
	}

	for i := range m.Measurements {
		m.Measurements[i].Value = r.FormValue(m.Measurements[i].Name)
	}

	if res, err := Calculate(m.Measurements); err != nil {
		m.Error = err
	} else {
		m.CalcResult = res
	}

	_ = tmpl.Execute(w, m)
}

func Calculate(measurements []MeasurementType) ([]MeasurementType, error) {
	i := make(map[string]float64)
	var o []MeasurementType

	for _, m := range measurements {
		if m.Value == "" {
			return o, fmt.Errorf("поле \"%s\" не заповнене", m.Label)
		}

		if val, err := strconv.ParseFloat(m.Value, 64); err != nil {
			return o, fmt.Errorf("поле \"%s\" містить невірне значення", m.Label)
		} else {
			i[m.Name] = val
		}

	}

	hg := i["hg"]
	cg := i["cg"]
	sg := i["sg"]
	og := i["og"]
	vg := i["vg"]
	wg := i["wg"]
	ag := i["ag"]
	const qidaf = 40.4

	cr := cg * (100 - wg - ag) / 100
	hr := hg * (100 - wg - ag) / 100
	or := og * (100 - wg - ag) / 100
	sr := sg * (100 - wg - ag) / 100
	ar := ag * (100 - wg) / 100
	vr := vg * (100 - wg) / 100
	qri := qidaf*(100-wg-ar)/100 - 0.025

	o = []MeasurementType{
		{"hr", "Hp", "%", fmt.Sprintf("%.2f", hr)},
		{"cr", "Cp", "%", fmt.Sprintf("%.2f", cr)},
		{"sr", "Sp", "%", fmt.Sprintf("%.2f", sr)},
		{"or", "Op", "%", fmt.Sprintf("%.2f", or)},
		{"vr", "Vp", "мг/кг", fmt.Sprintf("%.2f", vr)},
		{"ar", "Ap", "%", fmt.Sprintf("%.2f", ar)},
		{"qri", "Qri", "МДж/кг", fmt.Sprintf("%.2f", qri)},
	}

	return o, nil
}

func main() {
	http.HandleFunc("/", IndexHandler)

	fmt.Println("Server is listening...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
