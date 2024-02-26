## RA 작업 트레이 

특정 커맨드를 실행시키기 위한 용도

### 빌드 시 아래 명령어 실행

```bash
go build -ldflags -H=windowsgui
go build -ldflags -H=windowsgui -o RaTray.exe
```

### 참고사항 

- ra.exe 파일와 menu.json 파일은 동일 디렉토리에 존재해야 함
- 메뉴 수정 시 ra.exe 종료하고 menu.json 파일을 수정한 후 ra.exe 재 실행