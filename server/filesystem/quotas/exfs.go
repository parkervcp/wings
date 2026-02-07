package quotas

import (
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/parkervcp/fsquota"
	"github.com/pelican-dev/wings/config"
)

var exfsProjects []exfsProject

type exfsProject struct {
	ID       int
	Name     string
	BasePath string
}

const (
	projidTemplate = `{{ range . }}{{ .UUID }}:{{ .ID }}
{{ end }}`
	projectsTemplate = `{{ range . }}{{ .ID }}:{{ .BasePath }}/{{ .UUID }}
{{ end }}`

	projidFile  = `/etc/projid`
	projectFile = `/etc/projects`
)

// setQuota sets the quota in bytes for the specified server uuid
func (q exfsProject) setQuota(byteLimit uint64) (err error) {
	serverProject, err := fsquota.LookupProject(q.Name)
	if err != nil {
		return
	}

	serverDirPath := fmt.Sprintf("%s/%s", config.Get().System.Data, q.Name)
	projInfo, err := fsquota.GetProjectInfo(serverDirPath, serverProject)
	if err != nil {
		return
	}

	projInfo.Limits.Bytes.SetHard(byteLimit)

	if _, err = fsquota.SetProjectQuota(serverDirPath, serverProject, projInfo.Limits); err != nil {
		return
	}
	return
}

// getQuota gets the specified quotas and usage of a specified server uuid
func (q exfsProject) getQuota() (bytesUsed int64, err error) {
	serverProject, err := fsquota.LookupProject(q.Name)
	if err != nil {
		return
	}

	projInfo, err := fsquota.GetProjectInfo(q.BasePath, serverProject)
	if err != nil {
		return
	}

	// converts the uint64 to int64.
	// This should only be an issue in the terms of exabytes...
	return int64(projInfo.BytesUsed), nil
}

// enableEXFSQuota enables quotas on a specified directory
func (q exfsProject) enableEXFSQuota() (err error) {
	serverDir, err := os.OpenFile(fmt.Sprintf("%s/%s", q.BasePath, q.Name), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return
	}

	defer serverDir.Close()

	folderXattr, err := getXAttr(serverDir)
	if err != nil {
		return
	}

	// ensure project inherit flag is set
	if (folderXattr.XFlags & FS_XFLAG_PROJINHERIT) != 0 {
		if err = setXAttr(serverDir, fsXAttr{XFlags: FS_XFLAG_PROJINHERIT}); err != nil {
			return
		}
	}

	// ensure correct project id is set
	if folderXattr.ProjectID != uint32(q.ID) {
		if err = setXAttr(serverDir, fsXAttr{ProjectID: uint32(q.ID)}); err != nil {
			return
		}
	}

	return
}

func (q exfsProject) addProject() (err error) {
	basePath := config.Get().System.Data
	if strings.HasSuffix(basePath, "/") {
		basePath = strings.TrimSuffix(config.Get().System.Data, "/")
	}

	q.BasePath = basePath
	exfsProjects = append(exfsProjects, q)

	if err = writeEXFSProjects(); err != nil {
		return
	}

	if err = q.enableEXFSQuota(); err != nil {
		return
	}

	return
}

// removeProject drops a specified project from the
func (q exfsProject) removeProject() (err error) {
	for pos, project := range exfsProjects {
		if project.Name == q.Name {
			exfsProjects = append(exfsProjects[:pos], exfsProjects[pos+1:]...)
		}
	}

	err = writeEXFSProjects()
	return
}

func writeEXFSProjects() (err error) {
	// write out projid file
	idtmpl, err := template.New("projid").Parse(projidTemplate)
	if err != nil {
		return
	}

	if err = writeTemplate(idtmpl, projidFile, exfsProjects); err != nil {
		return
	}

	projtmpl, err := template.New("projects").Parse(projectsTemplate)
	if err != nil {
		return
	}

	if err = writeTemplate(projtmpl, projectFile, exfsProjects); err != nil {
		return
	}

	return
}

func writeTemplate(t *template.Template, file string, data interface{}) (err error) {
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	err = t.Execute(f, data)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}
