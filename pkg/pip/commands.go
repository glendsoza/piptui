package pip

import (
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

type PackageCache struct {
	data map[string]*Package
	lock *sync.RWMutex
}

func (pc *PackageCache) Add(pkg *Package) {
	pc.lock.Lock()
	defer pc.lock.Unlock()
	pc.data[pkg.Name] = pkg
}

func (pc *PackageCache) Get(name string) (*Package, bool) {
	pc.lock.RLock()
	defer pc.lock.RUnlock()
	v, ok := pc.data[name]
	return v, ok
}

type installStatus string

const (
	IN_PROGRESS installStatus = "In Progress"
	READY       installStatus = "Ready"
	INTERRUPTED installStatus = "Interrupted"
)

type PipHelper struct {
	cache         *PackageCache
	installStatus installStatus
	installCmd    *exec.Cmd
}

func NewPipHelper() *PipHelper {
	return &PipHelper{
		cache: &PackageCache{
			data: make(map[string]*Package),
			lock: &sync.RWMutex{},
		},
		installStatus: READY,
		installCmd:    nil,
	}
}

var matchSpaceRegex = regexp.MustCompile(`\s+`)

func (p *PipHelper) InstallInProgress() bool {
	return p.installCmd != nil
}

func (p *PipHelper) show(packages []string) ([]byte, error) {
	output, err := exec.Command("/bin/bash", "-c", fmt.Sprintf("pip show %s", strings.Join(packages, " "))).CombinedOutput()
	if err != nil {
		return nil, err
	}
	return output, nil
}

// starts the uninstall command and does not wait for it.
// output is written to the writer
func (p *PipHelper) Uninstall(name string) []byte {
	out, _ := exec.Command("/bin/bash", "-c", fmt.Sprintf("pip uninstall -y %s 2>&1", name)).Output()
	return out
}

// starts the install command and does not wait for it.
// output is written to the writer
func (p *PipHelper) InstallAndStream(nameAndVersion string) (io.Reader, error) {
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("pip install %s 2>&1", nameAndVersion))
	s, err := cmd.StdoutPipe()
	cmd.Start()
	p.installCmd = cmd
	go func() {
		cmd.Wait()
		p.installCmd = nil
	}()
	return s, err
}

func (p *PipHelper) CancelInstall() {
	if p.InstallInProgress() {
		p.installCmd.Process.Kill()
	}
}

func (p *PipHelper) Show(name string, version string) (*Package, error) {
	v, ok := p.cache.Get(name)
	if ok {
		return v, nil
	}
	raw, err := p.show([]string{name})
	if err != nil {
		return nil, err
	}
	pkg, err := ParseRawData(string(raw))
	if err != nil {
		return nil, err
	}
	p.cache.Add(pkg)
	return pkg, nil
}

func (p *PipHelper) List() (Packages, error) {
	output, err := exec.Command("/bin/bash", "-c", "pip list 2>&1").CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf(string(output))
	}
	splitData := strings.Split(matchSpaceRegex.ReplaceAllString(strings.TrimSpace(string(output)), " "), " ")[4:]
	var packages Packages
	var packagesNamesList []string
	for x := 0; x < len(splitData); x += 2 {
		s := splitData[x : x+2]
		packagesNamesList = append(packagesNamesList, s[0])
		if err != nil {
			return nil, err
		}
	}
	output, err = p.show(packagesNamesList)
	if err != nil {
		return nil, err
	}
	for _, data := range strings.Split(string(output), "---") {
		pkg, err := ParseRawData(data)
		if err != nil {
			return nil, err
		}
		packages = append(packages, pkg)
		p.cache.Add(pkg)
	}
	return packages, nil
}
