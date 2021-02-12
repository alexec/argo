// +build e2e

package e2e

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/argoproj/argo-workflows/v3/test/e2e/fixtures"
)

type PodTemplateSuite struct {
	fixtures.E2ESuite
}

func (s *PodTemplateSuite) SetupSuite() {
	s.E2ESuite.SetupSuite()
	s.Need(fixtures.PodTemplate)
}

func (s *PodTemplateSuite) TestPodTemplateWorkflow() {
	infos, err := ioutil.ReadDir("testdata/pod-template")
	assert.NoError(s.T(), err)
	for _, info := range infos {
		s.T().Run(info.Name(), func(t *testing.T) {
			s.Given().
				Workflow("@testdata/pod-template/" + info.Name()).
				When().
				SubmitWorkflow().
				WaitForWorkflow(fixtures.ToBeSucceeded)
		})
	}
}

func TestPodTemplateSuite(t *testing.T) {
	suite.Run(t, new(PodTemplateSuite))
}
