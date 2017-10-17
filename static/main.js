var sticky = document.querySelector('.sticky');
var origOffsetY = sticky.offsetTop;
var d = new Date();
var fd = [d.getFullYear(), d.getMonth()+1, d.getDate()]
if(fd[1].length==1){
    fd[1]="0"+fd[1]
}
window.scrollY >= origOffsetY ? sticky.classList.add('fixed') :
                                sticky.classList.remove('fixed');

function onScroll(e) {
  window.scrollY >= origOffsetY ? sticky.classList.add('fixed') :
                                  sticky.classList.remove('fixed');
}

document.addEventListener('scroll', onScroll);

window.onload=function(){
    $("#delete_box").attr("disabled", true)
    $("[data_date]").each(function(){
        var date = $(this).attr("data_date").match(/(\d+)/g);
        date.length = 3
        $(this).text(date.join("/"))
    });
    overallResults()
};

$("tr td[delete]").dblclick(deleteRecord);
$("tr").click(trClickHandler)

function trClickHandler(){
    if($("#id").attr("value") != $(this).attr("data_id")){
        // $(this).toggleClass("");
        $("form")[0].action="/edit/";
        var id = $(this).attr("data_id");
        var amount = $(this).children("[data_amount]").text();
        var comment = $(this).children("[data_comment]").text();
        var date = $(this).children("[data_date]").attr("data_date").match(/(\d+)/g);
        date.length=3
        $("#id").attr("value", id);
        $("#amount").attr("value", amount);
        $("#comment").attr("value", comment);
        $("#date").attr("value", date.join("-"))
        $("#submit").attr("value", "Изменить")
        $("#submit").removeClass("btn-outline-success")
        $("#submit").addClass("btn-outline-primary")
        // $("#delete_box").removeAttr("disabled")
    } else {
        clearForm()
    }

};
// $("*").click(clickHandler)
// function clickHandler(){
//     // if($("#delete_box").prop("checked")){
//     //     $("#submit").removeClass("btn-outline-primary")
//     //     $("#submit").addClass("btn-outline-danger")
//     //     $("#submit").attr("value", "Удалить")
//     //     $("form")[0].action="/delete/";
//     // } else
//     if ($("#id").attr("value")==0 || !$("#id").attr("value")) {
//         // $("#submit").removeClass("btn-outline-danger")
//         $("#submit").addClass("btn-outline-success")
//         $("#submit").attr("value", "Создать")
//         $("form")[0].action="/create/";
//     } else {
//         // $("#submit").removeClass("btn-outline-danger")
//         $("#submit").addClass("btn-outline-primary")
//         $("#submit").attr("value", "Изменить")
//         $("form")[0].action="/edit/";
//     }
// };
function clearForm() {
    $("form")[0].action="/create/";
    $("#id").attr("value", 0);
    $("#amount").removeAttr("value");
    $("#comment").removeAttr("value");
    $("#date").attr("value", fd.join("-"));
    $("#submit").attr("value", "Создать");
    $("#submit").removeClass("btn-outline-primary");
    $("#submit").addClass("btn-outline-success");
    $("#delete_box").prop("checked", false);
    $("#delete_box").attr("disabled", true);
};
function overallResults() {
    var sum=0
    var income=0
    var outcome=0
    $("[data_amount]").each(function(){
        i = parseInt(this.innerText)
        if(i > 0){
            income += i
        } else if (i < 0) {
            outcome += i
        }
    })
    sum = income + outcome
    $("#nav_income").text("Доход: " + income)
    $("#nav_outcome").text("Расход: " + (-outcome))
    $("#nav_balance").text("Итог: " + sum)
}
function deleteRecord() {
    var parent = $(this).parent()
    parent.hide()
    var xhttp = new XMLHttpRequest();
    var id = $(this).parent().attr("data_id");
    xhttp.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
            parent.remove();
            overallResults();
        }
    };
    xhttp.open("POST", "delete/", true);
    xhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
    xhttp.send("id=" + id);
}
function getRecords(limit, offset){
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
            updateTable(this.responseText)
        }
    };
    xhttp.open("GET", "get/?limit=" + limit + "&offset=" + offset, true);
    xhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
    xhttp.send();
}

function updateTable(json) {
    if(json != "null"){
        var rec = JSON.parse(json)
        rec.forEach(function(e) {
            var date = e.date.match(/(\d+)/g);
            date.length = 3
            date = date.join("/")
            var node = `<tr data_id="`+e.id+`">
              <th scope="row">`+e.id+`</th>
              <td data_amount>`+e.amount+`</td>
              <td data_comment>`+e.comment+`</td>
              <td data_date=`+e.date+`>` + date + `</td>
              <td delete class="btn-outline-danger">X</td>
            </tr>`
            $("#table tbody").append(node);
            $("#table tr:last-child").click(trClickHandler);
            $("#table tr:last-child td[delete]").dblclick(deleteRecord);
        });


        overallResults();

    } else {
        $("#more_records").attr("disabled", true);
        $("#more_records").text("Больше нет записей");
    }
}

function getMoreRecords(){
    getRecords(1, $("#table tbody tr").length)
}

$("#more_records").click(getMoreRecords)
