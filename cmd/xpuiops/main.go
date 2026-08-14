package main

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"go/format"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	maxJavaScriptSize = 32 << 20
	webPlayerURL      = "https://open.spotify.com/"
	userAgent         = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36"
)

var (
	descriptorRE = regexp.MustCompile(`(?i)["']([_a-z][_0-9a-z]*)["']\s*,\s*["'](query|mutation|subscription)["']\s*,\s*["']([0-9a-f]{64})["']\s*,`)
	bundleRE     = regexp.MustCompile(`https://[a-z0-9.-]+/cdn/build/web-player/[A-Za-z0-9~._-]+\.js`)
)

type extraction struct {
	sources       []string
	clientVersion string
	operations    []operation
}

type operation struct {
	name string
	kind string
	hash string
}

// collector merges operation descriptors from several bundles. Sources are
// applied in order and a later source overrides an earlier one, so the caller
// decides which build wins when the same operation carries different hashes.
type collector struct {
	found     map[string]operation
	conflicts []string
}

func newCollector() *collector {
	return &collector{found: make(map[string]operation)}
}

// scan pulls every descriptor out of a single JavaScript bundle. Conflicting
// descriptors inside one bundle are an error; conflicts across bundles are
// recorded and resolved in favour of the newer source.
func (c *collector) scan(origin string, contents []byte) error {
	seen := make(map[string]operation)
	for _, match := range descriptorRE.FindAllSubmatchIndex(contents, -1) {
		if !hasNullValue(bytes.TrimSpace(contents[match[1]:])) {
			return fmt.Errorf("%s contains an operation descriptor with a non-null value", origin)
		}
		op := operation{
			name: string(contents[match[2]:match[3]]),
			kind: strings.ToLower(string(contents[match[4]:match[5]])),
			hash: strings.ToLower(string(contents[match[6]:match[7]])),
		}
		if current, ok := seen[op.name]; ok && (current.kind != op.kind || current.hash != op.hash) {
			return fmt.Errorf("conflicting descriptors for operation %s in %s", op.name, origin)
		}
		seen[op.name] = op

		if current, ok := c.found[op.name]; ok && current.hash != op.hash {
			c.conflicts = append(c.conflicts, op.name)
		}
		c.found[op.name] = op
	}
	return nil
}

func (c *collector) operations() []operation {
	operations := make([]operation, 0, len(c.found))
	for _, op := range c.found {
		operations = append(operations, op)
	}
	sort.Slice(operations, func(i, j int) bool { return operations[i].name < operations[j].name })
	return operations
}

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "xpuiops:", err)
		os.Exit(1)
	}
}

func run(args []string, stdout io.Writer) error {
	flags := flag.NewFlagSet("xpuiops", flag.ContinueOnError)
	flags.SetOutput(stdout)
	archivePath := flags.String("archive", "", "path to xpui.spa")
	web := flags.Bool("web", false, "also scan the open.spotify.com web player bundles")
	outPath := flags.String("out", "", "output Go file")
	pkgName := flags.String("package", "pfrequest", "package name of the generated file")
	clientVersion := flags.String("client-version", "", "optional Spotify client version")
	check := flags.Bool("check", false, "fail if the generated file differs")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if *archivePath == "" && !*web {
		return errors.New("at least one of -archive or -web is required")
	}
	if *outPath == "" {
		return errors.New("-out is required")
	}
	if filepath.Ext(*outPath) != ".go" {
		return errors.New("-out must be a .go file")
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected arguments: %s", strings.Join(flags.Args(), " "))
	}

	result, err := extract(*archivePath, *clientVersion, *web, stdout)
	if err != nil {
		return err
	}
	generated, err := render(result, *pkgName)
	if err != nil {
		return err
	}

	if *check {
		current, err := os.ReadFile(*outPath)
		if err != nil {
			return err
		}
		if !bytes.Equal(current, generated) {
			return fmt.Errorf("%s is stale; rerun without -check", *outPath)
		}
		fmt.Fprintf(stdout, "%s is current (%d operations)\n", *outPath, len(result.operations))
		return nil
	}

	if err := os.WriteFile(*outPath, generated, 0o644); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "generated %d operations in %s\n", len(result.operations), *outPath)
	return nil
}

// extract merges both bundle sources. The web player is scanned first so that a
// locally installed desktop archive, whose descriptors match the App-Platform
// this library reports, overrides it on disagreement.
func extract(archivePath, clientVersion string, web bool, stdout io.Writer) (*extraction, error) {
	c := newCollector()
	var sources []string

	if web {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		source, err := extractWeb(ctx, c)
		if err != nil {
			return nil, err
		}
		sources = append(sources, source)
	}
	if archivePath != "" {
		source, err := extractArchive(archivePath, c)
		if err != nil {
			return nil, err
		}
		sources = append(sources, source)
	}

	if len(c.found) == 0 {
		return nil, errors.New("no persisted-query operation descriptors found")
	}
	if len(c.conflicts) > 0 {
		sort.Strings(c.conflicts)
		fmt.Fprintf(stdout, "note: later source overrode %d operation(s): %s\n",
			len(c.conflicts), strings.Join(slices.Compact(c.conflicts), ", "))
	}

	return &extraction{
		sources:       sources,
		clientVersion: clientVersion,
		operations:    c.operations(),
	}, nil
}

// extractWeb scans the open.spotify.com bundles, which carry the operations the
// desktop archive lazily fetches at runtime and therefore never ships.
func extractWeb(ctx context.Context, c *collector) (string, error) {
	page, err := fetch(ctx, webPlayerURL)
	if err != nil {
		return "", err
	}
	bundles := bundleRE.FindAllString(string(page), -1)
	if len(bundles) == 0 {
		return "", fmt.Errorf("no web player bundles referenced by %s", webPlayerURL)
	}

	sort.Strings(bundles)
	bundles = slices.Compact(bundles)
	for _, bundle := range bundles {
		contents, err := fetch(ctx, bundle)
		if err != nil {
			return "", err
		}
		if err := c.scan(bundle, contents); err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("%s (%s)", webPlayerURL, strings.Join(bundles, ", ")), nil
}

func fetch(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: %s", url, resp.Status)
	}

	contents, err := io.ReadAll(io.LimitReader(resp.Body, maxJavaScriptSize+1))
	if err != nil {
		return nil, err
	}
	if len(contents) > maxJavaScriptSize {
		return nil, fmt.Errorf("%s exceeds the %d-byte limit", url, maxJavaScriptSize)
	}
	return contents, nil
}

func extractArchive(archivePath string, c *collector) (string, error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return "", err
	}
	digest := sha256.New()
	if _, err := io.Copy(digest, file); err != nil {
		return "", err
	}
	archive, err := zip.NewReader(file, info.Size())
	if err != nil {
		return "", err
	}

	for _, asset := range archive.File {
		if !strings.HasSuffix(asset.Name, ".js") && !strings.HasSuffix(asset.Name, ".mjs") {
			continue
		}
		if asset.UncompressedSize64 > maxJavaScriptSize {
			return "", fmt.Errorf("%s exceeds the %d-byte JavaScript limit", asset.Name, maxJavaScriptSize)
		}

		contents, err := readZipFile(asset)
		if err != nil {
			return "", fmt.Errorf("read %s: %w", asset.Name, err)
		}
		if err := c.scan(asset.Name, contents); err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("%s SHA-256 %x", filepath.Base(archivePath), digest.Sum(nil)), nil
}

func render(result *extraction, pkgName string) ([]byte, error) {
	var source strings.Builder
	source.WriteString("// Code generated by xpuiops; DO NOT EDIT.\n")
	if result.clientVersion != "" {
		fmt.Fprintf(&source, "// Spotify client version: %s\n", strconv.Quote(result.clientVersion))
	}
	for _, src := range result.sources {
		fmt.Fprintf(&source, "// Source: %s\n", src)
	}
	source.WriteString("\n")
	fmt.Fprintf(&source, "package %s\n\n", pkgName)
	source.WriteString("func init() {\n")
	source.WriteString("generatedOperationHashes = map[Operation]string{\n")
	for _, op := range result.operations {
		fmt.Fprintf(&source, "\t%s: %s, // %s\n", strconv.Quote(op.name), strconv.Quote(op.hash), op.kind)
	}
	source.WriteString("}\n")
	source.WriteString("}\n")
	return format.Source([]byte(source.String()))
}

func readZipFile(file *zip.File) ([]byte, error) {
	reader, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	contents, err := io.ReadAll(io.LimitReader(reader, maxJavaScriptSize+1))
	if err != nil {
		return nil, err
	}
	if len(contents) > maxJavaScriptSize {
		return nil, errors.New("uncompressed data exceeds size limit")
	}
	return contents, nil
}

func hasNullValue(data []byte) bool {
	if !bytes.HasPrefix(data, []byte("null")) {
		return false
	}
	return len(data) == 4 || !((data[4] >= 'a' && data[4] <= 'z') || (data[4] >= 'A' && data[4] <= 'Z') || (data[4] >= '0' && data[4] <= '9') || data[4] == '_' || data[4] == '$')
}
