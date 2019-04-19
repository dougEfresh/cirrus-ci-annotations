package junit

import (
	"fmt"
	"github.com/cirruslabs/cirrus-ci-annotations/model"
	"github.com/cirruslabs/cirrus-ci-annotations/util"
	"github.com/joshdk/go-junit"
)

func ParseJUnitAnnotations(path string) (error, []model.Annotation) {
	suites, err := junit.IngestFile(path)
	if err != nil {
		return err, nil
	}
	result := make([]model.Annotation, 0)
	for _, suite := range suites {
		for _, test := range suite.Tests {
			fqn := fmt.Sprintf("%s.%s", test.Classname, test.Name)
			switch test.Status {
			case junit.StatusPassed:
				result = append(
					result,
					model.Annotation{
						Type:               model.TestResultAnnotationType,
						Level:              "notice",
						Message:            fqn,
						FullyQualifiedName: fqn,
					},
				)
			case junit.StatusFailed:
				result = append(
					result,
					model.Annotation{
						Type:               model.TestResultAnnotationType,
						Level:              "failure",
						Message:            fqn,
						FullyQualifiedName: fqn,
						RawDetails:         test.Error.Error(),
						Location: util.GuessLocationIgnored(
							test.Error.Error(),
							[]string{
								"junit",
							},
						),
					},
				)
			}
		}
	}
	return nil, result
}
