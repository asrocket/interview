package sort

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
)

const (
	MaxBuffSize = 2 * 1024 * 1024
)

type byLine []string

func (a byLine) Len() int           { return len(a) }
func (a byLine) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byLine) Less(i, j int) bool { return a[i] < a[j] }

type Sorter interface {
	Sort() error
}

type lineSorter struct {
	r        io.Reader
	w        io.Writer
	buffSize int
}

func (s *lineSorter) Sort() error {
	defer os.RemoveAll("tmp")

	partNum, err := s.splitFile()
	if err != nil {
		return err
	}
	if partNum == 1 {
		return s.writeResult(0)
	}
	return s.merge(0, 1, partNum-1, partNum)
}

func (s *lineSorter) splitFile() (int, error) {
	if err := os.MkdirAll("tmp", os.ModePerm); err != nil {
		return 0, err
	}

	var (
		lineLength int
		partNum    int
		buf        = make([]string, 0, 100)
		scanner    = bufio.NewScanner(s.r)
	)

	for scanner.Scan() {
		bytes := scanner.Bytes()
		if len(bytes)+lineLength > s.buffSize {
			sort.Sort(byLine(buf))

			file, err := os.Create(fmt.Sprintf("tmp/%d", partNum))
			if err != nil {
				return 0, err
			}

			for _, str := range buf {
				if _, err := file.WriteString(str); err != nil {
					return 0, err
				}
				if _, err := file.Write([]byte{'\n'}); err != nil {
					return 0, err
				}
			}

			buf = buf[:0]
			lineLength = 0
			_ = file.Close()
			partNum++
		}
		buf = append(buf, string(bytes))
		lineLength += len(bytes)
	}
	sort.Sort(byLine(buf))
	file, err := os.Create(fmt.Sprintf("tmp/%d", partNum))
	if err != nil {
		return 0, err
	}
	for _, str := range buf {
		if _, err := file.WriteString(str); err != nil {
			return 0, err
		}
		if _, err := file.Write([]byte{'\n'}); err != nil {
			return 0, err
		}
	}
	_ = file.Close()
	return partNum + 1, nil
}

func (s *lineSorter) merge(l, r, parts, partNum int) error {
	if r == partNum {
		return s.writeResult(l)
	}

	file, err := os.Create(fmt.Sprintf("tmp/%d", partNum))
	if err != nil {
		return err
	}

	lf, _ := os.Open(fmt.Sprintf("tmp/%d", l))
	rf, _ := os.Open(fmt.Sprintf("tmp/%d", r))

	lscanner := bufio.NewScanner(lf)
	rscanner := bufio.NewScanner(rf)

	scanl, scanr := lscanner.Scan(), rscanner.Scan()

	for {
		if scanl && scanr {
			switch {
			case lscanner.Text() < rscanner.Text():
				file.WriteString(lscanner.Text())
				file.Write([]byte{'\n'})
				scanl = lscanner.Scan()
				continue

			case lscanner.Text() > rscanner.Text():
				file.WriteString(rscanner.Text())
				file.Write([]byte{'\n'})
				scanr = rscanner.Scan()
				continue
			default:
				file.WriteString(lscanner.Text())
				file.Write([]byte{'\n'})
				file.WriteString(rscanner.Text())
				file.Write([]byte{'\n'})
				scanl, scanr = lscanner.Scan(), rscanner.Scan()
				continue
			}
		}

		if scanl {
			file.WriteString(lscanner.Text())
			file.Write([]byte{'\n'})
			scanl = lscanner.Scan()
			continue
		}

		if scanr {
			file.WriteString(rscanner.Text())
			file.Write([]byte{'\n'})
			scanr = rscanner.Scan()
			continue
		}
		break
	}

	file.Close()
	lf.Close()
	rf.Close()
	return s.merge(r+1, r+2, parts-1, partNum+1)
}

func (s *lineSorter) writeResult(part int) error {
	name := fmt.Sprintf("tmp/%d", part)
	lf, err := os.Open(name)
	if err != nil {
		return err
	}
	defer lf.Close()

	lscanner := bufio.NewScanner(lf)
	for lscanner.Scan() {
		_, _ = s.w.Write(lscanner.Bytes())
		_, _ = s.w.Write([]byte{'\n'})
	}

	return nil
}

func NewLineSorter(in io.Reader, out io.Writer) Sorter {
	return &lineSorter{
		r:        in,
		w:        out,
		buffSize: MaxBuffSize,
	}
}
