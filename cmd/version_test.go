package cmd

import (
	"testing"

	"github.com/franela/goblin"
	"github.com/gomicro/penname"
	. "github.com/onsi/gomega"
)

func TestVersionCmd(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Version", func() {
		g.It("should print the version set", func() {
			w := penname.New()
			c := NewVersionCmd(w, "test-version")

			err := c.Execute()
			Expect(err).To(BeNil())

			out := w.Written()
			Expect(len(out)).To(Equal(27))
			Expect(string(out)).To(ContainSubstring("test-version"))
		})
	})
}
