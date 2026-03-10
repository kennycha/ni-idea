package skill

import (
	"os"
	"path/filepath"
)

const SearchSkillContent = `---
name: ni-idea-search
description: 개인 지식 베이스에서 과거 문제 해결 기록, 기술 결정, 지식을 검색합니다. 에러 해결이나 기술 결정 시 관련 경험을 찾을 때 사용하세요.
argument-hint: "[검색 키워드]"
---

# ni-idea 탐색

## 사용 시점

- 에러 해결 시 과거 유사 문제 검색
- 기술 결정 시 이전 결정 사항 참조
- 특정 기술/도구 관련 지식 조회

## 명령어

### 검색
` + "```bash" + `
ni search "키워드"           # 문제/결정 검색 (기본)
ni search "키워드" --all     # 전체 타입 검색
ni search "키워드" --tag k8s # 태그 필터
` + "```" + `

### 조회
` + "```bash" + `
ni get problems/nextjs-caching  # 특정 노트 조회
ni list                         # 전체 목록
ni list --tag infra             # 태그별 목록
ni tags                         # 태그 목록
` + "```" + `

## 워크플로우

1. $ARGUMENTS에서 키워드 추출 (없으면 사용자에게 질문)
2. ` + "`ni search \"키워드\"`" + ` 실행
3. 관련 노트 발견 시 ` + "`ni get <path>`" + `로 상세 조회
4. 노트 내용을 참고하여 답변
`

const AddSkillContent = `---
name: ni-idea-add
description: 문제 해결 과정, 기술 결정, 새로 배운 지식을 개인 지식 베이스에 기록합니다. 작업 완료 후 경험을 정리할 때 사용하세요.
argument-hint: "[노트 타입: problem|decision|knowledge|practice]"
---

# ni-idea 노트 추가

## 사용 시점

- 에러 해결 후 해결 과정 기록
- 기술 결정 후 결정 사항 문서화
- 새로운 기술/개념 학습 후 정리
- 실습/튜토리얼 진행 중 로그 기록

## 노트 타입

| 타입 | 용도 |
|------|------|
| problem | 문제 해결 기록 |
| decision | 기술/아키텍처 결정 |
| knowledge | 개념/기술 문서 |
| practice | 실습 로그 |

## 명령어

` + "```bash" + `
ni add --type problem --title "제목"
ni add --type decision --title "제목" --tag infra,k8s
ni add --type knowledge --title "제목"
ni add --type practice --title "제목"
` + "```" + `

## 가이드라인

- 하나의 세션을 꼭 하나의 노트로 만들 필요는 없음
- 세션에서 여러 문제를 해결했거나, 여러 결정을 내렸다면 각각 별도 노트로 분리
- 작고 집중된 노트가 검색과 재사용에 유리함

## 워크플로우

1. $ARGUMENTS에서 타입 확인 (없으면 대화 내용에서 추론하거나 질문)
2. 대화 내용을 바탕으로 노트 내용 구성 (필요시 여러 노트로 분리)
3. ` + "`ni add --type <type> --title \"제목\" --tag tag1,tag2 --no-edit`" + ` 실행
4. 생성된 파일 경로 확인
5. 필요시 추가 내용 작성
`

func InstallSkills() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	skillsDir := filepath.Join(home, ".claude", "skills")

	// ni-idea-search
	searchDir := filepath.Join(skillsDir, "ni-idea-search")
	if err := os.MkdirAll(searchDir, 0755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(searchDir, "SKILL.md"), []byte(SearchSkillContent), 0644); err != nil {
		return err
	}

	// ni-idea-add
	addDir := filepath.Join(skillsDir, "ni-idea-add")
	if err := os.MkdirAll(addDir, 0755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(addDir, "SKILL.md"), []byte(AddSkillContent), 0644); err != nil {
		return err
	}

	return nil
}
