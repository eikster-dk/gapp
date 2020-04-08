package cli

import (
	"fmt"
	"github.com/eikc/gapp/internal/authentication"
	"github.com/google/go-github/v30/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

func repoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Short: "repo ",
		Use: "repo",
	}

	cmd.AddCommand(cloneByTagCmd())

	return cmd
}

func cloneByTagCmd() *cobra.Command {
	var cloneByTag = &cobra.Command{
		Short: "clone-tag will clone all repositories from the owner that matches the provided topic",
		Use:   "clone-tag [owner]",
		RunE:  cloneByTag,
	}
	cloneByTag.Flags().StringP("output", "o", "./", "Output is the location which repositories will be cloned to, defaults to current folder")
	cloneByTag.Flags().StringP("topic", "t", "", "clones repositories which has the provided topic")
	cloneByTag.MarkFlagRequired("topic")

	return cloneByTag
}

func cloneByTag(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	owner := args[0]
	if owner == "" {
		return fmt.Errorf("owner not provided")
	}

	topic,err := cmd.Flags().GetString("topic")
	if err != nil {
		return err
	}

	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}

	user, err := authentication.GetUser()
	if err != nil {
		return err
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: user.Token},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 50},
	}
	var allRepos []*github.Repository

	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, owner, opt)
		if err != nil {
			return err
		}

		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	var selected []*github.Repository

	for _, repo := range allRepos {
		for _, t := range repo.Topics {
			if t == topic {
				selected = append(selected, repo)
			}
		}
	}

	fmt.Println("number of repos: ", len(selected))

	for _, repo := range selected {
		fmt.Println("url: ", *repo.CloneURL)
		_, err := git.PlainClone(output, false, &git.CloneOptions{
			URL: *repo.CloneURL,
			Auth: &http.BasicAuth{
				Username: user.Username,
				Password: user.Token,
			},
			Progress: cmd.OutOrStdout(),
		})

		if err != nil {
			return err
		}
	}

	return nil
}
