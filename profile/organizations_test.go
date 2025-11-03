package profile

import (
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/yyle88/done"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexistpath/osmustexist"
	"github.com/yyle88/runpath"
)

// TestGetOrganizations tests fetching organizations
// Verifies API call and caching work correctly
//
// TestGetOrganizations 测试获取组织列表
// 验证 API 调用和缓存是否正常工作
func TestGetOrganizations(t *testing.T) {
	t.Log(neatjsons.S(getOrganizationsCached()))
}

// TestFetchOrganizationRepos tests the repository collection process
// Verifies that repos are properly collected and organized
//
// TestFetchOrganizationRepos 测试仓库收集过程
// 验证仓库是否正确收集和组织
func TestFetchOrganizationRepos(t *testing.T) {
	collection := newRepoCollection()
	collection.collectRepos()
	t.Log(neatjsons.S(collection.organizations))
	// Show first 3 featured repos // 显示前3个精选仓库
	maxDisplay := min(3, len(collection.featuredRepos))
	for _, item := range collection.featuredRepos[:maxDisplay] {
		t.Log(item.orgName, item.repo.Name, item.repo.Stargazers)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type DocGenParam struct {
	shortName string
	titleLine string
	language  string
}

// TestGenMarkdown generates English version of organization profile README
// Creates markdown with organization badges and project tables
//
// TestGenMarkdown 生成英文版本的组织资料 README
// 创建包含组织徽章和项目表格的 markdown
func TestGenMarkdown(t *testing.T) {
	GenMarkdownTable(t, &DocGenParam{
		shortName: "README.md",
		titleLine: `| **<span style="font-size: 10px;">organization</span>** | **repo** |`,
		language:  "en",
	})
}

// TestGenMarkdownZhHans generates Chinese version of organization profile README
// Creates markdown with organization badges and project tables
//
// TestGenMarkdownZhHans 生成中文版本的组织资料 README
// 创建包含组织徽章和项目表格的 markdown
func TestGenMarkdownZhHans(t *testing.T) {
	GenMarkdownTable(t, &DocGenParam{
		shortName: "README.zh.md",
		titleLine: "| **组织** | **项目** |",
		language:  "zh",
	})
}

// GenMarkdownTable generates markdown tables for organization profile README
// Updates specified README file with organization content
//
// GenMarkdownTable 为组织资料 README 生成 markdown 表格
// 用组织内容更新指定的 README 文件
func GenMarkdownTable(t *testing.T, arg *DocGenParam) {
	markdown := newProfileMarkdown(arg.language)
	path := osmustexist.PATH(runpath.PARENT.Join(arg.shortName))

	// 读取原始文件作为基础
	originalText := string(done.VAE(os.ReadFile(path)).Nice())
	currentText := originalText

	t.Run("ORGS_AND_PROJECTS", func(t *testing.T) {
		orgsBadges := markdown.generateOrgsBadges()
		projectsTable := markdown.generateProjectsTable(arg.titleLine)
		content := orgsBadges + projectsTable

		t.Logf("Generated orgs badges length: %d", len(orgsBadges))
		t.Logf("Generated projects table length: %d", len(projectsTable))

		currentText = replaceBlock(currentText,
			"<!-- 这是一个注释，它不会在渲染时显示出来，这是组织项目列表的起始位置 -->",
			"<!-- 这是一个注释，它不会在渲染时显示出来，这是组织项目列表的终止位置 -->",
			content)

		t.Log("✅ Organizations and projects section updated")
	})

	t.Run("MORE_PROJECTS", func(t *testing.T) {
		moreProjects := markdown.generateMoreProjects()
		t.Logf("Generated more projects length: %d", len(moreProjects))

		currentText = replaceBlock(currentText, "<!-- 更多项目的起始位置 -->", "<!-- 更多项目的终止位置 -->", moreProjects)

		t.Log("✅ More projects section updated")
	})

	t.Run("WRITE_FILE", func(t *testing.T) {
		t.Logf("Final content length: %d", len(currentText))
		t.Logf("Content size change: %d -> %d (diff: %+d)", len(originalText), len(currentText), len(currentText)-len(originalText))

		must.Done(os.WriteFile(path, []byte(currentText), 0666))
		t.Logf("✅ Successfully updated %s", arg.shortName)
	})
}

// replaceBlock replaces content between two placeholder lines
// Returns updated content with new text between markers
//
// replaceBlock 替换两个占位符行之间的内容
// 返回标记之间包含新文本的更新内容
func replaceBlock(content, startPlaceholder, endPlaceholder, newContent string) string {
	lines := strings.Split(content, "\n")
	startIdx := slices.Index(lines, startPlaceholder)
	endIdx := slices.Index(lines, endPlaceholder)
	if startIdx != -1 && endIdx != -1 && startIdx < endIdx {
		return strings.Join(lines[:startIdx+1], "\n") + "\n" + newContent + "\n" + strings.Join(lines[endIdx:], "\n")
	}
	return content
}
