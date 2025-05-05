{{ define "content" }}

    <h1>Field list</h1>
    <a href="/fields/create">Create new field</a>
    <table>
        <thead>
            <tr>
                <th>Name</th>
                <th>Alias</th>
                <th>Collection</th>
                <th>Actions</th>
            </tr>
        </thead>
        <tbody>
            {{ range .fields }}
            <tr>
                <td>{{ .Name }}</td>
                <td>{{ .Alias }}</td>
                <td>{{ .Collection.Name }}</td>
                <td><a href="/fields/edit/{{ .ID }}">Edit</a></td>
            </tr>
            {{ end }}
        </tbody>
    </table>

{{end}}