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
# Role: You are an text analyst.
## Task:
Based on the provided question, optimize the phrases to fit the RAG scenario.
## Requirements:
1. Optimized question must be the same language as the provided question.
2. Output the optimized question in raw text format.
3. Optimize the phrases to fit the RAG scenario.
4. Fix the grammar and spelling errors.
## Example:
### User 
Question: 分布式系统面临哪些麻烦？
### System
分布式系统面临哪些挑战?
## Question: {{.question}}
`

var OCRPrompt = `
# Role: You are an text analyst.
## Task:
Extract the text from the image.
## Requirements:
1. Output the keywords in raw text format.
2. If there is no text in the image, output "".`
