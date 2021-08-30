package router

import (
	"dingtalk-webhook/config"
	util "dingtalk-webhook/util"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func gitlab(c *gin.Context) {
	var content GitlabCallback
	if err := c.ShouldBindJSON(&content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if content.ObjectKind == "merge_request" {
		actionText := getActionText(content.ObjectAttributes.Action)
		title := fmt.Sprintf("%s %s %s/%d 合并请求 %s", content.User.Name, actionText, content.Repository.Name, content.ObjectAttributes.Iid, content.ObjectAttributes.Title)
		text := generateGitlabMsg(content)
		util.SendMarkdownMsg(config.AppConfig.Dingtalk.Webhook, config.AppConfig.Dingtalk.Secret, title, text, false, []string{}, []string{})
	}

	c.JSON(http.StatusOK, content)
}

func getActionText(action string) string {
	actionText := map[string]string{
		"open":      "创建",
		"close":     "关闭",
		"reopen":    "重新创建",
		"update":    "更新",
		"approve":   "通过",
		"unapprove": "否决",
		"merge":     "合并",
	}

	if val, ok := actionText[action]; ok {
		return val
	}

	return action
}

func generateGitlabMsg(content GitlabCallback) string {
	statusColors := map[string]string{
		"opened":     "#EA7D24",
		"closed":     "#CE0000",
		"reopened":   "#EA7D24",
		"updated":    "#EA7D24",
		"approved":   "#2A8735",
		"unapproved": "#EA4444",
		"merged":     "#55A557",
	}

	if val, ok := statusColors[content.ObjectAttributes.State]; ok {
		content.ObjectAttributes.State = fmt.Sprintf("<font color='%s'>%s</font>", val, content.ObjectAttributes.State)
	}

	assignees := ""

	if content.Assignees != nil && len(content.Assignees) > 0 {
		for _, assignee := range content.Assignees {
			if mobile, ok := config.AppConfig.UserMap.Users[assignee.Username]; ok {
				assignees = assignees + fmt.Sprintf("%s @%s ", assignee.Name, mobile)
			} else {
				assignees = assignees + fmt.Sprintf("%s ", assignee.Name)
			}
		}
	}

	actionText := getActionText(content.ObjectAttributes.Action)

	text := fmt.Sprintf("### %s 合并请求 [%s/%d](%s) \n", actionText, content.Repository.Name, content.ObjectAttributes.Iid, content.ObjectAttributes.Url) +
		"--- \n" +
		fmt.Sprintf("- 标题：[%s](%s) \n", content.ObjectAttributes.Title, content.ObjectAttributes.Url) +
		fmt.Sprintf("- 操作人：%s \n", content.User.Name) +
		fmt.Sprintf("- 指派人：%s \n", assignees) +
		fmt.Sprintf("- 项目：[%s](%s) \n", content.Repository.Name, content.Repository.Homepage) +
		fmt.Sprintf("- 分支：%s -> %s \n", content.ObjectAttributes.SourceBranch, content.ObjectAttributes.TargetBranch) +
		fmt.Sprintf("- 状态：%s \n", content.ObjectAttributes.State) +
		fmt.Sprintf("- 创建时间：%s \n", content.ObjectAttributes.CreatedAt.ToTime().In(time.FixedZone("CST", 8*3600)).Format("2006-01-02 15:04:05")) +
		fmt.Sprintf("- 最后提交：[%s](%s) \n", content.ObjectAttributes.LastCommit.Id[0:8], content.ObjectAttributes.LastCommit.Url) +
		fmt.Sprintf("- 提交人：%s \n", content.ObjectAttributes.LastCommit.Author.Name) +
		fmt.Sprintf("- 提交时间：%s \n", content.ObjectAttributes.LastCommit.Timestamp.In(time.FixedZone("CST", 8*3600)).Format("2006-01-02 15:04:05")) +
		fmt.Sprintf("- 备注：[%s](%s)", content.ObjectAttributes.LastCommit.Title, content.ObjectAttributes.LastCommit.Url)

	return text
}

// Please add the /print api to your gitlab webhook and get the callback structure.
// Or refer https://docs.gitlab.com/ee/user/project/integrations/webhooks.html
type GitlabCallback struct {
	ObjectKind       string                 `json:"object_kind"`
	EventType        string                 `json:"event_type"`
	User             GitlabUser             `json:"user"`
	Project          GitlabProject          `json:"project"`
	ObjectAttributes GitlabObjectAttributes `json:"object_attributes"`
	Repository       GitlabRepository       `json:"repository"`
	Assignees        []GitlabAssignee       `json:"assignees"`
}

type GitlabUser struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type GitlabProject struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	WebUrl        string `json:"web_url"`
	GitSshUrl     string `json:"git_ssh_url"`
	GitHttpUrl    string `json:"git_http_url"`
	Namespace     string `json:"namespace"`
	DefaultBranch string `json:"default_branch"`
	Homepage      string `json:"homepage"`
	Url           string `json:"url"`
	SshUrl        string `json:"ssh_url"`
	HttpUrl       string `json:"http_url"`
}

type GitlabObjectAttributes struct {
	Id           int              `json:"id"`
	Iid          int              `json:"iid"`
	MergStatus   string           `json:"merge_status"`
	SourceBranch string           `json:"source_branch"`
	TargetBranch string           `json:"target_branch"`
	CreatedAt    GitlabTime       `json:"created_at"`
	Title        string           `json:"title"`
	Description  string           `json:"description"`
	Url          string           `json:"url"`
	LastCommit   GitlabLastCommit `json:"last_commit"`
	State        string           `json:"state"`
	Action       string           `json:"action"`
}

type GitlabRepository struct {
	Name        string `json:"name"`
	Url         string `json:"url"`
	Description string `json:"description"`
	Homepage    string `json:"homepage"`
}

type GitlabAssignee struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type GitlabLastCommit struct {
	Id        string                 `json:"id"`
	Message   string                 `json:"message"`
	Title     string                 `json:"title"`
	Timestamp time.Time              `json:"timestamp"`
	Url       string                 `json:"url"`
	Author    GitlabLastCommitAuthor `json:"author"`
}

type GitlabLastCommitAuthor struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type GitlabTime time.Time

func (j *GitlabTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02 15:04:05 UTC", s)
	if err != nil {
		return err
	}
	*j = GitlabTime(t)
	return nil
}

func (j GitlabTime) ToTime() time.Time {
	return time.Time(j)
}
