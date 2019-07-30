package main

import (
	"{{ .ImportPath }}/resource"
	"github.com/cloudboss/ofcourse/ofcourse"
)

func main() {
	ofcourse.Out(&resource.Resource{})
}
