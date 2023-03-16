package cmd

import (
	"testing"

	"github.com/franela/goblin"
	"github.com/gomicro/penname"
	. "github.com/onsi/gomega"
)

func TestCompletionCmd(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Completion", func() {
		g.It("should output the default completion", func() {
			w := penname.New()
			c := NewCompletionCmd(w)

			err := c.Execute()
			Expect(err).To(BeNil())
			cmdOut := string(w.Written())

			w.Reset()

			err = rootCmd.GenZshCompletion(w)
			Expect(err).To(BeNil())
			expectedOut := string(w.Written())

			Expect(cmdOut).To(Equal(expectedOut))
		})

		g.It("should output the requested completion", func() {
			w := penname.New()
			c := NewCompletionCmd(w)

			err := c.Flag("shell").Value.Set("bash")
			Expect(err).To(BeNil())

			err = c.Execute()
			Expect(err).To(BeNil())
			cmdOut := string(w.Written())

			w.Reset()

			err = rootCmd.GenBashCompletion(w)
			Expect(err).To(BeNil())
			expectedOut := string(w.Written())

			Expect(cmdOut).To(Equal(expectedOut))
		})

		g.It("should return an error for an unknown shell", func() {
			w := penname.New()
			c := NewCompletionCmd(w)

			err := c.Flag("shell").Value.Set("foo")
			Expect(err).To(BeNil())

			err = c.Execute()
			Expect(err).NotTo(BeNil())
			Expect(err).To(MatchError(ErrUknownShell))
		})
	})
}
