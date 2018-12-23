//when load time called
window.onload = function() {
    var full_location = window.location.href;
    var call_location = full_location.substring(full_location.lastIndexOf("/"));

    $(".active").removeClass("active");
    if(call_location == "/board.do") {
        $("#board").addClass("active");

        $("#board_pagenation").val("");
        var cnt = $("#board_cnt").val();
        if(cnt > 1) {//show the pagenation elemenet
            var page = Number($("#board_page").val());
            if( page == "" ) page = 1;
            var i = Math.floor((page-1)/10)+1; //start number
            var end = i+9;
            if( end > cnt) end = cnt; //end number

            var prev_page = page-1;
            if( prev_page == 0 ) prev_page = 1;
            var next_page = page+1;
            if( next_page > cnt ) next_page = cnt;

            //add component
            $("#board_pagenation").append("<a onclick='call_board(1);'>&laquo;</a>");
            $("#board_pagenation").append("<a onclick='call_board("+prev_page+");'>&lsaquo;</a>");
            for( ; i<=end; i++){
                $("#board_pagenation").append("<a onclick='call_board("+i+");'>"+i+"</a>");
            }
            $("#board_pagenation").append("<a onclick='call_board("+next_page+");'>&rsaquo;</a>");
            $("#board_pagenation").append("<a onclick='call_board("+cnt+");'>&raquo;</a>");
            $("#board_pagenation a:eq("+(page+1)+")").addClass("active");
        }
    }
    else if( call_location.substr(0,6)=="/board" ) {
        $("#board").addClass("active");
    }
    else {
        $("#main").addClass("active");
    }
};

// login view
function login_view(){
    if($("#login").css('display') == 'none') {
        $("#login").show();
    }
    else {
        $("#login").hide();
    }
}

// login check
function login_check(){
    var password = $("#login input:eq(2)").val();
    if( $("#login input:eq(0)").val()!="" && password !="" ) {
        var encryptVal = SHA512(password)
        $("#login input:eq(1)").val(encryptVal);
        $.ajax({
            type: "post",
            url: "/login_check.do",
            data: $("#login_form").serialize(),
            contentType: "application/x-www-form-urlencoded",
            success: function(responseData, textStatus, jqXHR) {
                if(responseData == "0") {
                    location.reload();
                    //location.href="/";
                } else {
                    alert(responseData);
                }
            },
            error: function(jqXHR, textStatus, errorThrown) {
                console.log(errorThrown);
            }
        })
    }
    else {
        alert("Two passwords are different.");
    }
}
// logout
function logout(){
    $.ajax({
        type: "post",
        url: "/logout.do",
        data: "",
        contentType: "application/x-www-form-urlencoded",
        success: function(responseData, textStatus, jqXHR) {
            if(responseData == "0") {
                location.reload();
                //location.href="/";
            } else {
                alert(responseData);
            }
        },
        error: function(jqXHR, textStatus, errorThrown) {
            console.log(errorThrown);
        }
    })
}

// registration view
function registration_view(){
    if($("#registration").css('display') == 'none') {
        $("#registration").show();
    }
    else {
        $("#registration").hide();
    }
}

// registration check
function registration_check(){
    var is_error = false;
    var email =  $("#registration input:eq(0)").val();
    var password =  $("#registration input:eq(2)").val();
    var password_repeat = $("#registration input:eq(3)").val();
    var name =  $("#registration input:eq(4)").val();
    var telephone =  $("#registration input:eq(5)").val();

    if( email == "" ){
        is_error = true;
        alert("Please enter the Email");
    }
    else if( password == "" || password_repeat == ""){
        is_error = true;
        alert("Please enter the Password");
    }
    else if( name == "" ){
        is_error = true;
        alert("Please enter the Name");
    }
    else if( telephone != "" && isNaN(Number(telephone)) ){
        is_error = true;
        alert("Please enter the Telephone number(number only)");
    }

    //no error then execute.
    if(is_error == false) {
        if( password == password_repeat ) {
            var encryptVal = SHA512(password);
            $("#registration input:eq(1)").val(encryptVal);
            $.ajax({
                type: "post",
                url: "/registration_check.do",
                data: $("#registration_form").serialize(),
                contentType: "application/x-www-form-urlencoded",
                success: function(responseData, textStatus, jqXHR) {
                    if(responseData == "0") {
                        alert("Successfully restrated!");
                        registration_view();
                        login_view();
                    } else {
                        alert(responseData);
                    }
                },
                error: function(jqXHR, textStatus, errorThrown) {
                    console.log(errorThrown);
                }
            })
        }
        else {
            alert("Two passwords are different.");
        }
    }
}

