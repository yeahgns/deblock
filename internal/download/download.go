package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func File(url, dest string) error {
	client := &http.Client{Timeout: 10 * time.Minute}
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP response %d when downloading %s", resp.StatusCode, url)
	}

	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("couldn't create the file %s: %w", dest, err)
	}
	defer out.Close()

	counter := &progressWriter{total: resp.ContentLength}
	if _, err := io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	counter.finish()
	return nil
}

type progressWriter struct {
	total      int64
	downloaded int64
	lastPrint  time.Time
}

func (p *progressWriter) Write(b []byte) (int, error) {
	n := len(b)
	p.downloaded += int64(n)
	if time.Since(p.lastPrint) > 150*time.Millisecond {
		p.lastPrint = time.Now()
		p.print()
	}
	return n, nil
}

func (p *progressWriter) finish() {
	p.print()
	fmt.Println()
}

func (p *progressWriter) print() {
	if p.total > 0 {
		pct := float64(p.downloaded) / float64(p.total) * 100
		fmt.Printf("\r  downloading... %5.1f%% (%s / %s)   ", pct, humanBytes(p.downloaded), humanBytes(p.total))
	} else {
		fmt.Printf("\r  downloading... %s", humanBytes(p.downloaded))
	}
}

func humanBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}
