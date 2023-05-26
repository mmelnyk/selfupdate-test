package selfupdate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// Release collects data about a single release on GitHub.
type Release struct {
	Name        string    `json:"name"`
	TagName     string    `json:"tag_name"`
	Draft       bool      `json:"draft"`
	PreRelease  bool      `json:"prerelease"`
	PublishedAt time.Time `json:"published_at"`
	Assets      []Asset   `json:"assets"`
}

// Asset is a file uploaded and attached to a release.
type Asset struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (r Release) String() string {
	return fmt.Sprintf("%v %v, %d assets",
		r.TagName,
		r.PublishedAt.Local().Format("2006-01-02 15:04:05"),
		len(r.Assets))
}

const (
	githubAPITimeout        = 30 * time.Second
	gitRegex                = `^(https?|git)(:\/\/|@)([^\/:]+)[\/:]([^\/:]+)\/([^\/\.:]+)(|\.git)$`
	githubDomain            = "github.com"
	githubReleaseFormat     = "https://api.github.com/repos/%s/%s/releases/latest"
	githubAssetFormat       = "https://api.github.com/repos/%s/%s/releases/assets/%d"
	githubAPIAccept         = "application/vnd.github.v3+json"
	githubAPIContent        = "application/json"
	githubAPIAcceptBinaries = "application/octet-stream"
)

// githubError is returned by the GitHub API, e.g. for rate-limiting.
type githubError struct {
	Message string
}

// gitHubLatestRelease uses the GitHub API to get information about the latest
// release of a repository.
func githubLatestRelease(ctx context.Context, git string) (Release, error) {
	re := regexp.MustCompile(gitRegex)
	matches := re.FindStringSubmatch(git)
	if len(matches) < 6 {
		return Release{}, fmt.Errorf("invalid GitHub URL %q", git)
	}

	if matches[3] != githubDomain {
		return Release{}, fmt.Errorf("invalid GitHub domain %q", matches[3])
	}

	owner := matches[4]
	repo := matches[5]

	ctx, cancel := context.WithTimeout(ctx, githubAPITimeout)
	defer cancel()

	url := fmt.Sprintf(githubReleaseFormat, owner, repo)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Release{}, err
	}

	// pin API version 3
	req.Header.Set("Accept", githubAPIAccept)

	res, err := http.DefaultClient.Do(req.WithContext(ctx))
	// If we got an error, and the context has been canceled,
	// the context's error is probably more useful.
	if err != nil {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		default:
		}
		return Release{}, err
	}

	if res.StatusCode != http.StatusOK {
		content := res.Header.Get("Content-Type")
		if strings.Contains(content, githubAPIContent) {
			// try to decode error message
			var msg githubError
			jerr := json.NewDecoder(res.Body).Decode(&msg)
			if jerr == nil {
				return Release{}, fmt.Errorf("unexpected status %v (%v) returned, message:\n  %v", res.StatusCode, res.Status, msg.Message)
			}
		}

		_ = res.Body.Close()
		return Release{}, fmt.Errorf("unexpected status %v (%v) returned", res.StatusCode, res.Status)
	}

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		_ = res.Body.Close()
		return Release{}, err
	}

	err = res.Body.Close()
	if err != nil {
		return Release{}, err
	}

	var release Release
	err = json.Unmarshal(buf, &release)
	if err != nil {
		return Release{}, err
	}

	if release.TagName == "" {
		return Release{}, errors.New("tag name for latest release is empty")
	}

	return release, nil
}

func GetLatestVersion(giturl string) (string, error) {
	release, err := githubLatestRelease(context.Background(), giturl)

	if err != nil {
		return "", err
	}

	return release.TagName, nil
}

func DownloadLatestVersion(giturl string) error {
	return nil
}
