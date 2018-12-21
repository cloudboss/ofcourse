package main

import (
        "{{ .ImportPath }}/resource"
        "github.com/cloudboss/ofcourse/pkg/ofcourse"
)

func main() {
        ofcourse.Out(&resource.Resource{})
}
