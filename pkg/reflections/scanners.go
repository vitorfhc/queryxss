package reflections

import "github.com/sirupsen/logrus"

type Reflection struct {
	Url      string
	Severity string
	What     string
	Where    string
}

type ScanFunc func(client ScanHttpClient, url string, minLength uint) ([]*Reflection, error)

func SimpleScan(client ScanHttpClient, url string, minLength uint) ([]*Reflection, error) {
	logrus.Debug("running simple scan")
	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	reflections, err := FindReflectedQueryValues(url, res, minLength)
	if err != nil {
		return nil, err
	}

	logrus.Debugf("found %d reflections", len(reflections))

	return reflections, nil
}
