<!DOCTYPE html>
<html>
        <head>
<meta name="viewport" content="width=device-width, initial-scale=1">
<link rel="stylesheet" type="text/css" href="/assets/css/stylesheet.css">
<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css" integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous">
        <title>RJW Transportation System</title>
    </head>
    <body>
        
        <div class="container">
        <h1>Transportation management tool</h1>
        <h2>Enter the motor for transportation here</h2>
        <br />    
        <td><p id="error">{{.DisplayMessage}}</p></td>

    <div class="row">
        <div class="col-md-6 col-md-offset-3">

        <form action="/transport/create" method="POST"> 
            <div class="form-group">
                <label for ="action_flag">Action</label>
                <select class="form-control" name="action_flag">
                    <option value="C"{{if eq .ActionFlag "C"}} selected{{end}}>Collect</option>
                    <option value="D"{{if eq .ActionFlag "D"}} selected{{end}}>Deliver</option>
                </select>
            </div>
            <div class="form-group">
                <label for ="service_type">Service type</label>
                <select class="form-control" name="service_type">
                    <option value="R"{{if eq .ActionFlag "R"}} selected{{end}}>Regular</option>
                    <option value="E"{{if eq .ActionFlag "E"}} selected{{end}}>Emergency Service</option>
                </select>
            </div>
            <div class="form-group">
            <label for="job_no">Job Number</label>
            <input type="text" class="form-control" placeholder="E.g. 100252" name="job_no" id="spoon" value="{{.JobNo}}"> <button type="button" class="btn btn-default" value="Get job number" id="get-job">Find job</button>
            </div>
            <div class="form-group">
                <label for="department">Department</label>
                <select class="form-control" name="department" id="dept">
                    <option value="EL"{{if eq .Department "EL"}} selected{{end}}>Electrical</option>
                    <option value="ME"{{if eq .Department "ME"}} selected{{end}}>Mechanical</option>
                    <option value="PM"{{if eq .Department "PM"}} selected{{end}}>Electronics</option>
                    <option value="HI"{{if eq .Department "HI"}} selected{{end}}>Hi-Cycle</option>
                    <option value="BA"{{if eq .Department "BA"}} selected{{end}}>Balancing</option>
                    </select>
            </div>
            <div class="form-group">
                <label for="customer">Customer</label>
                <input type="text" class="form-control" name="customer" id="customer-name" value="{{.Customer}}">
            </div>
            <div class="form-group">
                <label for="job_description">Job Description</label>
                <input type="text" class="form-control" name="job_description" id="description" value="{{.JobDescription}}">
            </div>    
            <div class="form-group">
                <label for="date">Date required</label>
                <input type="date" class="form-control" name="date" value="{{.TransportDate}}">
            </div>
            <div class="form-group">    
                <button type="submit" class="btn btn-default" name="submit">Submit</button>
            </div>
  </form>

    
        <a href="/transport/complete">Complete or edit transport jobs</a> <br />
        <br />

            </div> <!-- end column -->
        </div> <!-- end row -->
    </div> <!-- end container -->

        <script src="/assets/js/jquery-3.3.1.min.js"></script>
        <script src="/assets/js/create.js"></script>
        
    </body>
</html>