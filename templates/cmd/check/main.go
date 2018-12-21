package main

import (
	"{{ .ImportPath }}/resource"
	"github.com/cloudboss/ofcourse/ofcourse"
)

func main() {
	ofcourse.Check(&resource.Resource{})
}
