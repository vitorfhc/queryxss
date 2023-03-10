package reflections

import "net/url"

func GetQueryValues(u string) ([]string, error) {
	up, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	values := []string{}
	for _, v := range up.Query() {
		values = append(values, v...)
	}
	return values, nil
}
