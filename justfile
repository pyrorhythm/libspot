set dotenv-load := true

test:
	go test ./...

gover:
	goveralls -repotoken ${GOVERALLS_TOKEN}

updsum SEMVER:
	sleep 3
	curl https://sum.golang.org/lookup/pyrorhythm.dev/libspot@{{SEMVER}}

[parallel]
upload-coverage-and-fetch SEMVER: gover (updsum SEMVER)

tag-push SEMVER:
	git tag {{SEMVER}}
	git push origin {{SEMVER}}

commit-push SEMVER:
    git add . ; git commit -m "release: {{SEMVER}}"
    git tag {{SEMVER}}
    git push ; git push origin {{SEMVER}}

release SEMVER: test (commit-push SEMVER) (upload-coverage-and-fetch SEMVER)
