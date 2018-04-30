<!DOCTYPE html>
<html>
<head>
<meta name="viewport" content="width=device-width, initial-scale=1">
<link rel="stylesheet" type="text/css" href="/assets/css/stylesheet.css">
<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css" integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous">

</head>
<body>

        <div class="container">

                <h1>Mark transport jobs as complete:</h1>
        
                <br />
        
                <form action="/transport/create" method="POST"> 
                <table class="table table-sm table-hover">
                    <thead>
                        <tr>
                            <th scope="col">Job No</th>
                            <th scope="col">Action</th>
                            <th scope="col">Company</th>
                            <th scope="col">Department</th>
                            <th scope="col">Job Description</th>
                            <th scope="col">Date</th>
                            <th scope="col">Complete?</th>
                            <th scope="col">Edit</th>
                        </tr>
                    </thead>
                    <!-- Start of Block -->
                    {{range .ReportsAll}}
                    <tbody>
                        <tr>
                            <td>{{.JobNo}}</td>
                            <td>{{.ActionFlag}}</td>
                            <td>{{.Customer}}</td>
                            <td>{{.Department}}</td>
                            <td>{{.JobDescription}}</td>
                            <td>{{.TransportDate}}</td>
                            <td><a href="/transport/complete/{{.TransportID}}">Click to complete</a></td>
                            <td><a href="/transport/edit/{{.TransportID}}">Click to edit</a></td>
                        </tr>
                    </tbody>
                    {{end}}
                    <!-- End of Block -->
                </table>
        
                <br />
                <a href="/transport/create">Create transport job</a>
                <br />
            </div>
</body>
</html>