package internal

import (
	"fmt"
	"strings"
)

type ModeDef struct {
	NeedsScreen        bool
	Small              bool
	TranscriptRequired bool
	SystemPrompt       func(contextBlock string, aiRules string) string
	BuildUserPrompt    func(turns []Turn, userText string) string
}

func FormatTranscript(turns []Turn, limit int) string {
	if len(turns) == 0 {
		return ""
	}
	recent := turns
	if limit > 0 && len(turns) > limit {
		recent = turns[len(turns)-limit:]
	}
	var lines []string
	for _, t := range recent {
		who := "面试官 (Them): "
		if t.Channel == "you" {
			who = "我 (You): "
		}
		lines = append(lines, who+t.Text)
	}
	return strings.Join(lines, "\n")
}

func appendAiRules(basePrompt string, aiRules string) string {
	if strings.TrimSpace(aiRules) == "" {
		return basePrompt
	}
	return basePrompt + "\n\n=== 用户自定义回复规则 (必须严格遵守) ===\n" + strings.TrimSpace(aiRules)
}

const BaseInstruction = "用自然、流利、自信的语气回答（若对方使用中文则用中文回答，若对方使用英文则用英文回答）。以第一人称直接输出候选人口头要说的话，不要有任何客套废话或引导句（如'你可以这样说：'或引号等）。"

var Modes = map[string]ModeDef{
	"assist": {
		NeedsScreen:        true,
		Small:              false,
		TranscriptRequired: false,
		SystemPrompt: func(contextBlock string, aiRules string) string {
			sys := "你是一款悬浮在用户屏幕上的实时AI面试与编码副驾驶（Copilot）。\n" +
				BaseInstruction + "\n\n" +
				"根据屏幕截屏和最近的对话记录，判断候选人当下最需要的内容并直接输出：\n" +
				"• 行为面试题 (BEHAVIORAL)：给出完整的 STAR 回答（情境、任务、行动、结果），融合候选人真实经历，清晰且有量化数据。\n" +
				"• 动机与求职意愿 (MOTIVATION)：给出真诚且具有吸引力的具体回答。\n" +
				"• 情境应对题 (SITUATIONAL)：展现结构化思维与果断的决策判断力。\n" +
				"• 项目与履历 (EXPERIENCE)：基于简历给出深入、自豪的技术细节解答。\n" +
				"• 技术/系统设计 (TECHNICAL)：条理清晰，先结论后论据，给出关键技术选型与权衡。\n" +
				"• 编程算法题 (LEETCODE)：给出思路要点 + 规范代码 + 复杂度。\n" +
				"• 薪酬意向/反问：给出自信得体的答复。"
			if contextBlock != "" {
				sys = contextBlock + "\n\n" + sys
			}
			return appendAiRules(sys, aiRules)
		},
		BuildUserPrompt: func(turns []Turn, userText string) string {
			t := FormatTranscript(turns, 14)
			if t == "" {
				t = "(暂无对话转录记录，请主要依据截屏内容)"
			}
			return fmt.Sprintf("最近对话记录：\n%s\n\n请结合屏幕截屏内容，给出我现在应该说出的最佳回答。", t)
		},
	},

	"say": {
		NeedsScreen:        false,
		Small:              false,
		TranscriptRequired: false,
		SystemPrompt: func(contextBlock string, aiRules string) string {
			sys := "你是一款在现场面试中向候选人实时提供提示的隐形副驾。\n" +
				BaseInstruction + "\n\n" +
				"针对面试官最近的提问或对话，提供一段自然、沉稳、高水准的第一人称回答。\n" +
				"直接输出实际口述的文本内容，无需多余解释，长度通常在 2~5 句话（约 60~120 字），直击痛点。"
			if contextBlock != "" {
				sys = contextBlock + "\n\n" + sys
			}
			return appendAiRules(sys, aiRules)
		},
		BuildUserPrompt: func(turns []Turn, userText string) string {
			t := FormatTranscript(turns, 16)
			if t == "" {
				t = "(尚未开始语音转录)"
			}
			return fmt.Sprintf("当前正在进行的面试对话记录：\n%s\n\n面试官刚说完，我接下来应该怎么自然自信地回答？", t)
		},
	},

	"recap": {
		NeedsScreen:        false,
		Small:              true,
		TranscriptRequired: true,
		SystemPrompt: func(contextBlock string, aiRules string) string {
			sys := "你是一款专业的面试与会议复盘专家。\n" +
				"请依据完整的对话转录记录，对本场对话进行清晰结构化的总结回顾：\n" +
				"• 讨论主题 (Topics Covered)：双方实际讨论的核心议题与关键信息。\n" +
				"• 面试官提问 (Questions Asked)：面试官提出的重要问题清单。\n" +
				"• 核心观点与承诺 (Key Points & Answers)：双方达成的共识、技术方案或候选人的核心阐述。\n" +
				"• 待跟进/薄弱点 (Open Threads)：需要进一步巩固、补充说明或后续跟进的事项。\n" +
				"采用加粗小标题 + 简练的项目符号，紧扣事实，避免空洞套话。"
			if contextBlock != "" {
				sys = contextBlock + "\n\n" + sys
			}
			return appendAiRules(sys, aiRules)
		},
		BuildUserPrompt: func(turns []Turn, userText string) string {
			t := FormatTranscript(turns, 0)
			return fmt.Sprintf("本次会议/面试的完整转录记录如下：\n%s\n\n请对此进行全面而精炼的回顾与复盘总结。", t)
		},
	},

	"ask": {
		NeedsScreen:        true,
		Small:              false,
		TranscriptRequired: false,
		SystemPrompt: func(contextBlock string, aiRules string) string {
			sys := "你是一款全能的实时面试与技术助手，拥有当前屏幕与对话上下文的感知能力。\n" +
				BaseInstruction + "\n\n" +
				"请直接、精准地回答候选人的输入问题。如果涉及个人经历，使用简历中的真实背景；如果是技术概念或代码，给出清晰透彻的解答与范例。"
			if contextBlock != "" {
				sys = contextBlock + "\n\n" + sys
			}
			return appendAiRules(sys, aiRules)
		},
		BuildUserPrompt: func(turns []Turn, userText string) string {
			t := FormatTranscript(turns, 12)
			var sb strings.Builder
			if t != "" {
				sb.WriteString("最近对话背景：\n" + t + "\n\n")
			}
			sb.WriteString("问题：" + userText)
			return sb.String()
		},
	},

	"answerThis": {
		NeedsScreen:        false,
		Small:              false,
		TranscriptRequired: false,
		SystemPrompt: func(contextBlock string, aiRules string) string {
			sys := "你是一款向候选人精准提示特定问题答案的面试副驾。\n" +
				BaseInstruction + "\n\n" +
				"用户提供了一句面试官的特定提问。请专注于针对该问题给出完美的第一人称回答：\n" +
				"• 突出优势与亮点，逻辑层次分明；\n" +
				"• 行为题使用 STAR 原则；技术题先给结论和最佳实践；\n" +
				"• 直接输出口述话术，2~5 句话，无需多余前言。"
			if contextBlock != "" {
				sys = contextBlock + "\n\n" + sys
			}
			return appendAiRules(sys, aiRules)
		},
		BuildUserPrompt: func(turns []Turn, userText string) string {
			return fmt.Sprintf("请针对面试官的这句具体提问给出高分回答：\n\"%s\"\n\n直接输出我口头要说的答案：", userText)
		},
	},

	"leetcode": {
		NeedsScreen:        true,
		Small:              false,
		TranscriptRequired: false,
		SystemPrompt: func(contextBlock string, aiRules string) string {
			// LeetCode mode stays strictly professional without personal chatter
			return "你是一位顶级算法竞赛专家与资深技术面试官。\n" +
				"截图中包含一道编程/算法题目。\n" +
				"请严格按以下四部分输出结构化解答：\n" +
				"1. 【题目重述】：一句话总结题目核心目标与边界限制；\n" +
				"2. 【解题思路】：阐述最优解法（如双指针、单调栈、动态规划等）与推导逻辑；\n" +
				"3. 【代码实现】：给出高效、规范、注释得当的代码块（若截屏中有指定语言，使用该语言，否则默认使用 Go 或 Python）；\n" +
				"4. 【复杂度分析】：明确给出时间复杂度 O(...) 与空间复杂度 O(...) 并作简要说明。"
		},
		BuildUserPrompt: func(turns []Turn, userText string) string {
			return "请解答截屏中所显示的算法与编程题目。"
		},
	},
}
