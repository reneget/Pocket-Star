package monitor

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

const (
	defaultMaxLogBytes int64 = 10 * 1024 * 1024 // 10 MB
	ringBufferSize           = 1000
)

type MetricsStore struct {
	mu       sync.RWMutex
	buf      [ringBufferSize]CheckResult
	head     int
	size     int

	filePath string
	maxBytes int64
	file     *os.File
}

func NewMetricsStore(dataDir string) (*MetricsStore, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("create monitor dir: %w", err)
	}

	path := filepath.Join(dataDir, "metrics.log")
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open metrics log: %w", err)
	}

	return &MetricsStore{
		filePath: path,
		maxBytes: defaultMaxLogBytes,
		file:     f,
	}, nil
}

func (s *MetricsStore) SetMaxBytes(n int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.maxBytes = n
}

func (s *MetricsStore) Push(r CheckResult) {
	s.mu.Lock()
	s.buf[s.head] = r
	s.head = (s.head + 1) % ringBufferSize
	if s.size < ringBufferSize {
		s.size++
	}
	s.mu.Unlock()

	s.appendToFile(r)
}

func (s *MetricsStore) Recent(n int) []CheckResult {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.size == 0 {
		return nil
	}
	if n > s.size {
		n = s.size
	}

	res := make([]CheckResult, n)
	start := s.head - n
	if start < 0 {
		start += ringBufferSize
	}
	for i := 0; i < n; i++ {
		res[i] = s.buf[(start+i)%ringBufferSize]
	}
	return res
}

func (s *MetricsStore) FilterRecent(n int, ct CheckType) []CheckResult {
	all := s.Recent(s.size)
	var filtered []CheckResult
	for _, r := range all {
		if r.Type == ct {
			filtered = append(filtered, r)
		}
		if len(filtered) >= n {
			break
		}
	}
	return filtered
}

func (s *MetricsStore) Last(ct CheckType) *CheckResult {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i := 0; i < s.size; i++ {
		idx := (s.head - 1 - i + ringBufferSize) % ringBufferSize
		if s.buf[idx].Type == ct {
			r := s.buf[idx]
			return &r
		}
	}
	return nil
}

func (s *MetricsStore) LastByLabel(ct CheckType, label string) *CheckResult {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i := 0; i < s.size; i++ {
		idx := (s.head - 1 - i + ringBufferSize) % ringBufferSize
		if s.buf[idx].Type == ct && s.buf[idx].Label == label {
			r := s.buf[idx]
			return &r
		}
	}
	return nil
}

func (s *MetricsStore) appendToFile(r CheckResult) {
	data, err := json.Marshal(r)
	if err != nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.file != nil {
		s.file.Write(data)
		s.file.Write([]byte("\n"))
		s.file.Sync()
		s.checkRotation()
	}
}

func (s *MetricsStore) checkRotation() {
	info, err := os.Stat(s.filePath)
	if err != nil {
		return
	}
	if info.Size() >= s.maxBytes {
		s.file.Close()
		s.file = nil
		backup := s.filePath + ".1"
		os.Remove(backup)
		os.Rename(s.filePath, backup)
		f, err := os.OpenFile(s.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			s.file = f
		}
	}
}

func (s *MetricsStore) ReadFileLines(n int) ([]CheckResult, error) {
	f, err := os.Open(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var results []CheckResult
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var r CheckResult
		if err := json.Unmarshal(scanner.Bytes(), &r); err != nil {
			continue
		}
		results = append(results, r)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(results) > n {
		results = results[len(results)-n:]
	}
	return results, nil
}

func (s *MetricsStore) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.file != nil {
		return s.file.Close()
	}
	return nil
}
