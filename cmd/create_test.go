package cmd

import (
	"fmt"
	"strings"
	"testing"

	"github.com/franela/goblin"
	"github.com/gomicro/penname"
	"github.com/gomicro/train/client/clienttest"
	"github.com/google/go-github/github"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestCreateCmd(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Create", func() {
		g.It("should create prs for org repos", func() {
			w := penname.New()

			cmd := NewCreateCmd(w)
			cmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
				clt = clienttest.New(&clienttest.Config{
					BaseBranchName: "release",
					Logins:         []string{"gomicro"},
					Repos: []*github.Repository{
						{
							Name: github.String("steward"),
							Owner: &github.User{
								Login: github.String("gomicro"),
							},
							DefaultBranch: github.String("master"),
						},
					},
				})

				dryRun = viper.GetBool("dryRun")
			}

			cmd.SetArgs([]string{"gomicro"})
			err := cmd.Execute()
			Expect(err).To(BeNil())
			cmdOut := string(w.Written())

			entityOut := "Entity: gomicro\n"
			Expect(strings.HasPrefix(cmdOut, entityOut)).To(BeTrue(), fmt.Sprintf("missing entity out line in output: got %s", cmdOut))
			cmdOut = strings.TrimPrefix(cmdOut, entityOut)

			baseOut := fmt.Sprintf("Base: %s\n", "release")
			Expect(strings.HasPrefix(cmdOut, baseOut)).To(BeTrue(), fmt.Sprintf("missing base out line in output: got %s", cmdOut))
			cmdOut = strings.TrimPrefix(cmdOut, baseOut)
			cmdOut = strings.TrimPrefix(cmdOut, "\n")

			Expect(cmdOut).To(Equal("\nRelease PRs Created:\n\nhttps://github.com/gomicro/steward/pull/0\n"))
		})
	})
}
