package model

var KeywordExtractionPrompt = `
# Role: You are an SEO expert and text analyst.
## Task:
Based on the provided question, extract the relevant keywords.
## Requirements:
1. Extract keywords in the same language as the question.
2. Output the keywords in JSON format.
## Example:
### User 
Question: 如何提高网站的转化率？
### System
["转化率", "网站优化", "用户行为分析"]
## Question: {{.question}}
`

var QuestionOptimizationPrompt = `
# 角色：你是一位提问分析师。
## 任务：
基于用户提供的问题，优化问题的措辞。
## 要求：
1. 优化后的问题必须与用户提供的问题相同的语言。
2. 输出优化后的问题，以原始文本格式输出。
3. 将用户的问题简化，只保留问题的关键部分。
4. 修复问题中的语法和拼写错误。
## 示例1：
### 用户
问题: 我现在有100万条数据需要存入数据库，该使用哪种数据库？
### 系统
大规模数据库有哪些？
## 示例2：
### 用户
问题：如果我要请假5天以上，需要向谁申请？
### 系统
公司请假申请流程是什么？
## 问题: 
{{.question}}
`

var OCRPrompt = `
# Role: You are an text analyst.
## Task:
Extract the text from the image.
## Requirements:
1. Output the keywords in raw text format.
2. If there is no text in the image, output "".`
