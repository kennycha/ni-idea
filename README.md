# ni-idea

개인 지식 베이스 CLI

문제 해결 기록, 기술 결정, 학습 내용을 로컬 마크다운으로 저장하고 빠르게 검색합니다.

## 설치

### macOS (Apple Silicon)

```bash
curl -L https://github.com/kennycha/ni-idea/releases/latest/download/ni_darwin_arm64.tar.gz | tar xz
sudo mv ni /usr/local/bin/
```

### macOS (Intel)

```bash
curl -L https://github.com/kennycha/ni-idea/releases/latest/download/ni_darwin_amd64.tar.gz | tar xz
sudo mv ni /usr/local/bin/
```

### Linux

```bash
curl -L https://github.com/kennycha/ni-idea/releases/latest/download/ni_linux_amd64.tar.gz | tar xz
sudo mv ni /usr/local/bin/
```

### From Source

```bash
go install github.com/kennycha/ni-idea/cmd/ni@latest
```

## 초기 설정

```bash
ni init
```

생성되는 구조:

```
~/.ni-idea/
├── config.yaml
├── notes/
│   ├── problems/
│   ├── decisions/
│   ├── knowledge/
│   └── practice/
└── templates/
```

Claude Code skill도 함께 설치할 수 있습니다 (선택).

## 사용법

### 노트 추가

```bash
ni add --type problem --title "Next.js 캐싱 이슈"
ni add --type decision --title "배포 전략" --tag infra,k8s
ni add --type knowledge --title "Go 동시성"
ni add --type practice --title "Kubernetes 실습"
```

### 검색

```bash
ni search "캐싱"              # problem, decision 검색
ni search "캐싱" --all        # 전체 타입 검색
ni search "캐싱" --tag nextjs # 태그 필터
```

### 조회

```bash
ni get problems/nextjs-caching  # 특정 노트
ni list                         # 전체 목록
ni list --type problem          # 타입별 목록
ni list --tag infra             # 태그별 목록
ni tags                         # 태그 목록
```

## 노트 타입

| 타입        | 용도                |
| ----------- | ------------------- |
| `problem`   | 에러/문제 해결 기록 |
| `decision`  | 기술/아키텍처 결정  |
| `knowledge` | 개념/기술 문서      |
| `practice`  | 실습/학습 로그      |

## Claude Code 연동

`ni init` 시 skill 설치를 선택하면 다른 프로젝트에서 Claude Code가 지식 베이스를 활용할 수 있습니다.

- `/ni-idea-search` - 과거 경험 검색
- `/ni-idea-add` - 새 지식 기록

## License

MIT
