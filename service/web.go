package service

import (
	"net/http"
	"text/template"

	"github.com/Sirupsen/logrus"

	"github.com/slok/khronos/job"
)

var (
	jobsListTemplate = template.Must(template.New("dash").Parse(dashHTML))
)

//JobsList is the context that JobsList endpoint will use to render the template
type JobsList struct {
	Jobs []*job.Job
}

// JobsList Processes a jobs list template
func (s *KhronosService) JobsList(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("Calling jobsList HTML endpoint")
	js, _ := s.Storage.GetJobs(0, 0)
	jobs := JobsList{
		Jobs: js,
	}
	jobsListTemplate.Execute(w, &jobs)
}

const dashHTML = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <title>Job lists</title>
  </head>
  <body>
	<h1>Jobs list</h1>
    <table>
        <tr>
            <th>ID</th>
            <th>Name</th>
            <th>Description</th>
            <th>When</th>
            <th>Active</th>
            <th>URL</th>
        </tr>
        {{ range .Jobs }}
        <tr>
            <td>{{ .ID }}</td>
            <td>{{ .Name }}</td>
            <td>{{ .Description }}</td>
            <td>{{ .When }}</td>
            <td>{{ .Active }}</td>
            <td>{{ .URL }}</td>
        </tr>
        {{ end }}
    </table>
  </body>
</html>`
