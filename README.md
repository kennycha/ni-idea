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
ni search "캐시" --fuzzy      # 퍼지 검색 (오타 허용)
```

### 조회

```bash
ni get problems/nextjs-caching  # 특정 노트
ni list                         # 전체 목록
ni list --type problem          # 타입별 목록
ni list --tag infra             # 태그별 목록
ni tags                         # 태그 목록
```

### 인덱스 관리

```bash
ni index status    # 인덱스 상태 확인
ni index rebuild   # 인덱스 재구축
```

### 원격 동기화

```bash
# 리모트 설정
ni remote add personal https://my-server.com
ni remote list
ni remote remove personal

# 동기화
ni push --all              # 전체 노트 업로드
ni push problems/my-note   # 특정 노트 업로드
ni pull                    # 리모트에서 다운로드
ni pull --theirs           # 충돌 시 리모트 버전 사용
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

## 서버 (선택)

여러 기기에서 노트를 동기화하거나 팀과 공유하려면 서버를 실행합니다.

### 서버 설치

```bash
# macOS (Apple Silicon)
curl -L https://github.com/kennycha/ni-idea/releases/latest/download/ni-server_darwin_arm64.tar.gz | tar xz

# macOS (Intel)
curl -L https://github.com/kennycha/ni-idea/releases/latest/download/ni-server_darwin_amd64.tar.gz | tar xz

# Linux
curl -L https://github.com/kennycha/ni-idea/releases/latest/download/ni-server_linux_amd64.tar.gz | tar xz

# From Source
go build -o ni-server ./server
```

### 서버 실행

```bash
NI_AUTH_TOKENS=your-secret-token ./ni-server --port 8080 --data ./data
```

### Docker / Kubernetes 배포

```bash
# Docker
docker build -t ni-server -f server/Dockerfile .
docker run -d -p 8080:8080 -e NI_AUTH_TOKENS=your-secret-token -v ni-data:/data ni-server

# microk8s
# 1. 토큰 수정: server/k8s/deployment.yaml의 Secret
# 2. 이미지 빌드
docker build -t ni-server:latest -f server/Dockerfile .
docker save ni-server:latest | microk8s ctr image import -

# 3. 배포
microk8s kubectl apply -f server/k8s/deployment.yaml

# 4. 확인
microk8s kubectl get pods -n ni-idea
```

### 클라이언트 연결

```bash
ni remote add myserver http://localhost:8080
# Token 입력 프롬프트

ni push --all   # 업로드
ni pull         # 다운로드
```

## License

MIT
