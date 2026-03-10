package template

import (
	"os"
	"path/filepath"
)

var Templates = map[string]string{
	"problem":   problemTemplate,
	"decision":  decisionTemplate,
	"knowledge": knowledgeTemplate,
	"practice":  practiceTemplate,
}

// GetTemplate reads template from templatesDir
func GetTemplate(noteType string, templatesDir string) (string, error) {
	path := filepath.Join(templatesDir, noteType+".md")
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// GetBuiltinTemplate returns the built-in template (fallback)
func GetBuiltinTemplate(noteType string) string {
	return Templates[noteType]
}

func WriteTemplates(templatesDir string) error {
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		return err
	}

	for noteType, content := range Templates {
		filename := noteType + ".md"
		path := filepath.Join(templatesDir, filename)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}

const problemTemplate = `---
title: ""
type: problem
tags: []
private: false
created: {{DATE}}
updated: {{DATE}}
---

## 환경

- OS / 런타임:
- 관련 스택 및 버전:
- 배포 환경: (로컬 / 스테이징 / 프로덕션)

## 증상

- 기대: 어떻게 동작해야 했는가
- 실제: 어떻게 동작했는가 (에러 메시지, 예상과 다른 동작 등)

## 실패한 시도들

- [ ] **시도한 것**
  - 선택 이유: 왜 이게 될 것 같았는가
  - 실패 이유: 왜 안 됐는가
  - 부작용: 이 시도로 인해 생긴 다른 문제 (없으면 생략)

## 원인

- 근본 원인:
- 트리거 조건: 어떤 상황에서 발생했는가 (항상 / 특정 조건에서만 등)

## 해결

- 해결 방법:
- 적용한 코드 / 명령어:
  ` + "```" + `

  ` + "```" + `
- 검증 방법: 어떻게 해결됐음을 확인했는가

## 교훈

- 핵심 원인 패턴: 이런 류의 문제가 왜 생기는가
- 다음에 먼저 확인할 것:
- 예방법: 같은 문제를 피하려면
`

const decisionTemplate = `---
title: ""
type: decision
tags: []
private: false
created: {{DATE}}
updated: {{DATE}}
---

## 배경

- 상황:
- 해결하려는 문제:
- 제약 조건: (시간 / 비용 / 팀 역량 / 기술 스택 등)

## 검토한 대안들

- [ ] **대안 A**
  - 제외 이유:
  - 참고: 성능 / 확장성 / 유지보수성 / 운영 부담 / 구현 난이도 / 리스크 / 비용

## 결정

- 선택:
- 선택 이유:
- 기대 효과:

## 트레이드오프

- 포기한 것:
- 감수한 리스크:
- 재검토 조건: 어떤 상황이 오면 이 결정을 다시 고려해야 하는가

## 교훈

- 이 결정에서 배운 것:
- 다음에 비슷한 결정을 할 때 먼저 확인할 것:
`

const knowledgeTemplate = `---
title: ""
type: knowledge
tags: []
private: false
created: {{DATE}}
updated: {{DATE}}
---

## 개요

- 한 줄 설명:
- 어떤 문제를 해결하는가:
- 언제 쓰는가:

## 학습 배경

- 계기: 어떤 상황에서 이 지식이 필요했는가
- 관련 작업 / 문제:

## 핵심 개념

- **개념 A**: 설명
- **개념 B**: 설명

## 사용법

- 기본 사용:
  ` + "```" + `

  ` + "```" + `
- 주요 옵션 / 패턴:
  ` + "```" + `

  ` + "```" + `

## 적용 사례

- **프로젝트 A**
  - 어떤 맥락에서 사용했는가:
  - 결과 / 효과:

## 주의사항

- **함정 A**: 설명
  - 증상:
  - 해결:
`

const practiceTemplate = `---
title: ""
type: practice
tags: []
private: false
created: {{DATE}}
updated: {{DATE}}
---

## 목표

- 무엇을 익히려 하는가:
- 완료 기준:

## 환경

- OS / 런타임:
- 관련 스택 및 버전:

## 실습 로그

### {{DATE}}

- 한 것:
- 결과:
- 막힌 것:

#### Q&A

- **Q**:
  - **A**:

## 최종 정리

- 배운 것:
- 아직 모르는 것:
`
