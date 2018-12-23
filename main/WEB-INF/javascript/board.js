//move page
function call_board_add(){
    location.href="/board_add.do";
}

function call_board_detail(seq){
    $("#board_detail_seq").val(seq);
    $("#board_detail_form").submit();
}

function call_board(page){
    $("#board_page").val(page);
    $("#board_form").submit();
}

function call_board_search(){
    var keyword = $("#board_search").val();
    $("#keyword").val(keyword);
    call_board(1);
}

function board_add_check(){
    $.ajax({
        type: "post",
        url: "/board_add_check.do",
        data: $("#board_add_form").serialize(),
        contentType: "application/x-www-form-urlencoded",
        success: function(responseData, textStatus, jqXHR) {
            if(responseData == "0") {
                location.href="/board.do";
            } else {
                alert(responseData);
            }
        },
        error: function(jqXHR, textStatus, errorThrown) {
            console.log(errorThrown);
        }
    });
}

function call_board_modify(){
    $("#board_modify_form").submit();
}

function call_board_modify_check(){
    $.ajax({
        type: "post",
        url: "/board_modify_check.do",
        data: $("#board_modify_check_form").serialize(),
        contentType: "application/x-www-form-urlencoded",
        success: function(responseData, textStatus, jqXHR) {
            if(responseData == "0") {
                alert("Succesfully modified");
                location.href="/board.do";
                //location.reload();
            } else {
                alert(responseData);
            }
        },
        error: function(jqXHR, textStatus, errorThrown) {
            console.log(errorThrown);
        }
    });
}

function board_delete(seq, author){
    $.ajax({
        type: "post",
        url: "/board_delete.do",
        data: JSON.parse('{"seq":'+seq+', "author":"'+author+'"}'),
        contentType: "application/x-www-form-urlencoded",
        success: function(responseData, textStatus, jqXHR) {
            if(responseData == "0") {
                alert("Succesfully deleted.");
                location.href="/board.do";
            } else {
                alert(responseData);
            }
        },
        error: function(jqXHR, textStatus, errorThrown) {
            console.log(errorThrown);
        }
    });
}

