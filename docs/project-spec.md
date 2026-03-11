# ni-idea 프로젝트 스펙

## 개요

지식(문제 해결, 인프라, 프론트엔드, 회사 도메인 등)을 로컬에 저장하고, CLI를 통해 빠르게 검색·조회할 수 있는 개인 지식 베이스 도구.

Claude Code가 `bash`를 통해 자연스럽게 호출하는 것을 1차 목표로 한다.

---

## 목표

- 단일 바이너리 CLI (`ni`) 로 어느 환경에서도 즉시 실행 가능
- 마크다운 기반 노트로 사람이 직접 읽고 쓰기 쉬운 구조 유지
- Claude Code가 `CLAUDE.md` 지시에 따라 자율적으로 `ni` 를 호출해 컨텍스트를 보강
- 점진적으로 기능 확장 가능한 구조 (검색 → 인덱싱 → 임베딩 등)

---

## 기술 스택

| 항목        | 선택                            | 이유                                            |
| ----------- | ------------------------------- | ----------------------------------------------- |
| 언어        | Go                              | 단일 바이너리, 크로스플랫폼, 런타임 의존성 없음 |
| 노트 포맷   | Markdown + YAML 프론트매터      | 사람이 읽기 쉽고 메타데이터 구조화 가능         |
| 검색 (초기) | 파일 시스템 grep                | 단순하고 의존성 없음                            |
| 검색 (확장) | bleve 또는 외부 인덱서          | 필요 시 도입                                    |
| 설정        | `~/.config/ni-idea/config.yaml` | XDG 규격 준수                                   |

---

## 디렉토리 구조

### 바이너리 / 소스

```
ni-idea/
├── cmd/
│   └── ni/
│       └── main.go
├── internal/
│   ├── search/         # 검색 로직
│   ├── store/          # 노트 읽기/쓰기
│   ├── formatter/      # stdout 출력 포맷
│   └── config/         # 설정 로딩
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

### 노트 저장소 (기본 경로: `~/notes` 또는 설정으로 변경)

```
notes/
├── problems/           # 문제 해결 기록 (기본 검색 대상)
├── decisions/          # 아키텍처·기술 결정 기록 (기본 검색 대상)
├── knowledge/          # 개념·기술 정리 (--all 옵션 시 검색)
├── practice/           # 실습 노트 (--all 옵션 시 검색)
└── domains/            # 회사별 도메인 지식
    ├── company-a/
    └── company-b/
```

---

## 노트 포맷

모든 노트는 YAML 프론트매터 + 마크다운 본문으로 구성한다. 타입별 템플릿은 아래와 같다.

---

### problem 템플릿

```markdown
---
title: ""
type: problem
tags: []
domain: general
private: false
created: YYYY-MM-DD
updated: YYYY-MM-DD
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
```

```
- 검증 방법: 어떻게 해결됐음을 확인했는가

## 교훈

- 핵심 원인 패턴: 이런 류의 문제가 왜 생기는가
- 다음에 먼저 확인할 것:
- 예방법: 같은 문제를 피하려면
```

---

### decision 템플릿

```markdown
---
title: ""
type: decision
tags: []
domain: general
private: false
created: YYYY-MM-DD
updated: YYYY-MM-DD
---

## 배경

- 상황:
- 해결하려는 문제:
- 제약 조건: (시간 / 비용 / 팀 역량 / 기술 스택 등)

## 검토한 대안들

- [ ] **대안 A**
  - 제외 이유:
  - 참고 (관련 있는 것만):
    - 성능 / 확장성 / 유지보수성 / 운영 부담 / 구현 난이도 / 리스크 / 비용

## 결정

- 선택:
- 선택 이유:
- 기대 효과:
- 평가 근거 (관련 있는 것만):
  - 성능:
    - 평가:
    - 근거:
  - 확장성:
    - 평가:
    - 근거:
  - 유지보수성:
    - 평가:
    - 근거:
  - 운영 부담:
    - 평가:
    - 근거:
  - 구현 난이도:
    - 정도: (낮음 / 중간 / 높음)
    - 근거:
  - 리스크:
    - 평가:
    - 근거:
  - 비용:
    - 평가:
    - 근거:

## 트레이드오프

- 포기한 것:
- 감수한 리스크:
- 재검토 조건: 어떤 상황이 오면 이 결정을 다시 고려해야 하는가

## 교훈

- 이 결정에서 배운 것:
- 다음에 비슷한 결정을 할 때 먼저 확인할 것:
```

---

### knowledge 템플릿

```markdown
---
title: ""
type: knowledge
tags: []
domain: general
private: false
created: YYYY-MM-DD
updated: YYYY-MM-DD
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
```

```
- 주요 옵션 / 패턴:
```

```
- 실제 사용 예시:
```

```

## 적용 사례

- **프로젝트 A**
- 어떤 맥락에서 사용했는가:
- 어떻게 적용했는가:
- 결과 / 효과:

## 주의사항

- **함정 A**: 설명
- 증상:
- 해결:
```

---

### practice 템플릿

단일 주제에 대해 여러 날짜에 걸쳐 이어서 작성한다.

```markdown
---
title: ""
type: practice
tags: []
domain: general
private: false
created: YYYY-MM-DD
updated: YYYY-MM-DD
---

## 목표

- 무엇을 익히려 하는가:
- 완료 기준:

## 환경

- OS / 런타임:
- 관련 스택 및 버전:

## 실습 로그

### YYYY-MM-DD

- 한 것:
- 결과:
- 막힌 것:

## 최종 정리

- 배운 것:
- 아직 모르는 것:
```

---

## CLI 명령어 스펙

### `ni search <query>`

키워드로 노트를 검색한다. 기본적으로 `problem`, `decision` 타입만 검색하며, `--all` 옵션으로 전체 타입을 포함한다.

```bash
ni search "Next.js 캐싱"
ni search "캐싱" --tag nextjs
ni search "배포" --domain company-a
ni search "k8s" --type problem
ni search "캐싱" --all         # knowledge, practice 포함 전체 검색
```

**옵션**

| 플래그              | 설명                                                         |
| ------------------- | ------------------------------------------------------------ |
| `--tag`             | 특정 태그로 필터                                             |
| `--domain`          | 도메인 필터 (general, company-a 등)                          |
| `--type`            | 노트 타입 직접 지정 (problem, decision, knowledge, practice) |
| `--all`             | knowledge, practice 타입 포함 전체 검색                      |
| `--limit`           | 결과 수 제한 (기본값: 10)                                    |
| `--include-private` | `private: true` 노트 포함 (기본값: 제외)                     |

**출력**: 매칭된 노트 목록 (제목, 경로, 태그, 매칭 라인 스니펫)

---

### `ni get <path>`

특정 노트의 전체 내용을 출력한다. 경로를 직접 지정하므로 `private: true` 노트도 항상 접근 가능.

```bash
ni get problems/frontend/nextjs-caching
ni get knowledge/infra/k8s-basics
```

**출력**: 마크다운 전문 (stdout)

---

### `ni list`

노트 목록을 출력한다.

```bash
ni list
ni list --tag infra
ni list --domain company-a
ni list problems/
ni list --include-private   # private 노트 포함
```

**출력**: 제목, 경로, 태그, 수정일 목록

---

### `ni add`

새 노트를 추가한다. 인터랙티브 모드 또는 플래그로 제어.

```bash
ni add                              # 인터랙티브
ni add --title "제목" --tag infra --type problem
```

에디터(`$EDITOR`)를 열어 직접 작성하게 한다.

---

### `ni tags`

사용 중인 태그 목록과 노트 수를 출력한다.

```bash
ni tags
```

---

## 출력 포맷 원칙

Claude Code가 소비하기 좋도록:

- 기본 출력은 **plain markdown** (stdout)
- 에러는 **stderr**
- `--json` 플래그로 JSON 출력 전환 가능 (선택적 구현)
- 결과 없을 때는 빈 출력 대신 명확한 메시지 출력

---

## 설정 파일

`~/.config/ni-idea/config.yaml`

```yaml
notes_dir: ~/notes # 노트 저장 경로
editor: nvim # ni add 시 사용할 에디터
default_domain: general # 기본 도메인
search:
  max_results: 10
```

---

## CLAUDE.md 연동 예시

프로젝트마다 `CLAUDE.md`에 아래와 같이 작성하면 Claude Code가 자동으로 `ni`를 활용한다.

```markdown
## 지식 베이스 (ni-idea)

관련 기술 결정이나 도메인 질문은 작업 전 반드시 `ni search`로 먼저 조회할 것.

- 기본 검색 (문제 해결·결정 선례): `ni search <키워드>`
- 이 회사 도메인: `ni search <키워드> --domain company-a`
- 레퍼런스 포함 전체 검색: `ni search <키워드> --all`

조회 결과가 있을 경우, 해당 맥락을 반영해서 설계 및 코드 작성.
```

---

## Agent Skills (npx skills 연동)

`ni-idea` 레포에 `SKILL.md`를 포함해두면, 다른 프로젝트에서 작업 중 Claude Code가 노트를 추가해야 할 때 올바른 포맷으로 `ni add`를 자율 호출할 수 있다.

### 설치

```bash
npx skills add kennycha/ni-idea --skill add-note -g -a claude-code
```

### 디렉토리 구조 추가

```
ni-idea/
└── skills/
    └── add-note/
        └── SKILL.md
```

### SKILL.md 초안

```markdown
---
name: add-note
description: >
  ni-idea 지식 베이스에 노트를 추가한다.
  문제를 해결했거나, 중요한 결정을 내렸거나, 30분 이상 삽질한 경우 이 스킬을 사용한다.
---

# ni-idea 노트 추가

작업 중 아래 상황에 해당하면 ni add로 지식 베이스에 노트를 남긴다.

## 트리거 조건

- 에러나 문제를 해결했을 때
- 아키텍처나 기술 결정을 내렸을 때
- 삽질이 30분 이상 걸렸을 때
- 나중에 다시 찾아볼 것 같은 패턴을 발견했을 때

## 노트 작성 방법

stdin으로 마크다운을 주입해 비인터랙티브하게 추가한다.

    ni add --title "<제목>" --tag <tag1>,<tag2> --type <type> [--domain <domain>] [--private] <<'EOF'
    ## 상황
    <어떤 맥락에서 발생했는가>

    ## 시도한 것들
    <뭘 해봤는가>

    ## 원인
    <왜 발생했는가>

    ## 해결
    <어떻게 해결했는가>

    ## 교훈
    <다음에 기억할 것>
    EOF

## 필드 가이드

| 필드      | 값                          | 설명                                |
| --------- | --------------------------- | ----------------------------------- |
| --type    | problem                     | 문제 해결 기록 (기본 검색 대상)     |
| --type    | decision                    | 아키텍처·기술 결정 (기본 검색 대상) |
| --type    | knowledge                   | 개념·기술 정리 (--all 시 검색)      |
| --type    | practice                    | 실습 노트 (--all 시 검색)           |
| --tag     | infra, frontend, backend 등 | 복수 지정 가능, 쉼표 구분           |
| --domain  | general, company-a 등       | 생략 시 general                     |
| --private | 플래그                      | 기본 검색에서 제외할 민감한 노트    |

## 주의사항

- --title 은 간결하게, 나중에 검색어가 될 키워드를 포함할 것
- 회사 내부 정보가 포함된 경우 반드시 --private 플래그 추가
- 노트 추가 후 ni search <키워드>로 정상 검색되는지 확인
```

---

## 구현 순서 (Phase)

### Phase 1 — MVP ✅

- [x] 설정 파일 로딩
- [x] `ni init` — 노트 디렉토리 + 타입별 템플릿 파일 생성
- [x] `ni search` — 파일명 + 본문 grep 기반 검색
- [x] `ni get` — 노트 전문 출력
- [x] `ni list` — 목록 출력
- [x] `ni add` — 에디터 연동 (인터랙티브) + stdin 주입 (비인터랙티브) 노트 추가
- [x] `ni tags` — 태그 목록
- [x] Claude Code skill 연동 (`ni-idea-search`, `ni-idea-add`)

### Phase 2 — 검색 고도화 ✅

- [x] 인덱싱 도입 (bleve)
- [x] 퍼지 검색 (`--fuzzy` 플래그)
- [ ] 임베딩 기반 의미 검색 (로컬 모델 또는 API) — 추후

### Phase 3 — 원격 저장소 (클라우드/팀 동기화) ✅

- [x] `ni remote add/list/remove` — 리모트 관리
- [x] `ni push` — 로컬 → 서버 업로드
- [x] `ni pull` — 서버 → 로컬 다운로드
- [x] 충돌 감지 및 해결 (`--force`, `--theirs`, `--ours`)
- [x] Go 서버 구현 (파일 기반 저장, 토큰 인증)
- [ ] 리모트 검색 (`--remote` 플래그) — 추후

**CLI 명령어**

```bash
# 리모트 관리
ni remote add <name> <url>
ni remote list
ni remote remove <name>

# 동기화
ni push <path>              # 특정 노트 업로드
ni push --all               # private 제외 전체 업로드
ni push --force             # 충돌 무시, 강제 업로드
ni pull                     # 리모트 노트 로컬로 동기화
ni pull --theirs            # 충돌 시 서버 버전 사용
ni pull --ours              # 충돌 시 로컬 버전 유지
```

**설정**

```yaml
remotes:
  - name: personal
    url: https://my-ni.vercel.app
    token: your-token
```

**서버 실행**

```bash
NI_AUTH_TOKENS=token1,token2 ./ni-server --port 8080 --data ./data
```

**규칙**

- `private: true` 노트는 명시적 `--include-private` 없이 push 불가
- 충돌 감지: `updated` 타임스탬프 비교
- 인덱스는 sync 제외 (pull 후 자동 재인덱싱)

---

## 비고

- 공개 가능한 일반 지식과 내부 도메인 지식은 폴더(주제 분류)와 `private` 필드(공개 여부)로 역할을 분리해 관리
- `private: true` 노트는 같은 폴더 안에 섞여 있어도 기본 검색·목록에서 자동 제외
- 민감도가 높은 폴더 전체를 숨기고 싶을 때는 `.gitignore` 병행 사용 가능
- 노트 저장소는 별도 git repo로 관리 권장 (dotfiles 또는 private repo)
