package internal

import (
	"fmt"
	"regexp"
	"strings"
)

var categoryPatterns = map[string][]*regexp.Regexp{
	"behavioral": {
		regexp.MustCompile(`(?i)(tell me about a time|give me an example|describe a situation|when you (had|have|faced|dealt|worked|led|managed|failed|struggled)|biggest (challenge|achievement|failure|mistake|success)|how did you handle|walk me through a time|have you ever|conflict with|difficult (coworker|colleague|manager|teammate)|under pressure|tight deadline|disagree(d)? with|took initiative|learned (quickly|fast|new)|proud of)`),
		regexp.MustCompile(`(?i)(说说你遇到的|举个例子|描述一次|遇到过的最大挑战|最成功的一件事|最有成就感|失败的经历|如何处理冲突|团队分歧|压力最大|主导的项目|在紧急情况下|复盘)`),
	},
	"motivation": {
		regexp.MustCompile(`(?i)(why (do you want|are you interested|this company|this role|us|here)|why (are you leaving|did you leave|move on)|what (attracted|draws|interests|excites|appeals) (you|to)|where do you see yourself|5 years|career goals|ideal (role|company|environment)|motivates you|why should we hire|what (do you|would you) bring|new opportunity)`),
		regexp.MustCompile(`(?i)(为什么想来我们公司|为什么考虑换工作|为什么离职|未来的职业规划|你对我们了解多少|你觉得你能带来什么|五年规划|求职动机|为什么选择我们)`),
	},
	"situational": {
		regexp.MustCompile(`(?i)(what would you do if|how would you (handle|approach|deal with)|imagine you|hypothetically|if you (joined|started|were)|prioritize|production (outage|incident|down)|codebase (is a mess|legacy)|disagree with (your manager|a decision)|first (30|60|90) days)`),
		regexp.MustCompile(`(?i)(如果你遇到|如果线上出现bug|系统宕机|线上故障|如果你和上级意见不合|需求延期|面对紧急突发|如果你接手遗留代码|入职前三个月)`),
	},
	"experience": {
		regexp.MustCompile(`(?i)(tell me about your (experience|background|role|work)|walk me through your (resume|background|career)|tell me about yourself|what were you responsible|biggest project|tech stack|day.to.day|what did you build|recent role|previous job)`),
		regexp.MustCompile(`(?i)(自我介绍|介绍一下你自己|介绍一下你的项目|讲讲你最近做过的项目|简历中|负责什么模块|核心职责|技术栈|在这个项目中扮演什么角色)`),
	},
	"compensation": {
		regexp.MustCompile(`(?i)(salary (expectation|requirement|range)|compensation|how much (are you|do you) (making|expect|want)|when can you start|notice period|start date|other (offer|interview)|questions for us|questions for me|do you have any questions)`),
		regexp.MustCompile(`(?i)(期望薪资|薪水要求|薪资预期|什么时候可以到岗|入职时间|有其他offer吗|还在看其他机会吗|你有什么想问我们的|有什么问题要问我)`),
	},
	"technical": {
		regexp.MustCompile(`(?i)(system design|design (a|an|the) (system|service|api|database)|tradeoff|difference between|how does .{2,40} work|complexity|algorithm|data structure|scale (this|to|it)|architecture|pros and cons|how would you design)`),
		regexp.MustCompile(`(?i)(系统设计|底层原理|源码分析|优缺点|区别是什么|如何设计一个|高并发|缓存击穿|雪崩|锁机制|线程池|索引原理|分库分表|架构设计|时空复杂度)`),
	},
}

func DetectCategory(turns []Turn) string {
	if len(turns) == 0 {
		return "general"
	}

	// Examine recent "them" (interviewer) speech
	var recentThem []string
	count := 0
	for i := len(turns) - 1; i >= 0 && count < 5; i-- {
		if turns[i].Channel == "them" {
			recentThem = append([]string{turns[i].Text}, recentThem...)
			count++
		}
	}
	combined := strings.Join(recentThem, " ")
	if strings.TrimSpace(combined) == "" {
		// If no them channel turns, check last user turns
		for i := len(turns) - 1; i >= 0 && count < 3; i-- {
			recentThem = append([]string{turns[i].Text}, recentThem...)
			count++
		}
		combined = strings.Join(recentThem, " ")
	}

	if combined == "" {
		return "general"
	}

	for cat, patterns := range categoryPatterns {
		for _, p := range patterns {
			if p.MatchString(combined) {
				return cat
			}
		}
	}

	return "general"
}

func clipText(s string, limit int) string {
	runes := []rune(s)
	if len(runes) <= limit {
		return s
	}
	return string(runes[:limit]) + "..."
}

func BuildInterviewContext(s Settings, mode string, turns []Turn) string {
	if mode == "leetcode" {
		return ""
	}

	category := DetectCategory(turns)
	var blocks []string

	// Resume context
	if strings.TrimSpace(s.ResumeText) != "" {
		limit := 1600
		if category == "behavioral" || category == "experience" {
			limit = 2400
		}
		blocks = append(blocks, fmt.Sprintf("=== 候选人个人经历与简历信息 ===\n%s", clipText(strings.TrimSpace(s.ResumeText), limit)))
	}

	// Job description
	if strings.TrimSpace(s.JobDescription) != "" {
		blocks = append(blocks, fmt.Sprintf("=== 目标岗位要求 (JD) ===\n%s", clipText(strings.TrimSpace(s.JobDescription), 800)))
	}

	// Category specific injections
	switch category {
	case "behavioral":
		if strings.TrimSpace(s.StarStories) != "" {
			blocks = append(blocks, fmt.Sprintf("=== 真实 STAR 故事库 (优先引用) ===\n%s\n提示：回答行为面试题时，采用情境(Situation) -> 任务(Task) -> 行动(Action) -> 结果(Result) 结构，包含具体数据与度量。", clipText(s.StarStories, 2000)))
		}
		if strings.TrimSpace(s.WorkStyle) != "" {
			blocks = append(blocks, fmt.Sprintf("工作理念与风格：\n%s", clipText(s.WorkStyle, 400)))
		}
	case "motivation":
		if strings.TrimSpace(s.WhyCompany) != "" {
			blocks = append(blocks, fmt.Sprintf("为什么选择贵司：\n%s", clipText(s.WhyCompany, 500)))
		}
		if strings.TrimSpace(s.WhyLeaving) != "" {
			blocks = append(blocks, fmt.Sprintf("离职原因与考虑：\n%s", clipText(s.WhyLeaving, 400)))
		}
		if strings.TrimSpace(s.WorkStyle) != "" {
			blocks = append(blocks, fmt.Sprintf("工作风格与偏好环境：\n%s", clipText(s.WorkStyle, 400)))
		}
	case "situational":
		if strings.TrimSpace(s.WorkStyle) != "" {
			blocks = append(blocks, fmt.Sprintf("决策风格与排期原则：\n%s", clipText(s.WorkStyle, 500)))
		}
		if strings.TrimSpace(s.StarStories) != "" {
			blocks = append(blocks, fmt.Sprintf("过往应对危机/排期的参考经历：\n%s", clipText(s.StarStories, 800)))
		}
	case "experience":
		if strings.TrimSpace(s.StarStories) != "" {
			blocks = append(blocks, fmt.Sprintf("代表作与项目亮点：\n%s", clipText(s.StarStories, 1000)))
		}
	case "compensation":
		if strings.TrimSpace(s.SalaryTarget) != "" {
			blocks = append(blocks, fmt.Sprintf("期望薪酬范围：\n%s", s.SalaryTarget))
		}
		if strings.TrimSpace(s.QuestionsToAsk) != "" {
			blocks = append(blocks, fmt.Sprintf("向面试官反问的问题清单：\n%s", clipText(s.QuestionsToAsk, 600)))
		}
	case "technical":
		// Resume & JD already included
	default:
		if strings.TrimSpace(s.StarStories) != "" {
			blocks = append(blocks, fmt.Sprintf("核心亮点与经验：\n%s", clipText(s.StarStories, 600)))
		}
	}

	if len(blocks) == 0 {
		return ""
	}

	res := strings.Join(blocks, "\n\n")
	if strings.TrimSpace(s.JobDescription) != "" {
		res += "\n\n注意：请结合上述目标岗位要求与候选人经历，针对性突出人岗匹配度。"
	}
	return res
}
