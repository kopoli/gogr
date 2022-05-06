// Generated with by licrep version v0.3.4
// https://github.com/kopoli/licrep
// Called with: licrep -o licenses.go

package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io"
	"strings"
)

// License is a representation of an embedded license.
type License struct {
	// Name of the license
	Name string

	// The text of the license
	Text string
}

// GetLicenses gets a map of Licenses where the keys are
// the package names.
func GetLicenses() (map[string]License, error) {
	type EncodedLicense struct {
		Name   string
		Text   string
		length int64
	}
	data := map[string]EncodedLicense{

		"github.com/OpenPeeDeeP/xdg": {
			Name: "NewBSD",
			Text: `
H4sIAAAAAAAC/5SST4/bNhDF7/oUDzklhbr9dyjQnGhpbBGQRZWk1vFRK9ExAUs0KHoX++0L0l6s0wYt
etJA5My89+NbqRK//Vic+stiUNvBzIvJssKdX739egz4OHzCrz//8nsOcTZza0xpTJux0wnpfIE3i/HP
ZnzIMmlGuwRvny7Buhn9PCJOtTMWd/GDSX+e7Nz7Vxycn5YcLzYc4Xz6ukvIJjfagx36OCBH7w3Oxk82
BDPi7N2zHc2IcOwDwtHg4E4n92LnrxjcPNrYtKSmyYQ/suwHfKtogTu8SRncaDBdlgBvQm/nNK9/cs/x
6M387IIdTI5wtEsGnOwS4oj7ZfP4NyWjXYZTbyfjH76nwM73BN4UnL0bL4P5NxEZkoz/KwI3a6MbLpOZ
QyKbIfb85DxcOBqPqQ/G2/60vDNOD5Ma7+QnR42xqSkezv1kophYvys+utNoPGb3fimhtyFSHNx8Hej8
gql/xZOJMRkRHMw8Or+YmIizd5MLBlc0YcFovH02Iw7eTW8wFncIL/HBb/nBcjZDDBDO3sZY+Rid+Rqi
Zbla0BVXUGKtd0wSuEIrxSMvqcRqD10RCtHuJd9UGpWoS5IKrClRiEZLvuq0kAofmAJXH7J4wJo96Esr
SSkICb5ta04ldkxK1mhOKgdvirorebPJseo0GqFR8y3XVEKLPC7N/tkGscaWZFGxRrMVr7neJyFrrpu4
ay0kGFomNS+6mkm0nWyFIjBJWclVUTO+pfIBvEEjQI/UaKiK1fV3XUbt33hcEWrOVjVlaVOzR8klFTra
ea8KXlKjWZ1DtVTwWNAX2rY1k/v8NlPRnx01mrM6K9mWbUjh438gaaUoOknbqFmsobqV0lx3mrARooyg
M0XykRekPqMWKtHqFOUomWZpcSvFmmv1OdarTvEEjTeapOxazUXzCZXY0SPJrGCdojLRFU2yqisSch+H
RgYJfo5dRboiGYEmUiwiUFryQt9dy4SEFlLfeURDm5pvqCkoqhFxyo4r+gQmuYoX+HXtju0humQ5vlGn
KEvlXWLz9JLga7DykUfZt8utUIrfcpKQFRWuuB+yvwIAAP//tNCNL+cFAAA=`,
			length: 1511,
		},
		"github.com/kopoli/appkit": {
			Name: "MIT",
			Text: `
H4sIAAAAAAAC/1xRzW7jNhC+8yk+5JQAQrrYY2+MRVtEJNKg6HV9pCU6YiuLhkg3yNsXIzu7zZ4Eceb7
HTt4NNKiDp2fksdjI+0TY6t4+ZjD25Dx2D3h+7fv3/DqxtHj1U3/uNkztvXzOaQU4oSQMPjZHz/wNrsp
+77AafYe8YRucPObL5Aj3PSBi59TnBCP2YUpTG9w6OLlg8UT8hASUjzldzd7uKmHSyl2wWXfo4/d9eyn
7DLpncLoEx7z4PHQ3hEPT4tI793IwgSafY7wHvIQrxmzT3kOHXEUCFM3Xnvy8DkewzncFQi+xE8sR1yT
LxafBc6xDyf6+iXW5XocQxoK9IGoj9fsCyR6XNosKMcfcUby48i6eAk+Ycn6y92yQ9YvVGi+V5To5X2I
569JQmKn6zyFNPgF00ekuCj+7btML7R+iuMY3ylaF6c+UKL0J2N0aneM//oly+26U8yhu9W9HODy66r3
URrcOOLo74X5HmGC+1+cmeRTdlMObsQlzove7zGfGbOVQKvXds+NgGyxNfqHLEWJB95Ctg8F9tJWemex
58ZwZQ/Qa3B1wKtUZQHx19aItoU2TDbbWoqygFSreldKtcHLzkJpi1o20ooSVoME71RStETWCLOquLL8
RdbSHgq2llYR51obcGy5sXK1q7nBdme2uhXgqoTSSqq1kWojGqHsM6SC0hA/hLJoK17XJMX4zlbakD+s
9PZg5KayqHRdCtPiRaCW/KUWNyl1wKrmsilQ8oZvxILSthKG0drNHfaVoCfS4wp8ZaVWFGOllTV8ZQtY
bexP6F62ogA3sqVC1kY3BaM69ZpWpCKcEjcWqhpfLqLN8r9rxU9ClILXUm1aAlPEz+Vn9l8AAAD//7MD
VDw4BAAA`,
			length: 1080,
		},
		"github.com/kopoli/gogr/lib": {
			Name: "MIT",
			Text: `
H4sIAAAAAAAC/1xRzW6rOBTeI/EOn7pqJdT5Wcxidm5wglWwI+PcTJYOOMUzBEfYmapvPzokvXd6Vwif
8/0eMzg0wqD2nZuiw2MjzFOe5dkqXD5m/zYkPHZP+P3X3/7Aqx1Hh1c7/WNnRztbN599jD5M8BGDm93x
A2+znZLrC5xm5xBO6AY7v7kCKcBOH7i4OYYJ4Zisn/z0BosuXD7yLJyQBh8Rwym929nBTj1sjKHzNrke
feiuZzclm0jw5EcX8ZgGh4f2jnh4WlR6Z8c88xNo+DnDu09DuCbMLqbZd0RSwE/deO3Jxed49Gd/lyD4
UkHMsxRwja5YrBY4h96f6OuWZJfrcfRxKNB74j5ekysQ6XEptaAkv4QZ0Y1jnnXh4l3EEveHv2WJ3F+o
1HSvKdLL+xDOX7P4mGen6zz5OLgF1AfEsGj+7bpEL7R/CuMY3ildF6beU6j4Jx2Orm6P4V+35LmdeQrJ
d7fSlzNcfhz3PoqDHUcc3b0118NPsP+PNJODmOyUvB1xCfMi+XPU58VCxdGqtdkzzSFabLX6Jkpe4oG1
EO1Dgb0wldoZ7JnWTJoD1BpMHvAqZFmA/7XVvG2hdJ6JZlsLXhYQclXvSiE3eNkZSGVQi0YYXsIokOKd
S/CW2BquVxWThr2IWphDkWdrYSSxrpUGw5ZpI1a7mmlsd3qrWg4mS0glhVxrITe84dI8Q0hIBf6NS4O2
YnVNWnnGdqZSmixipbYHLTaVQaXqkusWLxy1YC81v2nJA1Y1E02BkjVswxeUMhXXeUZ7N4PYV5zeSJFJ
sJURSlKSlZJGs5UpYJQ237F70fICTIuWOllr1RR5Rp2qNe0ISUDJbzTUN76cRenlf9fy74woOauF3LQE
XlJ+bj//FwAA//+j1dj1SwQAAA==`,
			length: 1099,
		},
	}

	decode := func(input string, length int64) (string, error) {
		data := &bytes.Buffer{}
		br := base64.NewDecoder(base64.StdEncoding, strings.NewReader(input))

		r, err := gzip.NewReader(br)
		if err != nil {
			return "", err
		}

		_, err = io.CopyN(data, r, length)
		if err != nil {
			return "", err
		}

		// Make sure the gzip is decoded successfully
		err = r.Close()
		if err != nil {
			return "", err
		}
		return data.String(), nil
	}

	ret := make(map[string]License)

	for k := range data {
		text, err := decode(data[k].Text, data[k].length)
		if err != nil {
			return nil, err
		}
		ret[k] = License{
			Name: data[k].Name,
			Text: text,
		}
	}

	return ret, nil
}
