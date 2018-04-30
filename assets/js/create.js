$(function() {
    console.log("Page loaded")

    $("#get-job").on("click", function(){
        console.log($("#spoon").val())
        $.ajax({
            url: "/transport/create/" + $("#spoon").val(), 
            method: "GET",
            dataType: "json",
            success: function(j){
                console.log(j)
                $("#dept").val(j.department)
                $("#customer-name").val(j.customerName)
                $("#description").val(j.jobDescription)
            },
            error: function(error){
                console.log(error)
            }
        })
    })

})