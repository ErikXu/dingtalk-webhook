package router

import (
	"dingtalk-webhook/config"
	util "dingtalk-webhook/util"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func jira(c *gin.Context) {
	var content JiraCallback
	if err := c.ShouldBindJSON(&content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if content.WebhookEvent == "jira:issue_updated" && content.Issue.Fields.Status.Name != "待验证" {
		c.JSON(http.StatusBadRequest, content)
		return
	}

	colors := map[string]string{
		"致命": "#CE0000",
		"严重": "#EA4444",
		"一般": "#EA7D24",
		"提示": "#2A8735",
		"建议": "#55A557",
	}

	if content.WebhookEvent == "jira:issue_created" {
		title := fmt.Sprintf("%s 新增 %s %s", content.Issue.Fields.Creator.DisplayName, content.Issue.Fields.Issuetype.Name, content.Issue.Key)
		text := generateCreatedMsg(content, colors, "")
		util.SendMarkdownMsg(config.AppConfig.Dingtalk.Webhook, config.AppConfig.Dingtalk.Secret, title, text, false, []string{}, []string{})
	}

	c.JSON(http.StatusBadRequest, content)
}

func generateCreatedMsg(content JiraCallback, colors map[string]string, mobile string) string {

	priorityText := content.Issue.Fields.Priority.Name
	if val, ok := colors[content.Issue.Fields.Priority.Name]; ok {
		priorityText = fmt.Sprintf("<font color='%s'>%s</font>", val, content.Issue.Fields.Priority.Name)
	}

	assigneeText := content.Issue.Fields.Assignee.DisplayName
	if len(mobile) > 0 {
		assigneeText = fmt.Sprintf("%s @%s", content.Issue.Fields.Assignee.DisplayName, mobile)
	}

	text := fmt.Sprintf("### 新增 %s [%s](%s/browse/%s) \n", content.Issue.Fields.Issuetype.Name, content.Issue.Key, config.AppConfig.Jira.Domain, content.Issue.Key) +
		"--- \n" +
		fmt.Sprintf("- 优先级：<font color='#52C41A'>%s</font> \n", priorityText) +
		fmt.Sprintf("- 创建人：%s \n", content.Issue.Fields.Creator.DisplayName) +
		fmt.Sprintf("- 指派人：%s \n", assigneeText) +
		fmt.Sprintf("- 创建时间：%s \n ", time.Unix(content.Timestamp/1000, 0).In(time.FixedZone("CST", 8*3600)).Format("2006-01-02 15:04:05")) +
		fmt.Sprintf("> [%s](%s/browse/%s)", content.Issue.Fields.Summary, config.AppConfig.Jira.Domain, content.Issue.Key)

	return text
}

// Please add the /print api to your jira webhook and get the callback structure.
// Or refer https://developer.atlassian.com/server/jira/platform/webhooks/#example--callback-for-an-issue-related-event
type JiraCallback struct {
	Timestamp          int64              `json:"timestamp"`
	WebhookEvent       string             `json:"webhookEvent"`
	IssueEventTypeName string             `json:"issue_event_type_name"`
	User               JiraUser           `json:"user"`
	Issue              JiraIssue          `json:"issue"`
	Changelog          JiraIssueChangeLog `json:"changelog"`
}

type JiraUser struct {
	Name         string `json:"name"`
	Key          string `json:"key"`
	EmailAddress string `json:"emailAddress"`
	DisplayName  string `json:"displayName"`
	Active       bool   `json:"active"`
}

type JiraIssue struct {
	Id     string     `json:"id"`
	Self   string     `json:"self"`
	Key    string     `json:"key"`
	Fields JiraFields `json:"fields"`
}

type JiraFields struct {
	Issuetype   JiraIssueType     `json:"issuetype"`
	Description string            `json:"description"`
	Summary     string            `json:"summary"`
	Creator     JiraIssueCreator  `json:"creator"`
	Priority    JiraIssuePriority `json:"priority"`
	Progress    JiraIssueProgress `json:"progress"`
	Comment     JiraIssueComment  `json:"comment"`
	Assignee    JiraIssueAssignee `json:"assignee"`
	Status      JiraIssueStatus   `json:"status"`
}

type JiraIssueType struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Self        string `json:"self"`
	Description string `json:"description"`
}

type JiraIssueCreator struct {
	Name         string `json:"name"`
	Key          string `json:"key"`
	EmailAddress string `json:"emailAddress"`
	DisplayName  string `json:"displayName"`
	Active       bool   `json:"active"`
}

type JiraIssuePriority struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Self string `json:"self"`
}

type JiraIssueProgress struct {
	Progress int `json:"progress"`
	Total    int `json:"total"`
}

type JiraIssueAssignee struct {
	Name         string `json:"name"`
	Key          string `json:"key"`
	EmailAddress string `json:"emailAddress"`
	DisplayName  string `json:"displayName"`
	Active       bool   `json:"active"`
}

type JiraIssueStatus struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type JiraIssueStatusCategory struct {
	Id   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}

type JiraIssueComment struct {
	Comments []JiraIssueCommentItem `json:"comments"`
	Total    int                    `json:"total"`
}

type JiraIssueCommentItem struct {
	Id           string                       `json:"id"`
	Self         string                       `json:"self"`
	Author       JiraIssueCommentAuthor       `json:"author"`
	Body         string                       `json:"body"`
	UpdateAuthor JiraIssueCommentUpdateAuthor `json:"updateAuthor"`
}

type JiraIssueCommentAuthor struct {
	Name         string `json:"name"`
	Key          string `json:"key"`
	EmailAddress string `json:"emailAddress"`
	DisplayName  string `json:"displayName"`
	Active       bool   `json:"active"`
}

type JiraIssueCommentUpdateAuthor struct {
	Name         string `json:"name"`
	Key          string `json:"key"`
	EmailAddress string `json:"emailAddress"`
	DisplayName  string `json:"displayName"`
	Active       bool   `json:"active"`
}

type JiraIssueChangeLog struct {
	Id    string                   `json:"id"`
	Items []JiraIssueChangeLogItem `json:"items"`
}

type JiraIssueChangeLogItem struct {
	Field      string `json:"field"`
	Fieldtype  string `json:"fieldtype"`
	From       string `json:"from"`
	FromString string `json:"fromString"`
	To         string `json:"to"`
	ToString   string `json:"toString"`
}
